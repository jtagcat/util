package scrape

import (
	"context"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
)

// sel: goquery selector
func (s *Scraper) Get(urlS, sel string) (_ *goquery.Selection, newURL string, _ error) { // by func(*chromedp.Selector)
	if err := s.limitLock(); err != nil {
		return nil, "", err
	}
	defer s.limitUnlock()

	ctx, cancel := chromedp.NewContext(s.Ctx)
	defer cancel()

	if s.Timeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, s.Timeout)
		defer cancel()
	}

	var gotHtml string
	actions := []chromedp.Action{
		chromedp.Navigate(urlS),
		chromedp.WaitReady(":root"),
		chromedp.OuterHTML(":root", &gotHtml),
		// chromedp.InnerHTML("document", &gotHtml, chromedp.ByJSPath),
		// chromedp.Nodes(sel, &nodes, by),
		chromedp.Location(&newURL),
	}
	if len(s.Cookies) > 0 {
		actions = append(
			[]chromedp.Action{network.SetCookies(s.Cookies)},
			actions...)
	}

	if err := chromedp.Run(ctx, actions...); err != nil {
		return nil, "", err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(gotHtml))
	if err != nil {
		return nil, "", err
	}

	// return nodes, nil
	return doc.Find(sel), newURL, nil
}

// based on https://github.com/chromedp/examples/blob/3384adb2158f6df7e6a48458875a3a5f24aea0c3/download_file/main.go
// timeout: 0 to disable
func (s *Scraper) DownloadFile(urlS, outdir string) (suggested, filename, newURL string, _ error) {
	// block until we take a spot in both queues or parent ctx cancelled
	if err := s.limitLock(); err != nil {
		return "", "", "", err
	}
	defer s.limitUnlock()

	select {
	case <-s.Ctx.Done():
		return "", "", "", context.Canceled
	case s.downloadsLimit <- false: // false is placeholder

	}
	defer func() {
		<-s.downloadsLimit // leave queue
	}()

	ctx, cancel := chromedp.NewContext(s.Ctx)
	defer cancel()

	if s.Timeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, s.Timeout)
		defer cancel()
	}

	done := make(chan string, 1)

	//## handle download event
	var requestID network.RequestID
	chromedp.ListenTarget(ctx, func(v interface{}) {
		switch ev := v.(type) {
		// opt a: browser renders (png)
		case *network.EventRequestWillBeSent:
			if ev.Request.URL == urlS {
				requestID = ev.RequestID
			}
		case *network.EventLoadingFinished:
			if ev.RequestID == requestID {
				close(done)
			}

		// opt b: direct download
		case *browser.EventDownloadWillBegin:
			suggested = ev.SuggestedFilename

		case *browser.EventDownloadProgress:
			if ev.State == browser.DownloadProgressStateCompleted {
				done <- ev.GUID
				close(done)
			}

		}
	})

	//## direct chrome interaction
	actions := []chromedp.Action{
		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllowAndName).
			WithDownloadPath(outdir).
			WithEventsEnabled(true),
		chromedp.Navigate(urlS),
		chromedp.Location(&newURL),
	}
	if len(s.Cookies) > 0 {
		actions = append(
			[]chromedp.Action{network.SetCookies(s.Cookies)},
			actions...)
	}

	if err := chromedp.Run(ctx, actions...); err != nil && !strings.Contains(err.Error(), "net::ERR_ABORTED") {
		// Upstream note: Ignoring the net::ERR_ABORTED page error is essential here
		// since downloads will cause this error to be emitted, although the
		// download will still succeed.
		return "", "", "", err
	}

	guid := <-done // blocks

	//## opt a: if browser rendered and is not direct donwload (eg. png)
	if guid == "" {
		// emulate a guid manually
		var guidPath string
		for {
			guid = uuid.New().String()
			guidPath = path.Join(outdir, guid)
			if _, err := os.Stat(guidPath); err != nil && !errors.Is(err, os.ErrNotExist) {
				return "", "", "", err
			} else if errors.Is(err, os.ErrNotExist) {
				break
			}
		}

		// get the downloaded bytes by request id
		var buf []byte
		if err := chromedp.Run(ctx,
			chromedp.ActionFunc(func(ctx context.Context) (err error) {
				buf, err = network.GetResponseBody(requestID).Do(ctx)
				return err
			})); err != nil {
			return "", "", "", err
		}

		if err := ioutil.WriteFile(guidPath, buf, os.ModePerm); err != nil {
			return "", "", "", err
		}
	}

	return suggested, guid, newURL, nil
}

func (s *Scraper) GetRaw(urlS string) ([]byte, error) {
	return s.DoRaw(urlS, "GET", nil)
}

func (s *Scraper) DoRaw(urlS, method string, data []byte) (body []byte, _ error) {
	if err := s.limitLock(); err != nil {
		return nil, err
	}
	defer s.limitUnlock()

	ctx, cancel := chromedp.NewContext(s.Ctx)
	defer cancel()

	if s.Timeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, s.Timeout)
		defer cancel()
	}

	done := make(chan string, 1)

	var requestID network.RequestID
	chromedp.ListenTarget(ctx, func(v interface{}) {
		switch ev := v.(type) {
		case *network.EventResponseReceived:
			requestID = ev.RequestID
			close(done)
		case *fetch.EventRequestPaused:
			go func() {
				fctx := cdp.WithExecutor(ctx, chromedp.FromContext(ctx).Target)
				f := fetch.ContinueRequest(ev.RequestID).WithMethod(method)
				switch method {
				case http.MethodPost, http.MethodPut, http.MethodDelete:
					if data != nil {
						f = f.WithPostData(base64.StdEncoding.EncodeToString(data))
					}
				}

				f.Do(fctx)
			}()
		}
	})

	actions := []chromedp.Action{
		network.Enable(),
		fetch.Enable(),
		chromedp.Navigate(urlS),
	}
	if len(s.Cookies) > 0 {
		actions = append(
			[]chromedp.Action{network.SetCookies(s.Cookies)},
			actions...)
	}

	if err := chromedp.Run(ctx, actions...); err != nil {
		return nil, err
	}

	<-done
	// get downloaded bytes for request id
	err := chromedp.Run(ctx,

		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			body, err = network.GetResponseBody(requestID).Do(ctx)
			return err
		}),
	)
	return body, err
}

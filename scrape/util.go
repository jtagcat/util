package scrape

import (
	"github.com/PuerkitoBio/goquery"
)

// useful with a for-range
func RawEach(s *goquery.Selection) (a []*goquery.Selection) {
	s.Each(func(_ int, s *goquery.Selection) {
		a = append(a, s)
	})
	return a
}

func RecursiveChildFiltered(s *goquery.Selection, filters ...string) *goquery.Selection {
	for _, f := range filters {
		s = s.ChildrenFiltered(f)
	}
	return s
}

// if only we could have `type Node cdp.Node` to use `func (n *Node)`
//
// if https://github.com/chromedp/cdproto/issues/20 is implemented, this func is deprecated

// Converts cdp.Node to goquery and filters children.
//
// panics: valid HTML almost always remains valid HTML
// func CdpFilterChildren(node *cdp.Node, sel string) *goquery.Selection {
// 	doc, err := goquery.NewDocumentFromReader(strings.NewReader(
// 		node.Dump("", "", false)))
// 	if err != nil {
// 		panic(fmt.Errorf("converting cdp.Node (%v) to goquery Doc: %e", node, err))
// 	}

// 	return doc.Selection.ChildrenFiltered(sel)
// }

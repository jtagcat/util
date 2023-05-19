module github.com/jtagcat/util

go 1.18

require (
	github.com/PuerkitoBio/goquery v1.8.0
	github.com/chromedp/cdproto v0.0.0-20230213000208-1903a0cd6c4c
	github.com/chromedp/chromedp v0.8.7
	github.com/fsnotify/fsnotify v1.6.0
	github.com/google/uuid v1.3.0
	github.com/stretchr/testify v1.8.1
	golang.org/x/exp v0.0.0-20230519143937-03e91628a987
	golang.org/x/sync v0.1.0
	k8s.io/apimachinery v0.26.1
)

require (
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/chromedp/sysutil v1.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.1.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/klog/v2 v2.90.0 // indirect
	k8s.io/utils v0.0.0-20230209194617-a36077c30491 // indirect
)

// https://github.com/kubernetes/kubernetes/pull/113398
// go get github.com/jtagcat/kubernetes/staging/src/k8s.io/apimachinery@ManagedExponentialBackoff
// when updating this, update the reference in retry package
replace k8s.io/apimachinery => github.com/jtagcat/kubernetes/staging/src/k8s.io/apimachinery v0.0.0-20221027124836-581f57977fff

module github.com/jtagcat/util

go 1.26.10

require (
	github.com/PuerkitoBio/goquery v1.13.0
	github.com/chromedp/cdproto v0.159.0
	github.com/chromedp/chromedp v0.21.0
	github.com/fsnotify/fsnotify v1.10.1
	github.com/gin-gonic/gin v1.12.0
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.12.1
	go.etcd.io/bbolt v1.5.0
	golang.org/x/sync v0.24.0
	k8s.io/apimachinery v0.37.2
)

require (
	github.com/andybalholm/cascadia v1.3.5 // indirect
	github.com/bytedance/gopkg v0.1.4 // indirect
	github.com/bytedance/sonic v1.15.5 // indirect
	github.com/bytedance/sonic/loader v0.5.3 // indirect
	github.com/chromedp/sysutil v1.1.0 // indirect
	github.com/cloudwego/base64x v0.1.7 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.16 // indirect
	github.com/gin-contrib/sse v1.2.0 // indirect
	github.com/go-json-experiment/json v0.0.0-20260820222146-c27c302e5fc3 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-playground/locales v0.14.2 // indirect
	github.com/go-playground/universal-translator v0.18.2 // indirect
	github.com/go-playground/validator/v10 v10.30.5 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.4.0 // indirect
	github.com/goccy/go-json v0.12.0 // indirect
	github.com/goccy/go-yaml v1.19.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.4.1 // indirect
	github.com/leodido/go-urn v1.5.0 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.4.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/quic-go/qpack v0.6.0 // indirect
	github.com/quic-go/quic-go v0.64.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.3.2 // indirect
	go.mongodb.org/mongo-driver/v2 v2.9.2 // indirect
	go.uber.org/mock v0.6.0 // indirect
	golang.org/x/arch v0.32.0 // indirect
	golang.org/x/crypto v0.58.0 // indirect
	golang.org/x/mod v0.41.0 // indirect
	golang.org/x/net v0.62.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	golang.org/x/tools v0.51.0 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/klog/v2 v2.140.0 // indirect
	k8s.io/utils v0.0.0-20261023123148-ca04ba58ecab // indirect
)

// https://github.com/kubernetes/kubernetes/pull/113398
// go get github.com/jtagcat/kubernetes/staging/src/k8s.io/apimachinery@ManagedExponentialBackoff
// when updating this, update the reference in retry package
replace k8s.io/apimachinery => github.com/jtagcat/kubernetes/staging/src/k8s.io/apimachinery v0.0.0-20221027124836-581f57977fff

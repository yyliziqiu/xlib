package xhttp

import (
	"net/http"
	"time"

	"github.com/yyliziqiu/xlib/xutil"
)

type Api struct {
	BaseURL  string
	TextHTTP *TextHTTP
	JsonHTTP *JSONHTTP
}

func (a Api) URL(path string) string {
	return xutil.JoinURL(a.BaseURL, path)
}

func newHTTPClientWithTimeout(d time.Duration) *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}

func GetBasicAuthRequestFunc(username string, password string) func(req *http.Request) {
	return func(req *http.Request) {
		req.SetBasicAuth(username, password)
	}
}

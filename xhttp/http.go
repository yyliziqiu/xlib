package xhttp

import (
	"net/http"
	"time"

	"github.com/yyliziqiu/xlib/xutil"
)

type API struct {
	BaseURL  string
	TextHTTP *TextHTTP
	JSONHTTP *JSONHTTP
}

func (a API) URL(path string) string {
	return xutil.JoinURL(a.BaseURL, path)
}

func NewHTTPClientWithTimeout(d time.Duration) *http.Client {
	return &http.Client{
		Timeout: d,
	}
}

func BasicAuthRequestFunc(username string, password string) func(req *http.Request) {
	return func(req *http.Request) {
		req.SetBasicAuth(username, password)
	}
}

package xhttp

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/yyliziqiu/xlib/xutil"
)

type API struct {
	BaseURL string
	Text    *Text
	JSON    *JSON
}

func (a API) URL(path string) string {
	return xutil.JoinURL(a.BaseURL, path)
}

func AppendQuery(rawURL string, query url.Values) (string, error) {
	uo, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	for k, v := range uo.Query() {
		for _, s := range v {
			query.Add(k, s)
		}
	}
	uo.RawQuery = query.Encode()

	return uo.String(), nil
}

func newClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}

func newRequest(method string, url string, query url.Values, header http.Header, body io.Reader, cb func(req *http.Request)) (*http.Request, error) {
	url, err := AppendQuery(url, query)
	if err != nil {
		return nil, fmt.Errorf("append query failed [%v]", err)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("new request failed [%v]", err)
	}

	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	cb(req)

	return req, nil
}

func basicAuth(username string, password string) func(req *http.Request) {
	return func(req *http.Request) {
		req.SetBasicAuth(username, password)
	}
}

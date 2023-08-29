package xhttp

import (
	"net/http"
	"net/url"

	"github.com/yyliziqiu/xlib/xutil"
)

type Api struct {
	Domain string
}

func (a Api) Url(path string) string {
	return xutil.JoinURL(a.Domain, path)
}

func (a Api) Get(path string, header http.Header, ao url.Values, bo interface{}) error {
	return Get2(a.Url(path), header, ao, bo)
}

func (a Api) PostForm(path string, header http.Header, ao url.Values, bo interface{}) error {
	return PostForm2(a.Url(path), header, ao, bo)
}

func (a Api) PostJSON(path string, header http.Header, ao interface{}, bo interface{}) error {
	return PostJSON2(a.Url(path), header, ao, bo)
}

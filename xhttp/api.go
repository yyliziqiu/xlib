package xhttp

import (
	"github.com/yyliziqiu/xlib/xutil"
)

type Api struct {
	BaseURL string
}

func (a Api) URL(path string) string {
	return xutil.JoinURL(a.BaseURL, path)
}

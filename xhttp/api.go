package xhttp

import (
	"github.com/yyliziqiu/xlib/xutil"
)

type Api struct {
	Domain string
}

func (a Api) Url(path string) string {
	return xutil.JoinURL(a.Domain, path)
}

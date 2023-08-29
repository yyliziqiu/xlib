package xapi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yyliziqiu/xlib/xapi/xresponse"
)

type CrosHeaders struct {
	MaxAge           string
	Origin           string
	ExposeHeaders    string
	AllowMethods     string
	AllowHeaders     string
	AllowCredentials string
}

var crosHeaders = &CrosHeaders{
	MaxAge:           "86400",
	Origin:           "*",
	ExposeHeaders:    "",
	AllowMethods:     "OPTIONS, HEAD, GET, POST, PUT, PATCH, DELETE",
	AllowHeaders:     "*",
	AllowCredentials: "false",
}

// CrosMiddleware
// 允许跨域，参考： https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers
// Only the 7 CORS-safelisted response headers are exposed: Cache-Control, Content-Language, Content-Length, Content-Type, Expires, Last-Modified, Pragma.
// CORS-safelisted request header is one of the following HTTP headers: Accept, Accept-Language, Content-Language, Content-Type.
func CrosMiddleware(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", crosHeaders.Origin)
	c.Header("Access-Control-Expose-Headers", crosHeaders.ExposeHeaders)
	c.Header("Access-Control-Allow-Credentials", crosHeaders.AllowCredentials)
	if c.Request.Method == http.MethodOptions {
		c.Header("Access-Control-Max-Age", crosHeaders.MaxAge)
		c.Header("Access-Control-Allow-Methods", crosHeaders.AllowMethods)
		c.Header("Access-Control-Allow-Headers", crosHeaders.AllowHeaders)
	}

	if c.Request.Method == http.MethodOptions {
		xresponse.Ok(c)
	}
}

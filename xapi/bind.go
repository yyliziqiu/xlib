package xapi

import (
	"github.com/gin-gonic/gin"

	"github.com/yyliziqiu/xlib/xapi/xresponse"
	"github.com/yyliziqiu/xlib/xerror"
)

var ParamError = xerror.NewError("A1000", "param error")

// BindForm 参数发生错误时，响应简单的错误信息
func BindForm(ctx *gin.Context, form interface{}) bool {
	err := ctx.ShouldBind(form)
	if err != nil {
		errorLogger.Warnf("Bind error, path: %s, form: %v, err: %v.", ctx.FullPath(), form, err)
		xresponse.Fail(ctx, ParamError)
		return false
	}
	return true
}

// BindFormVerbose 参数发生错误时，响应详细的错误信息
func BindFormVerbose(ctx *gin.Context, form interface{}) bool {
	err := ctx.ShouldBind(form)
	if err != nil {
		errorLogger.Warnf("Bind error, path: %s, form: %v, err: %v.", ctx.FullPath(), form, err)
		xresponse.Fail(ctx, ParamError.With(err))
		return false
	}
	return true
}

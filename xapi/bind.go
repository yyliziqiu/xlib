package xapi

import (
	"github.com/gin-gonic/gin"

	"github.com/yyliziqiu/xlib/xapi/xresponse"
	"github.com/yyliziqiu/xlib/xerror"
)

var ParamError = xerror.New("A1000", "param error")

func BindForm(ctx *gin.Context, form interface{}, verbose bool) bool {
	err := ctx.ShouldBind(form)
	if err != nil {
		if errorLogger != nil {
			errorLogger.Warnf("Bind failed, path: %s, form: %v, error: %v.", ctx.FullPath(), form, err)
		}
		if verbose {
			xresponse.Error(ctx, ParamError.Wrap(err))
		} else {
			xresponse.Error(ctx, ParamError)
		}
		return false
	}
	return true
}

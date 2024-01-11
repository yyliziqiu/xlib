package xresponse

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yyliziqiu/xlib/xerror"
)

var (
	BadRequestError          = xerror.New("A0400", "Bad Request")
	UnauthorizedError        = xerror.New("A0401", "Unauthorized")
	ForbiddenError           = xerror.New("A0403", "Forbidden")
	NotFoundError            = xerror.New("A0404", "Not Found")
	MethodNotAllowedError    = xerror.New("A0405", "Method Not Allowed")
	InternalServerErrorError = xerror.New("B0500", "Internal Server Error")
)

type ErrorResult struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewErrorResult(code string, message string) ErrorResult {
	return ErrorResult{
		Code:    code,
		Message: message,
	}
}

func NewErrorResultWithError(err error) ErrorResult {
	xerr, ok := err.(*xerror.Error)
	if ok {
		return NewErrorResult(xerr.Code, xerr.Message)
	}
	return NewErrorResult(BadRequestError.Code, err.Error())
}

func buildErrorResponse(err error) (int, ErrorResult) {
	var (
		statusCode = http.StatusBadRequest
		code       = BadRequestError.Code
		message    = err.Error()
	)

	xerr, ok := err.(*xerror.Error)
	if ok {
		if xerr.Code[0] != 'A' {
			statusCode = http.StatusInternalServerError
		}
		code = xerr.Code
		message = xerr.Message
	}

	return statusCode, NewErrorResult(code, message)
}

func Response(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.JSON(statusCode, data)
}

func ResponseError(ctx *gin.Context, statusCode int, code string, message string) {
	ctx.JSON(statusCode, NewErrorResult(code, message))
}

func OK(ctx *gin.Context) {
	ctx.String(http.StatusOK, "")
}

func Result(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}

func Error(ctx *gin.Context, err error) {
	statusCode, result := buildErrorResponse(err)
	ctx.JSON(statusCode, result)
}

func ErrorString(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, NewErrorResult(BadRequestError.Code, message))
}

func AbortOK(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusOK)
}

func AbortResult(ctx *gin.Context, data interface{}) {
	ctx.AbortWithStatusJSON(http.StatusOK, data)
}

func AbortError(ctx *gin.Context, err error) {
	statusCode, result := buildErrorResponse(err)
	ctx.AbortWithStatusJSON(statusCode, result)
}

func AbortErrorString(ctx *gin.Context, message string) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResult(BadRequestError.Code, message))
}

func AbortBadRequest(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, NewErrorResultWithError(BadRequestError))
}

func AbortUnauthorized(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, NewErrorResultWithError(UnauthorizedError))
}

func AbortForbidden(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusForbidden, NewErrorResultWithError(ForbiddenError))
}

func AbortNotFound(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusNotFound, NewErrorResultWithError(NotFoundError))
}

func AbortMethodNotAllowed(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, NewErrorResultWithError(MethodNotAllowedError))
}

func AbortInternalServerError(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewErrorResultWithError(InternalServerErrorError))
}

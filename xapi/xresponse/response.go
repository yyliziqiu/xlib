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

type responseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func newResponseError(err *xerror.Error) responseError {
	return responseError{
		Code:    err.Code,
		Message: err.Message,
	}
}

func buildErrorResponse(err error) (int, responseError) {
	e, ok := err.(*xerror.Error)
	if !ok {
		e = BadRequestError
	}

	statusCode := http.StatusBadRequest
	if e.Code[0] != 'A' {
		statusCode = http.StatusInternalServerError
	}

	return statusCode, newResponseError(e)
}

func OK(ctx *gin.Context) {
	ctx.String(http.StatusOK, "")
}

func Result(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}

func Error(ctx *gin.Context, err error) {
	statusCode, responseErr := buildErrorResponse(err)
	ctx.JSON(statusCode, responseErr)
}

func ErrorString(ctx *gin.Context, message string) {
	statusCode, responseErr := buildErrorResponse(BadRequestError.With(message))
	ctx.JSON(statusCode, responseErr)
}

func AbortOK(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusOK)
}

func AbortResult(ctx *gin.Context, data interface{}) {
	ctx.AbortWithStatusJSON(http.StatusOK, data)
}

func AbortError(ctx *gin.Context, err error) {
	statusCode, responseErr := buildErrorResponse(err)
	ctx.AbortWithStatusJSON(statusCode, responseErr)
}

func AbortErrorString(ctx *gin.Context, message string) {
	statusCode, responseErr := buildErrorResponse(BadRequestError.With(message))
	ctx.AbortWithStatusJSON(statusCode, responseErr)
}

func AbortBadRequest(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, newResponseError(BadRequestError))
}

func AbortUnauthorized(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, newResponseError(UnauthorizedError))
}

func AbortForbidden(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusForbidden, newResponseError(ForbiddenError))
}

func AbortNotFound(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusNotFound, newResponseError(NotFoundError))
}

func AbortMethodNotAllowed(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, newResponseError(MethodNotAllowedError))
}

func AbortInternalServerError(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, newResponseError(InternalServerErrorError))
}

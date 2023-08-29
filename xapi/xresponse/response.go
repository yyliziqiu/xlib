package xresponse

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yyliziqiu/xlib/xerror"
)

var (
	BadRequestError          = &xerror.Error{Code: "A0400", Message: "Bad Request"}
	UnauthorizedError        = &xerror.Error{Code: "A0401", Message: "Unauthorized"}
	ForbiddenError           = &xerror.Error{Code: "A0403", Message: "Forbidden"}
	NotFoundError            = &xerror.Error{Code: "A0404", Message: "Not Found"}
	MethodNotAllowedError    = &xerror.Error{Code: "A0405", Message: "Method Not Allowed"}
	InternalServerErrorError = &xerror.Error{Code: "B0500", Message: "Internal Server Error"}
)

func newErrorResponse(err *xerror.Error) errorResponse {
	return errorResponse{Code: err.Code, Message: err.Message}
}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func buildErrorResponse(err error) (statusCode int, res errorResponse) {
	e, ok := err.(*xerror.Error)
	if !ok {
		e = BadRequestError
	}

	switch e.Code[0] {
	case 'A':
		statusCode = http.StatusBadRequest
	case 'B':
		statusCode = http.StatusInternalServerError
	case 'C', 'D':
		statusCode = http.StatusServiceUnavailable
	default:
		statusCode = http.StatusBadRequest
	}

	return statusCode, newErrorResponse(e)
}

func Ok(c *gin.Context) {
	c.String(http.StatusOK, "")
}

func Result(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func Abort(c *gin.Context, err error) {
	sc, er := buildErrorResponse(err)
	c.AbortWithStatusJSON(sc, er)
}

func AbortString(c *gin.Context, message string) {
	sc, er := buildErrorResponse(BadRequestError.With(message))
	c.AbortWithStatusJSON(sc, er)
}

func Fail(c *gin.Context, err error) {
	sc, er := buildErrorResponse(err)
	c.JSON(sc, er)
}

func FailString(c *gin.Context, message string) {
	sc, er := buildErrorResponse(BadRequestError.With(message))
	c.JSON(sc, er)
}

func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, newErrorResponse(UnauthorizedError))
}

func Forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, newErrorResponse(ForbiddenError))
}

func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, newErrorResponse(NotFoundError))
}

func MethodNotAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, newErrorResponse(MethodNotAllowedError))
}

func InternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, newErrorResponse(InternalServerErrorError))
}

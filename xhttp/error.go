package xhttp

import (
	"fmt"
)

type ResponseError struct {
	StatusCode int
	Body       string
}

func NewResponseError(statusCode int, body string) *ResponseError {
	return &ResponseError{
		StatusCode: statusCode,
		Body:       body,
	}
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("status code [%d], body [%s]", e.StatusCode, e.Body)
}

func IsResponseError(err error) bool {
	_, ok := err.(*ResponseError)
	return ok
}

package xhttp

import (
	"fmt"
)

type ResponseError struct {
	StatusCode int
	Body       string
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("status code: %d, body: %s", e.StatusCode, e.Body)
}

func NewResponseError(statusBody int, body string) *ResponseError {
	return &ResponseError{
		StatusCode: statusBody,
		Body:       body,
	}
}

func IsResponseError(err error) bool {
	_, ok := err.(*ResponseError)
	return ok
}

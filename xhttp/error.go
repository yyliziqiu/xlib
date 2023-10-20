package xhttp

import (
	"fmt"
)

type StatusError struct {
	Code    int
	Message string
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("status code: %d, message: %s", e.Code, e.Message)
}

func NewStatusError(code int, message string) *StatusError {
	return &StatusError{
		Code:    code,
		Message: message,
	}
}

func IsStatusError(err error) bool {
	_, ok := err.(*StatusError)
	return ok
}

package xhttp

import (
	"fmt"
)

type ResponseError struct {
	StatusCode int
	BodyString string
	BodyStruct interface{}
}

func newResponseError(statusCode int, bodyString string, bodyStruct interface{}) *ResponseError {
	return &ResponseError{
		StatusCode: statusCode,
		BodyString: bodyString,
		BodyStruct: bodyStruct,
	}
}

func (e ResponseError) Error() string {
	return fmt.Sprintf("status code [%d], body [%s]", e.StatusCode, e.BodyString)
}

func IsResponseError(err error) bool {
	_, ok := err.(*ResponseError)
	return ok
}

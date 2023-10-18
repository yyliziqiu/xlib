package xerror

import "fmt"

func NewError(code string, message string) *Error {
	return &Error{Code: code, Message: message}
}

type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %s, message: %s", e.Code, e.Message)
}

func (e *Error) With(v interface{}) *Error {
	err := &Error{Code: e.Code}
	switch v.(type) {
	case error:
		err.Message = v.(error).Error()
	case string:
		err.Message = v.(string)
	default:
		err.Message = fmt.Sprintf("%v", v)
	}
	return err
}

func (e *Error) Withf(message string, a ...interface{}) *Error {
	return &Error{
		Code:    e.Code,
		Message: fmt.Sprintf(message, a...),
	}
}

func (e *Error) WithFields(a ...interface{}) *Error {
	return &Error{
		Code:    e.Code,
		Message: fmt.Sprintf(e.Message, a...),
	}
}

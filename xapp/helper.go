package xapp

import (
	"reflect"
)

func GetFieldValue(s any, fn string) (any, bool) {
	val := reflect.ValueOf(s)

	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	field := val.FieldByName(fn)
	if !field.IsValid() {
		return nil, false
	}

	return field.Interface(), true
}

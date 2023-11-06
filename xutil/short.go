package xutil

import (
	"golang.org/x/exp/constraints"
)

func If[T any](exp bool, a, b T) T {
	if exp {
		return a
	}
	return b
}

func IfNil(a, b any) any {
	if a == nil {
		return b
	}
	return a
}

func IfZero[T constraints.Integer | constraints.Float](a, b T) T {
	if a == 0 {
		return b
	}
	return a
}

func IfEmpty[T ~string | []any | map[any]any](a, b T) T {
	if len(a) == 0 {
		return b
	}
	return a
}

package xstring

import (
	"strings"
)

func Truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func Empty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

package xutil

import (
	"strings"
)

func JoinURL(segments ...string) string {
	if len(segments) == 0 {
		return ""
	}

	n := len(segments)
	sb := strings.Builder{}
	for i, segment := range segments {
		if i < n-1 {
			sb.WriteString(strings.Trim(segment, "/"))
			sb.WriteString("/")
		} else {
			sb.WriteString(strings.TrimLeft(segment, "/"))
		}
	}

	return sb.String()
}

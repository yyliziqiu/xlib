package xutil

import (
	"strings"
)

func JoinURL(segments ...string) string {
	if len(segments) == 0 {
		return ""
	}

	sb := strings.Builder{}
	for _, segment := range segments {
		sb.WriteString(strings.Trim(segment, "/"))
		sb.WriteString("/")
	}

	url := sb.String()
	if !strings.HasSuffix(segments[len(segments)-1], "/") {
		url = url[:len(url)-1]
	}

	return url
}

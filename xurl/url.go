package xurl

import (
	"strings"
)

func Join(segments ...string) string {
	if len(segments) == 0 {
		return ""
	}

	n := len(segments)
	sb := strings.Builder{}
	for i, segment := range segments {
		if segment == "" {
			continue
		} else if i == n-1 {
			sb.WriteString(strings.TrimLeft(segment, "/"))
		} else {
			sb.WriteString(strings.Trim(segment, "/"))
			sb.WriteString("/")
		}
	}

	return sb.String()
}

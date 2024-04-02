package xhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type HTTPError struct {
	status int
	errstr string
}

func newHTTPError(status int, errstr string) *HTTPError {
	return &HTTPError{
		status: status,
		errstr: errstr,
	}
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("status code [%d], message [%s]", e.status, e.errstr)
}

type JsonResponse interface {
	Failed() bool
}

func JoinURL(segments ...string) string {
	if len(segments) == 0 {
		return ""
	}

	url2 := segments[0]
	for i, segment := range segments {
		if i == 0 || segment == "" {
			continue
		} else {
			l := strings.HasSuffix(url2, "/")
			r := strings.HasPrefix(segment, "/")
			if l && r {
				url2 += segment[1:]
			} else if l || r {
				url2 += segment
			} else {
				url2 += "/" + segment
			}
		}
	}

	return url2
}

func AppendQuery(rawURL string, query url.Values) (string, error) {
	if len(query) == 0 {
		return rawURL, nil
	}

	uo, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	for k, v := range uo.Query() {
		for _, s := range v {
			query.Add(k, s)
		}
	}

	uo.RawQuery = query.Encode()

	return uo.String(), nil
}

func HeaderToString(header http.Header) string {
	if len(header) == 0 {
		return "{}"
	}

	m := make(map[string]string, len(header))
	for key := range header {
		m[key] = header.Get(key)
	}

	bs, _ := json.Marshal(m)

	return string(bs)
}

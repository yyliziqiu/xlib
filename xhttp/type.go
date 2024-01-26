package xhttp

import (
	"fmt"
	"net/http"
	"net/url"
)

type JsonResponse interface {
	IsError() bool
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

func H2S(header http.Header) string {
	if len(header) == 0 {
		return "[]"
	}
	str := ""
	for key := range header {
		value := header.Get(key)
		str += fmt.Sprintf("%s: %s; ", key, value)
	}
	return fmt.Sprintf("[%s]", str)
}

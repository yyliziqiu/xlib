package xhttp

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Option func(api *API)

func WithClient(client *http.Client) Option {
	return func(api *API) {
		api.client = client
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(api *API) {
		api.client = &http.Client{Timeout: timeout}
	}
}

func WithBaseURL(baseURL string) Option {
	return func(api *API) {
		api.baseURL = baseURL
	}
}

func WithResType(typ string) Option {
	return func(api *API) {
		api.ResType = typ
	}
}

func WithErrorFunc(f func(interface{}) bool) Option {
	return func(api *API) {
		api.errorFunc = f
	}
}

func WithErrorStruct(errorStruct interface{}) Option {
	return func(api *API) {
		api.errorStruct = errorStruct
	}
}

func WithRequestFunc(f func(r *http.Request)) Option {
	return func(api *API) {
		api.requestFunc = f
	}
}

func WithBasicAuth(username string, password string) Option {
	return func(api *API) {
		api.requestFunc = func(req *http.Request) {
			req.SetBasicAuth(username, password)
		}
	}
}

func WithLogger(logger *logrus.Logger) Option {
	return func(api *API) {
		api.logger = logger
	}
}

func WithMaxLogLength(n int) Option {
	return func(api *API) {
		api.maxLogLength = n
	}
}

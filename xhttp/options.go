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

func WithFormat(format string) Option {
	return func(api *API) {
		api.format = format
	}
}

func WithBaseURL(baseURL string) Option {
	return func(api *API) {
		api.baseURL = baseURL
	}
}

func WithError(error interface{}) Option {
	return func(api *API) {
		api.error = error
	}
}

func WithDumps(enabled bool) Option {
	return func(api *API) {
		api.dumps = enabled
	}
}

func WithLogger(logger *logrus.Logger) Option {
	return func(api *API) {
		api.logger = logger
	}
}

func WithLogLength(n int) Option {
	return func(api *API) {
		api.logLength = n
	}
}

func WithLogEscape(enabled bool) Option {
	return func(api *API) {
		api.logEscape = enabled
	}
}

func WithRequestBefore(f func(r *http.Request)) Option {
	return func(api *API) {
		api.requestBefore = f
	}
}

func WithBasicAuth(username string, password string) Option {
	return func(api *API) {
		api.requestBefore = func(req *http.Request) {
			req.SetBasicAuth(username, password)
		}
	}
}

func WithBearerToken(token string) Option {
	return func(api *API) {
		api.requestBefore = func(req *http.Request) {
			req.Header.Set("Authorization", "Bearer "+token)
		}
	}
}

func WithResponseAfter(f func(res *http.Response) error) Option {
	return func(api *API) {
		api.responseAfter = f
	}
}

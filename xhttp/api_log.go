package xhttp

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func (api *API) logPrepare(log string) string {
	if api.logLength <= 0 {
		return ""
	}
	if len(log) > api.logLength {
		log = log[:api.logLength]
	}
	// return strings.ReplaceAll(log, "\n", "\\n")
	return log
}

func (api *API) logInfo(format string, args ...interface{}) {
	if api.logger == nil {
		return
	}
	message := api.logPrepare(fmt.Sprintf(format, args...))
	api.logger.Info(message)
}

func (api *API) logWarn(format string, args ...interface{}) {
	if api.logger == nil {
		return
	}
	message := api.logPrepare(fmt.Sprintf(format, args...))
	api.logger.Warn(message)
}

func (api *API) dumpRequest(req *http.Request) {
	if !api.dump {
		return
	}

	bs, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		fmt.Printf("Dump request failed, error: %v\n", err)
		return
	}

	fmt.Println("\n---------- Request ----------")
	fmt.Printf(string(bs))
	fmt.Println("\n---------- Request End----------")
}

func (api *API) dumpResponse(res *http.Response) {
	if !api.dump {
		return
	}

	bs, err := httputil.DumpResponse(res, true)
	if err != nil {
		fmt.Printf("Dump response failed, error: %v", err)
		return
	}

	fmt.Println("\n---------- Response ----------")
	fmt.Printf(string(bs))
	fmt.Println("\n---------- Response End----------")
}

package xhttp

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"

	"github.com/yyliziqiu/xlib/xutil"
)

var logger *logrus.Logger

func SetLogger(lg *logrus.Logger) {
	logger = lg
}

var maxLength = 2000

func SetMaxLength(n int) {
	maxLength = n
}

func Get2(url string, header http.Header, query url.Values, out interface{}) error {
	timer := xutil.NewTimer()

	err := Get(url, header, query, out)

	if err != nil {
		logError(url, header, query, out, err, timer)
		return err
	}

	logInfo(url, header, query, out, timer)

	return nil
}

func PostForm2(url string, header http.Header, in url.Values, out interface{}) error {
	timer := xutil.NewTimer()

	err := PostForm(url, header, in, out)

	if err != nil {
		logError(url, header, in, out, err, timer)
		return err
	}

	logInfo(url, header, in, out, timer)

	return nil
}

func PostJSON2(url string, header http.Header, in interface{}, out interface{}) error {
	timer := xutil.NewTimer()

	err := PostJSON(url, header, in, out)

	if err != nil {
		logError(url, header, in, out, err, timer)
		return err
	}

	logInfo(url, header, in, out, timer)

	return nil
}

func logInfo(url string, header http.Header, in interface{}, out interface{}, timer xutil.Timer) {
	if logger == nil {
		return
	}

	hs, is, os := marshal(header, in, out)

	logger.Infof("Request succeed, url: %s, header: %s, request: %s, response: %s, cost: %s.", url, hs, is, os, timer.Stops())
}

func marshal(header http.Header, in interface{}, out interface{}) (string, string, string) {
	var hs, is, os []byte
	if header != nil {
		hs, _ = json.Marshal(header)
		hs = truncate(hs)
	}
	if in != nil {
		is, _ = json.Marshal(in)
		is = truncate(is)
	}
	if out != nil {
		os, _ = json.Marshal(out)
		os = truncate(os)
	}
	return string(hs), string(is), string(os)
}

func truncate(bs []byte) []byte {
	if len(bs) <= maxLength {
		return bs
	}
	return bs[:maxLength]
}

func logError(url string, header http.Header, in interface{}, out interface{}, err error, timer xutil.Timer) {
	if logger == nil {
		return
	}

	hs, is, os := marshal(header, in, out)

	logger.Infof("Request failed, url: %s, header: %s, request: %s, response: %s, error: %s, cost: %s.", url, hs, is, os, err.Error(), timer.Stops())
}

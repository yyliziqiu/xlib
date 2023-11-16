package xhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/yyliziqiu/xlib/xutil"
)

const (
	ResTypeJSON = "json"
	ResTypeText = "text"
)

type API struct {
	client       *http.Client
	baseURL      string
	ResType      string
	errorFunc    func(interface{}) bool // 当 status code 为 200 时有效，用来判断响应是否成功，出入参数为 out
	errorStruct  interface{}            // 不能是指针
	requestFunc  func(req *http.Request)
	logger       *logrus.Logger
	logLengthMax int
	dump         bool
}

func NewAPI(options ...Option) *API {
	api := &API{
		client:       &http.Client{Timeout: 5 * time.Second},
		baseURL:      "",
		ResType:      ResTypeJSON,
		errorStruct:  nil,
		requestFunc:  nil,
		logger:       nil,
		logLengthMax: 1024,
		dump:         false,
	}

	for _, option := range options {
		option(api)
	}

	return api
}

func (a *API) truncateLog(log string) string {
	if a.logLengthMax <= 0 {
		return ""
	}
	if len(log) <= a.logLengthMax {
		return log
	}
	return log[:a.logLengthMax]
}

func (a *API) logInfo(format string, args ...interface{}) {
	if a.logger == nil {
		return
	}
	log := fmt.Sprintf(format, args...)
	log = a.truncateLog(log)
	a.logger.Info(log)
}

func (a *API) logWarn(format string, args ...interface{}) {
	if a.logger == nil {
		return
	}
	log := fmt.Sprintf(format, args...)
	log = a.truncateLog(log)
	a.logger.Warn(log)
}

func (a *API) dumpRequest(req *http.Request) {
	if !a.dump {
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

func (a *API) dumpResponse(res *http.Response) {
	if !a.dump {
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

func (a *API) newRequest(method string, path string, query url.Values, header http.Header, body io.Reader) (*http.Request, error) {
	rawURL, err := AppendQuery(a.url(path), query)
	if err != nil {
		a.logWarn("Append query failed, URL: %s, Query: %v, Error: %v.", rawURL, query, err)
		return nil, fmt.Errorf("append query failed [%v]", err)
	}

	req, err := http.NewRequest(method, rawURL, body)
	if err != nil {
		a.logWarn("New Request failed, URL: %s, Error: %v.", rawURL, err)
		return nil, fmt.Errorf("new request failed [%v]", err)
	}

	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	if a.requestFunc != nil {
		a.requestFunc(req)
	}

	return req, nil
}

func (a *API) url(path string) string {
	if strings.HasPrefix(path, "http://") ||
		strings.HasPrefix(path, "https://") {
		return path
	}
	return xutil.JoinURL(a.baseURL, path)
}

func (a *API) doRequest(req *http.Request) (*http.Response, error) {
	a.dumpRequest(req)
	res, err := a.client.Do(req)
	if err != nil {
		a.logWarn("Do Request failed, URL: %s, Error: %v.", req.URL, err)
		return nil, err
	}
	return res, nil
}

func (a *API) handleResponse(res *http.Response, out interface{}) ([]byte, error) {
	a.dumpResponse(res)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed [%v]", err)
	}

	statusCode := res.StatusCode
	contentType := strings.ToLower(res.Header.Get("Content-Type"))

	switch a.ResType {
	case ResTypeJSON:
		if !strings.Contains(contentType, "application/json") {
			return nil, fmt.Errorf(
				"response content type is not application/json, status code [%d], content type [%s], content [%s]",
				res.StatusCode, res.Header.Get("Content-Type"), string(body),
			)
		}
		return body, a.handleJSONResponse(statusCode, body, out)
	default:
		return body, a.handleTextResponse(statusCode, body, out)
	}
}

func (a *API) handleJSONResponse(statusCode int, body []byte, out interface{}) error {
	if statusCode == http.StatusOK {
		if out != nil {
			err := json.Unmarshal(body, out)
			if err != nil {
				return fmt.Errorf("unmarshal response [%s] failed [%v]", string(body), err)
			}
			if a.errorFunc != nil {
				if a.errorFunc(out) {
					return newResponseError(statusCode, string(body), out)
				}
			}
		}
	} else {
		var ret interface{}
		if a.errorStruct != nil {
			ret = a.errorStruct
			err := json.Unmarshal(body, &ret)
			if err != nil {
				return fmt.Errorf("unmarshal response [%s] failed [%v]", string(body), err)
			}
		} else if out != nil {
			ret = out
			err := json.Unmarshal(body, out)
			if err != nil {
				return fmt.Errorf("unmarshal response [%s] failed [%v]", string(body), err)
			}
		}
		return newResponseError(statusCode, string(body), ret)
	}
	return nil
}

func (a *API) handleTextResponse(statusCode int, body []byte, out interface{}) error {
	if statusCode != http.StatusOK {
		return newResponseError(statusCode, string(body), nil)
	}

	out = &body

	return nil
}

func (a *API) Get(path string, query url.Values, header http.Header, out interface{}) error {
	//  创建请求
	req, err := a.newRequest(http.MethodGet, path, query, header, nil)
	if err != nil {
		return err
	}

	// 开启计时器
	timer := xutil.NewTimer()

	// 发送请求
	res, err := a.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// 解析并返回响应结果
	body, err := a.handleResponse(res, out)
	if err != nil {
		a.logWarn("Response failed, URL: %s, Header: %v, Error: %v, Cost: %s.", req.URL, header, err, timer.Stops())
	} else {
		a.logInfo("Response succeed, URL: %s, Header: %v, Response: %v, Cost: %s.", req.URL, header, body, timer.Stops())
	}

	return err
}

func (a *API) PostForm(path string, query url.Values, header http.Header, in url.Values, out interface{}) error {
	reqBody := in.Encode()
	logBody, err := url.QueryUnescape(reqBody)
	if err != nil {
		logBody = reqBody
	}

	// 创建请求
	req, err := a.newRequest(http.MethodPost, path, query, header, strings.NewReader(reqBody))
	if err != nil {
		return err
	}

	// 添加 header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 开启计时器
	timer := xutil.NewTimer()

	// 发送请求
	res, err := a.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// 解析并返回响应结果
	resBody, err := a.handleResponse(res, out)
	if err != nil {
		a.logWarn("Response failed, URL: %s, Header: %v, Request: %s, Error: %v, Cost: %s.", req.URL, header, logBody, err, timer.Stops())
	} else {
		a.logInfo("Response succeed, URL: %s, Header: %v, Request: %s, Response: %v, Cost: %s.", req.URL, header, logBody, string(resBody), timer.Stops())
	}

	return err
}

func (a *API) PostJSON(path string, query url.Values, header http.Header, in interface{}, out interface{}) error {
	// JSON 序列化
	if in == nil {
		in = struct{}{}
	}
	reqBody, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshal request body failed [%v]", err)
	}

	// 创建请求
	req, err := a.newRequest(http.MethodPost, path, query, header, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	// 添加 header
	req.Header.Set("Content-Type", "application/json")

	// 开启计时器
	timer := xutil.NewTimer()

	// 发送请求
	res, err := a.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// 解析并返回响应结果
	resBody, err := a.handleResponse(res, out)
	if err != nil {
		a.logWarn("Response failed, URL: %s, Header: %v, Request: %s, Error: %v, Cost: %s.", req.URL, header, string(reqBody), err, timer.Stops())
	} else {
		a.logInfo("Response succeed, URL: %s, Header: %v, Request: %s, Response: %v, Cost: %s.", req.URL, header, string(reqBody), string(resBody), timer.Stops())
	}

	return err
}

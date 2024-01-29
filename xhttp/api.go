package xhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/yyliziqiu/xlib/xurl"
	"github.com/yyliziqiu/xlib/xutil"
)

const (
	FormatJSON = "json"
	FormatText = "text"
)

type API struct {
	client    *http.Client
	format    string
	baseURL   string
	error     interface{}    // 不能是指针
	dumps     bool           // 将 HTTP 报文打印到控制台
	logger    *logrus.Logger // 如果为 nil，则不记录日志
	logLength int            // 日志最大长度
	logEscape bool           // 替换日志中的特殊字符

	requestBefore func(req *http.Request)        // 在发送请求前调用
	responseAfter func(res *http.Response) error // 在接收响应后调用
}

func New(options ...Option) *API {
	api := &API{
		client:        &http.Client{Timeout: 5 * time.Second},
		format:        FormatJSON,
		baseURL:       "",
		error:         nil,
		dumps:         false,
		logger:        nil,
		logLength:     1024,
		requestBefore: nil,
		responseAfter: nil,
	}

	for _, option := range options {
		option(api)
	}

	return api
}

func (api *API) newRequest(method string, path string, query url.Values, header http.Header, body io.Reader) (*http.Request, error) {
	rawURL, err := AppendQuery(api.url(path), query)
	if err != nil {
		api.logWarn("Append query failed, URL: %s, Query: %v, Error: %v.", rawURL, query, err)
		return nil, fmt.Errorf("append query failed [%v]", err)
	}

	req, err := http.NewRequest(method, rawURL, body)
	if err != nil {
		api.logWarn("New Request failed, URL: %s, Error: %v.", rawURL, err)
		return nil, fmt.Errorf("new request failed [%v]", err)
	}

	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	if api.requestBefore != nil {
		api.requestBefore(req)
	}

	return req, nil
}

func (api *API) url(path string) string {
	if strings.HasPrefix(path, "http://") ||
		strings.HasPrefix(path, "https://") {
		return path
	}
	if path == "" {
		return api.baseURL
	}
	return xurl.Join(api.baseURL, path)
}

func (api *API) doRequest(req *http.Request) (*http.Response, error) {
	api.dumpRequest(req)

	res, err := api.client.Do(req)
	if err != nil {
		api.logWarn("Do Request failed, URL: %s, Error: %v.", req.URL, err)
		return nil, err
	}

	return res, nil
}

func (api *API) handleResponse(res *http.Response, out interface{}) ([]byte, error) {
	api.dumpResponse(res)

	if api.responseAfter != nil {
		err := api.responseAfter(res)
		if err != nil {
			return nil, err
		}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed [%v]", err)
	}

	switch api.format {
	case FormatText:
		return body, api.handleTextResponse(res.StatusCode, body, out)
	default:
		return body, api.handleJSONResponse(res.StatusCode, body, out)
	}
}

func (api *API) handleJSONResponse(statusCode int, body []byte, out interface{}) error {
	if statusCode/100 == 2 {
		if out != nil {
			err := json.Unmarshal(body, out)
			if err != nil {
				return fmt.Errorf("unmarshal response [%s] failed [%v]", string(body), err)
			}
			if jr, ok := out.(JsonResponse); ok {
				if jr.IsError() {
					return newResponseError(statusCode, string(body), out)
				}
			}
		}
	} else {
		var ret interface{}
		if api.error != nil {
			ret = reflect.New(reflect.TypeOf(api.error)).Interface()
			err := json.Unmarshal(body, ret)
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

func (api *API) handleTextResponse(statusCode int, body []byte, out interface{}) error {
	if statusCode/100 != 2 {
		return newResponseError(statusCode, string(body), nil)
	}

	if out == nil {
		return nil
	}

	bs, ok := out.(*[]byte)
	if !ok {
		return fmt.Errorf("response receiver must *[]byte type")
	}
	*bs = body

	return nil
}

func (api *API) Get(path string, query url.Values, header http.Header, out interface{}) error {
	//  创建请求
	req, err := api.newRequest(http.MethodGet, path, query, header, nil)
	if err != nil {
		return err
	}

	// 开启计时器
	timer := xutil.NewTimer()

	// 发送请求
	res, err := api.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// 解析并返回响应结果
	body, err := api.handleResponse(res, out)
	if err != nil {
		api.logWarn("Response failed, URL: %s, header: %s, error: %v, cost: %s.", req.URL, H2S(header), err, timer.Stops())
	} else {
		api.logInfo("Response succeed, URL: %s, header: %s, response: %s, cost: %s.", req.URL, H2S(header), string(body), timer.Stops())
	}

	return err
}

func (api *API) PostForm(path string, query url.Values, header http.Header, in url.Values, out interface{}) error {
	reqBody := in.Encode()

	forlog, err := url.QueryUnescape(reqBody)
	if err != nil {
		forlog = reqBody
	}

	// 创建请求
	req, err := api.newRequest(http.MethodPost, path, query, header, strings.NewReader(reqBody))
	if err != nil {
		return err
	}

	// 添加 header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 开启计时器
	timer := xutil.NewTimer()

	// 发送请求
	res, err := api.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// 解析并返回响应结果
	resBody, err := api.handleResponse(res, out)
	if err != nil {
		api.logWarn("Response failed, URL: %s, header: %s, request: %s, error: %v, cost: %s.", req.URL, H2S(header), forlog, err, timer.Stops())
	} else {
		api.logInfo("Response succeed, URL: %s, header: %s, request: %s, response: %s, cost: %s.", req.URL, H2S(header), forlog, string(resBody), timer.Stops())
	}

	return err
}

func (api *API) PostJSON(path string, query url.Values, header http.Header, in interface{}, out interface{}) error {
	// JSON 序列化
	if in == nil {
		in = struct{}{}
	}
	reqBody, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshal request body failed [%v]", err)
	}

	// 创建请求
	req, err := api.newRequest(http.MethodPost, path, query, header, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	// 添加 header
	req.Header.Set("Content-Type", "application/json")

	// 开启计时器
	timer := xutil.NewTimer()

	// 发送请求
	res, err := api.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// 解析并返回响应结果
	resBody, err := api.handleResponse(res, out)
	if err != nil {
		api.logWarn("Response failed, URL: %s, header: %s, request: %s, error: %v, cost: %s.", req.URL, H2S(header), string(reqBody), err, timer.Stops())
	} else {
		api.logInfo("Response succeed, URL: %s, header: %s, request: %s, response: %s, cost: %s.", req.URL, H2S(header), string(reqBody), string(resBody), timer.Stops())
	}

	return err
}

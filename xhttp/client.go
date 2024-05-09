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

	"github.com/yyliziqiu/xlib/xutil"
)

const (
	FormatJSON = "json"
	FormatText = "text"
)

type Client struct {
	client    *http.Client
	format    string
	baseURL   string
	error     error          // 不能是指针
	dumps     bool           // 将 HTTP 报文打印到控制台
	logger    *logrus.Logger // 如果为 nil，则不记录日志
	logLength int            // 日志最大长度
	logEscape bool           // 替换日志中的特殊字符

	requestBefore func(req *http.Request)        // 在发送请求前调用
	responseAfter func(res *http.Response) error // 在接收响应后调用
}

func New(options ...Option) *Client {
	client := &Client{
		client:        &http.Client{Timeout: 5 * time.Second},
		format:        FormatJSON,
		baseURL:       "",
		error:         nil,
		dumps:         false,
		logger:        nil,
		logLength:     1024,
		logEscape:     false,
		requestBefore: nil,
		responseAfter: nil,
	}

	for _, option := range options {
		option(client)
	}

	return client
}

func (cli *Client) newRequest(method string, path string, query url.Values, header http.Header, body io.Reader) (*http.Request, error) {
	if !strings.HasPrefix(path, "http://") && !strings.HasPrefix(path, "https://") {
		path = JoinURL(cli.baseURL, path)
	}

	url2, err := AppendQuery(path, query)
	if err != nil {
		cli.logWarn("Append query failed, URL: %s, query: %s, error: %v.", url2, query.Encode(), err)
		return nil, fmt.Errorf("append query error [%v]", err)
	}

	req, err := http.NewRequest(method, url2, body)
	if err != nil {
		cli.logWarn("New request failed, URL: %s, error: %v.", url2, err)
		return nil, fmt.Errorf("new request error [%v]", err)
	}

	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	if cli.requestBefore != nil {
		cli.requestBefore(req)
	}

	return req, nil
}

func (cli *Client) doRequest(req *http.Request) (*http.Response, error) {
	cli.dumpRequest(req)

	res, err := cli.client.Do(req)
	if err != nil {
		cli.logWarn("Do request failed, URL: %s, error: %v.", req.URL, err)
		return nil, err
	}

	return res, nil
}

func (cli *Client) handleResponse(res *http.Response, out interface{}) ([]byte, error) {
	cli.dumpResponse(res)

	if cli.responseAfter != nil {
		err := cli.responseAfter(res)
		if err != nil {
			return nil, err
		}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response error [%v]", err)
	}

	switch cli.format {
	case FormatText:
		return body, cli.handleTextResponse(res.StatusCode, body, out)
	default:
		return body, cli.handleJSONResponse(res.StatusCode, body, out)
	}
}

func (cli *Client) handleJSONResponse(statusCode int, body []byte, out interface{}) error {
	if statusCode/100 == 2 {
		if out != nil {
			err := json.Unmarshal(body, out)
			if err != nil {
				return fmt.Errorf("unmarshal response error [%v]", err)
			}
			if jr, ok := out.(JsonResponse); ok {
				if jr.Failed() {
					err2, ok2 := out.(error)
					if ok2 {
						return err2
					}
					return newHTTPError(statusCode, string(body))
				}
			}
		}
	} else {
		if cli.error != nil {
			ret := reflect.New(reflect.TypeOf(cli.error)).Interface()
			err := json.Unmarshal(body, ret)
			if err != nil {
				return fmt.Errorf("unmarshal response error [%v]", err)
			}
			return ret.(error)
		} else if out != nil {
			err := json.Unmarshal(body, out)
			if err != nil {
				return fmt.Errorf("unmarshal response error [%v]", err)
			}
			err2, ok2 := out.(error)
			if ok2 {
				return err2
			}
			return newHTTPError(statusCode, string(body))
		}
	}
	return nil
}

func (cli *Client) handleTextResponse(statusCode int, body []byte, out interface{}) error {
	if statusCode/100 != 2 {
		return newHTTPError(statusCode, string(body))
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

func (cli *Client) Get(path string, query url.Values, header http.Header, out interface{}) error {
	req, err := cli.newRequest(http.MethodGet, path, query, header, nil)
	if err != nil {
		return err
	}

	timer := xutil.NewTimer()

	res, err := cli.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := cli.handleResponse(res, out)

	cli.logHTTP(req.URL, header, nil, body, err, timer.Stops())

	return err
}

func (cli *Client) PostForm(path string, query url.Values, header http.Header, in url.Values, out interface{}) error {
	reqBody := in.Encode()

	req, err := cli.newRequest(http.MethodPost, path, query, header, strings.NewReader(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	timer := xutil.NewTimer()

	res, err := cli.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := cli.handleResponse(res, out)

	reqBody, _ = url.QueryUnescape(reqBody)
	cli.logHTTP(req.URL, header, []byte(reqBody), resBody, err, timer.Stops())

	return err
}

func (cli *Client) PostJSON(path string, query url.Values, header http.Header, in interface{}, out interface{}) error {
	if in == nil {
		in = struct{}{}
	}
	reqBody, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshal request body error [%v]", err)
	}

	req, err := cli.newRequest(http.MethodPost, path, query, header, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	timer := xutil.NewTimer()

	res, err := cli.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := cli.handleResponse(res, out)

	cli.logHTTP(req.URL, header, reqBody, resBody, err, timer.Stops())

	return err
}

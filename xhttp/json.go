package xhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type JSON struct {
	Client      *http.Client
	ErrorStruct interface{} // must not point
	RequestFunc func(req *http.Request)
}

func NewJSON() *JSON {
	return &JSON{
		Client: http.DefaultClient,
	}
}

func NewJSON1(d time.Duration) *JSON {
	return &JSON{
		Client: newClient(d),
	}
}

func NewJSON2(d time.Duration, errorStruct interface{}) *JSON {
	return &JSON{
		Client:      newClient(d),
		ErrorStruct: errorStruct,
	}
}

func NewJSON3(d time.Duration, username string, password string) *JSON {
	return &JSON{
		Client:      newClient(d),
		RequestFunc: basicAuth(username, password),
	}
}

func NewJSON4(d time.Duration, errorStruct interface{}, username string, password string) *JSON {
	return &JSON{
		Client:      newClient(d),
		ErrorStruct: errorStruct,
		RequestFunc: basicAuth(username, password),
	}
}

func (h *JSON) Get(rawURL string, query url.Values, header http.Header, out interface{}) error {
	//  创建请求
	req, err := newRequest(http.MethodGet, rawURL, query, header, nil, h.RequestFunc)
	if err != nil {
		return err
	}

	// 发送请求
	res, err := h.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// 解析并返回响应结果
	return h.handleResponse(res, out)
}

func (h *JSON) doRequest(req *http.Request) (*http.Response, error) {
	if h.Client == nil {
		h.Client = http.DefaultClient
	}
	return h.Client.Do(req)
}

func (h *JSON) handleResponse(res *http.Response, out interface{}) error {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read response failed [%v]", err)
	}

	// 校验 content type
	if !strings.Contains(strings.ToLower(res.Header.Get("Content-Type")), "application/json") {
		return fmt.Errorf("response content type is not application/json, status code [%d], content type [%s], content [%s]",
			res.StatusCode, res.Header.Get("Content-Type"), string(body))
	}

	if res.StatusCode == http.StatusOK {
		if out != nil {
			err = json.Unmarshal(body, out)
			if err != nil {
				return fmt.Errorf("unmarshal successful response failed [%v]", err)
			}
		}
	} else {
		var ret interface{}
		if h.ErrorStruct != nil {
			ret = h.ErrorStruct
			err = json.Unmarshal(body, &ret)
			if err != nil {
				return fmt.Errorf("unmarshal fail response failed [%v]", err)
			}
		} else if out != nil {
			ret = out
			err = json.Unmarshal(body, out)
			if err != nil {
				return fmt.Errorf("unmarshal fail response failed [%v]", err)
			}
		}
		return newResponseError(res.StatusCode, string(body), ret)
	}

	return nil
}

func (h *JSON) PostForm(url string, query url.Values, header http.Header, in url.Values, out interface{}) error {
	// 创建请求
	req, err := newRequest(http.MethodPost, url, query, header, strings.NewReader(in.Encode()), h.RequestFunc)
	if err != nil {
		return err
	}

	// 添加 header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求
	res, err := h.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// 解析并返回响应结果
	return h.handleResponse(res, out)
}

func (h *JSON) PostJSON(url string, query url.Values, header http.Header, in interface{}, out interface{}) error {
	// json 序列化
	if in == nil {
		in = struct{}{}
	}
	body, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshal request body failed [%v]", err)
	}

	// 创建请求
	req, err := newRequest(http.MethodPost, url, query, header, bytes.NewReader(body), h.RequestFunc)
	if err != nil {
		return err
	}

	// 添加 header
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	res, err := h.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// 解析并返回响应结果
	return h.handleResponse(res, out)
}

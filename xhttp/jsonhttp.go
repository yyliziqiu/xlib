package xhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type JSONHTTP struct {
	Client      *http.Client
	ErrorStruct interface{} // must not point
	RequestFunc func(req *http.Request)
}

func (h *JSONHTTP) Get(rawURL string, header http.Header, query url.Values, out interface{}) (interface{}, error) {
	// 解析 url，合并参数
	uo, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	for k, v := range uo.Query() {
		for _, s := range v {
			query.Add(k, s)
		}
	}
	uo.RawQuery = query.Encode()

	//  创建请求
	req, err := h.newRequest(http.MethodGet, uo.String(), header, nil)
	if err != nil {
		return nil, err
	}

	// 发送请求
	res, err := h.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// 解析并返回响应结果
	return h.handleResponse(res, out)
}

func (h *JSONHTTP) newRequest(method string, url string, header http.Header, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("new request failed, %v", err)
	}
	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	h.RequestFunc(req)
	return req, nil
}

func (h *JSONHTTP) doRequest(req *http.Request) (*http.Response, error) {
	if h.Client == nil {
		h.Client = http.DefaultClient
	}
	return h.Client.Do(req)
}

func (h *JSONHTTP) handleResponse(res *http.Response, out interface{}) (interface{}, error) {
	outbs, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed, %v", err)
	}

	// 校验 content type
	if !strings.Contains(strings.ToLower(res.Header.Get("Content-Type")), "application/json") {
		return nil, fmt.Errorf("response content type is not application/json, status code: %d, content type: %s, content: %s",
			res.StatusCode, res.Header.Get("Content-Type"), string(outbs))
	}

	if res.StatusCode == http.StatusOK {
		if out != nil {
			err = json.Unmarshal(outbs, out)
			if err != nil {
				return nil, fmt.Errorf("unmarshal successful response failed, %s", err)
			}
		}
	} else {
		cpy := h.ErrorStruct
		if h.ErrorStruct != nil {
			err = json.Unmarshal(outbs, &cpy)
			if err != nil {
				return nil, fmt.Errorf("unmarshal fail response failed, %s", err)
			}
		} else if out != nil {
			err = json.Unmarshal(outbs, out)
			if err != nil {
				return nil, fmt.Errorf("unmarshal fail response failed, %s", err)
			}
		}
		return cpy, NewResponseError(res.StatusCode, string(outbs))
	}

	return nil, nil
}

func (h *JSONHTTP) PostForm(url string, header http.Header, in url.Values, out interface{}) (interface{}, error) {
	// 创建请求
	req, err := h.newRequest(http.MethodPost, url, header, strings.NewReader(in.Encode()))
	if err != nil {
		return nil, err
	}

	// 添加 header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求
	res, err := h.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// 解析并返回响应结果
	return h.handleResponse(res, out)
}

func (h *JSONHTTP) PostJSON(url string, header http.Header, in interface{}, out interface{}) (interface{}, error) {
	// json 序列化
	if in == nil {
		in = struct{}{}
	}
	inbs, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("marshal request body failed, %v", err)
	}

	// 创建请求
	req, err := h.newRequest(http.MethodPost, url, header, bytes.NewReader(inbs))
	if err != nil {
		return nil, err
	}

	// 添加 header
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	res, err := h.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// 解析并返回响应结果
	return h.handleResponse(res, out)
}

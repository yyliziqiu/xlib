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

type Text struct {
	Client      *http.Client
	RequestFunc func(req *http.Request)
}

func NewText() *Text {
	return &Text{
		Client: http.DefaultClient,
	}
}

func NewText1(d time.Duration) *Text {
	return &Text{
		Client: newClient(d),
	}
}

func NewText3(d time.Duration, username string, password string) *Text {
	return &Text{
		Client:      newClient(d),
		RequestFunc: basicAuth(username, password),
	}
}

func (h *Text) Get(rawURL string, query url.Values, header http.Header) ([]byte, error) {
	//  创建请求
	req, err := newRequest(http.MethodPost, rawURL, query, header, nil, h.RequestFunc)
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
	return h.handleResponse(res)
}

func (h *Text) doRequest(req *http.Request) (*http.Response, error) {
	if h.Client == nil {
		h.Client = http.DefaultClient
	}
	return h.Client.Do(req)
}

func (h *Text) handleResponse(res *http.Response) ([]byte, error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed [%v]", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, newResponseError(res.StatusCode, string(body), nil)
	}

	return body, nil
}

func (h *Text) PostForm(url string, query url.Values, header http.Header, in url.Values) ([]byte, error) {
	// 创建请求
	req, err := newRequest(http.MethodPost, url, query, header, strings.NewReader(in.Encode()), h.RequestFunc)
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
	return h.handleResponse(res)
}

func (h *Text) PostJSON(url string, query url.Values, header http.Header, in interface{}) ([]byte, error) {
	// json 序列化
	if in == nil {
		in = struct{}{}
	}
	body, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("marshal request body failed [%v]", err)
	}

	// 创建请求
	req, err := newRequest(http.MethodPost, url, query, header, bytes.NewReader(body), h.RequestFunc)
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
	return h.handleResponse(res)
}

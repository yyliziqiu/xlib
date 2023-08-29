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

var client *http.Client

func SetClient(c *http.Client) {
	client = c
}

func Get(rawURL string, header http.Header, query url.Values, out interface{}) error {
	// 解析 url，合并参数
	uo, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	for k, v := range uo.Query() {
		for _, s := range v {
			query.Add(k, s)
		}
	}
	uo.RawQuery = query.Encode()

	//  创建请求
	req, err := newRequest(http.MethodGet, uo.String(), header, nil)
	if err != nil {
		return err
	}

	// 发送请求
	res, err := getClient().Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// 解析并返回响应结果
	return handleResponse(res, out)
}

func newRequest(method string, url string, header http.Header, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("new request failed, %v", err)
	}
	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	return req, nil
}

func getClient() *http.Client {
	if client != nil {
		return client
	}
	return http.DefaultClient
}

func handleResponse(res *http.Response, out interface{}) error {
	outbs, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read response failed, %v", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("response error, status code: %d, content: %s", res.StatusCode, string(outbs))
	}
	if !strings.Contains(strings.ToLower(res.Header.Get("Content-Type")), "application/json") {
		return fmt.Errorf("content type error, content type: %s, content: %s", res.Header.Get("Content-Type"), string(outbs))
	}
	if out != nil {
		err = json.Unmarshal(outbs, out)
		if err != nil {
			return fmt.Errorf("json unmarshal failed, %v", err)
		}
	}
	return nil
}

func PostForm(url string, header http.Header, in url.Values, out interface{}) error {
	// 创建请求
	req, err := newRequest(http.MethodPost, url, header, strings.NewReader(in.Encode()))
	if err != nil {
		return err
	}

	// 添加 header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求
	res, err := getClient().Do(req)
	if err != nil {
		return fmt.Errorf("do request failed, %v", err)
	}
	defer res.Body.Close()

	// 解析并返回响应结果
	return handleResponse(res, out)
}

type emptyJSON struct{}

func PostJSON(url string, header http.Header, in interface{}, out interface{}) error {
	// json 序列化
	if in == nil {
		in = emptyJSON{}
	}
	inBytes, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshal http request failed, %v", err)
	}

	// 创建请求
	req, err := newRequest(http.MethodPost, url, header, bytes.NewReader(inBytes))
	if err != nil {
		return err
	}

	// 添加 header
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	res, err := getClient().Do(req)
	if err != nil {
		return fmt.Errorf("do request failed, %v", err)
	}
	defer res.Body.Close()

	// 解析并返回响应结果
	return handleResponse(res, out)
}

package xhttp

import (
	"net/url"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestAPI_Get(t *testing.T) {
	logger := logrus.New()

	api := NewAPI(
		WithResType(ResTypeText),
		WithTimeout(time.Second),
		WithLogger(logger),
	)

	var result []byte

	err := api.Get("https://www.baidu.com", nil, nil, &result)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(result))
}

func TestAPI_PostForm(t *testing.T) {
	logger := logrus.New()

	api := NewAPI(
		WithResType(ResTypeText),
		WithTimeout(time.Second),
		WithLogger(logger),
	)

	body := url.Values{}
	body.Set("param1", "1")
	body.Set("param2", "2")
	body.Set("param3", "3")

	err := api.PostForm("https://www.baidu.com", nil, nil, body, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAPI_PostJSON(t *testing.T) {
	logger := logrus.New()

	api := NewAPI(
		WithResType(ResTypeText),
		WithTimeout(time.Second),
		WithLogger(logger),
	)

	body := map[string]string{
		"name": "ylq",
		"pass": "xzc",
	}

	err := api.PostJSON("https://www.baidu.com", nil, nil, body, nil)
	if err != nil {
		t.Fatal(err)
	}
}

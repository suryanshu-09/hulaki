package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var username, password string

func SetParams(url *string, params map[string]string) {
	if len(params) != 0 {
		*url = fmt.Sprintf("%s?", *url)
		paramArr := make([]string, 0)
		for key, val := range params {
			paramArr = append(paramArr, strings.Join([]string{key, val}, "="))
		}
		finalParams := strings.Join(paramArr, "&")
		*url = fmt.Sprintf("%s%s", *url, finalParams)
	}
}

func SetHeaders(req *http.Request, headers map[string]string) {
	for key, val := range headers {
		req.Header.Set(key, val)
	}
}

func SetBasicAuth(usr, pass string) {
	username = usr
	password = pass
}

type (
	Args func(*Arg)
	Arg  struct {
		Body    io.Reader
		Params  map[string]string
		Headers map[string]string
	}
)

func WithBody(body io.Reader) Args {
	return func(arg *Arg) {
		arg.Body = body
	}
}

func WithParams(params map[string]string) Args {
	return func(arg *Arg) {
		arg.Params = params
	}
}

func WithHeaders(headers map[string]string) Args {
	return func(arg *Arg) {
		arg.Headers = headers
	}
}

func GetArgs(args []Args) (body io.Reader, params, headers map[string]string) {
	arg := Arg{Body: &bytes.Buffer{}, Params: make(map[string]string, 0), Headers: make(map[string]string, 0)}
	for _, a := range args {
		a(&arg)
	}
	body = arg.Body
	params = arg.Params
	headers = arg.Headers
	return
}

func HTTPGet(url string, args ...Args) (*http.Response, error) {
	body, params, headers := GetArgs(args)
	for range args {
	}
	SetParams(&url, params)
	req, err := http.NewRequest("GET", url, body)
	if err != nil {
		return nil, err
	}
	SetHeaders(req, headers)
	client := http.Client{}
	return client.Do(req)
}

func HTTPPost(url string, args ...Args) (*http.Response, error) {
	body, params, headers := GetArgs(args)
	SetParams(&url, params)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	SetHeaders(req, headers)
	client := http.Client{}
	return client.Do(req)
}

func HTTPPut(url string, args ...Args) (*http.Response, error) {
	body, params, headers := GetArgs(args)
	SetParams(&url, params)
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}
	SetHeaders(req, headers)
	client := http.Client{}
	return client.Do(req)
}

func HTTPPatch(url string, args ...Args) (*http.Response, error) {
	body, params, headers := GetArgs(args)
	SetParams(&url, params)
	req, err := http.NewRequest("PATCH", url, body)
	if err != nil {
		return nil, err
	}
	SetHeaders(req, headers)
	client := http.Client{}
	return client.Do(req)
}

func HTTPDelete(url string, args ...Args) (*http.Response, error) {
	body, params, headers := GetArgs(args)
	SetParams(&url, params)
	req, err := http.NewRequest("DELETE", url, body)
	if err != nil {
		return nil, err
	}
	SetHeaders(req, headers)
	client := http.Client{}
	return client.Do(req)
}

func HTTPHead(url string, args ...Args) (*http.Response, error) {
	body, params, headers := GetArgs(args)
	SetParams(&url, params)
	req, err := http.NewRequest("HEAD", url, body)
	if err != nil {
		return nil, err
	}
	SetHeaders(req, headers)
	client := http.Client{}
	return client.Do(req)
}

func HTTPOptions(url string, args ...Args) (*http.Response, error) {
	body, params, headers := GetArgs(args)
	SetParams(&url, params)
	req, err := http.NewRequest("OPTIONS", url, body)
	if err != nil {
		return nil, err
	}
	SetHeaders(req, headers)
	client := http.Client{}
	return client.Do(req)
}

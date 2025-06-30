/*
Package utils
Utility functions
*/
package utils

import (
	"net/http"
)

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

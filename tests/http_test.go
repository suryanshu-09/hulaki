package tests

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/suryanshu-09/hulaki/utils"
)

// http
// GET
// POST
// PUT
// PATCH
// DELETE
// HEAD
// OPTIONS

// PARAMS input box - key, value, description
// AUTH  buttload of auths, stick to basic auth, then bearer and then jwt
// HEADERS input box - key, value, description
// BODY form-data, x-www-form-urlencoded, raw, binary, graphQL
// SCRIPTS pre-request, post-response
// SETTINGS lots of stuff

func setupTestServer() *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Has("x") && r.URL.Query().Has("y") {
				x := r.URL.Query().Get("x")
				y := r.URL.Query().Get("y")
				w.Write([]byte(x + y))
				return
			}
			if r.Header.Get("getthis") == "got" && r.Header.Get("tryharder") == "tried harder" {
				w.Header().Set("getthis", "got")
				w.Header().Set("tryharder", "tried harder")
				return
			}
			body, _ := io.ReadAll(r.Body)
			if len(body) > 0 {
				w.Write(body)
				return
			}
			w.Write([]byte("this was a get request"))
		case http.MethodPost:
			w.Write([]byte("this was a post request"))
		case http.MethodPut:
			w.Write([]byte("this was a put request"))
		case http.MethodPatch:
			w.Write([]byte("this was a patch request"))
		case http.MethodDelete:
			w.Write([]byte("this was a delete request"))
		case http.MethodHead:
			w.Header().Set("works", "this was a head request")
		case http.MethodOptions:
			w.Header().Set("works", "this was a options request")
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	return httptest.NewServer(handler)
}

func TestGet(t *testing.T) {
	server := setupTestServer()
	defer server.Close()
	t.Run("test GET normally", func(t *testing.T) {
		resp, err := utils.HTTPGet(server.URL)
		if err != nil {
			t.Errorf("got an error:%s", err.Error())
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("got an error:%s", err.Error())
		}
		want := "this was a get request"
		got := string(body)
		if got != want {
			t.Errorf("got:%s\nwant:%s", got, want)
		}
	})
	t.Run("test GET params", func(t *testing.T) {
		params := map[string]string{"x": "16", "y": "12"}
		resp, err := utils.HTTPGet(server.URL, utils.WithParams(params))
		if err != nil {
			t.Errorf("got an error:%s", err.Error())
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("got an error:%s", err.Error())
		}
		want := "1612"
		got := string(body)
		if got != want {
			t.Errorf("got:%s\nwant:%s", got, want)
		}
	})
	t.Run("test GET headers", func(t *testing.T) {
		header := map[string]string{"getthis": "got", "tryharder": "tried harder"}
		resp, err := utils.HTTPGet(server.URL, utils.WithHeaders(header))
		if err != nil {
			t.Errorf("got an error:%s", err.Error())
		}
		defer resp.Body.Close()
		headers := resp.Header
		key1 := headers.Get("getthis")
		key2 := headers.Get("tryharder")
		want1 := "got"
		want2 := "tried harder"
		if key1 != want1 {
			t.Errorf("got:%s\nwant:%s", key1, want1)
		}
		if key2 != want2 {
			t.Errorf("got:%s\nwant:%s", key2, want2)
		}
	})
	t.Run("test GET body", func(t *testing.T) {
		buf := bytes.Buffer{}
		buf.WriteString("new body new test")
		resp, err := utils.HTTPGet(server.URL, utils.WithBody(&buf))
		if err != nil {
			t.Errorf("got an error:%s", err.Error())
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("got an error:%s", err.Error())
		}
		want := "new body new test"
		got := string(body)
		if got != want {
			t.Errorf("got:%s\nwant:%s", got, want)
		}
	})
}

func TestPost(t *testing.T) {
	server := setupTestServer()
	defer server.Close()
	t.Run("test POST normally", func(t *testing.T) {
		resp, err := utils.HTTPPost(server.URL)
		if err != nil {
			t.Errorf("got an error:%s", err.Error())
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("got an error:%s", err.Error())
		}
		want := "this was a post request"
		got := string(body)
		if got != want {
			t.Errorf("got:%s\nwant:%s", got, want)
		}
	})
}

func TestPut(t *testing.T) {
	server := setupTestServer()
	defer server.Close()
	t.Run("test PUT normally", func(t *testing.T) {
		resp, err := utils.HTTPPut(server.URL)
		if err != nil {
			t.Errorf("got an error:%s", err.Error())
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("got an error:%s", err.Error())
		}
		want := "this was a put request"
		got := string(body)
		if got != want {
			t.Errorf("got:%s\nwant:%s", got, want)
		}
	})
}

func TestPatch(t *testing.T) {
	server := setupTestServer()
	defer server.Close()
	t.Run("test PATCH normally", func(t *testing.T) {
		resp, err := utils.HTTPPatch(server.URL)
		if err != nil {
			t.Errorf("got an error:%s", err.Error())
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("got an error:%s", err.Error())
		}
		want := "this was a patch request"
		got := string(body)
		if got != want {
			t.Errorf("got:%s\nwant:%s", got, want)
		}
	})
}

func TestDelete(t *testing.T) {
	server := setupTestServer()
	defer server.Close()
	t.Run("test Delete normally", func(t *testing.T) {
		resp, err := utils.HTTPDelete(server.URL)
		if err != nil {
			t.Errorf("got an error:%s", err.Error())
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("got an error:%s", err.Error())
		}
		want := "this was a delete request"
		got := string(body)
		if got != want {
			t.Errorf("got:%s\nwant:%s", got, want)
		}
	})
}

func TestHead(t *testing.T) {
	server := setupTestServer()
	defer server.Close()
	t.Run("test Head normally", func(t *testing.T) {
		resp, err := utils.HTTPHead(server.URL)
		if err != nil {
			t.Errorf("got an error:%s", err.Error())
		}
		defer resp.Body.Close()
		head := resp.Header
		got := head.Get("works")
		if err != nil {
			t.Errorf("got an error:%s", err.Error())
		}
		want := "this was a head request"
		if got != want {
			t.Errorf("got:%s\nwant:%s", got, want)
		}
	})
}

func TestOptions(t *testing.T) {
	server := setupTestServer()
	defer server.Close()
	t.Run("test Options normally", func(t *testing.T) {
		resp, err := utils.HTTPOptions(server.URL)
		if err != nil {
			t.Errorf("got an error:%s", err.Error())
		}
		defer resp.Body.Close()
		head := resp.Header
		got := head.Get("works")
		if err != nil {
			t.Errorf("got an error:%s", err.Error())
		}
		want := "this was a options request"
		if got != want {
			t.Errorf("got:%s\nwant:%s", got, want)
		}
	})
}

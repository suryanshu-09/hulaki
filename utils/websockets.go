package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebsocketClient struct {
	URL    string
	WS     *websocket.Conn
	Header *http.Response
}

func NewWebsocketClient(url string, args ...Args) (*WebsocketClient, error) {
	_, params, headers := GetArgs(args)
	SetParams(&url, params)
	r, _ := http.NewRequest("GET", url, nil)
	SetHeaders(r, headers)
	ws, header, err := websocket.DefaultDialer.Dial(url, r.Header)
	if err != nil {
		return nil, err
	}
	return &WebsocketClient{WS: ws, URL: url, Header: header}, nil
}

func (ws *WebsocketClient) Read(out io.Writer) error {
	_, p, err := ws.WS.ReadMessage()
	if err != nil {
		return err
	}
	fmt.Fprint(out, string(p))
	return nil
}

func (ws *WebsocketClient) Write(mt int, msg io.Reader) error {
	data := new(bytes.Buffer)
	_, err := io.Copy(data, msg)
	if err != nil {
		return err
	}
	return ws.WS.WriteMessage(websocket.TextMessage, data.Bytes())
}

func (ws *WebsocketClient) Close() error {
	return ws.WS.Close()
}

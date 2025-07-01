package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebsocketClient struct {
	URL       string
	WS        *websocket.Conn
	Header    *http.Response
	ReadChan  chan string
	WriteChan chan string
	ErrorChan chan error
	CloseChan chan struct{}
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

	client := &WebsocketClient{
		WS:        ws,
		URL:       url,
		Header:    header,
		ReadChan:  make(chan string),
		WriteChan: make(chan string),
		ErrorChan: make(chan error),
		CloseChan: make(chan struct{}),
	}

	go client.startReader()
	go client.startWriter()

	go client.startReader()
	go client.startWriter()
	return client, nil
}

func (ws *WebsocketClient) Read(out io.Writer) error {
	select {
	case message := <-ws.ReadChan:
		fmt.Fprint(out, message)
		return nil
	case err := <-ws.ErrorChan:
		return err
	case <-ws.CloseChan:
		return fmt.Errorf("connection closed")
	}
}

func (ws *WebsocketClient) startReader() {
	for {
		select {
		case <-ws.CloseChan:
			return
		default:
			_, message, err := ws.WS.ReadMessage()
			if err != nil {
				ws.ErrorChan <- err
				return
			}
			ws.ReadChan <- string(message)
		}
	}
}

func (ws *WebsocketClient) startWriter() {
	for {
		select {
		case <-ws.CloseChan:
			return
		case message := <-ws.WriteChan:
			err := ws.WS.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				ws.ErrorChan <- err
				return
			}
		}
	}
}

func (ws *WebsocketClient) Write(mt int, msg io.Reader) error {
	data := new(bytes.Buffer)
	_, err := io.Copy(data, msg)
	if err != nil {
		return err
	}

	select {
	case ws.WriteChan <- data.String():
		return nil
	case err := <-ws.ErrorChan:
		return err
	case <-ws.CloseChan:
		return fmt.Errorf("connection closed")
	}
}

func (ws *WebsocketClient) Close() error {
	close(ws.CloseChan)
	close(ws.ReadChan)
	close(ws.WriteChan)
	close(ws.ErrorChan)
	return ws.WS.Close()
}

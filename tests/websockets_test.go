package tests

import (
	"io"
	"net/http"
	"net/http/httptest"

	"golang.org/x/net/websocket"
)

func setupWebsocketServer() {
	echosServer := func(ws *websocket.Conn) {
		io.Copy(ws, ws)
	}

	h := http.Handler(websocket.Handler(echosServer))

	httptest.NewServer(h)
}

package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/suryanshu-09/hulaki/utils"
)

var upgrader = websocket.Upgrader{}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(mt, msg)
	}
}

func TestWebSocketEcho(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(wsHandler))
	defer server.Close()

	wsURL := "ws" + server.URL[len("http"):]

	ws, err := utils.NewWebsocketClient(wsURL)
	if err != nil {
		t.Error(err)
	}
	msg := new(bytes.Buffer)
	msg.WriteString("hulaki")
	if err = ws.Write(websocket.TextMessage, msg); err != nil {
		t.Error(err)
	}
	got := new(bytes.Buffer)
	if err = ws.Read(got); err != nil {
		t.Error(err)
	}
	if got.String() != "hulaki" {
		t.Errorf("got:%s\nwant:hulaki", got.String())
	}
}

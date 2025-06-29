/*
Package main
main.go
*/
package main

import (
	"fmt"
	"log"

	"github.com/suryanshu-09/hulaki/cmd"
	"golang.org/x/net/websocket"
)

func main() {
	cmd.Execute()
}

func RunClient() {
	origin := "http://localhost/"
	url := "ws://localhost:12345/ws"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := ws.Write([]byte("hello, world!\n")); err != nil {
		log.Fatal(err)
	}
	msg := make([]byte, 512)
	var n int
	if n, err = ws.Read(msg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received: %s.\n", msg[:n])
}

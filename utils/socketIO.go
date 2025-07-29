package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type SocketIOClient struct {
	conn     *websocket.Conn
	url      string
	headers  map[string][]string
	messages chan SocketIOMessage
	done     chan bool
}

type SocketIOMessage struct {
	Type int    `json:"type"`
	Data string `json:"data"`
}

type SocketIOResponse struct {
	Event     string         `json:"event"`
	Data      any            `json:"data"`
	Error     *SocketIOError `json:"error,omitempty"`
	Status    string         `json:"status"`
	Connected bool           `json:"connected"`
}

type SocketIOError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SocketIORequest struct {
	Event string         `json:"event"`
	Data  map[string]any `json:"data,omitempty"`
}

func NewSocketIOClient(serverURL string, args ...Args) (*SocketIOClient, error) {
	_, params, headers := GetArgs(args)

	parsedURL, err := url.Parse(serverURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	if parsedURL.Scheme != "ws" && parsedURL.Scheme != "wss" && parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("invalid URL scheme: %s (must be ws, wss, http, or https)", parsedURL.Scheme)
	}

	if parsedURL.Host == "" {
		return nil, fmt.Errorf("invalid URL: missing host")
	}

	if parsedURL.Scheme == "http" {
		parsedURL.Scheme = "ws"
	} else if parsedURL.Scheme == "https" {
		parsedURL.Scheme = "wss"
	}

	if !strings.HasSuffix(parsedURL.Path, "/socket.io/") {
		if parsedURL.Path == "" || parsedURL.Path == "/" {
			parsedURL.Path = "/socket.io/"
		} else {
			parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/") + "/socket.io/"
		}
	}

	query := parsedURL.Query()
	query.Set("EIO", "4")
	query.Set("transport", "websocket")
	for key, value := range params {
		query.Set(key, value)
	}
	parsedURL.RawQuery = query.Encode()

	var reqHeaders map[string][]string
	if len(headers) > 0 {
		reqHeaders = make(map[string][]string)
		for key, value := range headers {
			reqHeaders[key] = []string{value}
		}
	}

	return &SocketIOClient{
		url:      parsedURL.String(),
		headers:  reqHeaders,
		messages: make(chan SocketIOMessage, 100),
		done:     make(chan bool),
	}, nil
}

func (c *SocketIOClient) Connect() error {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(c.url, c.headers)
	if err != nil {
		return fmt.Errorf("failed to connect to Socket.IO server: %w", err)
	}

	c.conn = conn
	return nil
}

func (c *SocketIOClient) Disconnect() {
	if c.conn != nil {
		c.conn.Close()
		c.done <- true
	}
}

func (c *SocketIOClient) Emit(event string, data any) error {
	if c.conn == nil {
		return fmt.Errorf("not connected")
	}

	eventData := []any{event}
	if data != nil {
		eventData = append(eventData, data)
	}

	payload, err := json.Marshal(eventData)
	if err != nil {
		return err
	}

	message := "42" + string(payload)
	return c.conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (c *SocketIOClient) IsConnected() bool {
	return c.conn != nil
}

func (c *SocketIOClient) ReadMessages() {
	if c.conn == nil {
		return
	}

	for {
		select {
		case <-c.done:
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				return
			}

			c.messages <- SocketIOMessage{
				Type: websocket.TextMessage,
				Data: string(message),
			}
		}
	}
}

func SocketIOConnect(serverURL string, args ...Args) (*SocketIOResponse, error) {
	client, err := NewSocketIOClient(serverURL, args...)
	if err != nil {
		return &SocketIOResponse{
			Error: &SocketIOError{
				Code:    "CONNECTION_ERROR",
				Message: err.Error(),
			},
			Status:    "error",
			Connected: false,
		}, nil
	}

	err = client.Connect()
	if err != nil {
		return &SocketIOResponse{
			Error: &SocketIOError{
				Code:    "CONNECTION_FAILED",
				Message: fmt.Sprintf("failed to connect to Socket.IO server: %v", err),
			},
			Status:    "error",
			Connected: false,
		}, nil
	}

	time.Sleep(100 * time.Millisecond)

	connected := client.IsConnected()
	client.Disconnect()

	return &SocketIOResponse{
		Event: "connect",
		Data: map[string]any{
			"url":       serverURL,
			"timestamp": time.Now().Unix(),
		},
		Status:    "success",
		Connected: connected,
	}, nil
}

func SocketIOEmit(serverURL, event string, args ...Args) (*SocketIOResponse, error) {
	client, err := NewSocketIOClient(serverURL, args...)
	if err != nil {
		return &SocketIOResponse{
			Error: &SocketIOError{
				Code:    "CONNECTION_ERROR",
				Message: err.Error(),
			},
			Status:    "error",
			Connected: false,
		}, nil
	}

	err = client.Connect()
	if err != nil {
		return &SocketIOResponse{
			Error: &SocketIOError{
				Code:    "CONNECTION_FAILED",
				Message: fmt.Sprintf("failed to connect to Socket.IO server: %v", err),
			},
			Status:    "error",
			Connected: false,
		}, nil
	}
	defer client.Disconnect()

	body, _, _ := GetArgs(args)
	var eventData map[string]any

	if body != nil {
		decoder := json.NewDecoder(body)
		if err := decoder.Decode(&eventData); err != nil && err != io.EOF {
			return &SocketIOResponse{
				Error: &SocketIOError{
					Code:    "INVALID_DATA",
					Message: fmt.Sprintf("failed to decode event data: %v", err),
				},
				Status:    "error",
				Connected: client.IsConnected(),
			}, nil
		}
	}

	time.Sleep(100 * time.Millisecond)

	err = client.Emit(event, eventData)
	if err != nil {
		return &SocketIOResponse{
			Error: &SocketIOError{
				Code:    "EMIT_FAILED",
				Message: fmt.Sprintf("failed to emit event: %v", err),
			},
			Status:    "error",
			Connected: client.IsConnected(),
		}, nil
	}

	return &SocketIOResponse{
		Event: event,
		Data: map[string]any{
			"event":     event,
			"data":      eventData,
			"url":       serverURL,
			"timestamp": time.Now().Unix(),
		},
		Status:    "success",
		Connected: client.IsConnected(),
	}, nil
}

func SocketIOListen(serverURL, event string, duration time.Duration, args ...Args) (*SocketIOResponse, error) {
	client, err := NewSocketIOClient(serverURL, args...)
	if err != nil {
		return &SocketIOResponse{
			Error: &SocketIOError{
				Code:    "CONNECTION_ERROR",
				Message: err.Error(),
			},
			Status:    "error",
			Connected: false,
		}, nil
	}

	err = client.Connect()
	if err != nil {
		return &SocketIOResponse{
			Error: &SocketIOError{
				Code:    "CONNECTION_FAILED",
				Message: fmt.Sprintf("failed to connect to Socket.IO server: %v", err),
			},
			Status:    "error",
			Connected: false,
		}, nil
	}
	defer client.Disconnect()

	go client.ReadMessages()

	receivedMessages := make([]string, 0)
	timeout := time.After(duration)

	for {
		select {
		case msg := <-client.messages:
			receivedMessages = append(receivedMessages, msg.Data)
		case <-timeout:
			goto done
		}
	}

done:
	return &SocketIOResponse{
		Event: event,
		Data: map[string]any{
			"event":           event,
			"messages":        receivedMessages,
			"message_count":   len(receivedMessages),
			"listen_duration": duration.String(),
			"url":             serverURL,
			"timestamp":       time.Now().Unix(),
		},
		Status:    "success",
		Connected: client.IsConnected(),
	}, nil
}

func WithNamespace(namespace string) Args {
	return WithParams(map[string]string{"namespace": namespace})
}

func WithSocketIOTimeout(timeout time.Duration) Args {
	return func(arg *Arg) {
	}
}

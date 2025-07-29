package tests

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/suryanshu-09/hulaki/utils"
)

func TestSocketIOConnect(t *testing.T) {
	t.Run("test Socket.IO connection error", func(t *testing.T) {
		// Test connection to non-existent server
		resp, err := utils.SocketIOConnect("ws://localhost:99999")
		if err != nil {
			t.Errorf("SocketIOConnect should not return error, got: %s", err.Error())
		}

		if resp.Error == nil {
			t.Errorf("Expected connection error, got nil")
		}

		if resp.Error.Code != "CONNECTION_ERROR" && resp.Error.Code != "CONNECTION_FAILED" {
			t.Errorf("Expected CONNECTION_ERROR or CONNECTION_FAILED, got: %s", resp.Error.Code)
		}
	})
}

func TestSocketIOEmit(t *testing.T) {
	t.Run("test Socket.IO emit connection error", func(t *testing.T) {
		// Test emit to non-existent server
		testData := map[string]any{"message": "test"}
		jsonData, _ := json.Marshal(testData)

		resp, err := utils.SocketIOEmit("ws://localhost:99999", "test-event",
			utils.WithBody(bytes.NewBuffer(jsonData)))
		if err != nil {
			t.Errorf("SocketIOEmit should not return error, got: %s", err.Error())
		}

		if resp.Error == nil {
			t.Errorf("Expected connection error, got nil")
		}

		if resp.Error.Code != "CONNECTION_ERROR" && resp.Error.Code != "CONNECTION_FAILED" {
			t.Errorf("Expected CONNECTION_ERROR or CONNECTION_FAILED, got: %s", resp.Error.Code)
		}
	})

	t.Run("test Socket.IO emit with invalid data", func(t *testing.T) {
		// Test emit with invalid JSON data
		invalidJSON := "invalid json {"

		resp, err := utils.SocketIOEmit("ws://localhost:3000", "test-event",
			utils.WithBody(bytes.NewBuffer([]byte(invalidJSON))))
		if err != nil {
			t.Errorf("SocketIOEmit should not return error, got: %s", err.Error())
		}

		// Should get either connection error or invalid data error
		if resp.Error == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestSocketIOListen(t *testing.T) {
	t.Run("test Socket.IO listen connection error", func(t *testing.T) {
		// Test listen on non-existent server
		resp, err := utils.SocketIOListen("ws://localhost:99999", "test-event", 100)
		if err != nil {
			t.Errorf("SocketIOListen should not return error, got: %s", err.Error())
		}

		if resp.Error == nil {
			t.Errorf("Expected connection error, got nil")
		}

		if resp.Error.Code != "CONNECTION_ERROR" && resp.Error.Code != "CONNECTION_FAILED" {
			t.Errorf("Expected CONNECTION_ERROR or CONNECTION_FAILED, got: %s", resp.Error.Code)
		}
	})
}

func TestSocketIOClient(t *testing.T) {
	t.Run("test Socket.IO client creation", func(t *testing.T) {
		client, err := utils.NewSocketIOClient("ws://localhost:3000")
		if err != nil {
			t.Errorf("NewSocketIOClient should not return error for valid URL, got: %s", err.Error())
		}

		if client == nil {
			t.Errorf("Expected non-nil client")
		}
	})

	t.Run("test Socket.IO client with invalid URL", func(t *testing.T) {
		_, err := utils.NewSocketIOClient("invalid-url")
		if err == nil {
			t.Errorf("Expected error for invalid URL")
		}
	})

	t.Run("test Socket.IO client with headers", func(t *testing.T) {
		headers := map[string]string{"authorization": "Bearer token"}
		client, err := utils.NewSocketIOClient("ws://localhost:3000", utils.WithHeaders(headers))
		if err != nil {
			t.Errorf("NewSocketIOClient with headers should not return error, got: %s", err.Error())
		}

		if client == nil {
			t.Errorf("Expected non-nil client")
		}
	})
}

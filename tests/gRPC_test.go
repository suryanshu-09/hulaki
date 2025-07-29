package tests

import (
	"testing"

	"github.com/suryanshu-09/hulaki/utils"
)

func TestGRPCCall(t *testing.T) {
	t.Run("test gRPC connection error", func(t *testing.T) {
		// Test connection to non-existent server
		resp, err := utils.GRPCCall("localhost:99999", "TestService", "TestMethod")
		if err != nil {
			t.Errorf("GRPCCall should not return error, got: %s", err.Error())
		}

		if resp.Error == nil {
			t.Errorf("Expected connection error, got nil")
		}

		if resp.Error.Code != "CONNECTION_ERROR" && resp.Error.Code != "CONNECTION_FAILED" {
			t.Errorf("Expected CONNECTION_ERROR or CONNECTION_FAILED, got: %s", resp.Error.Code)
		}
	})
}

func TestGRPCReflect(t *testing.T) {
	t.Run("test gRPC reflection connection error", func(t *testing.T) {
		// Test reflection on non-existent server
		resp, err := utils.GRPCReflect("localhost:99999")
		if err != nil {
			t.Errorf("GRPCReflect should not return error, got: %s", err.Error())
		}

		if resp.Error == nil {
			t.Errorf("Expected connection error, got nil")
		}

		if resp.Error.Code != "CONNECTION_ERROR" && resp.Error.Code != "CONNECTION_FAILED" {
			t.Errorf("Expected CONNECTION_ERROR or CONNECTION_FAILED, got: %s", resp.Error.Code)
		}
	})
}

func TestGRPCClient(t *testing.T) {
	t.Run("test gRPC client creation", func(t *testing.T) {
		client, err := utils.NewGRPCClient("localhost:50051")
		if err != nil {
			t.Errorf("NewGRPCClient should not return error for valid address, got: %s", err.Error())
		}

		if client == nil {
			t.Errorf("Expected non-nil client")
		}

		// Test close
		err = client.Close()
		if err != nil {
			t.Errorf("Close should not return error, got: %s", err.Error())
		}
	})

	t.Run("test gRPC client with headers", func(t *testing.T) {
		headers := map[string]string{"authorization": "Bearer token"}
		client, err := utils.NewGRPCClient("localhost:50051", utils.WithHeaders(headers))
		if err != nil {
			t.Errorf("NewGRPCClient with headers should not return error, got: %s", err.Error())
		}

		if client == nil {
			t.Errorf("Expected non-nil client")
		}

		client.Close()
	})
}

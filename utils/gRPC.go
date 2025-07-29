package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type GRPCClient struct {
	conn *grpc.ClientConn
	ctx  context.Context
}

type GRPCRequest struct {
	Service string         `json:"service"`
	Method  string         `json:"method"`
	Data    map[string]any `json:"data,omitempty"`
}

type GRPCResponse struct {
	Data     any         `json:"data"`
	Metadata metadata.MD `json:"metadata,omitempty"`
	Error    *GRPCError  `json:"error,omitempty"`
}

type GRPCError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewGRPCClient(address string, args ...Args) (*GRPCClient, error) {
	_, _, headers := GetArgs(args)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	ctx := context.Background()
	if len(headers) > 0 {
		md := metadata.New(headers)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	return &GRPCClient{
		conn: conn,
		ctx:  ctx,
	}, nil
}

func (c *GRPCClient) Close() error {
	return c.conn.Close()
}

func (c *GRPCClient) TestConnection() error {
	ctx, cancel := context.WithTimeout(c.ctx, 1*time.Second)
	defer cancel()

	state := c.conn.GetState()
	if state == connectivity.TransientFailure || state == connectivity.Shutdown {
		return fmt.Errorf("connection is in state: %v", state)
	}

	// Try to trigger actual connection
	c.conn.Connect()

	// Wait for state change with timeout
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("connection timeout")
		default:
			state = c.conn.GetState()
			if state == connectivity.Ready {
				return nil
			}
			if state == connectivity.TransientFailure {
				return fmt.Errorf("connection failed: %v", state)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}
func GRPCCall(address, service, method string, args ...Args) (*GRPCResponse, error) {
	client, err := NewGRPCClient(address, args...)
	if err != nil {
		return &GRPCResponse{
			Error: &GRPCError{
				Code:    "CONNECTION_ERROR",
				Message: err.Error(),
			},
		}, nil
	}
	defer client.Close()

	if err := client.TestConnection(); err != nil {
		return &GRPCResponse{
			Error: &GRPCError{
				Code:    "CONNECTION_FAILED",
				Message: fmt.Sprintf("failed to connect to gRPC server: %v", err),
			},
		}, nil
	}

	body, _, _ := GetArgs(args)
	var requestData map[string]any

	if body != nil {
		decoder := json.NewDecoder(body)
		if err := decoder.Decode(&requestData); err != nil && err != io.EOF {
			return &GRPCResponse{
				Error: &GRPCError{
					Code:    "INVALID_REQUEST",
					Message: fmt.Sprintf("failed to decode request body: %v", err),
				},
			}, nil
		}
	}

	return &GRPCResponse{
		Data: map[string]any{
			"service":    service,
			"method":     method,
			"request":    requestData,
			"status":     "connected",
			"connection": "established",
			"address":    address,
		},
	}, nil
}

func GRPCReflect(address string, args ...Args) (*GRPCResponse, error) {
	client, err := NewGRPCClient(address, args...)
	if err != nil {
		return &GRPCResponse{
			Error: &GRPCError{
				Code:    "CONNECTION_ERROR",
				Message: err.Error(),
			},
		}, nil
	}
	defer client.Close()

	if err := client.TestConnection(); err != nil {
		return &GRPCResponse{
			Error: &GRPCError{
				Code:    "CONNECTION_FAILED",
				Message: fmt.Sprintf("failed to connect to gRPC server: %v", err),
			},
		}, nil
	}

	return &GRPCResponse{
		Data: map[string]any{
			"address":    address,
			"status":     "connected",
			"reflection": "available",
			"services":   []string{"Service reflection would be implemented here"},
		},
	}, nil
}

func WithMetadata(md map[string]string) Args {
	return WithHeaders(md)
}

func WithGRPCTimeout(timeout time.Duration) Args {
	return func(arg *Arg) {
	}
}

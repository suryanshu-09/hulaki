package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/suryanshu-09/hulaki/utils"
)

func setupGraphQLServer() *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}

		var gqlReq utils.GraphQLRequest
		if err := json.Unmarshal(body, &gqlReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Mock responses based on query content
		if gqlReq.Query == "query GetUser($id: ID!) { user(id: $id) { name email } }" {
			response := utils.GraphQLResponse{
				Data: map[string]interface{}{
					"user": map[string]interface{}{
						"name":  "John Doe",
						"email": "john@example.com",
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		if gqlReq.Query == "mutation CreateUser($input: UserInput!) { createUser(input: $input) { id name } }" {
			response := utils.GraphQLResponse{
				Data: map[string]interface{}{
					"createUser": map[string]interface{}{
						"id":   "123",
						"name": "Jane Doe",
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		if gqlReq.Query == "invalid query" {
			response := utils.GraphQLResponse{
				Errors: []utils.GraphQLError{
					{
						Message: "Syntax Error: Unexpected Name \"invalid\"",
						Locations: []utils.GraphQLErrorLocation{
							{Line: 1, Column: 1},
						},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		// Default query response
		response := utils.GraphQLResponse{
			Data: map[string]interface{}{
				"hello": "world",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	return httptest.NewServer(handler)
}

func TestGraphQLQuery(t *testing.T) {
	server := setupGraphQLServer()
	defer server.Close()

	t.Run("test GraphQL query normally", func(t *testing.T) {
		query := "{ hello }"
		resp, err := utils.GraphQLQuery(server.URL, query)
		if err != nil {
			t.Errorf("got an error: %s", err.Error())
		}
		defer resp.Body.Close()

		gqlResp, err := utils.ParseGraphQLResponse(resp)
		if err != nil {
			t.Errorf("failed to parse GraphQL response: %s", err.Error())
		}

		data, ok := gqlResp.Data.(map[string]interface{})
		if !ok {
			t.Error("data is not a map")
		}

		hello, exists := data["hello"]
		if !exists {
			t.Error("hello field not found in response")
		}

		if hello != "world" {
			t.Errorf("got: %v, want: world", hello)
		}
	})

	t.Run("test GraphQL query with variables", func(t *testing.T) {
		query := "query GetUser($id: ID!) { user(id: $id) { name email } }"
		variables := map[string]interface{}{
			"id": "1",
		}

		resp, err := utils.GraphQLQuery(server.URL, query, utils.WithVariables(variables))
		if err != nil {
			t.Errorf("got an error: %s", err.Error())
		}
		defer resp.Body.Close()

		gqlResp, err := utils.ParseGraphQLResponse(resp)
		if err != nil {
			t.Errorf("failed to parse GraphQL response: %s", err.Error())
		}

		data, ok := gqlResp.Data.(map[string]interface{})
		if !ok {
			t.Error("data is not a map")
		}

		user, exists := data["user"].(map[string]interface{})
		if !exists {
			t.Error("user field not found in response")
		}

		if user["name"] != "John Doe" {
			t.Errorf("got name: %v, want: John Doe", user["name"])
		}

		if user["email"] != "john@example.com" {
			t.Errorf("got email: %v, want: john@example.com", user["email"])
		}
	})

	t.Run("test GraphQL query with headers", func(t *testing.T) {
		query := "{ hello }"
		headers := map[string]string{
			"Authorization": "Bearer test-token",
			"X-Custom":      "test-value",
		}

		resp, err := utils.GraphQLQuery(server.URL, query, utils.WithHeaders(headers))
		if err != nil {
			t.Errorf("got an error: %s", err.Error())
		}
		defer resp.Body.Close()

		// Verify request was processed correctly
		if resp.StatusCode != http.StatusOK {
			t.Errorf("got status code: %d, want: %d", resp.StatusCode, http.StatusOK)
		}
	})

	t.Run("test GraphQL query with params", func(t *testing.T) {
		query := "{ hello }"
		params := map[string]string{
			"version": "v1",
			"format":  "json",
		}

		resp, err := utils.GraphQLQuery(server.URL, query, utils.WithParams(params))
		if err != nil {
			t.Errorf("got an error: %s", err.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("got status code: %d, want: %d", resp.StatusCode, http.StatusOK)
		}
	})
}

func TestGraphQLMutation(t *testing.T) {
	server := setupGraphQLServer()
	defer server.Close()

	t.Run("test GraphQL mutation", func(t *testing.T) {
		mutation := "mutation CreateUser($input: UserInput!) { createUser(input: $input) { id name } }"
		variables := map[string]interface{}{
			"input": map[string]interface{}{
				"name":  "Jane Doe",
				"email": "jane@example.com",
			},
		}

		resp, err := utils.GraphQLMutation(server.URL, mutation, utils.WithVariables(variables))
		if err != nil {
			t.Errorf("got an error: %s", err.Error())
		}
		defer resp.Body.Close()

		gqlResp, err := utils.ParseGraphQLResponse(resp)
		if err != nil {
			t.Errorf("failed to parse GraphQL response: %s", err.Error())
		}

		data, ok := gqlResp.Data.(map[string]interface{})
		if !ok {
			t.Error("data is not a map")
		}

		createUser, exists := data["createUser"].(map[string]interface{})
		if !exists {
			t.Error("createUser field not found in response")
		}

		if createUser["id"] != "123" {
			t.Errorf("got id: %v, want: 123", createUser["id"])
		}

		if createUser["name"] != "Jane Doe" {
			t.Errorf("got name: %v, want: Jane Doe", createUser["name"])
		}
	})
}

func TestGraphQLErrorHandling(t *testing.T) {
	server := setupGraphQLServer()
	defer server.Close()

	t.Run("test GraphQL error response", func(t *testing.T) {
		query := "invalid query"

		resp, err := utils.GraphQLQuery(server.URL, query)
		if err != nil {
			t.Errorf("got an error: %s", err.Error())
		}
		defer resp.Body.Close()

		gqlResp, err := utils.ParseGraphQLResponse(resp)
		if err != nil {
			t.Errorf("failed to parse GraphQL response: %s", err.Error())
		}

		if len(gqlResp.Errors) == 0 {
			t.Error("expected GraphQL errors but got none")
		}

		firstError := gqlResp.Errors[0]
		expectedMessage := "Syntax Error: Unexpected Name \"invalid\""
		if firstError.Message != expectedMessage {
			t.Errorf("got error message: %s, want: %s", firstError.Message, expectedMessage)
		}

		if len(firstError.Locations) == 0 {
			t.Error("expected error location but got none")
		}

		if firstError.Locations[0].Line != 1 {
			t.Errorf("got error line: %d, want: 1", firstError.Locations[0].Line)
		}

		if firstError.Locations[0].Column != 1 {
			t.Errorf("got error column: %d, want: 1", firstError.Locations[0].Column)
		}
	})
}

func TestParseGraphQLResponse(t *testing.T) {
	t.Run("test parse valid GraphQL response", func(t *testing.T) {
		responseJSON := `{
			"data": {
				"user": {
					"name": "Test User",
					"email": "test@example.com"
				}
			}
		}`

		resp := &http.Response{
			Body: io.NopCloser(bytes.NewBufferString(responseJSON)),
		}

		gqlResp, err := utils.ParseGraphQLResponse(resp)
		if err != nil {
			t.Errorf("failed to parse response: %s", err.Error())
		}

		data, ok := gqlResp.Data.(map[string]interface{})
		if !ok {
			t.Error("data is not a map")
		}

		user, exists := data["user"].(map[string]interface{})
		if !exists {
			t.Error("user field not found")
		}

		if user["name"] != "Test User" {
			t.Errorf("got name: %v, want: Test User", user["name"])
		}
	})

	t.Run("test parse invalid JSON", func(t *testing.T) {
		invalidJSON := `{"data": invalid}`

		resp := &http.Response{
			Body: io.NopCloser(bytes.NewBufferString(invalidJSON)),
		}

		_, err := utils.ParseGraphQLResponse(resp)
		if err == nil {
			t.Error("expected error for invalid JSON but got none")
		}
	})
}

func TestWithVariables(t *testing.T) {
	t.Run("test WithVariables function", func(t *testing.T) {
		variables := map[string]interface{}{
			"id":   "123",
			"name": "Test User",
		}

		arg := &utils.Arg{}
		withVars := utils.WithVariables(variables)
		withVars(arg)

		if arg.Body == nil {
			t.Error("expected Body to be set but it was nil")
		}

		buf := new(bytes.Buffer)
		io.Copy(buf, arg.Body)

		var parsedVars map[string]interface{}
		err := json.Unmarshal(buf.Bytes(), &parsedVars)
		if err != nil {
			t.Errorf("failed to parse variables JSON: %s", err.Error())
		}

		if parsedVars["id"] != "123" {
			t.Errorf("got id: %v, want: 123", parsedVars["id"])
		}

		if parsedVars["name"] != "Test User" {
			t.Errorf("got name: %v, want: Test User", parsedVars["name"])
		}
	})
}

package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GraphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type GraphQLResponse struct {
	Data   any            `json:"data"`
	Errors []GraphQLError `json:"errors,omitempty"`
}

type GraphQLError struct {
	Message   string                 `json:"message"`
	Locations []GraphQLErrorLocation `json:"locations,omitempty"`
	Path      []any                  `json:"path,omitempty"`
}

type GraphQLErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

func GraphQLQuery(url string, query string, args ...Args) (*http.Response, error) {
	variables := make(map[string]any)
	_, params, headers := GetArgs(args)

	// Parse variables from args if provided
	for _, arg := range args {
		a := &Arg{}
		arg(a)
		if a.Body != nil {
			buf := new(bytes.Buffer)
			io.Copy(buf, a.Body)
			if buf.Len() > 0 {
				var vars map[string]any
				if err := json.Unmarshal(buf.Bytes(), &vars); err == nil {
					variables = vars
				}
			}
		}
	}

	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	SetParams(&url, params)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	SetHeaders(req, headers)

	client := http.Client{}
	return client.Do(req)
}

func GraphQLMutation(url string, mutation string, args ...Args) (*http.Response, error) {
	return GraphQLQuery(url, mutation, args...)
}

func ParseGraphQLResponse(resp *http.Response) (*GraphQLResponse, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var gqlResp GraphQLResponse
	if err := json.Unmarshal(body, &gqlResp); err != nil {
		return nil, fmt.Errorf("failed to parse GraphQL response: %w", err)
	}

	return &gqlResp, nil
}

func WithVariables(variables map[string]any) Args {
	return func(arg *Arg) {
		jsonVars, _ := json.Marshal(variables)
		arg.Body = bytes.NewBuffer(jsonVars)
	}
}

package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/suryanshu-09/hulaki/styles"
	"github.com/suryanshu-09/hulaki/utils"
)

// graphqlCmd represents the graphql command
var graphqlCmd = &cobra.Command{
	Use:   "graphql",
	Short: "Make a GraphQL request",
	Long: `The 'graphql' command sends GraphQL queries and mutations to a specified endpoint.
GraphQL is a query language for APIs that allows you to request exactly the data you need.
You can include variables and headers to customize the request.`,
	Example: `Examples:
1. Perform a basic GraphQL query:
   hulaki graphql https://api.example.com/graphql --query="{ users { name email } }"

2. Perform a GraphQL query with variables:
   hulaki graphql https://api.example.com/graphql --query="query GetUser($id: ID!) { user(id: $id) { name } }" --variables='{"id":"123"}'

3. Perform a GraphQL mutation:
   hulaki graphql https://api.example.com/graphql --query="mutation CreateUser($input: UserInput!) { createUser(input: $input) { id } }" --variables='{"input":{"name":"John"}}'

4. Perform a GraphQL request with custom headers:
   hulaki graphql https://api.example.com/graphql --query="{ users { name } }" --headers=Authorization=Bearer token`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a GraphQL endpoint URL")
		}

		url := args[0]
		query, err := cmd.Flags().GetString("query")
		if err != nil || query == "" {
			return errors.New("please provide a GraphQL query using --query flag")
		}

		variables, params, headers, err := graphQLIn(cmd)
		if err != nil {
			return err
		}

		var resp *http.Response
		if isMutation(query) {
			resp, err = utils.GraphQLMutation(url, query, utils.WithVariables(variables), utils.WithHeaders(headers), utils.WithParams(params))
		} else {
			resp, err = utils.GraphQLQuery(url, query, utils.WithVariables(variables), utils.WithHeaders(headers), utils.WithParams(params))
		}

		if err != nil {
			return err
		}

		return graphQLOut(cmd, resp)
	},
}

func init() {
	rootCmd.AddCommand(graphqlCmd)

	graphqlCmd.Flags().StringP("query", "q", "", "GraphQL query or mutation string (required)")
	graphqlCmd.Flags().String("variables", "", "GraphQL variables as JSON string")
	graphqlCmd.Flags().String("headers", "", "Custom headers for the GraphQL request, formatted as key=value pairs separated by commas")
	graphqlCmd.Flags().StringP("params", "p", "", "Query parameters for the GraphQL request, formatted as key=value pairs separated by commas")
	graphqlCmd.Flags().BoolP("less", "l", false, "Show only the response data, omitting headers and formatted output")
	graphqlCmd.Flags().Bool("raw", false, "Show raw JSON response without parsing GraphQL structure")
}

func graphQLIn(cmd *cobra.Command) (variables map[string]any, params map[string]string, headers map[string]string, err error) {
	// Parse variables
	variables = make(map[string]any)
	varsStr, err := cmd.Flags().GetString("variables")
	if err == nil && varsStr != "" {
		if varsStr == "-" {
			// Read from stdin
			buf := new(bytes.Buffer)
			io.Copy(buf, os.Stdin)
			varsStr = buf.String()
		}
		if err := json.Unmarshal([]byte(varsStr), &variables); err != nil {
			return nil, nil, nil, fmt.Errorf("invalid variables JSON: %w", err)
		}
	}

	// Parse params
	p, err := cmd.Flags().GetString("params")
	params = make(map[string]string)
	if err == nil && p != "" {
		if strings.Contains(p, ",") {
			paramsArr := strings.Split(p, ",")
			for _, param := range paramsArr {
				if strings.Contains(param, "=") {
					parts := strings.SplitN(param, "=", 2)
					params[parts[0]] = parts[1]
				}
			}
		} else if strings.Contains(p, "=") {
			parts := strings.SplitN(p, "=", 2)
			params[parts[0]] = parts[1]
		}
	}

	// Parse headers
	h, err := cmd.Flags().GetString("headers")
	headers = make(map[string]string)
	if err == nil && h != "" {
		if strings.Contains(h, ",") {
			headersArr := strings.Split(h, ",")
			for _, header := range headersArr {
				if strings.Contains(header, "=") {
					parts := strings.SplitN(header, "=", 2)
					headers[parts[0]] = parts[1]
				}
			}
		} else if strings.Contains(h, "=") {
			parts := strings.SplitN(h, "=", 2)
			headers[parts[0]] = parts[1]
		}
	}

	return variables, params, headers, nil
}

func graphQLOut(cmd *cobra.Command, resp *http.Response) error {
	defer resp.Body.Close()

	raw, _ := cmd.Flags().GetBool("raw")
	less, _ := cmd.Flags().GetBool("less")
	out := cmd.OutOrStdout()

	if raw {
		// Show raw response
		body := bytes.Buffer{}
		_, err := io.Copy(&body, resp.Body)
		if err != nil {
			return err
		}

		if !less {
			fmt.Fprintf(out, "%s\n", styles.Heading.Render("HEADERS"))
			for key, values := range resp.Header {
				for _, value := range values {
					fmt.Fprintf(out, "%s: %s\n", styles.Key.Render(key), value)
				}
			}
			fmt.Fprintf(out, "%s\n", styles.Heading.Render("BODY"))
		}

		fmt.Fprintf(out, "%s\n", styles.Content.Render(body.String()))
		return nil
	}

	// Parse GraphQL response
	gqlResp, err := utils.ParseGraphQLResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to parse GraphQL response: %w", err)
	}

	if !less {
		fmt.Fprintf(out, "%s\n", styles.Heading.Render("HEADERS"))
		for key, values := range resp.Header {
			for _, value := range values {
				fmt.Fprintf(out, "%s: %s\n", styles.Key.Render(key), value)
			}
		}

		// Show GraphQL errors if any
		if len(gqlResp.Errors) > 0 {
			fmt.Fprintf(out, "%s\n", styles.Heading.Render("ERRORS"))
			for _, gqlErr := range gqlResp.Errors {
				fmt.Fprintf(out, "%s\n", styles.Content.Render(gqlErr.Message))
				if len(gqlErr.Locations) > 0 {
					for _, loc := range gqlErr.Locations {
						fmt.Fprintf(out, "  Location: Line %d, Column %d\n", loc.Line, loc.Column)
					}
				}
			}
		}

		fmt.Fprintf(out, "%s\n", styles.Heading.Render("DATA"))
	}

	// Format and display data
	if gqlResp.Data != nil {
		dataJSON, err := json.MarshalIndent(gqlResp.Data, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format data: %w", err)
		}
		fmt.Fprintf(out, "%s\n", styles.Content.Render(string(dataJSON)))
	} else if len(gqlResp.Errors) > 0 && less {
		// In less mode, show errors if no data
		for _, gqlErr := range gqlResp.Errors {
			fmt.Fprintf(out, "Error: %s\n", gqlErr.Message)
		}
	}

	return nil
}

func isMutation(query string) bool {
	trimmed := strings.TrimSpace(strings.ToLower(query))
	return strings.HasPrefix(trimmed, "mutation")
}

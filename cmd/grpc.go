package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/suryanshu-09/hulaki/styles"
	"github.com/suryanshu-09/hulaki/utils"
)

// grpcCmd represents the grpc command
var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Make a gRPC request",
	Long: `The 'grpc' command connects to gRPC servers and makes service calls.
gRPC is a high-performance, language-agnostic RPC framework that uses HTTP/2 and Protocol Buffers.
You can connect to servers, call methods, and send data.`,
	Example: `Examples:
1. Connect to a gRPC server:
   hulaki grpc localhost:50051 --service=UserService --method=GetUser

2. Call a gRPC method with data:
   hulaki grpc localhost:50051 --service=UserService --method=CreateUser --data='{"name":"John","email":"john@example.com"}'

3. Call gRPC with custom headers:
   hulaki grpc localhost:50051 --service=UserService --method=GetUser --headers=authorization=Bearer token

4. Reflect on gRPC services:
   hulaki grpc localhost:50051 --reflect`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a gRPC server address (e.g., localhost:50051)")
		}

		address := args[0]
		reflect, _ := cmd.Flags().GetBool("reflect")

		if reflect {
			return grpcReflect(cmd, address)
		}

		service, err := cmd.Flags().GetString("service")
		if err != nil || service == "" {
			return errors.New("please provide a service name using --service flag")
		}

		method, err := cmd.Flags().GetString("method")
		if err != nil || method == "" {
			return errors.New("please provide a method name using --method flag")
		}

		data, params, headers, err := grpcIn(cmd)
		if err != nil {
			return err
		}

		var grpcArgs []utils.Args
		if len(headers) > 0 {
			grpcArgs = append(grpcArgs, utils.WithHeaders(headers))
		}
		if len(params) > 0 {
			grpcArgs = append(grpcArgs, utils.WithParams(params))
		}
		if data != nil {
			jsonData, _ := json.Marshal(data)
			grpcArgs = append(grpcArgs, utils.WithBody(bytes.NewBuffer(jsonData)))
		}

		resp, err := utils.GRPCCall(address, service, method, grpcArgs...)
		if err != nil {
			return err
		}

		return grpcOut(cmd, resp)
	},
}

func init() {
	rootCmd.AddCommand(grpcCmd)

	grpcCmd.Flags().StringP("service", "s", "", "gRPC service name (required)")
	grpcCmd.Flags().StringP("method", "m", "", "gRPC method name (required)")
	grpcCmd.Flags().String("data", "", "Request data as JSON string")
	grpcCmd.Flags().String("headers", "", "Custom headers/metadata for the gRPC request, formatted as key=value pairs separated by commas")
	grpcCmd.Flags().StringP("params", "p", "", "Query parameters for the gRPC request, formatted as key=value pairs separated by commas")
	grpcCmd.Flags().BoolP("less", "l", false, "Show only the response data, omitting headers and formatted output")
	grpcCmd.Flags().Bool("reflect", false, "Reflect on available gRPC services")
}

func grpcIn(cmd *cobra.Command) (data map[string]any, params map[string]string, headers map[string]string, err error) {
	// Parse data
	data = make(map[string]any)
	dataStr, err := cmd.Flags().GetString("data")
	if err == nil && dataStr != "" {
		if dataStr == "-" {
			// Read from stdin
			buf := new(bytes.Buffer)
			io.Copy(buf, os.Stdin)
			dataStr = buf.String()
		}
		if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
			return nil, nil, nil, fmt.Errorf("invalid data JSON: %w", err)
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

	return data, params, headers, nil
}

func grpcOut(cmd *cobra.Command, resp *utils.GRPCResponse) error {
	less, _ := cmd.Flags().GetBool("less")
	out := cmd.OutOrStdout()

	if resp.Error != nil {
		fmt.Fprintf(out, "%s\n", styles.Heading.Render("ERROR"))
		fmt.Fprintf(out, "%s: %s\n", styles.Key.Render("Code"), resp.Error.Code)
		fmt.Fprintf(out, "%s: %s\n", styles.Key.Render("Message"), resp.Error.Message)
		return nil
	}

	if !less {
		fmt.Fprintf(out, "%s\n", styles.Heading.Render("GRPC RESPONSE"))

		if resp.Metadata != nil && len(resp.Metadata) > 0 {
			fmt.Fprintf(out, "%s\n", styles.Heading.Render("METADATA"))
			for key, values := range resp.Metadata {
				for _, value := range values {
					fmt.Fprintf(out, "%s: %s\n", styles.Key.Render(key), value)
				}
			}
		}

		fmt.Fprintf(out, "%s\n", styles.Heading.Render("DATA"))
	}

	// Format and display data
	if resp.Data != nil {
		dataJSON, err := json.MarshalIndent(resp.Data, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format data: %w", err)
		}
		fmt.Fprintf(out, "%s\n", styles.Content.Render(string(dataJSON)))
	}

	return nil
}

func grpcReflect(cmd *cobra.Command, address string) error {
	resp, err := utils.GRPCReflect(address)
	if err != nil {
		return err
	}

	return grpcOut(cmd, resp)
}

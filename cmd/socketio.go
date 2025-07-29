package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/suryanshu-09/hulaki/styles"
	"github.com/suryanshu-09/hulaki/utils"
)

// socketioCmd represents the socketio command
var socketioCmd = &cobra.Command{
	Use:   "socketio",
	Short: "Connect to Socket.IO servers",
	Long: `The 'socketio' command connects to Socket.IO servers for real-time communication.
Socket.IO enables real-time bidirectional event-based communication between clients and servers.
You can connect, emit events, and listen for events.`,
	Example: `Examples:
1. Connect to a Socket.IO server:
   hulaki socketio ws://localhost:3000

2. Emit an event to the server:
   hulaki socketio ws://localhost:3000 --emit=message --data='{"text":"Hello World"}'

3. Listen for events from the server:
   hulaki socketio ws://localhost:3000 --listen=notification --duration=10s

4. Connect with custom headers:
   hulaki socketio ws://localhost:3000 --headers=authorization=Bearer token

5. Connect to a specific namespace:
   hulaki socketio ws://localhost:3000 --namespace=/chat`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a Socket.IO server URL (e.g., ws://localhost:3000)")
		}

		serverURL := args[0]
		emit, _ := cmd.Flags().GetString("emit")
		listen, _ := cmd.Flags().GetString("listen")
		durationStr, _ := cmd.Flags().GetString("duration")

		if emit != "" {
			return socketIOEmit(cmd, serverURL, emit)
		}

		if listen != "" {
			duration, err := time.ParseDuration(durationStr)
			if err != nil {
				duration = 5 * time.Second // Default duration
			}
			return socketIOListen(cmd, serverURL, listen, duration)
		}

		// Default: just connect and test
		return socketIOConnect(cmd, serverURL)
	},
}

func init() {
	rootCmd.AddCommand(socketioCmd)

	socketioCmd.Flags().String("emit", "", "Event name to emit to the server")
	socketioCmd.Flags().String("listen", "", "Event name to listen for from the server")
	socketioCmd.Flags().String("duration", "5s", "Duration to listen for events (e.g., 10s, 1m)")
	socketioCmd.Flags().String("data", "", "Event data as JSON string")
	socketioCmd.Flags().String("headers", "", "Custom headers for the Socket.IO connection, formatted as key=value pairs separated by commas")
	socketioCmd.Flags().String("namespace", "", "Socket.IO namespace to connect to")
	socketioCmd.Flags().StringP("params", "p", "", "Query parameters for the Socket.IO connection, formatted as key=value pairs separated by commas")
	socketioCmd.Flags().BoolP("less", "l", false, "Show only the response data, omitting headers and formatted output")
}

func socketioIn(cmd *cobra.Command) (data map[string]any, params map[string]string, headers map[string]string, err error) {
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

	// Add namespace to params if specified
	namespace, err := cmd.Flags().GetString("namespace")
	if err == nil && namespace != "" {
		params["namespace"] = namespace
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

func socketioOut(cmd *cobra.Command, resp *utils.SocketIOResponse) error {
	less, _ := cmd.Flags().GetBool("less")
	out := cmd.OutOrStdout()

	if resp.Error != nil {
		fmt.Fprintf(out, "%s\n", styles.Heading.Render("ERROR"))
		fmt.Fprintf(out, "%s: %s\n", styles.Key.Render("Code"), resp.Error.Code)
		fmt.Fprintf(out, "%s: %s\n", styles.Key.Render("Message"), resp.Error.Message)
		return nil
	}

	if !less {
		fmt.Fprintf(out, "%s\n", styles.Heading.Render("SOCKET.IO RESPONSE"))
		fmt.Fprintf(out, "%s: %s\n", styles.Key.Render("Status"), resp.Status)
		fmt.Fprintf(out, "%s: %t\n", styles.Key.Render("Connected"), resp.Connected)
		if resp.Event != "" {
			fmt.Fprintf(out, "%s: %s\n", styles.Key.Render("Event"), resp.Event)
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

func socketIOConnect(cmd *cobra.Command, serverURL string) error {
	_, params, headers, err := socketioIn(cmd)
	if err != nil {
		return err
	}

	var socketIOArgs []utils.Args
	if len(headers) > 0 {
		socketIOArgs = append(socketIOArgs, utils.WithHeaders(headers))
	}
	if len(params) > 0 {
		socketIOArgs = append(socketIOArgs, utils.WithParams(params))
	}

	resp, err := utils.SocketIOConnect(serverURL, socketIOArgs...)
	if err != nil {
		return err
	}

	return socketioOut(cmd, resp)
}

func socketIOEmit(cmd *cobra.Command, serverURL, event string) error {
	data, params, headers, err := socketioIn(cmd)
	if err != nil {
		return err
	}

	var socketIOArgs []utils.Args
	if len(headers) > 0 {
		socketIOArgs = append(socketIOArgs, utils.WithHeaders(headers))
	}
	if len(params) > 0 {
		socketIOArgs = append(socketIOArgs, utils.WithParams(params))
	}
	if data != nil && len(data) > 0 {
		jsonData, _ := json.Marshal(data)
		socketIOArgs = append(socketIOArgs, utils.WithBody(bytes.NewBuffer(jsonData)))
	}

	resp, err := utils.SocketIOEmit(serverURL, event, socketIOArgs...)
	if err != nil {
		return err
	}

	return socketioOut(cmd, resp)
}

func socketIOListen(cmd *cobra.Command, serverURL, event string, duration time.Duration) error {
	_, params, headers, err := socketioIn(cmd)
	if err != nil {
		return err
	}

	var socketIOArgs []utils.Args
	if len(headers) > 0 {
		socketIOArgs = append(socketIOArgs, utils.WithHeaders(headers))
	}
	if len(params) > 0 {
		socketIOArgs = append(socketIOArgs, utils.WithParams(params))
	}

	resp, err := utils.SocketIOListen(serverURL, event, duration, socketIOArgs...)
	if err != nil {
		return err
	}

	return socketioOut(cmd, resp)
}

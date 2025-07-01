package cmd

import (
	"bytes"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/v2/textarea"
	"github.com/charmbracelet/bubbles/v2/viewport"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"github.com/suryanshu-09/hulaki/styles"
	"github.com/suryanshu-09/hulaki/utils"
)

var wsCmd = &cobra.Command{
	Use:   "ws",
	Short: "Interactive WebSocket client",
	Long: `The "ws" command allows you to connect to a WebSocket server and interact with it in real-time.
You can send and receive messages using a terminal-based user interface.`,
	Example: `Examples:
1. Perform a basic GET request:
   hulaki http get https://example.com

2. Perform a GET request with query parameters:
   hulaki http get https://api.example.com/data --params=type=user,status=active

3. Perform a GET request with custom headers:
   hulaki http get https://api.example.com/data --headers=Authorization=BearerToken,Accept=application/json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params, headers, err := WebsocketIn(cmd, args)
		if err != nil {
			return err
		}
		url := args[0]
		ws, err := utils.NewWebsocketClient(url, utils.WithHeaders(headers), utils.WithParams(params))
		if err != nil {
			return err
		}
		defer ws.Close()

		wc := NewWebsocketCli(ws)

		go func() {
			for {
				recieve := new(bytes.Buffer)
				err := ws.Read(recieve)
				if err != nil {
					log.Println("Error reading message:", err)
					return
				}
				wc.AddMessage(recieve.String())
			}
		}()

		if _, err := tea.NewProgram(wc, tea.WithMouseAllMotion()).Run(); err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(wsCmd)

	wsCmd.Flags().String("headers", "", "Custom headers for the WebSocket connection, formatted as key=value pairs separated by commas")
	wsCmd.Flags().StringP("params", "p", "", "Query parameters for the WebSocket connection, formatted as key=value pairs separated by commas")
}

var termHeight, termWidth = 0, 0

type WebsocketCli struct {
	Input    textarea.Model
	ViewPort viewport.Model
	WS       *utils.WebsocketClient
	messages []string
}

func (wc *WebsocketCli) Init() tea.Cmd {
	return nil
}

func (wc *WebsocketCli) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		termHeight = msg.Height
		termWidth = msg.Width

		wc.Input.SetWidth(msg.Width - 4)
		wc.ViewPort.SetWidth(msg.Width - 2)
		wc.ViewPort.SetHeight(msg.Height / 5)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return wc, tea.Quit
		case "enter":
			send := new(bytes.Buffer)
			send.WriteString(wc.Input.Value())
			err := wc.WS.Write(websocket.TextMessage, send)
			if err != nil {
				log.Println("Error sending message:", err)
			}
			wc.Input.Reset()
		}
	}

	var cmds []tea.Cmd

	updatedInput, cmd := wc.Input.Update(msg)
	cmds = append(cmds, cmd)
	wc.Input = updatedInput

	updatedView, cmd := wc.ViewPort.Update(msg)
	cmds = append(cmds, cmd)
	wc.ViewPort = updatedView

	return wc, tea.Batch(cmds...)
}

func (wc *WebsocketCli) View() string {
	WsOutputStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Width(termWidth - 2).BorderForeground(lipgloss.Color("#8a2be2"))
	return lipgloss.JoinVertical(lipgloss.Left, styles.Key.Render("Send a message:"), WsOutputStyle.Render(wc.Input.View()), "\n", WsOutputStyle.Render(wc.ViewPort.View()))
}

func (wc *WebsocketCli) AddMessage(message string) {
	wc.messages = append(wc.messages, message)
	wc.ViewPort.SetContent(wc.formatMessages())
}

func (wc *WebsocketCli) formatMessages() string {
	var buffer bytes.Buffer
	for _, msg := range wc.messages {
		buffer.WriteString(msg + "\n")
	}
	return buffer.String()
}

func NewWebsocketCli(ws *utils.WebsocketClient) *WebsocketCli {
	ta := textarea.New()
	ta.CharLimit = 240
	ta.Styles = textarea.DefaultDarkStyles()
	ta.SetHeight(1)
	ta.VirtualCursor = true
	ta.ShowLineNumbers = false
	ta.Placeholder = "msg..."
	ta.Styles.Cursor.Shape = tea.CursorBlock
	ta.Styles.Cursor.Blink = true
	ta.Styles.Cursor.BlinkSpeed = 2 * time.Second
	ta.Styles.Cursor.Color = lipgloss.Color("#ff1493")
	ta.Focus()

	v := viewport.New()
	v.SetWidth(termWidth - 2)
	v.SetHeight(termHeight / 5)
	v.FillHeight = true
	v.MouseWheelEnabled = true
	v.SoftWrap = true
	v.KeyMap = viewport.KeyMap{}

	return &WebsocketCli{
		Input:    ta,
		ViewPort: v,
		WS:       ws,
		messages: []string{},
	}
}

func WebsocketIn(cmd *cobra.Command, args []string) (params map[string]string, headers map[string]string, err error) {
	p, err := cmd.Flags().GetString("params")
	params = make(map[string]string, 0)
	if err == nil {
		if i := strings.Index(p, ","); i != -1 {
			paramsArr := strings.SplitSeq(p, ",")
			for param := range paramsArr {
				if i := strings.Index(param, "="); i != -1 {
					temp := strings.Split(param, "=")
					params[temp[0]] = temp[1]
				}
			}
		} else {
			if i := strings.Index(p, "="); i != -1 {
				temp := strings.Split(p, "=")
				params[temp[0]] = temp[1]
			}
		}
	}

	h, _ := cmd.Flags().GetString("headers")
	headers = make(map[string]string, 0)
	if err == nil {
		if i := strings.Index(h, ","); i != -1 {
			headersArr := strings.SplitSeq(h, ",")
			for header := range headersArr {
				if i := strings.Index(header, "="); i != -1 {
					temp := strings.Split(header, "=")
					headers[temp[0]] = temp[1]
				}
			}
		} else {
			if i := strings.Index(h, "="); i != -1 {
				temp := strings.Split(h, "=")
				headers[temp[0]] = temp[1]
			}
		}
	}

	if len(args) < 1 {
		return nil, nil, errors.New("please provide a url")
	}
	return
}

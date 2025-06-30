package cmd

import (
	"bytes"
	"errors"
	"log"
	"sync"
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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide a url")
		}
		url := args[0]
		ws, err := utils.NewWebsocketClient(url)
		if err != nil {
			return err
		}
		defer ws.Close()

		wc := NewWebsocketCli(ws)
		wg := sync.WaitGroup{}

		wg.Add(1)
		go func() {
			defer wg.Done()
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

		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := tea.NewProgram(wc, tea.WithMouseAllMotion()).Run(); err != nil {
				log.Fatal(err)
			}
		}()

		wg.Wait()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(wsCmd)
}

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
		wc.Input.SetWidth(msg.Width - 2)
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
			wc.Input.SetValue("")
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
	WsOutputStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	return lipgloss.JoinVertical(lipgloss.Left, WsOutputStyle.Render(wc.Input.View()), "\n", WsOutputStyle.Render(wc.ViewPort.View()))
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
	ta.Prompt = styles.Key.Render("Send a message:")
	ta.Placeholder = "msg..."
	ta.SetHeight(1)
	ta.Styles.Cursor.Shape = tea.CursorBlock
	ta.Styles.Cursor.Blink = true
	ta.Styles.Cursor.BlinkSpeed = 2 * time.Second
	ta.Styles.Cursor.Color = lipgloss.Color("#ff1493")
	ta.Focus()

	v := viewport.New()
	v.FillHeight = true
	v.MouseWheelEnabled = true
	v.SoftWrap = true

	return &WebsocketCli{
		Input:    ta,
		ViewPort: v,
		WS:       ws,
		messages: []string{},
	}
}

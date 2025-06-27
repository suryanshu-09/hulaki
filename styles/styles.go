/*
Package styles
Styles for both cli and tui
*/
package styles

import "github.com/charmbracelet/lipgloss/v2"

var (
	Heading = lipgloss.NewStyle().Background(lipgloss.Color("#ff1493")).Bold(true).Margin(1).Padding(0, 1, 0, 1)
	Content = lipgloss.NewStyle().Margin(0, 0, 0, 1)
	Key     = lipgloss.NewStyle().Foreground(lipgloss.Color("#14ff82")).Bold(true).Margin(0, 0, 0, 1)
)

package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/MarquesCoding/DockerViewer/internal/styles"
)

type StatusBarModel struct {
	width   int
	section Section
	tab     int
	err     error
}

type Pane int

const (
	PaneLeft Pane = iota
	PaneRight
)

func NewStatusBar() StatusBarModel { return StatusBarModel{} }

func (m *StatusBarModel) SetWidth(w int)       { m.width = w }
func (m *StatusBarModel) SetSection(s Section)  { m.section = s }
func (m *StatusBarModel) SetTab(t int)          { m.tab = t }
func (m *StatusBarModel) SetError(err error)    { m.err = err }

func (m StatusBarModel) View() string {
	key := func(k, desc string) string {
		return styles.StyleKey.Render(k) +
			lipgloss.NewStyle().Foreground(styles.ColorHeader).Render(":"+desc)
	}

	hints := []string{
		key("←/→", "section"),
		key("↑↓", "navigate"),
		key("a/d", "panels"),
		key("s", "actions"),
	}
	if m.section == SectionContainers {
		hints = append(hints, key("[/]", "tabs"))
		if m.tab == 1 {
			hints = append(hints, key("f", "follow"))
		}
	}
	hints = append(hints, key("r", "refresh"), key("q", "quit"))

	left := strings.Join(hints, styles.StyleMuted.Render("  "))

	right := lipgloss.NewStyle().Foreground(styles.ColorFocused).Bold(true).Render("DockerViewer")
	if m.err != nil {
		right += lipgloss.NewStyle().Foreground(styles.ColorStopped).Render(" ✗ " + m.err.Error())
	}

	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right) - 2
	if gap < 1 {
		gap = 1
	}

	return styles.StyleStatusBar.Width(m.width).Render(left + strings.Repeat(" ", gap) + right)
}

package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/MarquesCoding/DockerViewer/internal/styles"
)

// Action represents one item in the popup menu.
type Action struct {
	Label       string
	ID          string // machine-readable key passed back to the caller
	Destructive bool
}

type PopupState int

const (
	PopupStateList PopupState = iota
	PopupStateConfirm
)

// PopupModel is an overlay action menu. It is embedded by value in the root
// model so no pointer aliasing is needed.
type PopupModel struct {
	title   string
	actions []Action
	cursor  int
	state   PopupState // list vs confirm
	width   int
	height  int
}

func NewPopup(title string, actions []Action) PopupModel {
	return PopupModel{title: title, actions: actions}
}

func (m *PopupModel) SetSize(w, h int) { m.width = w; m.height = h }

func (m *PopupModel) MoveUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m *PopupModel) MoveDown() {
	if m.cursor < len(m.actions)-1 {
		m.cursor++
	}
}

// Select returns the currently highlighted action and whether it needs confirmation.
func (m *PopupModel) Select() (Action, bool) {
	if m.cursor >= len(m.actions) {
		return Action{}, false
	}
	a := m.actions[m.cursor]
	return a, a.Destructive
}

// Confirm advances the popup into confirm state and returns the chosen action.
func (m *PopupModel) Confirm() Action {
	m.state = PopupStateConfirm
	return m.actions[m.cursor]
}

func (m PopupModel) NeedsConfirm() bool { return m.state == PopupStateConfirm }

func (m *PopupModel) ResetConfirm() { m.state = PopupStateList }

func (m PopupModel) View() string {
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorFocused).
		Padding(0, 1)

	titleStyle := lipgloss.NewStyle().
		Foreground(styles.ColorFocused).
		Bold(true).
		MarginBottom(1)

	if m.state == PopupStateConfirm {
		return m.confirmView(borderStyle, titleStyle)
	}
	return m.listView(borderStyle, titleStyle)
}

func (m PopupModel) listView(border, title lipgloss.Style) string {
	var b strings.Builder
	b.WriteString(title.Render(m.title) + "\n")

	for i, action := range m.actions {
		label := action.Label
		if action.Destructive {
			label = styles.StyleStopped.Render(label)
		} else {
			label = lipgloss.NewStyle().Foreground(styles.ColorText).Render(label)
		}

		if i == m.cursor {
			prefix := styles.StyleSelected.Render(" ▶ ")
			b.WriteString(prefix + label + "\n")
		} else {
			b.WriteString(styles.StyleMuted.Render("   ") + label + "\n")
		}
	}

	b.WriteString("\n" + styles.StyleMuted.Render("↑↓ navigate  enter select  esc cancel"))

	inner := border.Render(b.String())
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, inner)
}

func (m PopupModel) confirmView(border, title lipgloss.Style) string {
	action := m.actions[m.cursor]

	var b strings.Builder
	b.WriteString(title.Render("Confirm") + "\n")
	b.WriteString(lipgloss.NewStyle().Foreground(styles.ColorText).Render(
		"Really "+strings.ToLower(action.Label)+"?",
	) + "\n\n")
	b.WriteString(styles.StyleStopped.Bold(true).Render("  [enter] Yes") + "  ")
	b.WriteString(styles.StyleMuted.Render("[esc] Cancel"))

	inner := border.Render(b.String())
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, inner)
}

// ContainerActions returns the action list for a container given its current state.
func ContainerActions(state string) []Action {
	switch state {
	case "running":
		return []Action{
			{Label: "Stop", ID: "stop"},
			{Label: "Restart", ID: "restart"},
			{Label: "Pause", ID: "pause"},
			{Label: "Kill", ID: "kill", Destructive: true},
			{Label: "Remove (force)", ID: "remove", Destructive: true},
		}
	case "paused":
		return []Action{
			{Label: "Unpause", ID: "unpause"},
			{Label: "Stop", ID: "stop"},
			{Label: "Remove (force)", ID: "remove", Destructive: true},
		}
	case "exited", "created", "dead":
		return []Action{
			{Label: "Start", ID: "start"},
			{Label: "Remove", ID: "remove", Destructive: true},
		}
	default:
		return []Action{
			{Label: "Remove (force)", ID: "remove", Destructive: true},
		}
	}
}

func ImageActions() []Action {
	return []Action{
		{Label: "Remove", ID: "remove", Destructive: true},
	}
}

func NetworkActions() []Action {
	return []Action{
		{Label: "Remove", ID: "remove", Destructive: true},
	}
}

func VolumeActions() []Action {
	return []Action{
		{Label: "Remove", ID: "remove", Destructive: true},
	}
}

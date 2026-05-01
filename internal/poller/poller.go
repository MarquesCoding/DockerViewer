package poller

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	DefaultInterval = 3 * time.Second
	StatsInterval   = 2 * time.Second
)

type TickMsg struct{ T time.Time }

func TickCmd(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return TickMsg{T: t}
	})
}

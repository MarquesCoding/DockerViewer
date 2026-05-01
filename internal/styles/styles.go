package styles

import "github.com/charmbracelet/lipgloss"

var (
	ColorRunning  = lipgloss.Color("#00FF7F")
	ColorStopped  = lipgloss.Color("#FF6B6B")
	ColorPaused   = lipgloss.Color("#FFD700")
	ColorSelected = lipgloss.Color("#7D56F4")
	ColorBorder   = lipgloss.Color("#383838")
	ColorFocused  = lipgloss.Color("#7D56F4")
	ColorHeader   = lipgloss.Color("#888888")
	ColorMuted    = lipgloss.Color("#555555")
	ColorText     = lipgloss.Color("#DDDDDD")
	ColorKey      = lipgloss.Color("#7D56F4")

	StyleRunning = lipgloss.NewStyle().Foreground(ColorRunning)
	StyleStopped = lipgloss.NewStyle().Foreground(ColorStopped)
	StylePaused  = lipgloss.NewStyle().Foreground(ColorPaused)

	StyleSelected = lipgloss.NewStyle().
			Background(ColorSelected).
			Foreground(lipgloss.Color("#FFFFFF"))

	StylePanelBorder = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorBorder)

	StylePanelBorderFocused = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorFocused)

	StyleSectionHeader = lipgloss.NewStyle().
				Foreground(ColorHeader).
				Bold(true)

	StyleTabActive = lipgloss.NewStyle().
			Foreground(ColorFocused).
			Bold(true).
			Underline(true)

	StyleTabInactive = lipgloss.NewStyle().
				Foreground(ColorMuted)

	StyleStatusBar = lipgloss.NewStyle().
			Background(lipgloss.Color("#1a1a2e")).
			Foreground(lipgloss.Color("#888888")).
			Padding(0, 1)

	StyleKey = lipgloss.NewStyle().
			Foreground(ColorKey).
			Bold(true)

	StyleLabel = lipgloss.NewStyle().
			Foreground(ColorHeader).
			Width(12)

	StyleMuted = lipgloss.NewStyle().Foreground(ColorMuted)
)

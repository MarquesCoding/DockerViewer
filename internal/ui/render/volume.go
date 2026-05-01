package render

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/MarquesCoding/DockerViewer/internal/docker"
	"github.com/MarquesCoding/DockerViewer/internal/styles"
)

func Volume(detail docker.VolumeDetail, width int) string {
	var b strings.Builder
	info := detail.Info

	b.WriteString(styles.StyleSectionHeader.Render("  Volume Info") + "\n\n")

	row := func(label, value string) string {
		return styles.StyleLabel.Render(label) +
			lipgloss.NewStyle().Foreground(styles.ColorText).Render(value)
	}

	name := info.Name
	maxName := width - 16
	if maxName > 0 && len(name) > maxName {
		name = name[:maxName] + "..."
	}
	b.WriteString(row("Name:", name) + "\n")
	b.WriteString(row("Driver:", info.Driver) + "\n")
	b.WriteString(row("Scope:", info.Scope) + "\n")
	b.WriteString(row("Mountpoint:", info.Mountpoint) + "\n")

	if len(info.Labels) > 0 {
		b.WriteString("\n" + styles.StyleSectionHeader.Render("  Labels") + "\n")
		for k, v := range info.Labels {
			b.WriteString("  " + styles.StyleMuted.Render(k+"="+v) + "\n")
		}
	}

	_ = strings.Builder{}
	return b.String()
}

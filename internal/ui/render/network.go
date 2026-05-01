package render

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/MarquesCoding/DockerViewer/internal/docker"
	"github.com/MarquesCoding/DockerViewer/internal/styles"
)

func Network(detail docker.NetworkDetail, _ int) string {
	var b strings.Builder
	info := detail.Info

	b.WriteString(styles.StyleSectionHeader.Render("  Network Info") + "\n\n")

	row := func(label, value string) string {
		return styles.StyleLabel.Render(label) +
			lipgloss.NewStyle().Foreground(styles.ColorText).Render(value)
	}

	id := info.ID
	if len(id) > 12 {
		id = id[:12]
	}
	b.WriteString(row("Name:", info.Name) + "\n")
	b.WriteString(row("ID:", id) + "\n")
	b.WriteString(row("Driver:", info.Driver) + "\n")
	b.WriteString(row("Scope:", info.Scope) + "\n")

	if info.IPAM.Config != nil {
		for _, cfg := range info.IPAM.Config {
			if cfg.Subnet != "" {
				b.WriteString(row("Subnet:", cfg.Subnet) + "\n")
			}
			if cfg.Gateway != "" {
				b.WriteString(row("Gateway:", cfg.Gateway) + "\n")
			}
		}
	}

	if len(info.Containers) > 0 {
		b.WriteString("\n" + styles.StyleSectionHeader.Render("  Containers") + "\n\n")
		for _, ctr := range info.Containers {
			b.WriteString(fmt.Sprintf("  %s  %s\n",
				styles.StyleRunning.Render(ctr.Name),
				styles.StyleMuted.Render(ctr.IPv4Address),
			))
		}
	}

	_ = strings.Builder{}
	return b.String()
}

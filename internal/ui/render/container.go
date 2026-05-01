package render

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/MarquesCoding/DockerViewer/internal/docker"
	"github.com/MarquesCoding/DockerViewer/internal/styles"
)

func Container(detail docker.ContainerDetail, logs []string, activeTab int, width int) string {
	switch activeTab {
	case 1:
		return containerLogs(logs, width)
	case 2:
		return containerConfig(detail, width)
	default:
		return containerStats(detail, width)
	}
}

func containerStats(detail docker.ContainerDetail, width int) string {
	info := detail.Info
	st := detail.Stats

	var b strings.Builder
	b.WriteString(styles.StyleSectionHeader.Render("  Stats") + "\n\n")

	row := func(label, value string) string {
		return styles.StyleLabel.Render(label) +
			lipgloss.NewStyle().Foreground(styles.ColorText).Render(value)
	}

	state := info.State
	stateStr := "unknown"
	stateStyle := lipgloss.NewStyle()
	if state != nil {
		stateStr = state.Status
		switch state.Status {
		case "running":
			stateStyle = styles.StyleRunning
		case "exited":
			stateStyle = styles.StyleStopped
		case "paused":
			stateStyle = styles.StylePaused
		}
	}
	b.WriteString(row("Status:", stateStyle.Render(stateStr)) + "\n")
	b.WriteString(row("CPU:", fmt.Sprintf("%.2f%%", st.CPUPercent)) + "\n")
	b.WriteString(row("Memory:", fmt.Sprintf("%s / %s (%.1f%%)",
		formatBytes(st.MemUsed), formatBytes(st.MemMax), st.MemPercent)) + "\n")

	if info.NetworkSettings != nil {
		for name, net := range info.NetworkSettings.Networks {
			b.WriteString(row("Network:", fmt.Sprintf("%s  IP: %s", name, net.IPAddress)) + "\n")
		}
	}

	if state != nil && state.StartedAt != "" {
		if t, err := time.Parse(time.RFC3339Nano, state.StartedAt); err == nil {
			b.WriteString(row("Started:", t.Local().Format("2006-01-02 15:04:05")) + "\n")
		}
	}
	return b.String()
}

func containerLogs(logs []string, width int) string {
	if len(logs) == 0 {
		return styles.StyleMuted.Render("  No logs yet...")
	}
	var b strings.Builder
	for _, line := range logs {
		if len(line) > width-4 {
			line = line[:width-4]
		}
		b.WriteString(line + "\n")
	}
	return b.String()
}

func containerConfig(detail docker.ContainerDetail, _ int) string {
	info := detail.Info
	var b strings.Builder

	b.WriteString(styles.StyleSectionHeader.Render("  Config") + "\n\n")

	row := func(label, value string) string {
		return styles.StyleLabel.Render(label) +
			lipgloss.NewStyle().Foreground(styles.ColorText).Render(value)
	}

	b.WriteString(row("Name:", strings.TrimPrefix(info.Name, "/")) + "\n")
	b.WriteString(row("ID:", info.ID[:12]) + "\n")

	if info.Config != nil {
		b.WriteString(row("Image:", info.Config.Image) + "\n")
		if len(info.Config.Cmd) > 0 {
			b.WriteString(row("Command:", strings.Join(info.Config.Cmd, " ")) + "\n")
		}
		if len(info.Config.Entrypoint) > 0 {
			b.WriteString(row("Entrypoint:", strings.Join(info.Config.Entrypoint, " ")) + "\n")
		}
		b.WriteString(row("WorkingDir:", info.Config.WorkingDir) + "\n")

		if len(info.Config.Env) > 0 {
			b.WriteString("\n" + styles.StyleSectionHeader.Render("  Environment") + "\n")
			for _, env := range info.Config.Env {
				b.WriteString("  " + styles.StyleMuted.Render(env) + "\n")
			}
		}
	}

	if info.HostConfig != nil && len(info.HostConfig.PortBindings) > 0 {
		b.WriteString("\n" + styles.StyleSectionHeader.Render("  Ports") + "\n")
		for port := range info.HostConfig.PortBindings {
			for _, binding := range info.HostConfig.PortBindings[port] {
				b.WriteString(fmt.Sprintf("  %s -> %s:%s\n", port, binding.HostIP, binding.HostPort))
			}
		}
	}

	return b.String()
}

func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

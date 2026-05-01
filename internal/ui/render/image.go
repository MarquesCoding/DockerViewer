package render

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/MarquesCoding/DockerViewer/internal/docker"
	"github.com/MarquesCoding/DockerViewer/internal/styles"
)

func Image(detail docker.ImageDetail, width int) string {
	var b strings.Builder
	inspect := detail.Inspect

	b.WriteString(styles.StyleSectionHeader.Render("  Image Info") + "\n\n")

	row := func(label, value string) string {
		return styles.StyleLabel.Render(label) +
			lipgloss.NewStyle().Foreground(styles.ColorText).Render(value)
	}

	id := inspect.ID
	if len(id) > 19 {
		id = id[7:19]
	}
	b.WriteString(row("ID:", id) + "\n")
	if len(inspect.RepoTags) > 0 {
		b.WriteString(row("Tags:", strings.Join(inspect.RepoTags, ", ")) + "\n")
	}
	b.WriteString(row("Size:", formatBytes(uint64(inspect.Size))) + "\n")

	if inspect.Created != "" {
		if t, err := time.Parse(time.RFC3339Nano, inspect.Created); err == nil {
			b.WriteString(row("Created:", t.Local().Format("2006-01-02 15:04:05")) + "\n")
		}
	}

	if inspect.Config != nil {
		if len(inspect.Config.Cmd) > 0 {
			b.WriteString(row("Command:", strings.Join(inspect.Config.Cmd, " ")) + "\n")
		}
		if inspect.Config.WorkingDir != "" {
			b.WriteString(row("WorkDir:", inspect.Config.WorkingDir) + "\n")
		}
	}

	if len(detail.History) > 0 {
		b.WriteString("\n" + styles.StyleSectionHeader.Render("  Layers") + "\n\n")

		idW, sizeW := 14, 10
		colID := lipgloss.NewStyle().Width(idW).Foreground(styles.ColorMuted)
		colSize := lipgloss.NewStyle().Width(sizeW).Foreground(styles.ColorPaused)
		colCmd := lipgloss.NewStyle().Foreground(styles.ColorText)

		hdr := colID.Render("ID") + colSize.Render("SIZE") + colCmd.Render("COMMAND")
		b.WriteString("  " + styles.StyleSectionHeader.Render(hdr) + "\n")

		for _, layer := range detail.History {
			id := "<missing>"
			if layer.ID != "<missing>" && len(layer.ID) > 19 {
				id = layer.ID[7:19]
			} else if layer.ID != "<missing>" {
				id = layer.ID
			}

			cmd := layer.CreatedBy
			maxCmd := width - idW - sizeW - 6
			if maxCmd > 0 && len(cmd) > maxCmd {
				cmd = cmd[:maxCmd]
			}

			b.WriteString("  " + colID.Render(id) + colSize.Render(formatBytes(uint64(layer.Size))) + colCmd.Render(cmd) + "\n")
		}
	}

	_ = fmt.Sprintf
	return b.String()
}

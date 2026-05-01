package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/MarquesCoding/DockerViewer/internal/docker"
	"github.com/MarquesCoding/DockerViewer/internal/styles"
)

type Section int

const (
	SectionContainers Section = iota
	SectionImages
	SectionNetworks
	SectionVolumes
	sectionCount
)

func (s Section) String() string {
	return [...]string{"Containers", "Images", "Networks", "Volumes"}[s]
}

type LeftPanelModel struct {
	focused       bool
	width, height int
	activeSection Section
	cursor        int
	Containers    []docker.ContainerSummary
	Images        []docker.ImageSummary
	Networks      []docker.NetworkSummary
	Volumes       []docker.VolumeSummary
}

func NewLeftPanel() LeftPanelModel {
	return LeftPanelModel{focused: true}
}

func (m *LeftPanelModel) SetSize(w, h int) { m.width = w; m.height = h }
func (m *LeftPanelModel) SetFocused(f bool) { m.focused = f }

func (m *LeftPanelModel) NextSection() {
	m.activeSection = (m.activeSection + 1) % sectionCount
	m.cursor = 0
}

func (m *LeftPanelModel) PrevSection() {
	m.activeSection = (m.activeSection + sectionCount - 1) % sectionCount
	m.cursor = 0
}

func (m *LeftPanelModel) MoveDown() {
	if max := m.listLen() - 1; m.cursor < max {
		m.cursor++
	}
}

func (m *LeftPanelModel) MoveUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m *LeftPanelModel) ActiveSection() Section { return m.activeSection }

func (m *LeftPanelModel) SelectedContainerState() string {
	if m.cursor < len(m.Containers) {
		return m.Containers[m.cursor].State
	}
	return ""
}

func (m *LeftPanelModel) SelectedID() string {
	switch m.activeSection {
	case SectionContainers:
		if m.cursor < len(m.Containers) {
			return m.Containers[m.cursor].ID
		}
	case SectionImages:
		if m.cursor < len(m.Images) {
			return m.Images[m.cursor].ID
		}
	case SectionNetworks:
		if m.cursor < len(m.Networks) {
			return m.Networks[m.cursor].ID
		}
	case SectionVolumes:
		if m.cursor < len(m.Volumes) {
			return m.Volumes[m.cursor].Name
		}
	}
	return ""
}

func (m *LeftPanelModel) listLen() int {
	switch m.activeSection {
	case SectionContainers:
		return len(m.Containers)
	case SectionImages:
		return len(m.Images)
	case SectionNetworks:
		return len(m.Networks)
	case SectionVolumes:
		return len(m.Volumes)
	}
	return 0
}

func (m LeftPanelModel) View() string {
	innerW := m.width - 4
	if innerW < 1 {
		innerW = 1
	}

	var b strings.Builder

	tabs := make([]string, sectionCount)
	for i := Section(0); i < sectionCount; i++ {
		label := i.String()
		if i == m.activeSection {
			tabs[i] = styles.StyleTabActive.Render(label)
		} else {
			tabs[i] = styles.StyleTabInactive.Render(label)
		}
	}
	b.WriteString(strings.Join(tabs, styles.StyleMuted.Render(" │ ")) + "\n")
	b.WriteString(styles.StyleMuted.Render(strings.Repeat("─", innerW)) + "\n")

	availH := m.height - 4
	if availH < 1 {
		availH = 1
	}

	rows := m.renderRows(innerW)
	start := 0
	if m.cursor >= availH {
		start = m.cursor - availH + 1
	}
	end := start + availH
	if end > len(rows) {
		end = len(rows)
	}
	for i := start; i < end; i++ {
		b.WriteString(rows[i] + "\n")
	}

	borderStyle := styles.StylePanelBorder
	if m.focused {
		borderStyle = styles.StylePanelBorderFocused
	}
	return borderStyle.Width(m.width - 2).Height(m.height - 2).Render(b.String())
}

func (m LeftPanelModel) renderRows(w int) []string {
	switch m.activeSection {
	case SectionContainers:
		return m.renderContainers(w)
	case SectionImages:
		return m.renderImages(w)
	case SectionNetworks:
		return m.renderNetworks(w)
	case SectionVolumes:
		return m.renderVolumes(w)
	}
	return nil
}

func (m LeftPanelModel) renderContainers(w int) []string {
	if len(m.Containers) == 0 {
		return []string{styles.StyleMuted.Render("  No containers")}
	}
	rows := make([]string, len(m.Containers))
	for i, c := range m.Containers {
		stateStyle := styles.StyleStopped
		switch c.State {
		case "running":
			stateStyle = styles.StyleRunning
		case "paused":
			stateStyle = styles.StylePaused
		}
		stateLabel := stateStyle.Render(fmt.Sprintf("%-8s", c.State))
		name := c.Name
		if maxName := w - 10; maxName > 0 && len(name) > maxName {
			name = name[:maxName]
		}
		line := stateLabel + " " + name
		if i == m.cursor {
			line = styles.StyleSelected.Width(w).Render(line)
		}
		rows[i] = line
	}
	return rows
}

func (m LeftPanelModel) renderImages(w int) []string {
	if len(m.Images) == 0 {
		return []string{styles.StyleMuted.Render("  No images")}
	}
	rows := make([]string, len(m.Images))
	for i, img := range m.Images {
		tag := "<none>:<none>"
		if len(img.Tags) > 0 {
			tag = img.Tags[0]
		}
		size := formatBytesShort(img.Size)
		maxTag := w - len(size) - 2
		if maxTag > 0 && len(tag) > maxTag {
			tag = tag[:maxTag]
		}
		line := lipgloss.NewStyle().Width(w-len(size)-1).Render(tag) + styles.StyleMuted.Render(size)
		if i == m.cursor {
			line = styles.StyleSelected.Width(w).Render(tag + " " + size)
		}
		rows[i] = line
	}
	return rows
}

func (m LeftPanelModel) renderNetworks(w int) []string {
	if len(m.Networks) == 0 {
		return []string{styles.StyleMuted.Render("  No networks")}
	}
	rows := make([]string, len(m.Networks))
	for i, net := range m.Networks {
		name := net.Name
		if maxName := w - len(net.Driver) - 2; maxName > 0 && len(name) > maxName {
			name = name[:maxName]
		}
		line := lipgloss.NewStyle().Width(w-len(net.Driver)-1).Render(name) + styles.StyleMuted.Render(net.Driver)
		if i == m.cursor {
			line = styles.StyleSelected.Width(w).Render(net.Name + " " + net.Driver)
		}
		rows[i] = line
	}
	return rows
}

func (m LeftPanelModel) renderVolumes(w int) []string {
	if len(m.Volumes) == 0 {
		return []string{styles.StyleMuted.Render("  No volumes")}
	}
	rows := make([]string, len(m.Volumes))
	for i, vol := range m.Volumes {
		name := vol.Name
		if maxName := w - 2; maxName > 0 && len(name) > maxName {
			name = name[:maxName] + "…"
		}
		line := name
		if i == m.cursor {
			line = styles.StyleSelected.Width(w).Render(vol.Name)
		}
		rows[i] = line
	}
	return rows
}

func formatBytesShort(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%dB", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%ciB", float64(size)/float64(div), "KMGTPE"[exp])
}

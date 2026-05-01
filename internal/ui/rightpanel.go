package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/MarquesCoding/DockerViewer/internal/docker"
	"github.com/MarquesCoding/DockerViewer/internal/styles"
	"github.com/MarquesCoding/DockerViewer/internal/ui/render"
)

type RightPanelModel struct {
	focused       bool
	width, height int
	viewport      viewport.Model
	activeTab     int
	section       Section

	containerDetail docker.ContainerDetail
	imageDetail     docker.ImageDetail
	networkDetail   docker.NetworkDetail
	volumeDetail    docker.VolumeDetail
	detailLoaded    bool

	logs      []string
	logFollow bool
}

const maxLogLines = 500

func NewRightPanel() RightPanelModel {
	return RightPanelModel{viewport: viewport.New(80, 40)}
}

func (m *RightPanelModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.viewport.Width = w - 4
	m.viewport.Height = h - 6
	m.refresh()
}

func (m *RightPanelModel) SetFocused(f bool) { m.focused = f }

func (m *RightPanelModel) SetSection(s Section) {
	if m.section != s {
		m.section = s
		m.activeTab = 0
		m.logs = nil
		m.detailLoaded = false
		m.refresh()
	}
}

func (m *RightPanelModel) SetContainerDetail(d docker.ContainerDetail) {
	m.containerDetail = d
	m.detailLoaded = true
	m.refresh()
}

func (m *RightPanelModel) SetImageDetail(d docker.ImageDetail) {
	m.imageDetail = d
	m.detailLoaded = true
	m.refresh()
}

func (m *RightPanelModel) SetNetworkDetail(d docker.NetworkDetail) {
	m.networkDetail = d
	m.detailLoaded = true
	m.refresh()
}

func (m *RightPanelModel) SetVolumeDetail(d docker.VolumeDetail) {
	m.volumeDetail = d
	m.detailLoaded = true
	m.refresh()
}

func (m *RightPanelModel) AppendLog(line string) {
	m.logs = append(m.logs, line)
	if len(m.logs) > maxLogLines {
		m.logs = m.logs[len(m.logs)-maxLogLines:]
	}
	if m.section == SectionContainers && m.activeTab == 1 {
		m.refresh()
		m.viewport.GotoBottom()
	}
}

func (m *RightPanelModel) NextTab() {
	m.activeTab = (m.activeTab + 1) % m.tabCount()
	m.refresh()
}

func (m *RightPanelModel) PrevTab() {
	tc := m.tabCount()
	m.activeTab = (m.activeTab + tc - 1) % tc
	m.refresh()
}

func (m *RightPanelModel) ToggleFollow()  { m.logFollow = !m.logFollow }
func (m *RightPanelModel) ScrollDown()    { m.viewport.LineDown(3) }
func (m *RightPanelModel) ScrollUp()      { m.viewport.LineUp(3) }
func (m *RightPanelModel) PageDown()      { m.viewport.HalfViewDown() }
func (m *RightPanelModel) PageUp()        { m.viewport.HalfViewUp() }
func (m *RightPanelModel) ActiveTab() int { return m.activeTab }

func (m *RightPanelModel) tabCount() int {
	if m.section == SectionContainers {
		return 3
	}
	return 1
}

func (m *RightPanelModel) tabLabels() []string {
	if m.section == SectionContainers {
		return []string{"Stats", "Logs", "Config"}
	}
	return []string{"Info"}
}

func (m *RightPanelModel) refresh() {
	innerW := m.width - 4
	if innerW < 1 {
		innerW = 1
	}

	if !m.detailLoaded {
		m.viewport.SetContent(styles.StyleMuted.Render("  Loading..."))
		return
	}

	var content string
	switch m.section {
	case SectionContainers:
		content = render.Container(m.containerDetail, m.logs, m.activeTab, innerW)
	case SectionImages:
		content = render.Image(m.imageDetail, innerW)
	case SectionNetworks:
		content = render.Network(m.networkDetail, innerW)
	case SectionVolumes:
		content = render.Volume(m.volumeDetail, innerW)
	default:
		content = styles.StyleMuted.Render("  Select an item from the left panel")
	}
	m.viewport.SetContent(content)
}

func (m RightPanelModel) View() string {
	labels := m.tabLabels()
	tabParts := make([]string, len(labels))
	for i, label := range labels {
		if i == m.activeTab {
			tabParts[i] = styles.StyleTabActive.Render(label)
		} else {
			tabParts[i] = styles.StyleTabInactive.Render(label)
		}
	}

	innerW := m.width - 4
	if innerW < 1 {
		innerW = 1
	}

	var b strings.Builder
	b.WriteString(strings.Join(tabParts, styles.StyleMuted.Render(" │ ")) + "\n")
	b.WriteString(styles.StyleMuted.Render(strings.Repeat("─", innerW)) + "\n")
	b.WriteString(m.viewport.View())

	borderStyle := styles.StylePanelBorder
	if m.focused {
		borderStyle = styles.StylePanelBorderFocused
	}
	return borderStyle.Width(m.width - 2).Height(m.height - 2).Render(b.String())
}

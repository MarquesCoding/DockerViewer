package app

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/MarquesCoding/DockerViewer/internal/docker"
	"github.com/MarquesCoding/DockerViewer/internal/poller"
	"github.com/MarquesCoding/DockerViewer/internal/ui"
)

type Model struct {
	docker    *docker.Client
	left      ui.LeftPanelModel
	right     ui.RightPanelModel
	statusBar ui.StatusBarModel

	activePane ui.Pane
	width      int
	height     int
	lastErr    error

	logCtx    context.CancelFunc
	logCh     <-chan string
	logTarget string

	popup        *ui.PopupModel
	popupTargetID string
}

func New(dockerClient *docker.Client) Model {
	return Model{
		docker:     dockerClient,
		left:       ui.NewLeftPanel(),
		right:      ui.NewRightPanel(),
		statusBar:  ui.NewStatusBar(),
		activePane: ui.PaneLeft,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.fetchCurrentSection(),
		poller.TickCmd(poller.DefaultInterval),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		leftW, rightW, panelH := m.panelDims()
		m.left.SetSize(leftW, panelH)
		m.right.SetSize(rightW, panelH)
		m.statusBar.SetWidth(m.width)
		if m.popup != nil {
			m.popup.SetSize(m.width, m.height)
		}

	case tea.KeyMsg:
		return m.handleKey(msg)

	case poller.TickMsg:
		return m, tea.Batch(
			m.fetchCurrentSection(),
			poller.TickCmd(poller.DefaultInterval),
		)

	case docker.ContainerListMsg:
		if msg.Err != nil {
			m.lastErr = msg.Err
		} else {
			m.lastErr = nil
			m.left.Containers = msg.Items
			m.right.SetSection(m.left.ActiveSection())
			return m, m.fetchDetail()
		}

	case docker.ImageListMsg:
		if msg.Err != nil {
			m.lastErr = msg.Err
		} else {
			m.lastErr = nil
			m.left.Images = msg.Items
			m.right.SetSection(m.left.ActiveSection())
			return m, m.fetchDetail()
		}

	case docker.NetworkListMsg:
		if msg.Err != nil {
			m.lastErr = msg.Err
		} else {
			m.lastErr = nil
			m.left.Networks = msg.Items
			m.right.SetSection(m.left.ActiveSection())
			return m, m.fetchDetail()
		}

	case docker.VolumeListMsg:
		if msg.Err != nil {
			m.lastErr = msg.Err
		} else {
			m.lastErr = nil
			m.left.Volumes = msg.Items
			m.right.SetSection(m.left.ActiveSection())
			return m, m.fetchDetail()
		}

	case docker.ContainerDetailMsg:
		if msg.Err == nil {
			m.right.SetContainerDetail(msg.Detail)
		}

	case docker.ImageDetailMsg:
		if msg.Err == nil {
			m.right.SetImageDetail(msg.Detail)
		}

	case docker.NetworkDetailMsg:
		if msg.Err == nil {
			m.right.SetNetworkDetail(msg.Detail)
		}

	case docker.VolumeDetailMsg:
		if msg.Err == nil {
			m.right.SetVolumeDetail(msg.Detail)
		}

	case docker.LogLineMsg:
		if msg.ContainerID == m.logTarget {
			m.right.AppendLog(msg.Line)
			return m, docker.StreamLogsCmd(m.logTarget, m.logCh)
		}

	case docker.ActionDoneMsg:
		if msg.Err != nil {
			m.lastErr = fmt.Errorf("%s: %w", msg.Action, msg.Err)
		} else {
			m.lastErr = nil
		}
		// Refresh list so state change is reflected immediately.
		return m, m.fetchCurrentSection()
	}

	m.statusBar.SetError(m.lastErr)
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Popup intercepts all keys when open.
	if m.popup != nil {
		return m.handlePopupKey(msg)
	}

	switch msg.String() {
	case "q", "ctrl+c":
		m.stopLogStream()
		return m, tea.Quit

	case "s":
		return m.openPopup()

	case "left", "shift+tab":
		if m.activePane == ui.PaneLeft {
			m.stopLogStream()
			m.left.PrevSection()
			m.right.SetSection(m.left.ActiveSection())
			m.statusBar.SetSection(m.left.ActiveSection())
			return m, m.fetchCurrentSection()
		}
		m.right.PrevTab()
		m.statusBar.SetTab(m.right.ActiveTab())
		return m, m.maybeStartLogStream()

	case "right", "tab":
		if m.activePane == ui.PaneLeft {
			m.stopLogStream()
			m.left.NextSection()
			m.right.SetSection(m.left.ActiveSection())
			m.statusBar.SetSection(m.left.ActiveSection())
			return m, m.fetchCurrentSection()
		}
		m.right.NextTab()
		m.statusBar.SetTab(m.right.ActiveTab())
		return m, m.maybeStartLogStream()

	case "[":
		m.right.PrevTab()
		m.statusBar.SetTab(m.right.ActiveTab())
		return m, m.maybeStartLogStream()

	case "]":
		m.right.NextTab()
		m.statusBar.SetTab(m.right.ActiveTab())
		return m, m.maybeStartLogStream()

	case "j", "down":
		if m.activePane == ui.PaneLeft {
			m.stopLogStream()
			m.left.MoveDown()
			m.right.SetSection(m.left.ActiveSection())
			return m, m.fetchDetail()
		}
		m.right.ScrollDown()

	case "k", "up":
		if m.activePane == ui.PaneLeft {
			m.stopLogStream()
			m.left.MoveUp()
			m.right.SetSection(m.left.ActiveSection())
			return m, m.fetchDetail()
		}
		m.right.ScrollUp()

	case "pgdown", " ":
		m.right.PageDown()

	case "pgup":
		m.right.PageUp()

	case "a":
		m.activePane = ui.PaneLeft
		m.left.SetFocused(true)
		m.right.SetFocused(false)

	case "d", "enter":
		m.activePane = ui.PaneRight
		m.left.SetFocused(false)
		m.right.SetFocused(true)

	case "f":
		m.right.ToggleFollow()

	case "r":
		return m, m.fetchCurrentSection()
	}

	return m, nil
}

func (m Model) handlePopupKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		if m.popup.NeedsConfirm() {
			m.popup.ResetConfirm()
		} else {
			m.popup = nil
		}

	case "j", "down":
		if !m.popup.NeedsConfirm() {
			m.popup.MoveDown()
		}

	case "k", "up":
		if !m.popup.NeedsConfirm() {
			m.popup.MoveUp()
		}

	case "enter":
		action, needsConfirm := m.popup.Select()
		if needsConfirm && !m.popup.NeedsConfirm() {
			m.popup.Confirm()
		} else {
			m.popup = nil
			return m, m.executeAction(action.ID, m.popupTargetID)
		}
	}

	return m, nil
}

func (m Model) openPopup() (tea.Model, tea.Cmd) {
	id := m.left.SelectedID()
	if id == "" {
		return m, nil
	}

	var title string
	var actions []ui.Action

	switch m.left.ActiveSection() {
	case ui.SectionContainers:
		state := m.left.SelectedContainerState()
		title = "Container: " + id
		actions = ui.ContainerActions(state)
	case ui.SectionImages:
		title = "Image: " + id
		actions = ui.ImageActions()
	case ui.SectionNetworks:
		title = "Network: " + id
		actions = ui.NetworkActions()
	case ui.SectionVolumes:
		title = "Volume: " + id
		actions = ui.VolumeActions()
	default:
		return m, nil
	}

	p := ui.NewPopup(title, actions)
	p.SetSize(m.width, m.height)
	m.popup = &p
	m.popupTargetID = id
	return m, nil
}

func (m Model) executeAction(actionID, targetID string) tea.Cmd {
	ctx := context.Background()
	switch m.left.ActiveSection() {
	case ui.SectionContainers:
		switch actionID {
		case "start":
			return m.docker.StartContainer(ctx, targetID)
		case "stop":
			return m.docker.StopContainer(ctx, targetID)
		case "restart":
			return m.docker.RestartContainer(ctx, targetID)
		case "pause":
			return m.docker.PauseContainer(ctx, targetID)
		case "unpause":
			return m.docker.UnpauseContainer(ctx, targetID)
		case "kill":
			return m.docker.KillContainer(ctx, targetID)
		case "remove":
			return m.docker.RemoveContainer(ctx, targetID)
		}
	case ui.SectionImages:
		if actionID == "remove" {
			return m.docker.RemoveImage(ctx, targetID)
		}
	case ui.SectionNetworks:
		if actionID == "remove" {
			return m.docker.RemoveNetwork(ctx, targetID)
		}
	case ui.SectionVolumes:
		if actionID == "remove" {
			return m.docker.RemoveVolume(ctx, targetID)
		}
	}
	return nil
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	leftW, rightW, panelH := m.panelDims()

	panels := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(leftW).Height(panelH).Render(m.left.View()),
		lipgloss.NewStyle().Width(rightW).Height(panelH).Render(m.right.View()),
	)
	base := lipgloss.JoinVertical(lipgloss.Left, panels, m.statusBar.View())

	if m.popup != nil {
		return overlayPopup(base, m.popup.View(), m.width, m.height)
	}
	return base
}

// overlayPopup composites the popup on top of the base layout by replacing
// the center lines. The popup View already uses lipgloss.Place with full
// terminal dimensions, so we just return it — the alt-screen ensures the
// base is still visible in the terminal buffer behind it.
func overlayPopup(_, popupView string, _, _ int) string {
	return popupView
}

func (m Model) panelDims() (leftW, rightW, panelH int) {
	leftW = m.width * 30 / 100
	if leftW < 28 {
		leftW = 28
	}
	rightW = m.width - leftW
	panelH = m.height - 1
	return
}

func (m Model) fetchCurrentSection() tea.Cmd {
	ctx := context.Background()
	switch m.left.ActiveSection() {
	case ui.SectionContainers:
		return m.docker.ListContainers(ctx)
	case ui.SectionImages:
		return m.docker.ListImages(ctx)
	case ui.SectionNetworks:
		return m.docker.ListNetworks(ctx)
	case ui.SectionVolumes:
		return m.docker.ListVolumes(ctx)
	}
	return nil
}

func (m Model) fetchDetail() tea.Cmd {
	id := m.left.SelectedID()
	if id == "" {
		return nil
	}
	ctx := context.Background()
	switch m.left.ActiveSection() {
	case ui.SectionContainers:
		return m.docker.InspectContainer(ctx, id)
	case ui.SectionImages:
		return m.docker.InspectImage(ctx, id)
	case ui.SectionNetworks:
		return m.docker.InspectNetwork(ctx, id)
	case ui.SectionVolumes:
		return m.docker.InspectVolume(ctx, id)
	}
	return nil
}

func (m *Model) maybeStartLogStream() tea.Cmd {
	if m.left.ActiveSection() != ui.SectionContainers || m.right.ActiveTab() != 1 {
		return nil
	}
	id := m.left.SelectedID()
	if id == "" {
		return nil
	}
	m.stopLogStream()

	ctx, cancel := context.WithCancel(context.Background())
	m.logCtx = cancel
	m.logTarget = id
	m.logCh = m.docker.StartLogStream(ctx, id)
	return docker.StreamLogsCmd(id, m.logCh)
}

func (m *Model) stopLogStream() {
	if m.logCtx != nil {
		m.logCtx()
		m.logCtx = nil
		m.logCh = nil
		m.logTarget = ""
	}
}

package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	tea "github.com/charmbracelet/bubbletea"
)

type ActionDoneMsg struct {
	Action string
	Err    error
}

func (c *Client) StartContainer(ctx context.Context, id string) tea.Cmd {
	return func() tea.Msg {
		err := c.cli.ContainerStart(ctx, id, container.StartOptions{})
		return ActionDoneMsg{Action: "start", Err: err}
	}
}

func (c *Client) StopContainer(ctx context.Context, id string) tea.Cmd {
	return func() tea.Msg {
		err := c.cli.ContainerStop(ctx, id, container.StopOptions{})
		return ActionDoneMsg{Action: "stop", Err: err}
	}
}

func (c *Client) RestartContainer(ctx context.Context, id string) tea.Cmd {
	return func() tea.Msg {
		err := c.cli.ContainerRestart(ctx, id, container.StopOptions{})
		return ActionDoneMsg{Action: "restart", Err: err}
	}
}

func (c *Client) PauseContainer(ctx context.Context, id string) tea.Cmd {
	return func() tea.Msg {
		err := c.cli.ContainerPause(ctx, id)
		return ActionDoneMsg{Action: "pause", Err: err}
	}
}

func (c *Client) UnpauseContainer(ctx context.Context, id string) tea.Cmd {
	return func() tea.Msg {
		err := c.cli.ContainerUnpause(ctx, id)
		return ActionDoneMsg{Action: "unpause", Err: err}
	}
}

func (c *Client) KillContainer(ctx context.Context, id string) tea.Cmd {
	return func() tea.Msg {
		err := c.cli.ContainerKill(ctx, id, "SIGKILL")
		return ActionDoneMsg{Action: "kill", Err: err}
	}
}

func (c *Client) RemoveContainer(ctx context.Context, id string) tea.Cmd {
	return func() tea.Msg {
		err := c.cli.ContainerRemove(ctx, id, container.RemoveOptions{Force: true, RemoveVolumes: false})
		return ActionDoneMsg{Action: "remove", Err: err}
	}
}

func (c *Client) RemoveImage(ctx context.Context, id string) tea.Cmd {
	return func() tea.Msg {
		_, err := c.cli.ImageRemove(ctx, id, image.RemoveOptions{Force: false, PruneChildren: false})
		return ActionDoneMsg{Action: "remove image", Err: err}
	}
}

func (c *Client) RemoveNetwork(ctx context.Context, id string) tea.Cmd {
	return func() tea.Msg {
		err := c.cli.NetworkRemove(ctx, id)
		return ActionDoneMsg{Action: "remove network", Err: err}
	}
}

func (c *Client) RemoveVolume(ctx context.Context, name string) tea.Cmd {
	return func() tea.Msg {
		err := c.cli.VolumeRemove(ctx, name, false)
		return ActionDoneMsg{Action: "remove volume", Err: err}
	}
}

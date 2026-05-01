package docker

import (
	"context"

	"github.com/docker/docker/api/types/network"
	tea "github.com/charmbracelet/bubbletea"
)

type NetworkSummary struct {
	ID     string
	Name   string
	Driver string
	Scope  string
}

type NetworkDetail struct {
	Info network.Inspect
}

func (c *Client) ListNetworks(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		networks, err := c.cli.NetworkList(ctx, network.ListOptions{})
		if err != nil {
			return NetworkListMsg{Err: err}
		}

		summaries := make([]NetworkSummary, 0, len(networks))
		for _, net := range networks {
			id := net.ID
			if len(id) > 12 {
				id = id[:12]
			}
			summaries = append(summaries, NetworkSummary{
				ID:     id,
				Name:   net.Name,
				Driver: net.Driver,
				Scope:  net.Scope,
			})
		}
		return NetworkListMsg{Items: summaries}
	}
}

func (c *Client) InspectNetwork(ctx context.Context, id string) tea.Cmd {
	return func() tea.Msg {
		info, err := c.cli.NetworkInspect(ctx, id, network.InspectOptions{})
		if err != nil {
			return NetworkDetailMsg{Err: err}
		}
		return NetworkDetailMsg{Detail: NetworkDetail{Info: info}}
	}
}

package docker

import (
	"context"

	"github.com/docker/docker/api/types/volume"
	tea "github.com/charmbracelet/bubbletea"
)

type VolumeSummary struct {
	Name       string
	Driver     string
	Mountpoint string
	Scope      string
}

type VolumeDetail struct {
	Info volume.Volume
}

func (c *Client) ListVolumes(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		resp, err := c.cli.VolumeList(ctx, volume.ListOptions{})
		if err != nil {
			return VolumeListMsg{Err: err}
		}

		summaries := make([]VolumeSummary, 0, len(resp.Volumes))
		for _, vol := range resp.Volumes {
			summaries = append(summaries, VolumeSummary{
				Name:       vol.Name,
				Driver:     vol.Driver,
				Mountpoint: vol.Mountpoint,
				Scope:      vol.Scope,
			})
		}
		return VolumeListMsg{Items: summaries}
	}
}

func (c *Client) InspectVolume(ctx context.Context, name string) tea.Cmd {
	return func() tea.Msg {
		info, err := c.cli.VolumeInspect(ctx, name)
		if err != nil {
			return VolumeDetailMsg{Err: err}
		}
		return VolumeDetailMsg{Detail: VolumeDetail{Info: info}}
	}
}

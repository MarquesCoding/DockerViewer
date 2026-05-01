package docker

import (
	"context"

	"github.com/docker/docker/api/types/image"
	tea "github.com/charmbracelet/bubbletea"
)

type ImageSummary struct {
	ID      string
	Tags    []string
	Size    int64
	Created int64
}

type ImageDetail struct {
	History []image.HistoryResponseItem
	Inspect image.InspectResponse
}

func (c *Client) ListImages(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		images, err := c.cli.ImageList(ctx, image.ListOptions{})
		if err != nil {
			return ImageListMsg{Err: err}
		}

		summaries := make([]ImageSummary, 0, len(images))
		for _, img := range images {
			id := img.ID
			if len(id) > 19 {
				id = id[7:19]
			}
			summaries = append(summaries, ImageSummary{
				ID:      id,
				Tags:    img.RepoTags,
				Size:    img.Size,
				Created: img.Created,
			})
		}
		return ImageListMsg{Items: summaries}
	}
}

func (c *Client) InspectImage(ctx context.Context, id string) tea.Cmd {
	return func() tea.Msg {
		inspect, err := c.cli.ImageInspect(ctx, id)
		if err != nil {
			return ImageDetailMsg{Err: err}
		}
		history, err := c.cli.ImageHistory(ctx, id)
		if err != nil {
			history = nil
		}
		return ImageDetailMsg{Detail: ImageDetail{
			Inspect: inspect,
			History: history,
		}}
	}
}

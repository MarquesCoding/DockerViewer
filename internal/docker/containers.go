package docker

import (
	"context"
	"encoding/json"
	"io"

	"github.com/docker/docker/api/types/container"
	tea "github.com/charmbracelet/bubbletea"
)

type ContainerSummary struct {
	ID      string
	Name    string
	Image   string
	Status  string
	State   string
	CPU     float64
	Memory  float64
	MemUsed uint64
	MemMax  uint64
}

type ContainerDetail struct {
	Info  container.InspectResponse
	Stats StatsSnapshot
}

type StatsSnapshot struct {
	CPUPercent float64
	MemUsed    uint64
	MemMax     uint64
	MemPercent float64
}

func (c *Client) ListContainers(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		containers, err := c.cli.ContainerList(ctx, container.ListOptions{All: true})
		if err != nil {
			return ContainerListMsg{Err: err}
		}

		summaries := make([]ContainerSummary, 0, len(containers))
		for _, ctr := range containers {
			name := ""
			if len(ctr.Names) > 0 {
				name = ctr.Names[0]
				if len(name) > 0 && name[0] == '/' {
					name = name[1:]
				}
			}
			summaries = append(summaries, ContainerSummary{
				ID:     ctr.ID[:12],
				Name:   name,
				Image:  ctr.Image,
				Status: ctr.Status,
				State:  ctr.State,
			})
		}
		return ContainerListMsg{Items: summaries}
	}
}

func (c *Client) InspectContainer(ctx context.Context, id string) tea.Cmd {
	return func() tea.Msg {
		info, err := c.cli.ContainerInspect(ctx, id)
		if err != nil {
			return ContainerDetailMsg{Err: err}
		}
		stats, err := c.fetchStats(ctx, id)
		if err != nil {
			stats = StatsSnapshot{}
		}
		return ContainerDetailMsg{Detail: ContainerDetail{Info: info, Stats: stats}}
	}
}

func (c *Client) FetchStats(ctx context.Context, id string) tea.Cmd {
	return func() tea.Msg {
		stats, err := c.fetchStats(ctx, id)
		if err != nil {
			return ContainerStatsMsg{Err: err}
		}
		return ContainerStatsMsg{Stats: stats}
	}
}

func (c *Client) fetchStats(ctx context.Context, id string) (StatsSnapshot, error) {
	resp, err := c.cli.ContainerStats(ctx, id, false)
	if err != nil {
		return StatsSnapshot{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return StatsSnapshot{}, err
	}

	var s container.StatsResponse
	if err := json.Unmarshal(body, &s); err != nil {
		return StatsSnapshot{}, err
	}

	cpuDelta := float64(s.CPUStats.CPUUsage.TotalUsage - s.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(s.CPUStats.SystemUsage - s.PreCPUStats.SystemUsage)
	numCPUs := float64(s.CPUStats.OnlineCPUs)
	if numCPUs == 0 {
		numCPUs = float64(len(s.CPUStats.CPUUsage.PercpuUsage))
	}

	var cpuPercent float64
	if systemDelta > 0 && cpuDelta > 0 {
		cpuPercent = (cpuDelta / systemDelta) * numCPUs * 100.0
	}

	memUsed := s.MemoryStats.Usage
	memMax := s.MemoryStats.Limit
	var memPercent float64
	if memMax > 0 {
		memPercent = float64(memUsed) / float64(memMax) * 100.0
	}

	return StatsSnapshot{
		CPUPercent: cpuPercent,
		MemUsed:    memUsed,
		MemMax:     memMax,
		MemPercent: memPercent,
	}, nil
}

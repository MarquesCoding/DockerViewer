package docker

import (
	"bufio"
	"context"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	tea "github.com/charmbracelet/bubbletea"
)

type LogLineMsg struct {
	ContainerID string
	Line        string
}

// StreamLogsCmd returns a tea.Cmd that reads one log line from ch and returns it
// as a LogLineMsg. Re-issue this command in Update to keep receiving lines.
func StreamLogsCmd(containerID string, ch <-chan string) tea.Cmd {
	return func() tea.Msg {
		line, ok := <-ch
		if !ok {
			return nil
		}
		return LogLineMsg{ContainerID: containerID, Line: line}
	}
}

// StartLogStream spawns a goroutine that writes lines to the returned channel.
// Cancel ctx to stop it.
func (c *Client) StartLogStream(ctx context.Context, containerID string) <-chan string {
	ch := make(chan string, 64)
	go func() {
		defer close(ch)
		opts := container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
			Tail:       "100",
		}
		reader, err := c.cli.ContainerLogs(ctx, containerID, opts)
		if err != nil {
			return
		}
		defer reader.Close()

		pr, pw := io.Pipe()
		go func() {
			defer pw.Close()
			stdcopy.StdCopy(pw, pw, reader)
		}()

		scanner := bufio.NewScanner(pr)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case ch <- scanner.Text():
			}
		}
	}()
	return ch
}

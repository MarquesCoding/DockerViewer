package docker

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
)

type Client struct {
	cli *client.Client
}

func New() (*Client, error) {
	// If DOCKER_HOST is set, use it directly (covers Colima, remote, etc.)
	if os.Getenv("DOCKER_HOST") != "" {
		return newClient(client.FromEnv)
	}

	// Try standard socket, then common Colima socket path.
	candidates := []string{
		"/var/run/docker.sock",
		filepath.Join(os.Getenv("HOME"), ".colima/default/docker.sock"),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return newClient(client.WithHost("unix://" + path))
		}
	}

	// Nothing found — fall back to FromEnv and let Docker SDK produce the error.
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("could not connect to Docker (tried %v): %w", candidates, err)
	}
	return &Client{cli: cli}, nil
}

func newClient(opt client.Opt) (*Client, error) {
	cli, err := client.NewClientWithOpts(opt, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Client{cli: cli}, nil
}

func (c *Client) Close() {
	c.cli.Close()
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.cli.Ping(ctx)
	return err
}

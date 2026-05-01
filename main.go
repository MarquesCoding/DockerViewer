package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/MarquesCoding/DockerViewer/internal/app"
	"github.com/MarquesCoding/DockerViewer/internal/docker"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("dockerviewer", version)
		return
	}
	dockerClient, err := docker.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to Docker: %v\n", err)
		os.Exit(1)
	}
	defer dockerClient.Close()

	model := app.New(dockerClient)
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

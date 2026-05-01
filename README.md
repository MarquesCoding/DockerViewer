# DockerViewer

A terminal UI for Docker — think lazydocker, but built from scratch as a personal project to learn Go.

<img width="1642" height="554" alt="image" src="https://github.com/user-attachments/assets/b7be171b-88e1-4ffb-8199-f47103d11288" />

This is early-stage and actively being developed. Expect rough edges.

## What it does

- Browse containers, images, networks, and volumes in a split-pane TUI
- View container stats, live logs, and config
- Run actions on resources (start, stop, restart, pause, kill, remove)
- Auto-refreshes every 3 seconds
- Works with Docker Desktop and Colima

## Installation

### Homebrew

```sh
brew install MarquesCoding/dockerviewer/dockerviewer
```

### From source

Requires Go 1.21+.

```sh
git clone https://github.com/MarquesCoding/DockerViewer
cd dockerviewer
go build -o dockerviewer .
./dockerviewer
```

## Usage

```sh
dockerviewer
```

If you're using Colima or a non-standard Docker socket, set `DOCKER_HOST` first:

```sh
export DOCKER_HOST=unix://$HOME/.colima/default/docker.sock
dockerviewer
```

## Keybindings

| Key | Action |
|-----|--------|
| `←` / `→` | Switch section (Containers / Images / Networks / Volumes) |
| `↑` / `↓` | Navigate list |
| `a` / `d` | Switch panel focus |
| `s` | Open actions menu |
| `[` / `]` | Cycle detail tabs |
| `f` | Toggle log follow |
| `r` | Force refresh |
| `q` | Quit |

## Status

Early stages — built primarily to learn Go and the [Bubbletea](https://github.com/charmbracelet/bubbletea) TUI framework. Not production-ready.

Contributions and feedback welcome.

## License

MIT

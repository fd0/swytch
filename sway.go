package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/mdirkse/i3ipc-go"
)

type Window struct {
	Program   string
	Title     string
	Workspace string
	Active    bool

	ID int64
}

func getAllWindows(socket *i3ipc.IPCSocket) ([]Window, error) {
	root, err := socket.GetTree()
	if err != nil {
		return nil, fmt.Errorf("get tree: %w", err)
	}

	windows := []Window{}

	for _, node := range root.Leaves() {
		if node.Name == "" {
			continue
		}

		program := node.Window_Properties.Class
		if program == "" {
			program = node.AppID
		}

		window := Window{
			Title:     node.Name,
			Program:   program,
			Workspace: strings.Split(node.Workspace().Name, ":")[0],
			Active:    node.Focused,

			ID: node.ID,
		}

		windows = append(windows, window)
	}

	return windows, nil
}

func FocusWindow(id string) error {
	selector := fmt.Sprintf("[con_id=%s]", id)
	cmd := exec.Command("swaymsg", selector, "focus")

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("focus window: %w", err)
	}

	return nil
}

func MoveWindowToCurrentWorkspace(id string) error {
	selector := fmt.Sprintf("[con_id=%s]", id)
	cmd := exec.Command("swaymsg", selector, "move", "container", "to", "workspace", "current")

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("focus window: %w", err)
	}

	return nil
}

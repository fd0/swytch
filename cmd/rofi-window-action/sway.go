package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/joshuarubin/go-sway"
)

type Window struct {
	Program   string
	Title     string
	Workspace string
	Active    bool

	ID int64
}

func newWindow(workspace *sway.Node, node *sway.Node) Window {
	var program string

	if node.WindowProperties != nil {
		program = node.WindowProperties.Class
	}

	if program == "" && node.AppID != nil {
		program = *node.AppID
	}

	window := Window{
		Title:     node.Name,
		Program:   program,
		Workspace: strings.Split(workspace.Name, ":")[0],
		Active:    node.Focused,

		ID: node.ID,
	}

	return window
}

func traverseNodes(workspace *sway.Node, node *sway.Node, list []Window) []Window {
	if node.Type == sway.NodeCon || node.Type == sway.NodeFloatingCon {
		// add windows to the list, ignore containers
		if node.Name != "" {
			list = append(list, newWindow(workspace, node))
		}
	}

	for _, n := range node.Nodes {
		list = traverseNodes(workspace, n, list)
	}

	for _, n := range node.FloatingNodes {
		list = traverseNodes(workspace, n, list)
	}

	return list
}

func getAllWindows(ctx context.Context, client sway.Client) ([]Window, error) {
	root, err := client.GetTree(ctx)
	if err != nil {
		return nil, fmt.Errorf("get tree: %w", err)
	}

	var list []Window

	for _, output := range root.Nodes {
		if output.Type != sway.NodeOutput {
			continue
		}

		for _, workspace := range output.Nodes {
			if workspace.Name == "__i3_scratch" {
				continue
			}

			list = traverseNodes(workspace, workspace, list)
		}
	}

	return list, nil
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

func KillWindow(id string) error {
	selector := fmt.Sprintf("[con_id=%s]", id)
	cmd := exec.Command("swaymsg", selector, "kill")

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("kill window: %w", err)
	}

	return nil
}

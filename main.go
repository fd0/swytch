package main

import (
	"fmt"
	"html"
	"os"
	"strconv"
	"strings"

	"github.com/mdirkse/i3ipc-go"
	"github.com/spf13/pflag"
)

type Options struct {
	WorkspaceColors []string
}

func main() {
	var opts Options

	flags := pflag.NewFlagSet("swytch", pflag.ExitOnError)
	flags.StringSliceVar(&opts.WorkspaceColors, "workspace-colors", []string{"lightblue", "lightgreen", "orange", "red", "magenta", "cyan"}, "Use workspace `color`")

	err := flags.Parse(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse flags: %v\n", err)
		os.Exit(10)
	}

	if v, ok := os.LookupEnv("SWYTCH_WORKSPACE_COLORS"); ok {
		opts.WorkspaceColors = strings.Split(v, ",")
	}

	err = run(opts, flags.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(11)
	}
}

func run(opts Options, args []string) error {
	if _, ok := os.LookupEnv("ROFI_RETV"); !ok {
		// just run rofi with the right parameters
		return RunRofi("SWYTCH_WORKSPACE_COLORS=" + strings.Join(opts.WorkspaceColors, ","))
	}

	retv, err := strconv.Atoi(os.Getenv("ROFI_RETV"))
	if err != nil {
		return fmt.Errorf("invalid value passed in ROFI_RETV: %w", err)
	}

	info := os.Getenv("ROFI_INFO")

	// called by rofi
	if info != "" {
		// item selected
		fmt.Fprintf(os.Stderr, "startup, selected item %v (%v) by rofi: %v\n", info, retv, args[1])

		switch retv {
		case 1:
			fmt.Fprintf(os.Stderr, "focus window %v\n", info)

			return FocusWindow(info)

		case 10:
			fmt.Fprintf(os.Stderr, "move window %v to current workspace\n", info)

			return MoveWindowToCurrentWorkspace(info)

		case 11:
			fmt.Fprintf(os.Stderr, "kill window %v\n", info)

			return KillWindow(info)

		default:
			return fmt.Errorf("unknown keyboard shortcut %v received from rofi via $ROFI_RETV", retv)
		}
	}

	// build menu
	socket, err := i3ipc.GetIPCSocket()
	if err != nil {
		return fmt.Errorf("connect window manager: %w", err)
	}

	windows, err := getAllWindows(socket)
	if err != nil {
		return fmt.Errorf("get windows: %w", err)
	}

	// configure rofi
	dispOpts := DisplayOptions{
		Prompt:     "window",
		NoCustom:   true,
		UseHotKeys: true,
		MarkupRows: true,
	}

	fmt.Print(dispOpts.ConfigString())

	workspace := ""
	color := ""
	colorIndex := 0

	for _, w := range windows {
		if w.Workspace != workspace {
			color = opts.WorkspaceColors[colorIndex%len(opts.WorkspaceColors)]
			colorIndex++
			workspace = w.Workspace
		}

		active := ""
		if w.Active {
			active = ` font_weight="bold"`
		}

		text := fmt.Sprintf("<span foreground=%q>[%s]</span>", color, html.EscapeString(w.Workspace))
		text += fmt.Sprintf("\t<span%s>%s\t%s</span>", active, html.EscapeString(w.Program), html.EscapeString(w.Title))

		row := Row{
			Text: text,
			Info: fmt.Sprintf("%d", w.ID),
		}
		fmt.Print(row.ConfigString())
	}

	return nil
}

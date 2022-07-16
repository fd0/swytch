package main

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
	"strconv"

	"github.com/mdirkse/i3ipc-go"
	"github.com/spf13/pflag"
)

type Options struct {
	WorkspaceColors []string
	Debug           bool
}

func main() {
	var opts Options

	flags := pflag.NewFlagSet("swytch", pflag.ExitOnError)
	flags.StringSliceVar(&opts.WorkspaceColors, "workspace-colors", []string{"lightblue", "lightgreen", "orange", "red", "magenta", "cyan"}, "Use workspace `color`")
	flags.BoolVar(&opts.Debug, "debug", false, "Print debug messages")

	err := flags.Parse(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse flags: %v\n", err)
		os.Exit(10)
	}

	// get opts from environment if passed from there, overrides cli flags
	if v, ok := os.LookupEnv("ROFI_WINDOW_ACTION_OPTS"); ok {
		err = json.Unmarshal([]byte(v), &opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "json decode opts from $ROFI_WINDOW_ACTION_OPTS: %v\n", err)
			os.Exit(12)
		}
	}

	err = run(opts, flags.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(11)
	}
}

func run(opts Options, args []string) error {
	if _, ok := os.LookupEnv("ROFI_RETV"); !ok {
		// just run rofi with the right parameters, pass all opts as json in environment
		buf, err := json.Marshal(opts)
		if err != nil {
			return fmt.Errorf("json encode: %w", err)
		}

		return RunRofi("ROFI_WINDOW_ACTION_OPTS=" + string(buf))
	}

	retv, err := strconv.Atoi(os.Getenv("ROFI_RETV"))
	if err != nil {
		return fmt.Errorf("invalid value passed in ROFI_RETV: %w", err)
	}

	info := os.Getenv("ROFI_INFO")

	// called by rofi
	if info != "" {
		// item selected
		if opts.Debug {
			fmt.Fprintf(os.Stderr, "selected item %v (%v) by rofi: %v\n", info, retv, args[1])
		}

		switch retv {
		case 1:
			if opts.Debug {
				fmt.Fprintf(os.Stderr, "focus window %v\n", info)
			}

			return FocusWindow(info)

		case 10:
			if opts.Debug {
				fmt.Fprintf(os.Stderr, "move window %v to current workspace\n", info)
			}

			return MoveWindowToCurrentWorkspace(info)

		case 11:
			if opts.Debug {
				fmt.Fprintf(os.Stderr, "kill window %v\n", info)
			}

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

package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mdirkse/i3ipc-go"
)

func main() {
	err := run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(10)
	}
}

func run(args []string) error {
	if _, ok := os.LookupEnv("ROFI_RETV"); !ok {
		// just run rofi with the right parameters
		return RunRofi()
	}

	retv, err := strconv.Atoi(os.Getenv("ROFI_RETV"))
	if err != nil {
		return fmt.Errorf("invalid value passed in ROFI_RETV: %w", err)
	}

	info := os.Getenv("ROFI_INFO")

	// called by rofi
	if len(args) > 1 {
		// item selected
		fmt.Fprintf(os.Stderr, "startup, selected item by rofi via %v (%v): %v\n", retv, info, args[1])

		switch retv {
		case 1:
			fmt.Fprintf(os.Stderr, "focus window %v\n", info)

			return FocusWindow(info)

		case 10:
			fmt.Fprintf(os.Stderr, "move window %v to current workspace\n", info)

			return MoveWindowToCurrentWorkspace(info)

		default:
			return fmt.Errorf("unknown keyboard shortcut %v received from rofi via $ROFI_RETV", retv)
		}
	}

	// build menu
	fmt.Fprintf(os.Stderr, "startup, show rofi menu\n")

	socket, err := i3ipc.GetIPCSocket()
	if err != nil {
		return fmt.Errorf("connect window manager: %w", err)
	}

	windows, err := getAllWindows(socket)
	if err != nil {
		return fmt.Errorf("get windows: %w", err)
	}

	// configure rofi
	opts := DisplayOptions{
		Prompt:     "window",
		NoCustom:   true,
		UseHotKeys: true,
	}

	fmt.Print(opts.ConfigString())

	for _, window := range windows {
		row := Row{
			Text: window.String(),
			Info: fmt.Sprintf("%d", window.ID),
		}
		fmt.Print(row.ConfigString())
	}

	return nil
}

package main

import (
	"fmt"
	"os"

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
	socket, err := i3ipc.GetIPCSocket()
	if err != nil {
		return fmt.Errorf("connect window manager: %w", err)
	}

	version, err := socket.GetVersion()
	if err != nil {
		return fmt.Errorf("get version: %w", err)
	}

	fmt.Printf("version: %v\n", version)

	return nil
}

// icon   Program     Title
//        code-oss    main.go
//        chromium    Wayland support
//        Gvim        notes-x

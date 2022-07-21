package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

func RunRofi(ctx context.Context, additionalEnvironment ...string) error {
	program, err := os.Executable()
	if err != nil {
		return fmt.Errorf("find executable: %w", err)
	}

	cmd := exec.CommandContext(ctx,
		"rofi", "-modi", "swytch:"+program, "-show", "swytch",
		"-kb-accept-alt", "", // disable alternative accept
		"-kb-custom-1", "Shift+Return", // set custom keybinding 1 to shift+return
		"-kb-custom-2", "Control+c", // kill window
	)
	cmd.Env = append(os.Environ(), additionalEnvironment...)
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("run rofi: %w", err)
	}

	return nil
}

type DisplayOptions struct {
	Prompt     string
	Message    string
	MarkupRows bool
	NoCustom   bool
	UseHotKeys bool
}

func formatOpt(name, value string) string {
	return "\x00" + name + "\x1f" + value + "\n"
}

func (opts DisplayOptions) ConfigString() string {
	res := ""

	if opts.Prompt != "" {
		res += formatOpt("prompt", opts.Prompt)
	}

	if opts.Message != "" {
		res += formatOpt("message", opts.Message)
	}

	if opts.MarkupRows {
		res += formatOpt("markup-rows", "true")
	}

	if opts.NoCustom {
		res += formatOpt("no-custom", "true")
	}

	if opts.UseHotKeys {
		res += formatOpt("use-hot-keys", "true")
	}

	return res
}

type Row struct {
	Text string
	Info string
}

func (r Row) ConfigString() string {
	extra := ""

	if r.Info != "" {
		extra += "\x00" + "info" + "\x1f" + r.Info
	}

	return r.Text + extra + "\n"
}

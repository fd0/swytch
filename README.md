# Introduction

This program is a replacement for the window selection built into rofi, which
does not work with Sway. It extracts the windows via `swaymsg` instead of
querying the X server.

It is inspired by [`swytch.sh`](https://github.com/wilecoyote2015/Swytch),
but implemented in Go instead of Shell for speed.

## Building

Install Go, then run `go build`.


## Keybindings

 * `Return` switches focus to the window
 * `Shift+Return` moves the window to the current workspace
 * `Control+C` kill the window

# Dependency

It includes a fork of the library `i3ipc-go` from [here](https://github.com/emcconville/i3ipc-go)
(branch `sway_support`) for Sway support (at least until
https://github.com/mdirkse/i3ipc-go/pull/9 is merged). Additionally, we've
added the `app_id` field used for native Wayland windows.

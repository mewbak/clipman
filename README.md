# Clipman

A basic clipboard manager for Wayland, with support for persisting copy buffers after an application exits.

## Installing

Requirements:

- a windows manager that uses `wlr-data-control`, like Sway and other wlroots-based WMs.
- wl-clipboard >= 2.0
- dmenu or rofi

[Install go](https://golang.org/doc/install), add `$GOPATH/bin` to your path, then run `go get github.com/yory8/clipman` OR run `go install` inside this folder.

Archlinux users can find a PKGBUILD [here](https://aur.archlinux.org/packages/clipman/).

## Usage

Run the binary in your Sway session by adding `exec wl-paste -t text --watch clipman store` (or `exec wl-paste -t text --watch clipman store 1>> PATH/TO/LOGFILE 2>&1 &` to log errors) at the beginning of your config.
For primary clipboard support, also add `exec wl-paste -p -t text --watch clipman store --histpath="~/.local/share/clipman-primary.json`.

To query the history and select items, run the binary as `clipman pick`. You can assign it to a keybinding: `bindsym $mod+h exec clipman pick`.
For primary clipboard support, `clipman pick --histpath="~/.local/share/clipman-primary.json`.

To remove items from history, `clipman clear` and `clipman clear --all`.

For more options: `clipman -h`.

## Versions

This projects follows SemVer conventions.

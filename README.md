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

Run the binary in your Sway session by adding `exec clipman -d` (or `exec clipman -d 1>> PATH/TO/LOGFILE 2>&1 &` to log errors) at the beginning of your config.

To query the history and select items, run the binary as `clipman -s`. You can assign it to a keybinding: `bindsym $mod+h exec clipman -s`.

For more options: `clipman -h`.

## Versions

This projects follows SemVer conventions.

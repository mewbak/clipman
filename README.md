# Clipman

A basic clipboard manager for Wayland, with support for persisting copy buffers after an application exits.

## Installing

Requirements:

- a windows manager that uses `wlr-data-control`, like Sway and other wlroots-based WMs.
- wl-clipboard >= 2.0
- dmenu, rofi or wofi

[Install go](https://golang.org/doc/install), add `$GOPATH/bin` to your path, then run `go get github.com/yory8/clipman` OR run `go install` inside this folder.

Archlinux users can find a PKGBUILD [here](https://aur.archlinux.org/packages/clipman/).

## Usage

Run the binary in your Sway session by adding `exec wl-paste -t TEXT --watch clipman store` (or `exec wl-paste -t TEXT --watch clipman store 1>> PATH/TO/LOGFILE 2>&1 &` to log errors) at the beginning of your config.
For primary clipboard support, also add `exec wl-paste -p -t TEXT --watch clipman store --histpath="~/.local/share/clipman-primary.json`.

To query the history and select items, run the binary as `clipman pick`. You can assign it to a keybinding: `bindsym $mod+h exec clipman pick`.
For primary clipboard support, `clipman pick --histpath="~/.local/share/clipman-primary.json`.

To remove items from history, `clipman clear` and `clipman clear --all`.

To serve the last history item at startup, add `exec clipman restore` to your Sway config.

For more options: `clipman -h`.

## Known Issues

### Loss of rich text

- All items stored in history are treated as plain text.

- By default, we continue serving the last copied item even after its owner has exited. The trade-off is that we *always immediately* lose rich content: for example, if you copy some bold text in LibreOffice, when you paste it right after it will be unformatted text; or, if you copy a bookmark in Firefox, you won't be able to paste it in another bookmark folder. To disable this behaviour, you must give up persistency-after-exit by passing the `-P` option to `clipman store`. (Items manually picked from history will still be just plain text.)

## Versions

This projects follows SemVer conventions.

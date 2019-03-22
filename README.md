# Clipman

A basic clipboard manager for Wayland.

## Requirements

- a windows manager that uses `wlr-data-control`, like Sway and other wlroots-bases WMs.
- wl-clipboard from latest master (NOT v1.0)
- dmenu

## Usage

The demon that collects history is written in Go. Install it in your path, then run it in your Sway session by adding `exec clipman` at the beginning of your config.

You can configure how many history items to preserve (default: 15) by editing directly the source.

To query the history and select items, run the provided python script (`clipman.py`). You can assign it to a keybinding.

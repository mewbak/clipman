# Clipman

A basic clipboard manager for Wayland.

## Requirements

- a windows manager that uses `wlr-data-control`, like Sway and other wlroots-bases WMs.
- wl-clipboard from latest master (NOT v1.0)
- dmenu

## Usage

Install the binary in your path, then run it in your Sway session by adding `exec clipman -d` at the beginning of your config.

To query the history and select items, run the binary as `clipman -s`. You can assign it to a keybinding: `bindsym $mod+h exec clipman -s`.

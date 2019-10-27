# Next

**New features**

- `restore` command to serve the last history item, useful when run at startup.
- `--tool-args` argument to pass additional args to dmenu/rofi/etc.
- rofi and wofi now display a prompt hint to remind you whether you are picking or clearing

**Notable Bug fixes**

- we don't leak our clipboard to `ps` anymore

# 1.1

**New features**

- add support for wofi selector, a native wayland rofi clone
- serve next-to-last item when clearing last item

# 1.0

**Breaking changes**:

- switch from flags to subcommands: `wl-paste -t text --watch clipman store` instead than `clipman -d` and `clipman pick` instead than `clipman -s`
- switch demon from polling to event-driven: requires wl-clipboard >= 2.0
- rename "selector" flag to "tool"

**New features**:

- primary clipboard support: `wl-paste -p -t text --watch clipman store --histpath="~/.local/share/clipman-primary.json` and `clipman pick --histpath="~/.local/share/clipman-primary.json`
- new `clear` command for removing item(s) from history
- STDOUT tool for querying history through external tools (fzf, etc)

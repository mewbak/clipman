# 1.0

**Breaking changes**:

- switch from flags to subcommands: `wl-paste -t text --watch clipman store` instead than `clipman -d` and `clipman pick` instead than `clipman -s`
- switch demon from polling to event-driven: requires wl-clipboard >= 2.0
- rename "selector" flag to "tool"

**New features**:

- primary clipboard support: `wl-paste -p -t text --watch clipman store --histpath="~/.local/share/clipman-primary.json` and `clipman pick --histpath="~/.local/share/clipman-primary.json`
- new `clear` command for removing item(s) from history
- STDOUT tool for querying history through external tools (fzf, etc)

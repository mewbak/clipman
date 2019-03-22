#!/usr/bin/python
import json
import os
import subprocess
import sys

MAX = 15


def selector(input_, max_):
    """Send list to dmenu for selection. Display max items."""
    cmd = [
        "dmenu",
        "-b",
        "-fn",
        "-misc-dejavu sans mono-medium-r-normal--17-120-100-100-m-0-iso8859-16",
        "-l",
        str(max_),
    ]
    chosen = (
        subprocess.run(
            cmd, input="\n".join(input_).encode("utf-8"), capture_output=True
        )
        .stdout.decode()
        .strip()
    )
    return chosen


def main():
    try:
        with open(os.path.expanduser("~/.local/share/clipman.json")) as f:
            history = json.load(f)
    except FileNotFoundError:
        sys.exit("No history available")

    history = [repr(x) for x in reversed(history)]  # don't expand newlines

    selected = selector(history, MAX)

    subprocess.run(["wl-copy", selected])


if __name__ == "__main__":
    main()

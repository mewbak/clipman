package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func selector(history []string, max int, tool string) error {
	// reverse the history
	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}

	selected, err := dmenu(history, max, tool)
	if err != nil {
		return err
	}

	// serve selection to the OS
	err = exec.Command("wl-copy", []string{"--", selected}...).Run()

	return err
}

func dmenu(list []string, max int, tool string) (string, error) {
	if len(list) == 0 {
		return "", nil
	}

	if tool == "-" {
		escaped, _ := preprocessHistory(list, false)
		os.Stdout.WriteString(strings.Join(escaped, "\n"))
		return "", nil
	}

	bin, err := exec.LookPath(tool)
	if err != nil {
		return "", fmt.Errorf("%s is not installed", tool)
	}

	var args []string
	switch tool {
	case "dmenu":
		args = []string{"dmenu", "-b",
			"-fn",
			"-misc-dejavu sans mono-medium-r-normal--17-120-100-100-m-0-iso8859-16",
			"-l",
			strconv.Itoa(max)}
	case "rofi":
		args = []string{"rofi", "-dmenu",
			"-lines",
			strconv.Itoa(max)}
	default:
		return "", fmt.Errorf("Unsupported tool")
	}

	escaped, guide := preprocessHistory(list, true)
	input := strings.NewReader(strings.Join(escaped, "\n"))

	cmd := exec.Cmd{Path: bin, Args: args, Stdin: input}
	selected, err := cmd.Output()
	if err != nil {
		if err.Error() == "exit status 1" {
			// dmenu exits with this error when no selection done
			return "", nil
		}
		return "", err
	}
	trimmed := selected[:len(selected)-1] // drop newline added by dmenu

	sel, ok := guide[string(trimmed)]
	if !ok {
		return "", fmt.Errorf("couldn't recover original string; please report this bug along with a copy of your clipman.json")
	}

	return sel, nil
}

func preprocessHistory(list []string, cutting bool) ([]string, map[string]string) {
	// dmenu will break if items contain newlines, so we must pass them as literals.
	// however, when it sends them back, we need a way to restore them
	var escaped []string
	guide := make(map[string]string)

	for _, original := range list {
		repr := fmt.Sprintf("%#v", original)
		max := len(repr) - 1 // drop right quote

		// dmenu will split lines longer than 1200 something; we cut at 400 to spare memory
		if cutting {
			maxChars := 400
			if max > maxChars {
				max = maxChars
			}
		}

		repr = repr[1:max] // drop left quote
		guide[repr] = original
		escaped = append(escaped, repr)
	}

	return escaped, guide
}

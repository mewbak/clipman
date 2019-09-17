package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func selector(data []string, max int, tool string) (string, error) {
	if len(data) == 0 {
		return "", errors.New("No history available")
	}

	// output to stdout and return
	if tool == "STDOUT" {
		escaped, _ := preprocessData(data, false)
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
		return "", fmt.Errorf("Unsupported tool: %s", tool)
	}

	processed, guide := preprocessData(data, true)

	cmd := exec.Cmd{Path: bin, Args: args, Stdin: strings.NewReader(strings.Join(processed, "\n"))}
	b, err := cmd.CombinedOutput()
	if err != nil {
		if err.Error() == "exit status 1" {
			// dmenu/rofi exits with this error when no selection done
			return "", nil
		}
		return "", err
	}
	selected := string(b[:len(b)-1]) // drop newline added by dmenu/rofi

	sel, ok := guide[selected]
	if !ok {
		return "", errors.New("couldn't recover original string")
	}

	return sel, nil
}

// preprocessData:
// - reverses the data
// - escapes \n (it would break external selectors)
// - optionally it cuts items longer than 400 bytes (dmenu doesn't allow more than ~1200)
// A guide is created to allow restoring the selected item.
func preprocessData(data []string, cutting bool) ([]string, map[string]string) {
	var escaped []string
	guide := make(map[string]string)

	for i := len(data) - 1; i >= 0; i-- { // reverse slice
		original := data[i]

		// escape newlines
		repr := strings.ReplaceAll(original, "\\n", "\\\\n") // preserve literal \n
		repr = strings.ReplaceAll(repr, "\n", "\\n")

		// optionally cut to maxChars
		const maxChars = 400
		if cutting && len(repr) > maxChars {
			repr = repr[:maxChars]
		}

		guide[repr] = original
		escaped = append(escaped, repr)
	}

	return escaped, guide
}

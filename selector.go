package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func selector(history []string, max int) error {
	// reverse the history
	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}

	selected, err := dmenu(history, max)
	if err != nil {
		// dmenu exits with error when no selection done
		return nil
	}

	// serve selection to the OS
	if err := exec.Command("wl-copy", []string{"--", selected}...).Run(); err != nil {
		return err
	}

	return nil
}

func dmenu(list []string, max int) (string, error) {
	if len(list) == 0 {
		return "", nil
	}

	args := []string{"dmenu", "-b",
		"-fn",
		"-misc-dejavu sans mono-medium-r-normal--17-120-100-100-m-0-iso8859-16",
		"-l",
		strconv.Itoa(max)}

	// dmenu will break if items contain newlines, so we must pass them as literals.
	// however, when it sends them back, we need a way to restore them to non literals
	guide := make(map[string]string)
	reprList := []string{}
	for _, original := range list {
		repr := fmt.Sprintf("%#v", original)
		max := len(repr) - 1 // drop right quote
		maxChars := 400
		// dmenu will split lines longer than 1200 something; we cut at 400 to spare memory
		if max > maxChars {
			max = maxChars
		}
		repr = repr[1:max] // drop left quote
		guide[repr] = original
		reprList = append(reprList, repr)
	}

	input := strings.NewReader(strings.Join(reprList, "\n"))

	cmd := exec.Cmd{Path: "/usr/bin/dmenu", Args: args, Stdin: input}
	selected, err := cmd.Output()
	if err != nil {
		return "", err
	}
	trimmed := selected[:len(selected)-1] // drop newline

	sel, ok := guide[string(trimmed)]
	if !ok {
		return "", fmt.Errorf("couldn't recover original string; please report this bug along with a copy of your clipman.json")
	}

	return sel, nil
}

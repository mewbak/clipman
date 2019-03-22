package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func selector(history []string) error {

	// reverse the history
	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}

	selected, err := dmenu(history, max)
	if err != nil {
		// dmenu exits with error when no selection done
		return nil
	}

	if err := exec.Command("wl-copy", selected).Run(); err != nil {
		return err
	}

	return nil
}

func dmenu(list []string, max int) (string, error) {
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
		repr = repr[1 : len(repr)-1] // drop quotes
		guide[repr] = original
		reprList = append(reprList, repr)
	}

	input := strings.NewReader(strings.Join(reprList, "\n"))

	cmd := exec.Cmd{Path: "/usr/bin/dmenu", Args: args, Stdin: input}
	selected, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return guide[string(selected)], nil
}

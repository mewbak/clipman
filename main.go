package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app       = kingpin.New("clipman", "A clipboard manager for Wayland")
	histpath  = app.Flag("histpath", "Path of history file").Default("~/.local/share/clipman.json").String()
	demon     = app.Command("listen", "Run as a demon to record clipboard events")
	picker    = app.Command("pick", "Pick an item from clipboard history")
	noPersist = demon.Flag("no-persist", "Don't persist a copy buffer after a program exits").Short('P').Default("false").Bool()
	maxDemon  = demon.Flag("max-items", "history size").Default("15").Int()
	maxPicker = picker.Flag("max-items", "scrollview length").Default("15").Int()
	tool      = picker.Flag("selector", "Which selector to use: dmenu/rofi/-").Default("dmenu").String()
)

func main() {
	app.HelpFlag.Short('h')
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case "listen":
		persist := !*noPersist
		histfile, history, err := getHistory()
		if err != nil {
			log.Fatal(err)
		}
		listen(history, histfile, persist, *maxDemon)
	case "pick":
		_, history, err := getHistory()
		if err != nil {
			log.Fatal(err)
		}

		if err := selector(history, *maxPicker, *tool); err != nil {
			log.Fatal(err)
		}
	}
}

func getHistory() (string, []string, error) {
	// set histfile
	histfile := *histpath
	if strings.HasPrefix(histfile, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", nil, err
		}
		histfile = strings.Replace(histfile, "~", home, 1)
	}

	// read existing history
	var history []string
	b, err := ioutil.ReadFile(histfile)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", nil, fmt.Errorf("Failure reading history file: %s", err)
		}
	} else {
		if err := json.Unmarshal(b, &history); err != nil {
			return "", nil, fmt.Errorf("Failure parsing history: %s", err)
		}
	}

	return histfile, history, nil
}

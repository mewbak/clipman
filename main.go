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
	app        = kingpin.New("clipman", "A clipboard manager for Wayland")
	asDemon    = app.Flag("demon", "Run as a demon to record clipboard events").Short('d').Default("false").Bool()
	asSelector = app.Flag("select", "Select an item from clipboard history").Short('s').Default("false").Bool()
	noPersist  = app.Flag("no-persist", "Don't persist a copy buffer after a program exits").Short('P').Default("false").Bool()
	max        = app.Flag("max-items", "history size (with -d) or scrollview length (with -s)").Default("15").Int()
	tool       = app.Flag("selector", "Which selector to use: dmenu/rofi/-").Default("dmenu").String()
	histpath   = app.Flag("histpath", "Directory where to save history").Default("~/.local/share/clipman.json").String()
)

func main() {
	app.HelpFlag.Short('h')
	kingpin.MustParse(app.Parse(os.Args[1:]))
	modeCount := 0
	if *asDemon {
		modeCount++
	}
	if *asSelector {
		modeCount++
	}
	if modeCount != 1 {
		fmt.Println("Missing or incompatible options. You must provide exactly one of these:")
		fmt.Println("  -d, --demon")
		fmt.Println("  -s, --select")
		fmt.Println("See -h/--help for info")
		os.Exit(1)
	}

	// set histfile
	histfile := *histpath
	if strings.HasPrefix(histfile, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		histfile = strings.Replace(histfile, "~", home, 1)
	}

	// read existing history
	var history []string
	b, err := ioutil.ReadFile(histfile)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("Failure reading history file: %s", err)
		}
	} else {
		if err := json.Unmarshal(b, &history); err != nil {
			log.Fatalf("Failure parsing history: %s", err)
		}
	}

	if *asDemon {
		persist := !*noPersist
		listen(history, histfile, persist, *max)
	} else if *asSelector {
		if len(history) == 0 {
			log.Fatal("No history available")
		}
		if err := selector(history, *max, *tool); err != nil {
			log.Fatal(err)
		}
	}
}

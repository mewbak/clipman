package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app        = kingpin.New("clipman", "A clipboard manager for Wayland")
	asDemon    = app.Flag("demon", "Run as a demon to record clipboard events").Short('d').Default("false").Bool()
	asSelector = app.Flag("select", "Select an item from clipboard history").Short('s').Default("false").Bool()
	noPersist  = app.Flag("no-persist", "Don't persist a copy buffer after a program exits").Short('P').Default("false").Bool()
	max        = app.Flag("max-items", "items to store in history (with -d) or display before scrolling (with -s)").Default("15").Int()
	tool       = app.Flag("selector", "Which selector to use: dmenu/rofi").Default("dmenu").String()
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

	h, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	histfile := path.Join(h, ".local/share/clipman.json")

	var history []string
	b, err := ioutil.ReadFile(histfile)
	if err == nil {
		if err := json.Unmarshal(b, &history); err != nil {
			log.Fatalf("Failure unmarshaling history (main.38): %s", err)
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

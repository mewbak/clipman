package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"

	"gopkg.in/alecthomas/kingpin.v2"
)

const max = 15

var (
	app        = kingpin.New("clipman", "A clipboard manager for Wayland")
	asDemon    = app.Flag("demon", "Run as a demon to record clipboard events").Short('d').Default("false").Bool()
	asSelector = app.Flag("select", "Select an item from clipboard history").Short('s').Default("false").Bool()
)

var (
	histfile string
	history  []string
)

func main() {
	app.HelpFlag.Short('h')
	kingpin.MustParse(app.Parse(os.Args[1:]))
	if (*asDemon && *asSelector) || (!*asDemon && !*asSelector) {
		log.Fatal("Missing or incompatible options. See -h/--help for info")
	}

	h, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	histfile = path.Join(h, ".local/share/clipman.json")

	b, err := ioutil.ReadFile(histfile)
	if err == nil {
		if err := json.Unmarshal(b, &history); err != nil {
			log.Fatal(err)
		}
	}

	if *asDemon {
		if err := listen(history, histfile); err != nil {
			log.Fatal(err)
		}
	} else if *asSelector {
		if err := selector(history); err != nil {
			log.Fatal(err)
		}
	}
}

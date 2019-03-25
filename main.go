package main

import (
	"encoding/json"
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
	histfile := path.Join(h, ".local/share/clipman.json")

	var history []string
	b, err := ioutil.ReadFile(histfile)
	if err == nil {
		if err := json.Unmarshal(b, &history); err != nil {
			log.Fatal(err)
		}
	}

	if *asDemon {
		persist := !*noPersist
		if err := listen(history, histfile, persist, *max); err != nil {
			log.Fatal(err)
		}
	} else if *asSelector {
		if err := selector(history, *max); err != nil {
			log.Fatal(err)
		}
	}
}

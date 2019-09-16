package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app        = kingpin.New("clipman", "A clipboard manager for Wayland")
	histpath   = app.Flag("histpath", "Path of history file").Default("~/.local/share/clipman.json").String()
	storer     = app.Command("store", "Run from `wl-paste --watch` to record clipboard events")
	picker     = app.Command("pick", "Pick an item from clipboard history")
	clearer    = app.Command("clear", "Remove an item from history")
	noPersist  = storer.Flag("no-persist", "Don't persist a copy buffer after a program exits").Short('P').Default("false").Bool()
	maxDemon   = storer.Flag("max-items", "history size").Default("15").Int()
	maxPicker  = picker.Flag("max-items", "scrollview length").Default("15").Int()
	pickTool   = picker.Flag("selector", "Which selector to use: dmenu/rofi/-").Default("dmenu").String()
	clearTool  = clearer.Flag("selector", "Which selector to use: dmenu/rofi/-").Default("dmenu").String()
	maxClearer = clearer.Flag("max-items", "scrollview length").Default("15").Int()
	clearAll   = clearer.Flag("all", "Remove all items").Short('a').Default("false").Bool()
)

func main() {
	app.HelpFlag.Short('h')
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case "store":
		persist := !*noPersist

		histfile, history, err := getHistory()
		if err != nil {
			log.Fatal(err)
		}

		// read copy from stdin
		var stdin []string
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			stdin = append(stdin, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatal("Error getting input from stdin.")
		}
		text := strings.Join(stdin, "\n")

		store(text, history, histfile, *maxDemon, persist)
	case "pick":
		_, history, err := getHistory()
		if err != nil {
			log.Fatal(err)
		}

		selection, err := selector(history, *maxPicker, *pickTool)
		if err != nil {
			log.Fatal(err)
		}

		// serve selection to the OS
		err = exec.Command("wl-copy", []string{"--", selection}...).Run()
	case "clear":
		histfile, history, err := getHistory()
		if err != nil {
			log.Fatal(err)
		}

		if *clearAll {
			if err := os.Remove(histfile); err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}

		selection, err := selector(history, *maxClearer, *clearTool)
		if err != nil {
			log.Fatal(err)
		}

		if err := write(filter(history, selection), histfile); err != nil {
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

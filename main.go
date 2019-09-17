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
	app      = kingpin.New("clipman", "A clipboard manager for Wayland")
	histpath = app.Flag("histpath", "Path of history file").Default("~/.local/share/clipman.json").String()

	storer    = app.Command("store", "Record clipboard events (run as argument to `wl-paste --watch`)")
	maxDemon  = storer.Flag("max-items", "history size").Default("15").Int()
	noPersist = storer.Flag("no-persist", "Don't persist a copy buffer after a program exits").Short('P').Default("false").Bool()

	picker    = app.Command("pick", "Pick an item from clipboard history")
	maxPicker = picker.Flag("max-items", "scrollview length").Default("15").Int()
	pickTool  = picker.Flag("tool", "Which selector to use: dmenu/rofi/STDOUT").Short('t').Default("dmenu").String()

	clearer    = app.Command("clear", "Remove item(s) from history")
	maxClearer = clearer.Flag("max-items", "scrollview length").Default("15").Int()
	clearTool  = clearer.Flag("tool", "Which selector to use: dmenu/rofi/STDOUT").Short('t').Default("dmenu").String()
	clearAll   = clearer.Flag("all", "Remove all items").Short('a').Default("false").Bool()
)

func main() {
	app.HelpFlag.Short('h')
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case "store":
		histfile, history, err := getHistory(*histpath)
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

		persist := !*noPersist
		if err := store(text, history, histfile, *maxDemon, persist); err != nil {
			log.Fatal(err)
		}
	case "pick":
		_, history, err := getHistory(*histpath)
		if err != nil {
			log.Fatal(err)
		}

		selection, err := selector(history, *maxPicker, *pickTool)
		if err != nil {
			log.Fatal(err)
		}

		if selection != "" {
			// serve selection to the OS
			if err := exec.Command("wl-copy", []string{"--", selection}...).Run(); err != nil {
				log.Fatal(err)
			}
		}
	case "clear":
		histfile, history, err := getHistory(*histpath)
		if err != nil {
			log.Fatal(err)
		}

		// remove all history
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

		if selection != "" {
			if selection == history[len(history)-1] {
				// it's the latest item
				// in this case, wl-copy is still serving the copy, so wipe it
				if err := exec.Command("wl-copy", "-c").Run(); err != nil {
					log.Fatal(err)
				}
			}
			if err := write(filter(history, selection), histfile); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func getHistory(rawPath string) (string, []string, error) {
	// set histfile; expand user home
	histfile := rawPath
	if strings.HasPrefix(histfile, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", nil, err
		}
		histfile = strings.Replace(histfile, "~", home, 1)
	}

	// read history if it exists
	var history []string
	b, err := ioutil.ReadFile(histfile)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", nil, fmt.Errorf("failure reading history file: %s", err)
		}
	} else {
		if err := json.Unmarshal(b, &history); err != nil {
			return "", nil, fmt.Errorf("failure parsing history: %s", err)
		}
	}

	return histfile, history, nil
}

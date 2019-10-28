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

const version = "1.1.0"

var (
	app      = kingpin.New("clipman", "A clipboard manager for Wayland")
	histpath = app.Flag("histpath", "Path of history file").Default("~/.local/share/clipman.json").String()

	storer    = app.Command("store", "Record clipboard events (run as argument to `wl-paste --watch`)")
	maxDemon  = storer.Flag("max-items", "history size").Default("15").Int()
	noPersist = storer.Flag("no-persist", "Don't persist a copy buffer after a program exits").Short('P').Default("false").Bool()

	picker       = app.Command("pick", "Pick an item from clipboard history")
	maxPicker    = picker.Flag("max-items", "scrollview length").Default("15").Int()
	pickTool     = picker.Flag("tool", "Which selector to use: dmenu/rofi/wofi/STDOUT").Short('t').Default("dmenu").String()
	pickToolArgs = picker.Flag("tool-args", "Extra arguments to pass to the --tool").Short('T').Default("").String()

	clearer       = app.Command("clear", "Remove item(s) from history")
	maxClearer    = clearer.Flag("max-items", "scrollview length").Default("15").Int()
	clearTool     = clearer.Flag("tool", "Which selector to use: dmenu/rofi/wofi/STDOUT").Short('t').Default("dmenu").String()
	clearToolArgs = clearer.Flag("tool-args", "Extra arguments to pass to the --tool").Short('T').Default("").String()
	clearAll      = clearer.Flag("all", "Remove all items").Short('a').Default("false").Bool()

	restorer = app.Command("restore", "Serve the last recorded item from history")
)

func main() {
	app.Version(version)
	app.HelpFlag.Short('h')
	app.VersionFlag.Short('v')
	action := kingpin.MustParse(app.Parse(os.Args[1:]))

	histfile, history, err := getHistory(*histpath)
	if err != nil {
		log.Fatal(err)
	}

	switch action {
	case "store":
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
		selection, err := selector(history, *maxPicker, *pickTool, "pick", *pickToolArgs)
		if err != nil {
			log.Fatal(err)
		}

		if selection != "" {
			// serve selection to the OS
			serveTxt(selection)
		}
	case "restore":
		if len(history) == 0 {
			log.Println("Nothing to restore")
			return
		}

		serveTxt(history[len(history)-1])
	case "clear":
		// remove all history
		if *clearAll {
			if err := wipeAll(histfile); err != nil {
				log.Fatal(err)
			}
			return
		}

		selection, err := selector(history, *maxClearer, *clearTool, "clear", *clearToolArgs)
		if err != nil {
			log.Fatal(err)
		}

		if selection == "" {
			return
		}

		if len(history) < 2 {
			// there was only one possible item we could select, and we selected it,
			// so wipe everything
			if err := wipeAll(histfile); err != nil {
				log.Fatal(err)
			}
			return
		}

		if selection == history[len(history)-1] {
			// wl-copy is still serving the copy, so replace with next latest
			// note: we alread exited if less than 2 items
			serveTxt(history[len(history)-2])
		}

		if err := write(filter(history, selection), histfile); err != nil {
			log.Fatal(err)
		}
	}
}

func wipeAll(histfile string) error {
	// clear WM's clipboard
	if err := exec.Command("wl-copy", "-c").Run(); err != nil {
		return err
	}

	if err := os.Remove(histfile); err != nil {
		return err
	}

	return nil
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

func serveTxt(s string) {
	bin, err := exec.LookPath("wl-copy")
	if err != nil {
		log.Printf("couldn't find wl-copy: %v\n", err)
	}

	cmd := exec.Cmd{Path: bin, Stdin: strings.NewReader(s)}
	if err := cmd.Run(); err != nil {
		log.Printf("error running wl-copy: %s\n", err)
	}
}

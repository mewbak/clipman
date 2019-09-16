package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os/exec"
)

func store(text string, history []string, histfile string, max int, persist bool) {
	if text == "" {
		return
	}

	l := len(history)
	if l > 0 {
		if history[l-1] == text {
			return
		}

		if l >= max {
			// usually just one item, but more if we reduce our --max-items value
			history = history[l-max+1:]
		}

		// remove duplicates
		history = filter(history, text)
	}

	history = append(history, text)

	// dump history to file so that other apps can query it
	if err := write(history, histfile); err != nil {
		log.Fatalf("Fatal error writing history: %s", err)
	}

	if persist {
		// make the copy buffer available to all applications,
		// even when the source has disappeared
		if err := exec.Command("wl-copy", []string{"--", text}...).Run(); err != nil {
			log.Printf("Error running wl-copy: %s", err)
		}
	}

	return
}

func filter(history []string, text string) []string {
	var (
		found bool
		idx   int
	)

	for i, el := range history {
		if el == text {
			found = true
			idx = i
			break
		}
	}

	if found {
		// we know that idx can't be the last element, because
		// we never get to call this function if that's the case
		history = append(history[:idx], history[idx+1:]...)
	}

	return history
}

func write(history []string, histfile string) error {
	histlog, err := json.Marshal(history)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(histfile, histlog, 0644)

	return err
}

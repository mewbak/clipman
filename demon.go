package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type historyBuf struct {
	buf      []string // field name as required by io.Writer, don't change
	histfile string
	max      int
	persist  bool
}

func (hb *historyBuf) Write(p []byte) (n int, err error) {
	hb.buf = store(string(p), hb.buf, hb.histfile, hb.max, hb.persist)
	return len(p), err // signature as required by io.Writer, don't change
}

func write(history []string, histfile string) error {
	histlog, err := json.Marshal(history)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(histfile, histlog, 0644)

	return err
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

func store(text string, history []string, histfile string, max int, persist bool) []string {
	l := len(history)
	if l > 0 {
		if history[l-1] == text {
			return history
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

	return history
}

func listen(history []string, histfile string, persist bool, max int) {
	cmd := exec.Command("wl-paste", "-t", "text", "--watch", "cat")
	cmd.Stdout = &historyBuf{history, histfile, max, persist}
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("Error running wl-paste (cmd.Start): %s", err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatalf("Error running wl-paste (cmd.Wait): %s", err)
	}
}

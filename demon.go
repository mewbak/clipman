package main

import (
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"time"
)

const sleep = 1 * time.Second

func write(history []string, histfile string) error {
	histlog, err := json.Marshal(history)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(histfile, histlog, 0644)
	if err != nil {
		return err
	}
	return nil
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

func listen(history []string, histfile string) error {

	for {

		t, err := exec.Command("wl-paste", []string{"-n", "-t", "text"}...).Output()
		text := string(t)
		if err != nil || text == "" {
			// there's nothing to select, so we sleep.
			time.Sleep(sleep)
			continue
		}

		l := len(history)

		if l > 0 {

			// wl-paste will always give back the last copied text
			// (as long as the place we copied from is still running)
			if history[l-1] == text {
				time.Sleep(sleep)
				continue
			}

			if l == max {
				// we know that at any given time len(history) cannot be bigger than max,
				// so it's enough to drop the first element
				history = history[1:]
			}

			// remove duplicates
			// consider doing this in the frontend, for sparing resources
			history = filter(history, text)

		}

		history = append(history, text)

		// dump history to file so that other apps can query it
		err = write(history, histfile)
		if err != nil {
			return err
		}

		time.Sleep(sleep)
	}

}

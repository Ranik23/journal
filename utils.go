package main

import (
	"errors"
	"fmt"
	Log "log"
	"os"
	"strings"

	"github.com/jroimartin/gocui"
)

var options []string
var cursor = 0


func loadServices() (*[]string, error) {
	dirs := []string{"/etc/systemd/system", "/lib/systemd/system", "/usr/lib/systemd/system"}
	var newSlice []string

	for _, dir := range dirs {
		file, err := os.Open(dir)
		if err != nil {
			Log.Println(err)
			continue
		}

		slice, err := file.Readdirnames(-1)
		if err != nil {
			return nil, err
		}

		for _, str := range slice {
			if strings.Contains(str, ".service") {
				newSlice = append(newSlice, str)
			}
		}
	}

	return &newSlice, nil
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}


func handleLogError(err error, logView *gocui.View) error {
	if errors.Is(err, ErrNoLogsFound) {
		fmt.Fprintln(logView, "No logs found matching criteria.")
		return nil
	}
	if errors.Is(err, ErrEOF) {
		fmt.Fprintln(logView, "No logs left")
		return nil
	}
	if errors.Is(err, ErrSeekTail) {
		fmt.Fprintln(logView, "Error seek tail")
		return nil
	}
	fmt.Fprintln(logView, err.Error())
	return err
}

package main

import (
	"errors"
	"fmt"
	Log "log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
	"journal/logs"
)

var options []string
var cursor = 0

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("dropdown", 0, 0, maxX/5, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Dropdown"
		drawOptions(v)
		if _, err := setCurrentViewOnTop(g, "dropdown"); err != nil {
			return err
		}
	}

	if v, err := g.SetView("service", maxX/5+1, 0, maxX-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Chosen Service"
	}

	if v, err := g.SetView("Logs", maxX/5+1, 3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Logs"
		v.Wrap = true
	}
	return nil
}

func drawOptions(v *gocui.View) {
	v.Clear()
	green := color.New(color.FgGreen).SprintFunc()
	for i, option := range options {
		if i == cursor {
			fmt.Fprintf(v, "%s\n", green("-> " + option))
		} else {
			fmt.Fprintf(v, "-> %s\n", option)
		}
	}
}
func nextOption(g *gocui.Gui, v *gocui.View) error {
	if cursor < len(options)-1 {
		cursor++
	}
	_, oy := v.Origin()
	v.SetOrigin(0, oy+1)
	return updateDropdown(g)
}

func prevOption(g *gocui.Gui, v *gocui.View) error {
	if cursor > 0 {
		cursor--
	}
	_, oy := v.Origin()
	v.SetOrigin(0, oy-1)
	return updateDropdown(g)
}

func updateDropdown(g *gocui.Gui) error {
	v, err := g.View("dropdown")
	if err != nil {
		return err
	}
	drawOptions(v)
	return nil
}

func selectOption(g *gocui.Gui, v *gocui.View) error {

	selected := options[cursor]

	logView, err := g.View("Logs")
	if err != nil {
		return err
	}

	logView.Clear()

	serviceView, err := g.View("service")
	if err != nil {
		return err
	}
	serviceView.Clear()
	fmt.Fprintf(serviceView, "%s", selected)

	array, err := log.FetchLogs(selected, "", 50)
	if err != nil {
		if errors.Is(err, log.ErrNoLogsFound) {
			fmt.Fprintln(logView, "No logs found matching criteria.")
			return nil
		}
		if errors.Is(err, log.ErrEOF) {
			fmt.Fprintln(logView, "No logs left")
			return nil
		}

		if errors.Is(err, log.ErrSeekTail) {
			fmt.Fprintln(logView, "Error seek tail")
			return nil
		}

		fmt.Println(logView, err.Error())
		return err
	}

	for _, entity := range *array {
		fmt.Fprintln(logView, entity)
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

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

func main() {
	services, err := loadServices()
	if err != nil {
		Log.Panicln(err)
	}

	options = *services

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		Log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Mouse = true
	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("dropdown", gocui.KeyArrowDown, gocui.ModNone, nextOption); err != nil {
		Log.Panicln(err)
	}
	if err := g.SetKeybinding("dropdown", gocui.KeyArrowUp, gocui.ModNone, prevOption); err != nil {
		Log.Panicln(err)
	}
	if err := g.SetKeybinding("dropdown", gocui.KeyEnter, gocui.ModNone, selectOption); err != nil {
		Log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		Log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		Log.Panicln(err)
	}
}

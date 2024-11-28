package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
)

func enterServiceName(g *gocui.Gui, v *gocui.View) error {
	serviceView, _ := g.View("service")
	serviceView.Clear()
	v.SetCursor(0, 0)
	text := v.Buffer()
	flag := false

	for i, service := range options {
		if service == text {
			cursor = i
			flag = true
			break
		}
	}

	dropdownView, _ := g.View("dropdown")

	if !flag {
		fmt.Fprint(serviceView, "Service Not Found: ", text)
	} else {
		selectOption(g, dropdownView)

		if _, err := setCurrentViewOnTop(g, "dropdown"); err != nil {
			return err
		}
	}

	v.Clear()

	return nil
}

func drawOptions(v *gocui.View) {
	v.Clear()
	green := color.New(color.FgGreen).SprintFunc()
	for i, option := range options {
		if i == cursor {
			fmt.Fprintf(v, "%s\n", green("-> "+option))
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
	if err := v.SetOrigin(0, oy+1); err != nil {
		return err
	}
	drawOptions(v)
	return nil
}

func prevOption(g *gocui.Gui, v *gocui.View) error {
	if cursor > 0 {
		cursor--
	}
	_, oy := v.Origin()
	if err := v.SetOrigin(0, oy-1); err != nil {
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

	array, err := FetchLogs(selected, "", 50)
	if err != nil {
		logView.Title = fmt.Sprintf("Logs(%d)", 0)
		return handleLogError(err, logView)
	} else {
		logView.Title = fmt.Sprintf("Logs(%d)", len(*array))
	}

	for _, entity := range *array {
		fmt.Fprintln(logView, entity)
	}

	return nil
}

func quitOption(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func switchToDropDown(g *gocui.Gui, v *gocui.View) error {
	if _, err := setCurrentViewOnTop(g, "dropdown"); err != nil {
		return err
	}
	return nil
}

func switchToSearch(g *gocui.Gui, v *gocui.View) error {
	if _, err := setCurrentViewOnTop(g, "search"); err != nil {
		return err
	}
	return nil
}

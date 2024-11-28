package main

import (
	Log "log"
	"github.com/jroimartin/gocui"
)

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
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quitOption); err != nil {
		Log.Panicln(err)
	}

	if err := g.SetKeybinding("search", gocui.KeyEnter, gocui.ModNone, enterServiceName); err != nil {
		Log.Panicln(err)
	}

	if err := g.SetKeybinding("search", gocui.KeyCtrlQ, gocui.ModNone, switchToDropDown); err != nil {
		Log.Panicln(err)
	}

	if err := g.SetKeybinding("dropdown", gocui.KeyCtrlQ, gocui.ModNone, switchToSearch); err != nil {
		Log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		Log.Panicln(err)
	}
}

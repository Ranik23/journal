package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("dropdown", 0, 3, maxX/5, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Dropdown"
		drawOptions(v)
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
		v.Wrap = true
	}

	if v, err := g.SetView("search", 0, 0, maxX/5, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Title = "Search Service"
		v.Overwrite = true

		if _, err := setCurrentViewOnTop(g, "search"); err != nil {
			return err
		}
	}

	if v, err := g.SetView("text", 0, maxY/2+3, maxX/5, maxY/2+5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Commands"
		fmt.Fprint(v, "Use Ctrl+Q To Switch\n")
	}
	return nil
}

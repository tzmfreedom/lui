package main

import (
	"fmt"

	"github.com/tzmfreedom/gocui"
)

type Menu struct {
	x, y, w, h int
}

func (w *Menu) Layout(g *gocui.Gui) error {
	if v, err := g.SetView("Menu", w.x, w.y, w.x+w.w, w.y+w.h); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		menues := []string{
			"Execute Query",
			"Describe SObject",
			"Execute Anonymous",
		}
		for _, menu := range menues {
			fmt.Fprintln(v, menu)
		}
		v.Title = "Menu"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.SetCursor(0, 0)

		if err := g.SetKeybinding("Menu", gocui.KeyArrowUp, gocui.ModNone, up(0)); err != nil {
			return err
		}
		if err := g.SetKeybinding("Menu", 'k', gocui.ModNone, up(0)); err != nil {
			return err
		}
		if err := g.SetKeybinding("Menu", gocui.KeyArrowDown, gocui.ModNone, down(2)); err != nil {
			return err
		}
		if err := g.SetKeybinding("Menu", 'j', gocui.ModNone, down(2)); err != nil {
			return err
		}
		if err := g.SetKeybinding("Menu", gocui.MouseLeft, gocui.ModNone, setCurrentView); err != nil {
			return err
		}

		g.SetCurrentView("Menu")
	}
	return nil
}

func newMenu(x, y, w, h int) *Menu {
	return &Menu{x, y, w, h}
}

package main

import (
	"fmt"

	"github.com/tzmfreedom/gocui"
)

type DebugView struct {
	x, y, w, h int
}

func (w *DebugView) Layout(g *gocui.Gui) error {
	if _, err := g.SetView("Debug", w.x, w.y, w.x+w.w, w.y+w.h); err != gocui.ErrUnknownView {
		return err
	}
	return nil
}

func newDebugView(x, y, w, h int) *DebugView {
	return &DebugView{x, y, w, h}
}

func debug(g *gocui.Gui, text string) error {
	v, err := g.SetCurrentView("Debug")
	if err != nil {
		return err
	}
	fmt.Fprintf(v, text)
	return nil
}

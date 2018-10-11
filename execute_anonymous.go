package main

import "github.com/tzmfreedom/gocui"

type ExecuteAnonymous struct {
	x, y, w, h int
}

func (w *ExecuteAnonymous) Layout(g *gocui.Gui) error {
	if v, err := g.SetView("ExecuteAnonymous", w.x, w.y, w.x+w.w, w.y+w.h); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Execute Anonymous"
		v.Editable = true
		v.SetCursor(0, 0)

		if err := g.SetKeybinding("ExecuteAnonymous", gocui.MouseRelease, gocui.ModNone, setCurrentViewForEditor); err != nil {
			return err
		}
	}
	return nil
}

func newExecuteAnonymous(x, y, w, h int) *ExecuteAnonymous {
	return &ExecuteAnonymous{x, y, w, h}
}

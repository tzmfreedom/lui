package main

import "github.com/tzmfreedom/gocui"

type SoqlEditor struct {
	x, y, w, h int
}

func (w *SoqlEditor) Layout(g *gocui.Gui) error {
	if v, err := g.SetView("SoqlEditor", w.x, w.y, w.x+w.w, w.y+w.h); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "SOQL Editor"
		v.Editable = true
		v.SetCursor(0, 0)

		if err := g.SetKeybinding("SoqlEditor", gocui.MouseRelease, gocui.ModNone, setCurrentViewForEditor); err != nil {
			return err
		}
	}
	return nil
}

func newSoqlEditor(x, y, w, h int) *SoqlEditor {
	return &SoqlEditor{x, y, w, h}
}

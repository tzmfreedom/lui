package main

import (
	"github.com/tzmfreedom/go-soapforce"
	"github.com/tzmfreedom/gocui"
)

type DescribeView struct {
	x, y, w, h  int
	View        *gocui.View
	Records     []*soapforce.SObject
	SObjectType string

	RecordView *RecordView
}

type Describe struct {
	x, y, w, h int
}

func (w *Describe) Layout(g *gocui.Gui) error {
	v, err := g.SetView("Describe", w.x, w.y, w.x+w.w, w.y+w.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "List View"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.SetCursor(0, 1)

		if err := g.SetKeybinding("Describe", gocui.KeyArrowUp, gocui.ModNone, up(1)); err != nil {
			return err
		}
		if err := g.SetKeybinding("Describe", 'k', gocui.ModNone, up(1)); err != nil {
			return err
		}
		if err := g.SetKeybinding("Describe", gocui.KeyArrowDown, gocui.ModNone, down(w.h-4)); err != nil {
			return err
		}
		if err := g.SetKeybinding("Describe", 'j', gocui.ModNone, down(w.h-2)); err != nil {
			return err
		}
		if err := g.SetKeybinding("Describe", gocui.MouseLeft, gocui.ModNone, setCurrentView); err != nil {
			return err
		}
	}
	return nil
}

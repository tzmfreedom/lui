package main

import (
	"fmt"

	"github.com/tzmfreedom/go-soapforce"
	"github.com/tzmfreedom/gocui"
)

type ListView struct {
	x, y, w, h int
	Records    []*soapforce.SObject
}

func (w *ListView) Layout(g *gocui.Gui) error {
	if v, err := g.SetView("ListView", w.x, w.y, w.x+w.w, w.y+w.h); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		headers := []string{
			"LastName",
		}
		for _, h := range headers {
			fmt.Fprintln(v, h)
		}
		for i, r := range w.Records {
			for _, h := range headers {
				fmt.Fprintln(v, r.Fields[h])
			}
			if i >= w.h-4 {
				break
			}
		}
		v.Title = "List View"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.SetCursor(0, 1)

		if err := g.SetKeybinding("ListView", gocui.KeyArrowUp, gocui.ModNone, up(1)); err != nil {
			return err
		}
		if err := g.SetKeybinding("ListView", 'k', gocui.ModNone, up(1)); err != nil {
			return err
		}
		if err := g.SetKeybinding("ListView", gocui.KeyArrowDown, gocui.ModNone, down(w.h-4)); err != nil {
			return err
		}
		if err := g.SetKeybinding("ListView", 'j', gocui.ModNone, down(w.h-2)); err != nil {
			return err
		}
		if err := g.SetKeybinding("ListView", gocui.MouseLeft, gocui.ModNone, setCurrentView); err != nil {
			return err
		}
	}
	return nil
}

func newListView(x, y, w, h int, r []*soapforce.SObject) *ListView {
	return &ListView{x, y, w, h, r}
}

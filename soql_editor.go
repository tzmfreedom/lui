package main

import (
	"fmt"

	"bytes"

	"github.com/tzmfreedom/gocui"
)

type SoqlEditor struct {
	x, y, w, h int
	listView   *ListView
}

func (w *SoqlEditor) Layout(g *gocui.Gui) error {
	if v, err := g.SetView("SoqlEditor", w.x, w.y, w.x+w.w, w.y+w.h); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "SOQL Editor"
		v.Editable = true
		v.SetCursor(0, 0)
		fmt.Fprintf(v, "SELECT Id, LastName, FirstName FROM Contact")

		if err := g.SetKeybinding("SoqlEditor", gocui.MouseRelease, gocui.ModNone, setCurrentViewForEditor); err != nil {
			return err
		}
		if err := g.SetKeybinding("SoqlEditor", gocui.KeyCtrlE, gocui.ModNone, execSOQL(w.listView)); err != nil {
			return err
		}
		if err := g.SetKeybinding("SoqlEditor", gocui.KeyCtrlS, gocui.ModNone, saveSOQL); err != nil {
			return err
		}
		if err := g.SetKeybinding("SoqlEditor", gocui.KeyCtrlC, gocui.ModNone, copySOQL); err != nil {
			return err
		}
		if err := g.SetKeybinding("SoqlEditor", gocui.KeyCtrlV, gocui.ModNone, copySOQL); err != nil {
			return err
		}
		g.SetCurrentView("SoqlEditor")
	}
	return nil
}

func readSOQL(v *gocui.View) (string, error) {
	v.Rewind()
	buf := &bytes.Buffer{}
	_, err := buf.ReadFrom(v)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func execSOQL(l *ListView) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		soql, err := readSOQL(v)
		if err != nil {
			return err
		}
		g.SetCurrentView("ListView")
		return l.Render(soql)
	}
}

func saveSOQL(g *gocui.Gui, v *gocui.View) error {
	_, err := readSOQL(v)
	if err != nil {
		return err
	}
	return nil
}

func copySOQL(g *gocui.Gui, v *gocui.View) error {
	_, err := readSOQL(v)
	if err != nil {
		return err
	}
	return nil
}

func pasteSOQL(g *gocui.Gui, v *gocui.View) error {
	_, err := readSOQL(v)
	if err != nil {
		return err
	}
	return nil
}

func newSoqlEditor(x, y, w, h int, listView *ListView) *SoqlEditor {
	return &SoqlEditor{x, y, w, h, listView}
}

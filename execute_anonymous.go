package main

import (
	"bytes"

	"github.com/tzmfreedom/go-soapforce"
	"github.com/tzmfreedom/gocui"
)

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
		if err := g.SetKeybinding("ExecuteAnonymous", gocui.KeyCtrlE, gocui.ModNone, execApexCode); err != nil {
			return err
		}
		if err := g.SetKeybinding("ExecuteAnonymous", gocui.KeyCtrlS, gocui.ModNone, saveSOQL); err != nil {
			return err
		}
		if err := g.SetKeybinding("ExecuteAnonymous", gocui.KeyCtrlC, gocui.ModNone, copySOQL); err != nil {
			return err
		}
		if err := g.SetKeybinding("ExecuteAnonymous", gocui.KeyCtrlV, gocui.ModNone, copySOQL); err != nil {
			return err
		}
		g.SetCurrentView("SoqlEditor")
	}
	return nil
}

func execApexCode(g *gocui.Gui, v *gocui.View) error {
	code, err := readApexCode(v)
	if err != nil {
		return err
	}
	categories := []*soapforce.LogInfo{
		{
			Category: soapforce.LogCategoryApex_code,
			Level:    soapforce.LogCategoryLevelFinest,
		},
	}
	client.SetDebuggingHeader(categories)
	res, err := client.ExecuteAnonymous(code)
	if err != nil {
		return err
	}
	debug(res)
	return nil
}

func readApexCode(v *gocui.View) (string, error) {
	v.Rewind()
	buf := &bytes.Buffer{}
	_, err := buf.ReadFrom(v)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func newExecuteAnonymous(x, y, w, h int) *ExecuteAnonymous {
	return &ExecuteAnonymous{x, y, w, h}
}

package main

import (
	"fmt"
	"strings"

	"github.com/tzmfreedom/go-soapforce"
	"github.com/tzmfreedom/gocui"
)

type DescribeField struct {
	x, y, w, h int
	Fields     []*soapforce.Field
}

func (w *DescribeField) Render(g *gocui.Gui) error {
	v, err := g.SetView("DescribeField", w.x, w.y, w.x+w.w, w.y+w.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Describe Field"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.SetCursor(0, 0)

		for _, f := range w.Fields {
			vals := []string{
				display(f.Label, describeColWidth),
				display(f.Name, describeColWidth),
				display(string(*f.Type_), describeColWidth),
			}
			fmt.Fprintln(v, strings.Join(vals, "|"))
		}

		if err := g.SetKeybinding("DescribeField", gocui.KeyArrowUp, gocui.ModNone, up(0)); err != nil {
			return err
		}
		if err := g.SetKeybinding("DescribeField", 'k', gocui.ModNone, up(0)); err != nil {
			return err
		}
		if err := g.SetKeybinding("DescribeField", gocui.KeyArrowDown, gocui.ModNone, down(len(w.Fields))); err != nil {
			return err
		}
		if err := g.SetKeybinding("DescribeField", 'j', gocui.ModNone, down(len(w.Fields))); err != nil {
			return err
		}
		if err := g.SetKeybinding("DescribeField", 'G', gocui.ModNone, toBottom); err != nil {
			return err
		}
		if err := g.SetKeybinding("DescribeField", 'g', gocui.ModNone, toTop); err != nil {
			return err
		}
		if err := g.SetKeybinding("DescribeField", gocui.MouseLeft, gocui.ModNone, setCurrentView); err != nil {
			return err
		}
		if err := g.SetKeybinding("DescribeField", 'q', gocui.ModNone, backToDescribe); err != nil {
			return err
		}
		g.SetCurrentView("DescribeField")
	}
	return nil
}

func backToDescribe(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybindings("DescribeField")
	g.DeleteView("DescribeField")
	g.SetCurrentView("Describe")
	return nil
}

func moveToDescribe(g *gocui.Gui, v *gocui.View) error {
	g.SetCurrentView("Describe")
	return nil
}

func newDescribeField(x, y, w, h int, fields []*soapforce.Field) *DescribeField {
	return &DescribeField{x, y, w, h, fields}
}

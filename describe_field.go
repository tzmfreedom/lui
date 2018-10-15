package main

import (
	"fmt"
	"strings"

	"github.com/tzmfreedom/go-soapforce"
	"github.com/tzmfreedom/gocui"
)

type DescribeField struct {
	x, y, w, h    int
	DisplayFields []*soapforce.Field
	Fields        []*soapforce.Field
	View          *gocui.View
}

func (w *DescribeField) Layout(g *gocui.Gui) error {
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
		w.View = v

		w.Render()

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
		if err := g.SetKeybinding("DescribeField", gocui.KeyCtrlF, gocui.ModNone, w.showSearchBox); err != nil {
			return err
		}
		g.SetCurrentView("DescribeField")
	}
	return nil
}

func (w *DescribeField) Render() {
	w.View.Clear()
	for _, f := range w.DisplayFields {
		vals := []string{
			display(f.Label, describeColWidth),
			display(f.Name, describeColWidth),
			display(string(*f.Type_), describeColWidth),
		}
		fmt.Fprintln(w.View, strings.Join(vals, "|"))
	}
}

func (w *DescribeField) showSearchBox(g *gocui.Gui, v *gocui.View) error {
	fieldSearchBox := newFieldSearchBox(w.x+20, w.y-2, w.w-20, 2, w)
	return fieldSearchBox.Render(g)
}

func backToDescribe(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybindings("DescribeField")
	g.DeleteView("DescribeField")
	g.DeleteKeybindings("FieldSearchBox")
	g.DeleteView("FieldSearchBox")
	g.SetCurrentView("Describe")
	return nil
}

func moveToDescribe(g *gocui.Gui, v *gocui.View) error {
	g.SetCurrentView("Describe")
	return nil
}

func newDescribeField(x, y, w, h int, fields []*soapforce.Field) *DescribeField {
	return &DescribeField{x, y, w, h, fields, fields, nil}
}

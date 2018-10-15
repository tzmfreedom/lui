package main

import (
	"fmt"
	"strings"

	"github.com/tzmfreedom/go-soapforce"
	"github.com/tzmfreedom/gocui"
)

type DescribeView struct {
	x, y, w, h int
	SObjects   []*soapforce.DescribeGlobalSObjectResult
	View       *gocui.View

	dx, dy, dw, dh int
}

var describeColWidth = 20

func (w *DescribeView) Layout(g *gocui.Gui) error {
	v, err := g.SetView("Describe", w.x, w.y, w.x+w.w, w.y+w.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Describe Global"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.SetCursor(0, 0)
		w.View = v

		//headers := []string{
		//	display("Label", describeColWidth),
		//	display("Name", describeColWidth),
		//}
		//fmt.Fprintln(v, strings.Join(headers, "|"))
		//fmt.Fprintln(v, strings.Repeat("─", w.w-2))
		for _, sobj := range w.SObjects {
			vals := []string{
				display(sobj.Label, describeColWidth),
				display(sobj.Name, describeColWidth),
			}
			fmt.Fprintln(v, strings.Join(vals, "|"))
		}

		if err := g.SetKeybinding("Describe", gocui.KeyArrowUp, gocui.ModNone, up(0)); err != nil {
			return err
		}
		if err := g.SetKeybinding("Describe", 'k', gocui.ModNone, up(0)); err != nil {
			return err
		}
		if err := g.SetKeybinding("Describe", gocui.KeyArrowDown, gocui.ModNone, down(len(w.SObjects))); err != nil {
			return err
		}
		if err := g.SetKeybinding("Describe", 'j', gocui.ModNone, down(len(w.SObjects))); err != nil {
			return err
		}
		if err := g.SetKeybinding("Describe", 'G', gocui.ModNone, toBottom); err != nil {
			return err
		}
		if err := g.SetKeybinding("Describe", 'g', gocui.ModNone, toTop); err != nil {
			return err
		}
		if err := g.SetKeybinding("Describe", gocui.MouseLeft, gocui.ModNone, setCurrentView); err != nil {
			return err
		}
		if err := g.SetKeybinding("Describe", gocui.KeyEnter, gocui.ModNone, w.showFields); err != nil {
			return err
		}
		if err := g.SetKeybinding("Describe", gocui.KeyCtrlF, gocui.ModNone, w.showSearchBox); err != nil {
			return err
		}
	}
	return nil
}

func (w *DescribeView) Render() error {
	w.View.Clear()
	for _, m := range w.SObjects {
		vals := []string{
			display(m.Label, describeColWidth),
			display(m.Name, describeColWidth),
		}
		fmt.Fprintln(w.View, strings.Join(vals, "|"))
	}
	return nil
}

func (w *DescribeView) showFields(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	sobj, err := getDescribeSObjectResult(w.SObjects[cy].Name)
	if err != nil {
		return err
	}
	field := newDescribeField(w.dx, w.dy, w.dw, w.dh, sobj.Fields)
	return field.Render(g)
}

func (w *DescribeView) showSearchBox(g *gocui.Gui, v *gocui.View) error {
	searchBox := newSearchBox(w.x+20, w.y-2, w.w-20, 2, w)
	return searchBox.Render(g)
}

func newDescribeView(x, y, w, h int, sobjects []*soapforce.DescribeGlobalSObjectResult, dx, dy, dw, dh int) *DescribeView {
	return &DescribeView{x, y, w, h, sobjects, nil, dx, dy, dw, dh}
}

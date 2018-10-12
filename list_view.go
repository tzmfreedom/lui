package main

import (
	"fmt"

	"strings"

	"github.com/tzmfreedom/go-soapforce"
	"github.com/tzmfreedom/gocui"
)

type ListView struct {
	x, y, w, h int
	View       *gocui.View
	Records    []*soapforce.SObject
}

func (w *ListView) Layout(g *gocui.Gui) error {
	v, err := g.SetView("ListView", w.x, w.y, w.x+w.w, w.y+w.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		w.View = v
		v.Title = "List View"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.SetCursor(0, 1)
		w.Render("")

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

func (w *ListView) Render(soql string) error {
	w.View.Clear()
	if soql == "" {
		return nil
	}
	soql, err := expandSOQL(soql)
	if err != nil {
		return err
	}

	headers := []string{
		"LastName",
	}
	for _, h := range headers {
		fmt.Fprintln(w.View, h)
	}

	r, err := client.Query(soql)
	if err != nil {
		return err
	}
	w.Records = r.Records

	if w.Records != nil {
		for i, r := range w.Records {
			for _, h := range headers {
				fmt.Fprintln(w.View, r.Fields[h])
			}
			if i >= w.h-4 {
				break
			}
		}
	}
	return nil
}

func expandSOQL(soql string) (string, error) {
	sobject := "Contact"
	res, err := client.DescribeSObject(sobject)
	if err != nil {
		return "", err
	}
	fields := make([]string, len(res.Fields))
	for i, f := range res.Fields {
		fields[i] = f.Name
	}
	f := strings.Join(fields, ", ")
	return strings.Replace(soql, "*", f, 0), nil
}

func newListView(x, y, w, h int) *ListView {
	return &ListView{x, y, w, h, nil, nil}
}

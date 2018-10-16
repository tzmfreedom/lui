package main

import (
	"fmt"
	"strings"

	"github.com/sahilm/fuzzy"
	"github.com/tzmfreedom/go-soapforce"
	"github.com/tzmfreedom/gocui"
)

type FieldSearchBox struct {
	x, y, w, h    int
	DescribeField *DescribeField
}

func (w *FieldSearchBox) Render(g *gocui.Gui) error {
	if v, err := g.SetView("FieldSearchBox", w.x, w.y, w.x+w.w, w.y+w.h); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "FieldSearchBox"
		v.Editable = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.SetCursor(0, 0)
		v.Editor = gocui.EditorFunc(w.searchBox)

		if err := g.SetKeybinding("FieldSearchBox", gocui.MouseLeft, gocui.ModNone, setCurrentViewForEditor); err != nil {
			return err
		}
		if err := g.SetKeybinding("FieldSearchBox", gocui.KeyEsc, gocui.ModNone, backToDescribe); err != nil {
			return err
		}
		if err := g.SetKeybinding("FieldSearchBox", gocui.KeyEnter, gocui.ModNone, moveToDescribeField); err != nil {
			return err
		}
	}
	g.SetCurrentView("FieldSearchBox")
	g.Cursor = true
	return nil
}

func (s *FieldSearchBox) search(v *gocui.View) error {
	pattern, err := readText(v)
	if err != nil {
		return err
	}
	pattern = strings.TrimSpace(pattern)
	df := s.DescribeField
	if pattern == "" {
		df.DisplayFields = df.Fields
		df.Render()
		return nil
	}
	names := make([]string, len(df.Fields))
	mapToField := map[string]*soapforce.Field{}
	for i, f := range df.Fields {
		target := fmt.Sprintf("%s %s", f.Label, f.Name)
		names[i] = target
		mapToField[target] = f
	}
	matches := fuzzy.Find(pattern, names)
	displayFields := make([]*soapforce.Field, len(matches))
	for i, m := range matches {
		displayFields[i] = mapToField[m.Str]
	}
	df.DisplayFields = displayFields
	df.Render()
	return nil
}

func (s *FieldSearchBox) searchBox(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
		s.search(v)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
		s.search(v)
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
		s.search(v)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
		s.search(v)
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	}
}

func moveToDescribeField(g *gocui.Gui, v *gocui.View) error {
	g.SetCurrentView("DescribeField")
	return nil
}

func newFieldSearchBox(x, y, w, h int, df *DescribeField) *FieldSearchBox {
	return &FieldSearchBox{x, y, w, h, df}
}

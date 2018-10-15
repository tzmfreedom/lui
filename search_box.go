package main

import (
	"bytes"

	"strings"

	"github.com/sahilm/fuzzy"
	"github.com/tzmfreedom/go-soapforce"
	"github.com/tzmfreedom/gocui"
)

type SearchBox struct {
	x, y, w, h int
	Describe   *DescribeView
}

func (w *SearchBox) Render(g *gocui.Gui) error {
	if v, err := g.SetView("SearchBox", w.x, w.y, w.x+w.w, w.y+w.h); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "SearchBox"
		v.Editable = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.SetCursor(0, 0)
		v.Editor = gocui.EditorFunc(w.searchBox)

		if err := g.SetKeybinding("SearchBox", gocui.MouseLeft, gocui.ModNone, setCurrentViewForEditor); err != nil {
			return err
		}
		if err := g.SetKeybinding("SearchBox", gocui.KeyEsc, gocui.ModNone, backToDescribe); err != nil {
			return err
		}
		if err := g.SetKeybinding("SearchBox", gocui.KeyEnter, gocui.ModNone, moveToDescribe); err != nil {
			return err
		}
	}
	g.SetCurrentView("SearchBox")
	return nil
}

func (s *SearchBox) search(v *gocui.View) error {
	pattern, err := readText(v)
	if err != nil {
		return err
	}
	pattern = strings.TrimSpace(pattern)
	dv := s.Describe
	if pattern == "" {
		dv.SObjects = descGlobalResults
		return dv.Render()
	}
	names := make([]string, len(descGlobalResults))
	for i, sobj := range descGlobalResults {
		names[i] = sobj.Name
	}
	mapToSObject := map[string]*soapforce.DescribeGlobalSObjectResult{}
	for _, sobj := range descGlobalResults {
		mapToSObject[sobj.Name] = sobj
	}
	matches := fuzzy.Find(pattern, names)
	newSObjects := make([]*soapforce.DescribeGlobalSObjectResult, len(matches))
	for i, m := range matches {
		newSObjects[i] = mapToSObject[m.Str]
	}
	dv.SObjects = newSObjects
	return dv.Render()
}

func readText(v *gocui.View) (string, error) {
	v.Rewind()
	buf := &bytes.Buffer{}
	_, err := buf.ReadFrom(v)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (s *SearchBox) searchBox(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
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

func newSearchBox(x, y, w, h int, d *DescribeView) *SearchBox {
	return &SearchBox{x, y, w, h, d}
}

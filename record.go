package main

import (
	"fmt"
	"regexp"

	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/tzmfreedom/go-soapforce"
	"github.com/tzmfreedom/gocui"
)

type RecordView struct {
	x, y, w, h  int
	SObjectType string
	Record      *soapforce.SObject
}

func (w *RecordView) Render(g *gocui.Gui) error {
	recordView, err := g.SetView("Record", w.x, w.y, w.x+w.w, w.y+w.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		recordView.Title = "Record Detail"
		recordView.Highlight = true
		recordView.SelBgColor = gocui.ColorGreen
		recordView.SelFgColor = gocui.ColorBlack
		recordView.SetCursor(0, 1)

		recordView.Clear()
		r, err := getRecordDetail(w.SObjectType, w.Record.Id)
		if err != nil {
			return err
		}
		sobj, err := getDescribeSObjectResult(w.SObjectType)
		max := maxFieldLabelLength(sobj.Fields)

		fmt.Fprintln(recordView, fmt.Sprintf("%s | %s", runewidth.FillRight("ID", max), r.Id))
		if err != nil {
			return err
		}

		for _, f := range sobj.Fields {
			value, ok := r.Fields[f.Name]
			if !ok {
				value = ""
			}
			fmt.Fprintln(recordView, fmt.Sprintf("%s | %s", runewidth.FillRight(f.Label, max), value))
		}

		if err := g.SetKeybinding("Record", gocui.KeyArrowUp, gocui.ModNone, up(0)); err != nil {
			return err
		}
		if err := g.SetKeybinding("Record", 'k', gocui.ModNone, up(0)); err != nil {
			return err
		}
		if err := g.SetKeybinding("Record", gocui.KeyArrowDown, gocui.ModNone, down(len(sobj.Fields)-2)); err != nil {
			return err
		}
		if err := g.SetKeybinding("Record", 'j', gocui.ModNone, down(len(sobj.Fields)-2)); err != nil {
			return err
		}
		if err := g.SetKeybinding("Record", gocui.MouseLeft, gocui.ModNone, setCurrentView); err != nil {
			return err
		}
		if err := g.SetKeybinding("Record", 'q', gocui.ModNone, backToList); err != nil {
			return err
		}
		g.SetCurrentView("Record")
	}
	return nil
}

func getRecordDetail(sobjectType, id string) (*soapforce.SObject, error) {
	sobj, err := getDescribeSObjectResult(sobjectType)
	if err != nil {
		return nil, err
	}
	fields := make([]string, len(sobj.Fields))
	for i, f := range sobj.Fields {
		fields[i] = f.Name
	}
	result, err := client.Query(fmt.Sprintf("SELECT %s FROM %s WHERE id = '%s'", strings.Join(fields, ", "), sobjectType, id))
	if err != nil {
		return nil, err
	}
	return result.Records[0], nil
}

func maxFieldLabelLength(results []*soapforce.Field) int {
	max := 0
	for _, result := range results {
		l := runewidth.StringWidth(result.Label)
		if l > max {
			max = l
		}
	}
	return max
}

func display(label string, max int) string {
	l := max - runewidth.StringWidth(label)
	if l > 0 {
		return runewidth.FillRight(label, max)
	}
	return runewidth.Truncate(label, max, "...")
}

func backToList(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybindings("Record")
	g.DeleteView("Record")
	g.SetCurrentView("ListView")
	return nil
}

func expandSOQL(soql string) (string, error) {
	sobjectType, err := getSobjectFromSoql(soql)
	if err != nil {
		return "", err
	}
	res, err := getDescribeSObjectResult(sobjectType)
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

func getSobjectFromSoql(soql string) (string, error) {
	r, err := regexp.Compile(`(?i)^SELECT\s+.*\s+FROM\s+([\d_a-zA-Z]+)`)
	if err != nil {
		return "", err
	}
	matches := r.FindAllStringSubmatch(soql, -1)
	return matches[0][1], nil
}

func getFields(soql string) ([]string, error) {
	r, err := regexp.Compile(`(?i)^SELECT\s+(.*)\s+FROM\s+([\d_a-zA-Z]+)`)
	if err != nil {
		return nil, err
	}
	matches := r.FindAllStringSubmatch(soql, -1)
	unprocessedFields := strings.Split(matches[0][1], ",")
	fields := make([]string, len(unprocessedFields))
	for i, f := range unprocessedFields {
		fields[i] = strings.TrimSpace(f)
	}
	return fields, nil
}

func newRecordView(x, y, w, h int, sobjectType string, record *soapforce.SObject) *RecordView {
	return &RecordView{x, y, w, h, sobjectType, record}
}

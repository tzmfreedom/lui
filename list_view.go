package main

import (
	"fmt"
	"os/exec"
	"runtime"

	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/tzmfreedom/go-soapforce"
	"github.com/tzmfreedom/gocui"
)

type ListView struct {
	x, y, w, h  int
	View        *gocui.View
	Records     []*soapforce.SObject
	SObjectType string

	RecordView *RecordView
}

type RecordView struct {
	x, y, w, h int
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
		if err := g.SetKeybinding("ListView", gocui.KeyEnter, gocui.ModNone, w.ShowRecord); err != nil {
			return err
		}
		if err := g.SetKeybinding("ListView", 'o', gocui.ModNone, w.OpenRecordDetail); err != nil {
			return err
		}
		if err := g.SetKeybinding("ListView", 'e', gocui.ModNone, w.OpenRecordEdit); err != nil {
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
	w.SObjectType = "Contact"

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

func (w *ListView) OpenRecordDetail(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	r := w.Records[cy-1]
	url, err := getRecordDetailUrl(w.SObjectType, r)
	if err != nil {
		return err
	}
	return openBrowser(url)
}

func (w *ListView) OpenRecordEdit(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	r := w.Records[cy-1]
	url, err := getRecordEditUrl(w.SObjectType, r)
	if err != nil {
		return err
	}
	return openBrowser(url)
}

func getRecordDetailUrl(sobjectType string, r *soapforce.SObject) (string, error) {
	res, err := getDescribeSObjectResult(sobjectType)
	if err != nil {
		return "", nil
	}
	return strings.Replace(res.UrlDetail, "{ID}", r.Id, -1), nil
}

func getRecordEditUrl(sobjectType string, r *soapforce.SObject) (string, error) {
	res, err := getDescribeSObjectResult(sobjectType)
	if err != nil {
		return "", nil
	}
	return strings.Replace(res.UrlEdit, "{ID}", r.Id, -1), nil
}

func getDescribeSObjectResult(sobjectType string) (*soapforce.DescribeSObjectResult, error) {
	var sobj *soapforce.DescribeSObjectResult
	var err error
	if v, ok := descSObjectResults[sobjectType]; ok {
		sobj = v
	} else {
		sobj, err = client.DescribeSObject(sobjectType)
		if err != nil {
			return nil, err
		}
		descSObjectResults[sobjectType] = sobj
	}
	return sobj, nil
}

func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url,dll,FiileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

func (w *ListView) ShowRecord(g *gocui.Gui, v *gocui.View) error {
	recordView, err := g.SetView("Record", w.RecordView.x, w.RecordView.y, w.RecordView.x+w.RecordView.w, w.RecordView.y+w.RecordView.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		recordView.Title = "Record Detail"
		recordView.Highlight = true
		recordView.SelBgColor = gocui.ColorGreen
		recordView.SelFgColor = gocui.ColorBlack
		recordView.SetCursor(0, 1)

		_, cy := v.Cursor()
		recordView.Clear()
		r, err := getRecordDetail(w.SObjectType, w.Records[cy].Id)
		if err != nil {
			return err
		}
		sobj, err := getDescribeSObjectResult(w.SObjectType)
		max := maxFieldLabelLength(sobj.Fields)

		fmt.Fprintln(recordView, fmt.Sprintf("%s | %s", paddingLabel("ID", max), r.Id))
		if err != nil {
			return err
		}

		for _, f := range sobj.Fields {
			value, ok := r.Fields[f.Name]
			if !ok {
				value = ""
			}
			fmt.Fprintln(recordView, fmt.Sprintf("%s | %s", paddingLabel(f.Label, max), value))
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

func paddingLabel(label string, max int) string {
	l := max - runewidth.StringWidth(label)
	for l > 0 {
		label += " "
		l--
	}
	return label
}

func backToList(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybindings("Record")
	g.DeleteView("Record")
	g.SetCurrentView("ListView")
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

func newListView(x, y, w, h int, rv *RecordView) *ListView {
	return &ListView{x, y, w, h, nil, nil, "", rv}
}

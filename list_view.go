package main

import (
	"fmt"
	"os/exec"
	"runtime"

	"strings"

	"github.com/tzmfreedom/go-soapforce"
	"github.com/tzmfreedom/gocui"
)

type ListView struct {
	x, y, w, h     int
	rx, ry, rw, rh int
	View           *gocui.View
	Records        []*soapforce.SObject
	SObjectType    string
}

var colWidth = 20

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
		v.SetCursor(0, 2)
		w.Render("")

		if err := g.SetKeybinding("ListView", gocui.KeyArrowUp, gocui.ModNone, up(2)); err != nil {
			return err
		}
		if err := g.SetKeybinding("ListView", 'k', gocui.ModNone, up(2)); err != nil {
			return err
		}
		if err := g.SetKeybinding("ListView", gocui.KeyArrowDown, gocui.ModNone, w.downDynamically); err != nil {
			return err
		}
		if err := g.SetKeybinding("ListView", 'j', gocui.ModNone, w.downDynamically); err != nil {
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

	headers, err := getFields(soql)
	if err != nil {
		return err
	}
	displayHeaders := make([]string, len(headers))
	for i, h := range headers {
		displayHeaders[i] = display(h, colWidth)
	}
	fmt.Fprintln(w.View, strings.Join(displayHeaders, "|"))
	fmt.Fprintln(w.View, strings.Repeat("─", w.w-2))

	r, err := client.Query(soql)
	if err != nil {
		return err
	}
	w.Records = r.Records
	w.SObjectType, err = getSobjectFromSoql(soql)
	if err != nil {
		return err
	}

	if w.Records != nil {
		for _, r := range w.Records {
			values := make([]string, len(headers))
			for i, h := range headers {
				if h == "Id" {
					values[i] = display(r.Id, colWidth)
				} else {
					value := r.Fields[h]
					if value != nil {
						values[i] = display(value.(string), colWidth)
					} else {
						values[i] = display("", colWidth)
					}
				}
			}
			fmt.Fprintln(w.View, strings.Join(values, "|"))
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

func (w *ListView) downDynamically(g *gocui.Gui, v *gocui.View) error {
	maxY := len(w.Records)
	if v != nil {
		cx, cy := v.Cursor()
		ox, oy := v.Origin()
		if err := v.SetCursor(cx, cy+1); err != nil && (oy+cy) < maxY {
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
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
	_, cy := v.Cursor()
	rv := newRecordView(w.rx, w.ry, w.rw, w.rh, w.SObjectType, w.Records[cy])
	return rv.Render(g)
}

func newListView(x, y, w, h, rx, ry, rw, rh int) *ListView {
	return &ListView{x, y, w, h, rx, ry, rw, rh, nil, nil, ""}
}

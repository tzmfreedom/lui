package main

import (
	"fmt"
	"os"

	"github.com/tzmfreedom/go-soapforce"
	"github.com/tzmfreedom/gocui"
)

var client *soapforce.Client
var descGlobalResults []*soapforce.DescribeGlobalSObjectResult
var descSObjectResults = map[string]*soapforce.DescribeSObjectResult{}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	client = soapforce.NewClient()
	username := os.Getenv("SALESFORCE_USERNAME")
	password := os.Getenv("SALESFORCE_PASSWORD")
	result, err := client.Login(username, password)
	if err != nil {
		panic(err)
	}
	descResult, err := client.DescribeGlobal()
	if err != nil {
		panic(err)
	}
	descGlobalResults = descResult.Sobjects

	maxX, maxY := g.Size()
	m := newMenu(0, 0, 25, 7)
	uinfo := newUserInfo(maxX/2, 0, maxX/2-1, maxY/2-1, result.UserInfo)
	rv := &RecordView{maxX / 2, 0, maxX/2 - 1, maxY - 1}
	lv := newListView(0, maxY/2, maxX-1, maxY/2-1, rv)
	soql := newSoqlEditor(26, 0, maxX/2-26-1, 7, lv)
	ea := newExecuteAnonymous(0, 8, maxX/2-1, maxY/2-9)
	// d := newDebugView(0, maxY-3, maxX-1, 2)

	g.Mouse = true
	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen
	g.SetManager(m, uinfo, ea, soql, lv)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, quit); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, moveToNext); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModAlt, moveTo("SoqlEditor")); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlE, gocui.ModAlt, moveTo("ExecuteAnonymous")); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlL, gocui.ModAlt, moveTo("ListView")); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlM, gocui.ModAlt, moveTo("Menu")); err != nil {
		panic(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}

func up(minY int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if v != nil {
			_, cy := v.Cursor()
			fmt.Fprintf(v, string(cy))
			if cy > minY {
				err := v.SetCursor(0, cy-1)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func down(maxY int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if v != nil {
			_, cy := v.Cursor()
			if cy < maxY {
				err := v.SetCursor(0, cy+1)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func setCurrentView(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return err
	}
	g.Cursor = false
	return nil
}

func setCurrentViewForEditor(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return err
	}
	g.Cursor = true
	lines := v.ViewBufferLines()
	cx, cy := v.Cursor()
	if len(lines) > cy {
		line := lines[cy]
		if len(line) > cx {
			return nil
		}
		return v.SetCursor(len(line), cy)
	}
	var x, y int
	if len(lines) == 0 {
		x = 0
		y = 0
	} else {
		y = len(lines) - 1
		x = len(lines[y])
		if x < 0 {
			x = 0
		}
	}
	v.SetCursor(x, y)
	return nil
}

var menuOrder = []string{
	"Menu",
	"SoqlEditor",
	"ExecuteAnonymous",
	"ListView",
}

func moveToNext(g *gocui.Gui, v *gocui.View) error {
	current := g.CurrentView().Name()
	for i, menu := range menuOrder {
		if menu == current {
			var next string
			if i+1 == len(menuOrder) {
				next = menuOrder[0]
			} else {
				next = menuOrder[i+1]
			}
			g.SetCurrentView(next)
			return nil
		}
	}
	return nil
}

func moveTo(name string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		g.SetCurrentView(name)
		return nil
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

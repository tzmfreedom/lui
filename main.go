package main

import (
	"fmt"
	"os"

	"github.com/tzmfreedom/go-soapforce"
	"github.com/tzmfreedom/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	client := soapforce.NewClient()
	username := os.Getenv("SALESFORCE_USERNAME")
	password := os.Getenv("SALESFORCE_PASSWORD")
	result, err := client.Login(username, password)
	qr, err := client.Query("SELECT Id, LastName, FirstName FROM Contact")
	if err != nil {
		panic(err)
	}

	maxX, maxY := g.Size()
	m := newMenu(0, 0, 25, 7)
	uinfo := newUserInfo(maxX/2, 0, maxX/2-1, maxY/2-1, result.UserInfo)
	lv := newListView(0, maxY/2, maxX-1, maxY/2-1, qr.Records)
	soql := newSoqlEditor(26, 0, maxX/2-26-1, 7)
	ea := newExecuteAnonymous(0, 8, maxX/2-1, maxY/2-9)
	// d := newDebugView(0, maxY-3, maxX-1, 2)

	g.Mouse = true
	g.SetManager(m, uinfo, ea, soql, lv)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
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
	cx, cy := v.Cursor()
	lines := v.ViewBufferLines()
	if len(lines) > cy {
		line := lines[cy]
		if len(line) > cx {
			fmt.Fprintf(v, string(line[cy]))
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

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

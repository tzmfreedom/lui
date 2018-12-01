package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/Songmu/prompter"
	"github.com/tzmfreedom/go-soapforce"
	"github.com/tzmfreedom/gocui"
)

var client *soapforce.Client
var descGlobalResults []*soapforce.DescribeGlobalSObjectResult
var descSObjectResults = map[string]*soapforce.DescribeSObjectResult{}

var (
	Version string
)

type option struct {
	Username string
	Password string
	Endpoint string
}

func main() {
	option := parse()
	client = soapforce.NewClient()
	if option.Endpoint != "" {
		client.LoginUrl = option.Endpoint
	}
	result, err := client.Login(option.Username, option.Password)
	if err != nil {
		panic(err)
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if runtime.GOOS == "windows" {
		g.ASCII = true
	}
	if err != nil {
		panic(err)
	}
	defer g.Close()

	descResult, err := client.DescribeGlobal()
	if err != nil {
		panic(err)
	}
	descGlobalResults = descResult.Sobjects

	maxX, maxY := g.Size()
	// m := newMenu(0, 0, 25, 7)
	uinfo := newUserInfo(maxX/2, 0, maxX/2-1, 7, result.UserInfo)
	dv := newDescribeView(maxX/2, 8, maxX/2-1, maxY/2-9, descGlobalResults, maxX/2, maxY/2, maxX/2-1, maxY/2-1)
	lv := newListView(0, maxY/2, maxX-1, maxY/2, maxX/2, maxY/2, maxX/2-1, maxY/2)
	soql := newSoqlEditor(0, 0, maxX/2-1, 7, lv)
	ea := newExecuteAnonymous(0, 8, maxX/2-1, maxY/2-9)

	g.Mouse = true
	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen
	g.SetManager(soql, uinfo, ea, dv, lv)

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

func parse() *option {
	flg := flag.NewFlagSet("lui", flag.ExitOnError)
	username := flg.String("u", "", "username")
	endpoint := flg.String("e", "", "endpoint")
	version := flg.Bool("v", false, "version")
	flg.Usage = func() {
		fmt.Printf(`NAME:
   lui - Lightining plaform terminal UI

USAGE:
   lui [options]

VERSION:
   %s

OPTIONS:
   -u  username
   -e  endpoint (e.g. test.salesforce.com)
   -v  print the version
`, Version)
		os.Exit(0)
	}
	err := flg.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}
	if *version {
		fmt.Printf("%s\n")
		os.Exit(0)
		return nil
	}

	if *username == "" && os.Getenv("SALESFORCE_USERNAME") != "" {
		*username = os.Getenv("SALESFORCE_USERNAME")
	}
	if *username == "" {
		*username = prompter.Prompt("Enter username", "")
	}
	if *endpoint == "" && os.Getenv("SALESFORCE_ENDPOINT") != "" {
		*endpoint = os.Getenv("SALESFORCE_ENDPOINT")
	}
	if *endpoint == "" {
		*endpoint = prompter.Prompt("Enter endpoint (e.g. test.salesforce.com)", "login.salesforce.com")
	}

	password := os.Getenv("SALESFORCE_PASSWORD")
	if password == "" {
		password = prompter.Password("Enter password")
	}

	return &option{
		Username: *username,
		Password: password,
		Endpoint: *endpoint,
	}
}

func up(minY int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if v != nil {
			cx, cy := v.Cursor()
			ox, oy := v.Origin()
			if (cy + oy) > minY {
				if err := v.SetCursor(cx, cy-1); err != nil {
					if err := v.SetOrigin(ox, oy-1); err != nil {
						return err
					}
				}
			} else if (cy+oy) == 2 && oy != 0 {
				v.SetCursor(cx, cy+1)
				v.SetOrigin(ox, oy-1)
			}
		}
		return nil
	}
}

func down(maxY int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
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
}

func toTop(g *gocui.Gui, v *gocui.View) error {
	return nil
}

func toBottom(g *gocui.Gui, v *gocui.View) error {
	return nil
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
	"SoqlEditor",
	"ExecuteAnonymous",
	"Describe",
	"DescribeField",
	"ListView",
	"Record",
}

func moveToNext(g *gocui.Gui, v *gocui.View) error {
	current := g.CurrentView().Name()
	for i, name := range menuOrder {
		if name == current {
			for {
				i++
				if i == len(menuOrder) {
					i = 0
					break
				}
				nextMenu := menuOrder[i]
				if _, err := g.View(nextMenu); err == nil {
					break
				}
			}
			g.SetCurrentView(menuOrder[i])
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

package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/tzmfreedom/go-soapforce"
	"github.com/tzmfreedom/gocui"
)

var cx int = 0
var cy int = 0

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)
	g.Mouse = true

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		panic(err)
	}
	up := func (g *gocui.Gui, v *gocui.View) error {
		if v != nil {
			if cy != 0 {
				cy--
			}
			v.SetCursor(0, cy+1)
		}
		return nil
	}

	down := func (g *gocui.Gui, v *gocui.View) error {
		if v != nil {
			if cy <= maxRecordLength {
				cy++
				v.SetCursor(0, cy+1)
			}
		}
		return nil
	}

	setCurrentView := func (g *gocui.Gui, v *gocui.View) error {
		if _, err := g.SetCurrentView(v.Name()); err != nil {
			return err
		}
		return nil
	}

	if err := g.SetKeybinding("Menu", 'k', gocui.ModNone, up); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("Menu", 'j', gocui.ModNone, down); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("Records", gocui.KeyArrowUp, gocui.ModNone, up); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("Records", 'k', gocui.ModNone, up); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("Records", gocui.KeyArrowDown, gocui.ModNone, down); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("Records", 'j', gocui.ModNone, down); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("Menu", gocui.MouseLeft, gocui.ModNone, setCurrentView); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("Records", gocui.MouseLeft, gocui.ModNone, setCurrentView); err != nil {
		panic(err)
	}

	client := soapforce.NewClient()
	username := os.Getenv("SALESFORCE_USERNAME")
	password := os.Getenv("SALESFORCE_PASSWORD")
	result, err := client.Login(username, password)
	userInfo = result.UserInfo
	qr, err := client.Query("SELECT Id, LastName, FirstName FROM Contact")
	if err != nil {
		panic(err)
	}
	records = qr.Records

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}

var client *soapforce.Client
var records []*soapforce.SObject
var userInfo *soapforce.GetUserInfoResult
var maxRecordLength = 20

func layout(g *gocui.Gui) error {

	maxX, maxY := g.Size()
	if v, err := g.SetView("Menu", 0, 0, 20, 5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		menues := []string{
			"Execute Query",
			"Describe SObject",
			"Execute Anonymous",
		}
		for _, menu := range menues {
			fmt.Fprintln(v, menu)
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.SetCursor(0, 0)
	}
	if v, err := g.SetView("User Info", maxX/2+1, 0, maxX-1, 10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fields := map[string]string{
			"ID":       userInfo.UserId,
			"Name":     userInfo.UserName,
			"FullName": userInfo.UserFullName,
			"Email":    userInfo.UserEmail,
			"OrgID":    userInfo.OrganizationId,
			"OrgName":  userInfo.OrganizationName,
		}
		max := 0
		for k, _ := range fields {
			l := len(k)
			if l > max {
				max = l
			}
		}
		for k, val := range fields {
			fmt.Fprintln(v, fmt.Sprintf("%"+strconv.Itoa(max+2)+"s | %s", k, val))
		}
	}

	if v, err := g.SetView("Records", 0, maxY/2+1, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		headers := []string{
			"LastName",
		}
		for _, h := range headers {
			fmt.Fprintln(v, h)
		}
		for i, r := range records {
			for _, h := range headers {
				fmt.Fprintln(v, r.Fields[h])
			}
			if i > maxRecordLength {
				break
			}
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.SetCursor(0, 1)
	}
	g.SetCurrentView("Menu")
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

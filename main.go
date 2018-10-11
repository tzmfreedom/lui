package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/tzmfreedom/go-soapforce"
	"github.com/tzmfreedom/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	g.Cursor = true
	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		panic(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}

}

func layout(g *gocui.Gui) error {
	client := soapforce.NewClient()
	username := os.Getenv("SALESFORCE_USERNAME")
	password := os.Getenv("SALESFORCE_PASSWORD")
	u, err := client.Login(username, password)
	if err != nil {
		return err
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("User Info", 0, 0, maxX/2, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		info := u.UserInfo
		fields := map[string]string{
			"ID":       info.UserId,
			"Name":     info.UserName,
			"FullName": info.UserFullName,
			"Email":    info.UserEmail,
			"OrgID":    info.OrganizationId,
			"OrgName":  info.OrganizationName,
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

	qr, err := client.Query("SELECT Id, LastName, FirstName FROM Contact")
	if err != nil {
		return err
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
		for i, r := range qr.Records {
			v.Highlight = true
			v.SelBgColor = gocui.ColorGreen
			v.SelFgColor = gocui.ColorBlack
			v.SetCursor(0, 1)
			fmt.Fprintln(v, "aaa bbb ccc")
			for _, h := range headers {
				fmt.Fprintln(v, r.Fields[h])
			}
			if i > 5 {
				break
			}
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

package main

import (
	"fmt"
	"strconv"

	"github.com/tzmfreedom/go-soapforce"
	"github.com/tzmfreedom/gocui"
)

type UserInfo struct {
	x, y, w, h int
	Info       *soapforce.GetUserInfoResult
}

func (w *UserInfo) Layout(g *gocui.Gui) error {
	if v, err := g.SetView("UserInfo", w.x, w.y, w.x+w.w, w.y+w.h); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fields := [][]string{
			{"ID", w.Info.UserId},
			{"Name", w.Info.UserName},
			{"FullName", w.Info.UserFullName},
			{"Email", w.Info.UserEmail},
			{"OrgID", w.Info.OrganizationId},
			{"OrgName", w.Info.OrganizationName},
		}
		max := 0
		for _, tuple := range fields {
			l := len(tuple[0])
			if l > max {
				max = l
			}
		}
		for _, tuple := range fields {
			fmt.Fprintln(v, fmt.Sprintf("%"+strconv.Itoa(max+2)+"s | %s", tuple[0], tuple[1]))
		}
		v.Title = "User Info"
	}
	return nil
}

func newUserInfo(x, y, w, h int, uinfo *soapforce.GetUserInfoResult) *UserInfo {
	return &UserInfo{x, y, w, h, uinfo}
}

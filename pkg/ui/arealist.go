package ui

import (
	"fmt"
	"github.com/askovpen/goated/pkg/msgapi"
	"github.com/askovpen/gocui"
	"log"
	"strconv"
)

func quitEnter(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	if cy == 1 {
		ActiveWindow = "AreaList"
		g.DeleteView("QuitMsg")
	} else {
		return gocui.ErrQuit
	}
	return nil
}
func quitUp(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	v.SetCursor(cx, 1-cy)
	return nil
}
func quitAreaList(g *gocui.Gui, v *gocui.View) error {
	v, _ = g.SetView("QuitMsg", 2, 1, 17, 4)
	v.Title = "Quit goAtEd?"
	v.TitleFgColor = gocui.ColorYellow | gocui.AttrBold
	v.FrameFgColor = gocui.ColorRed | gocui.AttrBold
	v.FrameBgColor = gocui.ColorBlack
	fmt.Fprintf(v, "     Yes!     \n      No       ")
	v.Highlight = true
	v.SelBgColor = gocui.ColorBlue
	v.SelFgColor = gocui.ColorWhite | gocui.AttrBold
	ActiveWindow = "QuitMsg"
	g.SetCurrentView("QuitMsg")
	return nil
}

func getAreaNew(m msgapi.AreaPrimitive) string {
	if m.GetCount()-m.GetLast() > 0 {
		return "\033[37;1m+\033[0m"
	}
	return " "
}
func areaPgup(g *gocui.Gui, v *gocui.View) error {
	_, oy := v.Origin()
	_, maxY := v.Size()
	if oy-maxY+1 < 0 {
		if oy > 0 {
			v.SetOrigin(0, 0)
		} else {
			v.SetCursor(0, 1)
		}
	} else {
		v.SetOrigin(0, oy-maxY+1)
	}
	return nil
}
func areaPgdn(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	_, oy := v.Origin()
	_, maxY := v.Size()
	if len(msgapi.Areas)-oy < maxY-1 && cy != maxY {
		v.SetCursor(0, len(msgapi.Areas)-oy)
	} else if cy < maxY-1 {
		v.SetCursor(0, maxY-1)
	} else if len(msgapi.Areas)-oy > maxY-1 {
		v.SetOrigin(0, oy+maxY-1)
		if len(msgapi.Areas)-(oy+maxY-1) < maxY-1 {
			v.SetCursor(0, len(msgapi.Areas)-(oy+maxY-1))
		}
	}
	return nil
}
func areaNext(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if cy == len(msgapi.Areas) {
			return nil
		}
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if cy+oy == len(msgapi.Areas) {
				return nil
			}
			StatusLine = fmt.Sprintf(" %s: %d msgs, %d unread",
				msgapi.Areas[cy+oy].GetName(),
				msgapi.Areas[cy+oy].GetCount(),
				msgapi.Areas[cy+oy].GetCount()-msgapi.Areas[cy+oy].GetLast())
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		} else {
			StatusLine = fmt.Sprintf(" %s: %d msgs, %d unread",
				msgapi.Areas[cy].GetName(),
				msgapi.Areas[cy].GetCount(),
				msgapi.Areas[cy].GetCount()-msgapi.Areas[cy].GetLast())
		}
	}
	return nil
}

func areaPrev(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if cy > 1 {
			StatusLine = fmt.Sprintf(" %s: %d msgs, %d unread",
				msgapi.Areas[cy+oy-2].GetName(),
				msgapi.Areas[cy+oy-2].GetCount(),
				msgapi.Areas[cy+oy-2].GetCount()-msgapi.Areas[cy+oy-2].GetLast())
			if err := v.SetCursor(cx, cy-1); err != nil {
				log.Print(err)
				return err
			}
		} else if oy > 0 {
			StatusLine = fmt.Sprintf(" %s: %d msgs, %d unread",
				msgapi.Areas[cy+oy-2].GetName(),
				msgapi.Areas[cy+oy-2].GetCount(),
				msgapi.Areas[cy+oy-2].GetCount()-msgapi.Areas[cy+oy-2].GetLast())
			if err := v.SetOrigin(ox, oy-1); err != nil {
				log.Print(err)
				return err
			}
		}
	}
	return nil
}

func viewArea(g *gocui.Gui, v *gocui.View) error {
	_, oy := v.Origin()
	_, cy := v.Cursor()
	//log.Printf("view %d", oy+cy)
	err := viewMsg(cy+oy-1, msgapi.Areas[cy+oy-1].GetLast())
	if err != nil {
		errorMsg(err.Error(), "AreaList")
		return nil
	}
	g.SetCurrentView("MsgBody")
	ActiveWindow = "MsgBody"
	return nil
}

// CreateAreaList create arealist
func CreateAreaList(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView("AreaList", 0, 0, maxX-1, maxY-2)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	v.Wrap = false
	v.Highlight = true
	v.SelBgColor = gocui.ColorBlue
	v.SelFgColor = gocui.ColorWhite | gocui.AttrBold
	v.FgColor = gocui.ColorWhite
	v.FrameFgColor = gocui.ColorBlue | gocui.AttrBold
	v.FrameBgColor = gocui.ColorBlack
	v.Clear()
	fmt.Fprintf(v, "\033[33;1m Area %-"+strconv.FormatInt(int64(maxX-23), 10)+"s %6s %6s \033[0m\n",
		"EchoID", "Msgs", "New")
	for i, a := range msgapi.Areas {
		fmt.Fprintf(v, "%4d%s %-"+strconv.FormatInt(int64(maxX-23), 10)+"s %6d %6d \n",
			i+1,
			getAreaNew(a),
			a.GetName(),
			a.GetCount(),
			a.GetCount()-a.GetLast())
	}
	_, cy := v.Cursor()
	if cy == 0 {
		areaNext(App, v)
	}
	return nil
}

func answerMsgAreaList(g *gocui.Gui, v *gocui.View) error {
	newMsgType = newMsgTypeAnswer | newMsgTypeAnswerNewArea
	inlineAreaList(g, "answerMsgAreaList")
	return nil
}

func forwardAreaList(g *gocui.Gui, v *gocui.View) error {
	newMsgType = newMsgTypeForward
	inlineAreaList(g, "forwardMsg")
	return nil
}

func inlineAreaList(g *gocui.Gui, t string) error {
	maxX, maxY := g.Size()
	title := ""
	switch t {
	case "answerMsgAreaList":
		title = "Answer In Area:"
	case "forwardMsg":
		title = "Forward To Area:"
	}
	v, _ := g.SetView("iAreaList", 0, 5, maxX-1, maxY-1)
	v.Wrap = false
	v.Title = title
	v.Highlight = true
	v.SelBgColor = gocui.ColorBlue
	v.SelFgColor = gocui.ColorWhite | gocui.AttrBold
	v.FgColor = gocui.ColorWhite
	v.FrameFgColor = gocui.ColorBlue | gocui.AttrBold
	v.FrameBgColor = gocui.ColorBlack
	ActiveWindow = "iAreaList"
	fmt.Fprintf(v, "\033[33;1m Area %-"+strconv.FormatInt(int64(maxX-23), 10)+"s %6s %6s \033[0m\n",
		"EchoID", "Msgs", "New")
	for i, a := range msgapi.Areas {
		fmt.Fprintf(v, "%4d%s %-"+strconv.FormatInt(int64(maxX-23), 10)+"s %6d %6d \n",
			i+1,
			getAreaNew(a),
			a.GetName(),
			a.GetCount(),
			a.GetCount()-a.GetLast())
	}
	_, cy := v.Cursor()
	if cy == 0 {
		areaNext(g, v)
	}
	return nil
}

func answerMsgAreaListEscape(g *gocui.Gui, v *gocui.View) error {
	newMsgType = 0
	g.DeleteView("iAreaList")
	ActiveWindow = "MsgBody"
	return nil
}
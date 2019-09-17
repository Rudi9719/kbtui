package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"strconv"
	"time"
)

func main() {
	kbtui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer kbtui.Close()

	kbtui.SetManagerFunc(layout)

	go testAsync(kbtui)
	if err := kbtui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := kbtui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
func testAsync(kbtui *gocui.Gui) {
	for i := 0; i < 50; i++ {
		printToView(kbtui, "Chat", "Message #" + strconv.Itoa(i) + "\n")
		time.Sleep(1 * time.Second)
	}
	clearView(kbtui, "Chat")
}
func clearView(kbtui *gocui.Gui, viewName string) {
	kbtui.Update(func(g *gocui.Gui) error {
		inputView, err := kbtui.View(viewName)
		if err != nil {
			return err
		} else {
			inputView.Clear()
		}
		return nil
	})
}

func printToView(kbtui *gocui.Gui, viewName string, message string) {
	kbtui.Update(func(g *gocui.Gui) error {
		updatingView, err := kbtui.View(viewName)
		if err != nil {
			return err
		} else {
			_, _ = fmt.Fprintf(updatingView, message)
		}
		return nil
	})
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if feedView, err := g.SetView("Feed", 12, 0, maxX-1, maxY/5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		feedView.Autoscroll = true
		_, _ = fmt.Fprintln(feedView, "Feed Window")
	}
	if chatView, err2 := g.SetView("Chat", 12, maxY/5+1, maxX-1, maxY-5); err2 != nil {
		if err2 != gocui.ErrUnknownView {
			return err2
		}
		chatView.Autoscroll = true
		_, _ = fmt.Fprintln(chatView, "Chat Window")
	}
	if inputView, err3 := g.SetView("Input", 12, maxY-4, maxX-1, maxY-1); err3 != nil {
		if err3 != gocui.ErrUnknownView {
			return err3
		}
		inputView.Editable = true
		_, _ = fmt.Fprintln(inputView, "Input Window")
	}
	if listView, err4 := g.SetView("List", 0, 0, 10, maxY-1); err4 != nil {
		if err4 != gocui.ErrUnknownView {
			return err4
		}
		_, _ = fmt.Fprintf(listView, "Lists\nWindow")
	}
	return nil
}

func Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	cx, _ := v.Cursor()
	ox, _ := v.Origin()
	limit := ox+cx+1 > 255
	switch {
	case ch != 0 && mod == 0 && !limit:
		v.EditWrite(ch)
	case key == gocui.KeySpace && !limit:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	}
}
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

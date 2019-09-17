package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
)

func main() {
	kbtui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer kbtui.Close()

	kbtui.SetManagerFunc(layout)

	if err := kbtui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := kbtui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if feedView, err := g.SetView("Feed", 12, 0, maxX-1, maxY/5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(feedView, "Feed Window")
	}
	if chatView, err2 := g.SetView("Chat", 12, maxY/5+1, maxX-1, maxY-5); err2 != nil {
		if err2 != gocui.ErrUnknownView {
			return err2
		}
		fmt.Fprintln(chatView, "Chat Window")
	}
	if inputView, err3 := g.SetView("Input", 12, maxY-4, maxX-1, maxY-1); err3 != nil {
		if err3 != gocui.ErrUnknownView {
			return err3
		}
		fmt.Fprintln(inputView, "Input Window")
	}
	if listView, err4 := g.SetView("List", 0, 0, 10, maxY-1); err4 != nil {
		if err4 != gocui.ErrUnknownView {
			return err4
		}
		fmt.Fprintln(listView, "Lists")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

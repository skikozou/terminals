package main

import "github.com/nsf/termbox-go"

func drawText(x, y int, msg string, fg, bg termbox.Attribute) {
	for i, ch := range msg {
		termbox.SetCell(x+i, y, ch, fg, bg)
	}
}

func drawBox(x, y, width, height int, fg, bg termbox.Attribute) {
	for w := 0; w < width; w++ {
		for h := 0; h < height; h++ {
			termbox.SetCell(x+w, y+h, ' ', fg, bg)
		}
	}
}

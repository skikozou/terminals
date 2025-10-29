package main

import (
	"github.com/nsf/termbox-go"
)

type WindowHandle struct {
	X, Y, Width, Height int
	Title, TrimTitle    string
	Fg, Bg              termbox.Attribute
	isDrag              bool
	DragPosX, DragPosY  int
	Content             string
}

type Screen struct {
	Windows  []*WindowHandle
	Debug    bool
	Dragging bool
}

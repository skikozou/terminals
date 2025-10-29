package main

import (
	"strings"
)

func (h *WindowHandle) drawWindow(isfocus bool) {
	f, b := h.Fg, h.Bg
	if isfocus {
		//f, b = termbox.ColorBlack, termbox.ColorWhite
	}
	// draw box
	drawBox(h.X, h.Y+1, h.Width, h.Height, h.Fg, h.Bg)
	// write titlebar
	h.TrimTitle = h.Title
	if h.Width-3 < len(h.Title) { // 3 is {2 padding, x button}
		h.TrimTitle = h.Title[0:h.Width-(3+4)] + "..." // 3+4 is slice num and "..."
	}
	paddinglen := h.Width - (len(h.TrimTitle) + 3)
	padding := ""
	if paddinglen > 0 {
		padding = strings.Repeat(" ", paddinglen)
	}
	titlebar := h.TrimTitle + padding + " x "
	drawText(h.X, h.Y, titlebar, b, f)
	// write content
	h.drawContent()
}

func (h *WindowHandle) drawContent() {
	lines := strings.Split(h.Content, "\n")
	for i := 0; i < h.Height && i < len(lines); i++ {
		limit := min(len(lines[i]), h.Width)
		drawText(h.X, h.Y+i+1, lines[i][0:limit], h.Fg, h.Bg)
	}
}

func (h *WindowHandle) isTitlebarArea(posX, posY int) bool {
	return h.X <= posX && posX < h.X+h.Width-3 && h.Y == posY
}

func (h *WindowHandle) isExitButtonArea(posX, posY int) bool {
	return h.X+h.Width-3 <= posX && posX < h.X+h.Width && posY == h.Y
}

func (h *WindowHandle) onWindow(posX, posY int) bool {
	return h.X <= posX && posX < h.X+h.Width && h.Y <= posY && posY < h.Y+h.Height+1
}

func (s *Screen) closeWindow(index int) {
	ws := s.Windows
	copy(ws[index:], ws[index+1:])
	ws[len(ws)-1] = nil
	ws = ws[:len(ws)-1]
	s.Windows = ws
}

func (s *Screen) focusWindow(index int) {
	ws := s.Windows
	ws = append(ws, ws[index])
	ws[index] = nil
	ws = append(ws[:index], ws[index+1:]...)
	s.Windows = ws
}

func (s *Screen) getfocus() *WindowHandle {
	return s.Windows[0]
}

// 1 2 3 4 5 want 3
// 1 2 4 5 3

// 1 2 3 4 5
// 1 2 3 4 5 3
// 1 2 |3| 4 5 3
// 1 2 4 5 3

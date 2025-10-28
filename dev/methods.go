package main

import (
	"strings"

	"github.com/nsf/termbox-go"
)

func (h *WindowHandle) drawWindow(isfocus bool) {
	f, b := h.Fg, h.Bg
	if isfocus {
		f, b = termbox.ColorBlack, termbox.ColorWhite
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

<<<<<<< HEAD
func (h *WindowHandle) onWindow(posX, PosY int) bool {

}

func (s *Screen) closeWindow(window *WindowHandle) {
=======
func (h *WindowHandle) onWindow(posX, posY int) bool {
        return h.X <= posX && posX < h.X+h.Width && h.Y <= posY && posY < h.Y+h.Height+1
}

func (s *Screen) closeWindow(index int) {
>>>>>>> a21cbb1 (uh h)
	if len(s.Windows) != 1 {
		s.Focus = s.Windows[len(s.Windows)-2]
	} else {
		s.Focus = nil
	}
	s.Focusindex--
	ws := s.Windows
	copy(ws[index:], ws[index+1:])
	ws[len(ws)-1] = nil
	ws = ws[:len(ws)-1]
	s.Windows = ws
}

func (s *Screen) focusWindow(index int) {
	s.Focus = s.Windows[index]
	s.Focusindex = index
}

func (s *Screen) focus(index int) {
	
}

// 1 2 3 4 5 want 3
// 1 2 4 5 3

// 1 2 3 4 5 3
// 1 2 |3| 4 5 3
// 1 2 4 5 3

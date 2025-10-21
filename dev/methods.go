package main

import "strings"

func (h *WindowHandle) drawWindow() {
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
	drawText(h.X, h.Y, titlebar, h.Bg, h.Fg)
	// write content
	h.drawContent()
}

func (h *WindowHandle) drawContent() {
	lines := strings.Split(h.Content, "\n")
	for i := 0; i < h.Height && i < len(lines); i++ {
		limit := len(lines[i])
		if limit > h.Width {
			limit = h.Width
		}
		drawText(h.X, h.Y+i+1, lines[i][0:limit], h.Fg, h.Bg)
	}
}

func (h *WindowHandle) IsTitlebarArea(posX, posY int) bool {
	return h.X <= posX && posX < h.X+h.Width-3 && h.Y == posY
}

func (h *WindowHandle) IsExitButtonArea(posX, posY int) bool {
	return h.X+h.Width-3 <= posX && posX < h.X+h.Width && posY == h.Y
}

func (h *WindowHandle) closeWindow(Windows *[]*WindowHandle) {
	ws := *Windows
	for i, w := range ws {
		if w == h {
			copy(ws[i:], ws[i+1:])
			ws[len(ws)-1] = nil
			ws = ws[:len(ws)-1]
			break
		}
	}
	*Windows = ws
}

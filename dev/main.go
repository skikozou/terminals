package main

import (
	"fmt"
	"strings"

	"github.com/nsf/termbox-go"
	"github.com/sirupsen/logrus"
)

type WindowHandle struct {
	X, Y, Width, Height int
	Title, TrimTitle    string
	Fg, Bg              termbox.Attribute
	isDrag              bool
	DragPosX, DragPosY  int
	Content             string
	Closing             func(*[]*WindowHandle, *WindowHandle)
}

type AppState struct {
	Windows  []*WindowHandle
	Debug    bool
	Dragging bool
}

func InitLogrus() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:            true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})
}

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
	arrContent := strings.Split(h.Content, "\n")
	for i := 0; i < h.Height && i < len(arrContent); i++ {
		limit := len(arrContent[i])
		if limit > h.Width {
			limit = h.Width
		}
		drawText(h.X, h.Y+i+1, arrContent[i][0:limit], h.Fg, h.Bg)
	}
}

func (h *WindowHandle) IsTitlebarArea(posX, posY int) bool {
	return h.X <= posX && posX < h.X+h.Width-3 && h.Y == posY
}

func (h *WindowHandle) IsExitButtonArea(posX, posY int) bool {
	return h.X+h.Width-3 <= posX && posX < h.X+h.Width && posY == h.Y
}

func closeWindow(Windows *[]*WindowHandle, target *WindowHandle) {
	ws := *Windows
	for i, w := range ws {
		if w == target {
			copy(ws[i:], ws[i+1:])
			ws[len(ws)-1] = nil
			ws = ws[:len(ws)-1]
			break
		}
	}
	*Windows = ws
}

func main() {
	InitLogrus()
	err := termbox.Init()
	if err != nil {
		logrus.Fatal(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputAlt | termbox.InputMouse)

	state := &AppState{
		Windows:  make([]*WindowHandle, 0),
		Debug:    false,
		Dragging: false,
	}

	state.Windows = append(state.Windows, &WindowHandle{
		X:         0,
		Y:         0,
		Width:     50,
		Height:    7,
		Title:     "test window",
		TrimTitle: "trimed",
		Fg:        termbox.ColorRed,
		Bg:        termbox.ColorWhite,
		Content:   "",
		Closing:   closeWindow,
	})

	for {
		//window test
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		for _, w := range state.Windows {
			w.drawWindow()
		}

		if state.Debug {
			drawText(0, 0, fmt.Sprintf("Windows: %d ", len(state.Windows)), termbox.ColorLightGray, termbox.ColorDarkGray)
		}

		termbox.Flush()

		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc {
				return
			} else if ev.Ch == 'n' {
				newH := &WindowHandle{
					X:         0,
					Y:         0,
					Width:     50,
					Height:    7,
					Title:     "test window",
					TrimTitle: "trimed",
					Fg:        termbox.ColorRed,
					Bg:        termbox.ColorWhite,
					Content:   "Hello, Golang!\n;)",
					Closing:   closeWindow,
				}

				state.Windows = append(state.Windows, newH)
			} else if ev.Key == termbox.KeyTab {
				state.Debug = !state.Debug
			}

		case termbox.EventMouse:
			var toclose *WindowHandle
			for _, w := range state.Windows {
				switch ev.Key {
				case termbox.MouseLeft:
					if w.IsExitButtonArea(ev.MouseX, ev.MouseY) && !w.isDrag && !state.Dragging {
						toclose = w
					}

					if w.isDrag {
						w.X = ev.MouseX + w.DragPosX
						w.Y = ev.MouseY
					}

					if w.IsTitlebarArea(ev.MouseX, ev.MouseY) && !state.Dragging {
						w.isDrag = true
						state.Dragging = true
						w.DragPosX = w.X - ev.MouseX
						w.DragPosY = w.Y
					}
					termbox.Flush()
				case termbox.MouseRelease:
					w.isDrag = false
					state.Dragging = false
				}
			}
			if toclose != nil {
				closeWindow(&state.Windows, toclose)
			}
		}
	}
}

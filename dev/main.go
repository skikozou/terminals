package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
	"github.com/sirupsen/logrus"
)

func InitLogrus() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:            true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})
}

func main() {
	InitLogrus()
	/* debug area */
	// logrus.Info(termbox.RGBToAttribute(0, 0, 255))
	// logrus.Info(termbox.ColorBlue)
	// return
	// end area */
	err := termbox.Init()
	if err != nil {
		logrus.Fatal(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputAlt | termbox.InputMouse)
	termbox.SetOutputMode(termbox.OutputRGB | termbox.Output256)

	state := &Screen{
		Windows:    make([]*WindowHandle, 0),
		Focus:      nil,
		Focusindex: -1,
		Debug:      false,
		Dragging:   false,
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
	})

	state.Focus = state.Windows[0]
	state.Focusindex = 0

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
					Fg:        termbox.RGBToAttribute(0, 100, 0),
					Bg:        termbox.RGBToAttribute(100, 0, 0),
					Content:   "Hello, Golang!\n;)",
				}

				state.Focus = newH
				state.Focusindex = len(state.Windows)
				state.Windows = append(state.Windows, newH)
			} else if ev.Key == termbox.KeyTab {
				state.Debug = !state.Debug
			}

		case termbox.EventMouse:
			var toclose *WindowHandle
			for i, w := range state.Windows {
				switch ev.Key {
				case termbox.MouseLeft:
					if w.IsExitButtonArea(ev.MouseX, ev.MouseY) && !w.isDrag && !state.Dragging && i == state.Focusindex {
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
				toclose.closeWindow(&state.Windows)
			}
		}
	}
}

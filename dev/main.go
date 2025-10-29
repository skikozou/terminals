package main

import (
	"fmt"
	"math/rand"

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

func randomColor() int {
	return rand.Intn(16)
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
	//termbox.SetOutputMode(termbox.OutputRGB)

	state := &Screen{
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
		Fg:        termbox.ColorLightGray,
		Bg:        termbox.ColorLightGray,
		Content:   "",
	})

	for {
		//window test
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		for i, w := range state.Windows {
			w.drawWindow(i == len(state.Windows)-1)
		}

		if state.Debug {
			drawText(0, 0, fmt.Sprintf("Windows: %v", state.Windows), termbox.ColorWhite, termbox.ColorBlack)
		}

		termbox.Flush()

		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				return
			case termbox.KeyTab:
				state.Debug = !state.Debug
			}

			switch ev.Ch {
			case 'n':
				newH := &WindowHandle{
					X:         0,
					Y:         0,
					Width:     50,
					Height:    7,
					Title:     "test window",
					TrimTitle: "trimed",
					Fg:        termbox.Attribute(randomColor()),
					Bg:        termbox.Attribute(randomColor()),
					Content:   "Hello, Golang!\n;)",
				}

				state.Windows = append(state.Windows, newH)
			}

		case termbox.EventMouse:
			var (
				tofocus = -1
				toclose = -1
			)
			for i, w := range state.Windows {
				switch ev.Key {
				case termbox.MouseLeft:
					if w.onWindow(ev.MouseX, ev.MouseY) && !state.Dragging {
						tofocus = i
					}

					if w.isExitButtonArea(ev.MouseX, ev.MouseY) && !w.isDrag && !state.Dragging && i == len(state.Windows)-1 {
						toclose = i
					}

					if w.isDrag {
						w.X = ev.MouseX + w.DragPosX
						w.Y = ev.MouseY
					}

					if w.isTitlebarArea(ev.MouseX, ev.MouseY) && !state.Dragging {
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

			if tofocus != -1 {
				state.focusWindow(tofocus)
			}

			if toclose != -1 {
				state.closeWindow(toclose)
			}
		}
	}
}

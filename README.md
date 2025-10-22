# TUI window
### Manage screen
```go
type Screen struct {
	Windows  []*WindowHandle
	Debug    bool
	Dragging bool
}
```
### Manage window
```go
type WindowHandle struct {
	X, Y, Width, Height int
	Title, TrimTitle    string
	Fg, Bg              termbox.Attribute
	isDrag              bool
	DragPosX, DragPosY  int
	Content             string
}
```

# Task
- focus system
- arrow control
- min/normal size control
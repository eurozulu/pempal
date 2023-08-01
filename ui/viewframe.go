package ui

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"strings"
)

type ViewFrame interface {
	Offset() ViewOffset
	Position() *ViewOffset
	ResetPosition() ViewFrame
	Print(msg ...string)
	Println(msg ...string)
	Clear()
	ClearLine()
	WithColour(c ViewColours) ViewFrame
	WithRelativeOffset(x, y int) ViewFrame
}

type viewFrame struct {
	offset   ViewOffset
	position *ViewOffset
	colours  ViewColours
}

func (f viewFrame) Offset() ViewOffset {
	return f.offset
}

func (f viewFrame) Position() *ViewOffset {
	return f.position
}

func (f viewFrame) Print(msg ...string) {
	offset := f.Position()
	s := strings.Join(msg, "")
	for _, c := range s {
		termbox.SetCell(offset.X, offset.Y, c, f.colours.Foreground, f.colours.Background)
		offset.X += runewidth.RuneWidth(c)
	}
}

func (f viewFrame) Println(msg ...string) {
	f.Print(strings.Join(msg, ""))
	f.position.X = f.offset.X
	f.position.Y++
}

func (f *viewFrame) Clear() {
	termbox.Clear(DefaultColours.Background, DefaultColours.Background)
	f.ResetPosition()
}

func (f *viewFrame) ClearLine() {
	w, _ := termbox.Size()
	f.position.X = f.offset.X
	f.Print(strings.Repeat(" ", w))
	f.position.X = f.offset.X
}

func (f viewFrame) WithColour(c ViewColours) ViewFrame {
	return &viewFrame{
		offset:   f.offset,
		position: f.position,
		colours:  f.colours.MergeColours(c),
	}
}

func (f viewFrame) WithRelativeOffset(x, y int) ViewFrame {
	fr := &viewFrame{
		offset: ViewOffset{
			X: f.position.X + x,
			Y: f.position.Y + y,
		},
		colours: f.colours,
	}
	fr.ResetPosition()
	return fr
}

func (f *viewFrame) ResetPosition() ViewFrame {
	f.position = &ViewOffset{
		X: f.offset.X,
		Y: f.offset.Y,
	}
	return f
}

func newViewFrame(offset ViewOffset) ViewFrame {
	return &viewFrame{
		offset: offset,
		position: &ViewOffset{
			X: offset.X,
			Y: offset.Y,
		},
		colours: textColours,
	}
}

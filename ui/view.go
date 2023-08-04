package ui

import (
	"fmt"
)

var ErrAborted = fmt.Errorf("aborted")

type View interface {
	Render(frame ViewFrame)
	Label() string
	fmt.Stringer
	Colours() ViewColours
}

type HiddenView interface {
	View
	IsHidden() bool
	SetHidden(hidden bool)
}

type view struct {
	label, text string
	colours     ViewColours
	hidden      bool
}

func (v view) Colours() ViewColours {
	return v.colours
}

func (v view) Render(frame ViewFrame) {
	if v.IsHidden() {
		return
	}
	v.renderLabel(frame)
	v.renderText(frame)
}

func (v view) renderLabel(frame ViewFrame) {
	if v.label != "" {
		f := frame.WithColour(v.colours)
		f.Print(v.label, ": ")
	}
}

func (v view) renderText(frame ViewFrame) {
	if v.text != "" {
		f := frame.WithColour(v.colours)
		f.Print(v.text)
	}
}

func (v view) Label() string {
	return v.label
}

func (v view) String() string {
	return v.text
}

func (v view) IsHidden() bool {
	return v.hidden && v.text == ""
}

func (v *view) SetHidden(hidden bool) {
	v.hidden = hidden
}

func NewLabelView(label, text string) *view {
	return &view{
		label:   label,
		text:    text,
		colours: DefaultColours,
	}
}

package ui

import (
	"github.com/nsf/termbox-go"
	"unicode"
)

// TextView extends the View interface to make the Views Text value mutable.
type TextView interface {
	View
	AppendText(ch rune)
	SetText(text string)
	GetText() string
	SetColours(colours ViewColours)
}

type textView struct {
	view
}

func (tv textView) String() string {
	return truncate(tv.text, 25)
}

func (tv *textView) GetText() string {
	return tv.text
}

func (tv *textView) AppendText(ch rune) {
	if ch == rune(termbox.KeyBackspace) || ch == rune(termbox.KeyBackspace2) {
		if tv.text != "" {
			tv.text = tv.text[:len(tv.text)-1]
		}
		return
	} else if unicode.IsPrint(ch) && !unicode.IsSymbol(ch) {
		tv.text = string(append([]rune(tv.text), ch))
	}
}

func (tv *textView) SetText(text string) {
	tv.text = text
}

func (tv *textView) SetColours(colours ViewColours) {
	tv.colours = tv.colours.MergeColours(colours)
}

func truncate(s string, width int) string {
	if len(s) <= width {
		return s
	}
	return s[:width-3] + "..."
}

func NewTextView(label, text string) TextView {
	return &textView{view: *NewLabelView(label, text)}
}

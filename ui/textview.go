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
	SetColours(colours ViewColours)
}

type textView struct {
	view
}

func (m *textView) AppendText(ch rune) {
	if ch == rune(termbox.KeyBackspace) || ch == rune(termbox.KeyBackspace2) {
		if m.text != "" {
			m.text = m.text[:len(m.text)-1]
		}
		return
	} else if unicode.IsPrint(ch) {
		m.text = string(append([]rune(m.text), ch))
	}
}

func (m *textView) SetText(text string) {
	m.text = text
}

func (m *textView) SetColours(colours ViewColours) {
	m.colours = m.colours.MergeColours(colours)
}

func NewTextView(label, text string) TextView {
	return &textView{view: *NewLabelView(label, text)}
}

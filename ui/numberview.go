package ui

import (
	"github.com/nsf/termbox-go"
	"strconv"
	"unicode"
)

type NumberView interface {
	View
	AppendText(ch rune)
	SetText(text string)
}

type numberView struct {
	textView
}

func (m *numberView) AppendText(ch rune) {
	if ch == rune(termbox.KeyBackspace) || ch == rune(termbox.KeyBackspace2) {
		if m.text != "" {
			m.text = m.text[:len(m.text)-1]
		}
		return
	} else if unicode.IsNumber(ch) {
		m.text = string(append([]rune(m.text), ch))
	}
}

func (m *numberView) SetText(text string) {
	if text != "" {
		_, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			text = text + err.Error()
		}
	}
	m.text = text
}

func NewNumberView(label string, value int64) NumberView {
	tv := NewTextView(label, strconv.FormatInt(value, 10)).(*textView)
	return &numberView{*tv}
}

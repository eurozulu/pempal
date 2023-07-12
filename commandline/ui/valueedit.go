package ui

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"strings"
	"unicode"
)

var ERRAborted = fmt.Errorf("aborted")

type ValueEdit struct {
	Name      string
	ValueType InputType
	Options   []string
}

func (ed *ValueEdit) Edit(offset ViewOffset, value string) (string, error) {
	isRoot := !termbox.IsInit
	if isRoot {
		if err := termbox.Init(); err != nil {
			return "", err
		}
		defer termbox.Close()
	}

	selected := ed.optionIndex(value)
	list := NewItemListOfValues(ed.Options)

	for {
		if isRoot {
			termbox.Clear(ColourBackground.ToAttribute(), ColourBackground.ToAttribute())
		}
		os := offset
		ed.renderValue(&os, value)
		os.YOffset++
		os.XOffset -= 25
		list.renderList(os, selected)

		ev, err := nextKeyEvent()
		if err != nil {
			return "", err
		}
		switch ev.Key {
		case termbox.KeyArrowUp: // Previous line
			if selected > 0 {
				selected--
				value = list[selected].Value
			}
		case termbox.KeyArrowDown: // next line
			if selected < len(list)-1 {
				selected++
				value = list[selected].Value
			}
		case termbox.KeyEnter: // Selected last item
			return value, nil

		case termbox.KeyEsc: // aborted
			return "", ERRAborted

		default:
			value = ed.handleKeyInput(*ev, value)
		}
	}
}

func (ed *ValueEdit) CanEdit() bool {
	return ed.ValueType != InputTypeNone
}

func (ed *ValueEdit) renderValue(offset *ViewOffset, value string) {
	fg := ColourForeground.ToAttribute()
	bg := ColourBackgroundEdit.ToAttribute()
	tbprint(offset, fg, bg, ed.Name)
	tbprint(offset, fg, bg, ": ")
	tbprint(offset, fg, bg, padValue(value, 25))
}

func (ed *ValueEdit) handleKeyInput(event termbox.Event, value string) string {
	switch event.Key {
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		if len(value) > 0 {
			return value[:len(value)-1]
		}
	case termbox.KeySpace:
		return ed.handleCharInput(' ', value)
	default:
		if event.Ch > 0 {
			return ed.handleCharInput(event.Ch, value)
		}
	}
	return value
}

func (ed *ValueEdit) handleCharInput(ch rune, value string) string {
	switch ed.ValueType {
	case InputTypeNumbers:
		if !unicode.IsNumber(ch) && !unicode.IsDigit(ch) {
			return value
		}
	case InputTypeLetters:
		if !unicode.IsLetter(ch) {
			return value
		}
	case InputTypePrintable:
		if !unicode.IsPrint(ch) && !unicode.IsSpace(ch) && !unicode.IsPunct(ch) {
			return value
		}
	default:
		// InputTypeNone
		return value
	}
	return string(append([]rune(value), ch))
}

func (ed *ValueEdit) optionIndex(s string) int {
	for i, v := range ed.Options {
		if s == v {
			return i
		}
	}
	return -1
}

func padValue(s string, width int) string {
	w := width - len(s)
	if w <= 0 {
		return s
	}
	return strings.Join([]string{s, strings.Repeat(" ", w)}, "")
}

func nextKeyEvent() (*termbox.Event, error) {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			return &ev, nil
		case termbox.EventError:
			return nil, ev.Err
		default:
			continue
		}
	}
}

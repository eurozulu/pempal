package ui

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"strings"
	"unicode"
)

var ERRAborted = fmt.Errorf("aborted")

const valuePadWidth = 20

const (
	// InputTypeNone only the Options of the field may be selected
	InputTypeNone InputType = iota
	// InputTypeNumbers only the options and any numeric value may be entered
	InputTypeNumbers
	// InputTypeLetters only the options and any string of letters may be entered
	InputTypeLetters
	// InputTypePrintable only the options and any printable string of characters may be entered
	InputTypePrintable
)

type InputType int

type EditItem struct {
	Name      string
	ValueType InputType
	Options   []string
}

func (ed *EditItem) Edit(offset ViewOffset, value string) (string, error) {
	selected := ed.optionIndex(value)
	list := newItemListOfValues(ed.Options)

	for {
		os := offset
		ed.renderValue(&os, value)
		os.YOffset++
		os.XOffset -= 25
		list.render(os, selected)
		if err := termbox.Flush(); err != nil {
			return "", err
		}

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
			if ed.ValueType != InputTypeNone {
				value = ed.handleKeyInput(*ev, value)
			}
		}
	}
}

func (ed *EditItem) CanEdit() bool {
	return ed.ValueType != InputTypeNone
}

func (ed *EditItem) renderValue(offset *ViewOffset, value string) {
	fg := ColourForeground.toAttribute()
	bg := ColourBackgroundEdit.toAttribute()
	tbprint(offset, fg, bg, ed.Name)
	tbprint(offset, fg, bg, ": ")
	tbprint(offset, fg, bg, padValue(value, valuePadWidth))
}

func (ed *EditItem) handleKeyInput(event termbox.Event, value string) string {
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

func (ed *EditItem) handleCharInput(ch rune, value string) string {
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

func (ed *EditItem) optionIndex(s string) int {
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

func newItemListOfValues(values []string) ItemList {
	items := make([]ListItem, len(values))
	for i, v := range values {
		items[i] = ListItem{Value: v}
	}
	return items
}

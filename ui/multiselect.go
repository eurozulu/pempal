package ui

import (
	"github.com/nsf/termbox-go"
	"strings"
)

type multiSelect struct {
	TextList
	selected []string
}

func (ms *multiSelect) OnChildUpdate(child View) {
	ms.SetText(ms.getSelectedValues())
}

func (ms multiSelect) renderChild(frame ViewFrame, child View, selected bool) {
	if selected || ms.isTextSelected(child.String()) {
		frame = frame.WithColour(selectedColour)
	}
	frame.Print(padText(child.Label(), 20))
	frame.Print(truncateOrPadValue(child.String(), 25))
}

func (ms *multiSelect) AppendText(ch rune) {
	switch ch {
	case rune(termbox.KeyArrowUp):
		ms.setSelectedIndex(-1)
	case rune(termbox.KeyArrowDown):
		ms.setSelectedIndex(1)
	case rune(termbox.KeySpace):
		ms.toggleSelected()
	default:
		if ms.allowInput {
			ms.textView.AppendText(ch)
		}
	}
}

func (ms multiSelect) getSelectedValues() string {
	// ensure current selection is in selectedlist
	if ms.text != "" && ms.selectedIndex(ms.text) < 0 {
		ms.selected = append(ms.selected, ms.text)
	}
	return strings.Join(ms.selected, ",")
}

func (ms *multiSelect) toggleSelected() {
	if ms.text == "" {
		// nothing selected
		return
	}
	i := ms.selectedIndex(ms.text)
	if i < 0 {
		// not already selected, add to list
		ms.selected = append(ms.selected, ms.text)
	} else {
		// already selected, remove from list
		ss := ms.selected[:i]
		if i+1 < len(ms.selected) {
			ss = append(ss, ms.selected[i+1:]...)
		}
		ms.selected = ss
	}
}

func (ms multiSelect) isTextSelected(s string) bool {
	return ms.selectedIndex(s) >= 0
}

func (ms multiSelect) selectedIndex(s string) int {
	for i, sz := range ms.selected {
		if sz == s {
			return i
		}
	}
	return -1
}

func NewMultiSelectHidden(label, text string, choices ...string) *multiSelect {
	ms := NewMultiSelect(label, text, choices...)
	ms.hidden = true
	return ms
}

func NewMultiSelect(label, text string, choices ...string) *multiSelect {
	tl := NewTextList(label, text, choices...)
	return &multiSelect{
		TextList: *tl,
	}
}

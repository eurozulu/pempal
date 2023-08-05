package ui

import (
	"github.com/nsf/termbox-go"
	"strings"
)

type ParentView interface {
	TextView
	ChildViews() []View
	SelectedIndex() int
	ChildByLabel(label string) View
}

type MutableParentView interface {
	SetSelectedIndex(index int)
	SetChildViews(children []View)
}

type parentView struct {
	textView
	children      []View
	allowInput    bool
	selectedindex int
}

func (pv parentView) ChildViews() []View {
	return pv.children
}

func (pv *parentView) SetChildViews(children []View) {
	pv.children = children
}

func (pv parentView) SelectedIndex() int {
	return pv.selectedindex
}

func (pv parentView) ChildByLabel(label string) View {
	for _, c := range pv.children {
		if c.Label() == label {
			return c
		}
	}
	return nil
}

func (pv parentView) Render(frame ViewFrame) {
	if pv.textView.IsHidden() {
		return
	}
	pv.textView.renderLabel(frame)
	pv.textView.renderText(frame.WithColour(allowEditColour))
	if len(pv.children) > 0 {
		frame.Println()
		pv.renderChildren(frame)
	}
}

func (pv *parentView) AppendText(ch rune) {
	switch ch {
	case rune(termbox.KeyArrowUp):
		pv.SetSelectedIndex(pv.selectedindex - 1)
	case rune(termbox.KeyArrowDown):
		pv.SetSelectedIndex(pv.selectedindex + 1)
	default:
		if pv.allowInput {
			pv.textView.AppendText(ch)
			pv.SetSelectedIndexByText(pv.text)
		}
	}
}

func (pv *parentView) SetText(text string) {
	pv.text = text
}

func (pv parentView) renderChildren(frame ViewFrame) {
	childFrame := frame.WithRelativeOffset(len(pv.label)+2, 0)
	for i, v := range pv.children {
		childFrame.ClearLine()
		pv.renderChild(childFrame, v, pv.selectedindex == i)
		if childFrame.Position().X > childFrame.Offset().X {
			// only move to next line if child rendered something.
			childFrame.Println()
		}
	}
}

func (pv parentView) renderChild(frame ViewFrame, child View, selected bool) {
	if selected {
		frame = frame.WithColour(selectedColour)
	}
	if hv, ok := child.(HiddenView); ok && hv.IsHidden() {
		return
	}
	if child.Label() != "" {
		frame.Print(padText(child.Label(), 20))
	}
	frame = frame.WithColour(child.Colours())
	frame.Print(padText(child.String(), 25))
}

func (pv *parentView) SetSelectedIndex(index int) {
	increment := 1
	if index < pv.selectedindex {
		increment = -1
	}
	for {
		if index < 0 || index >= len(pv.children) {
			return
		}
		if hv, ok := pv.children[index].(HiddenView); ok && hv.IsHidden() {
			// if hidden view, move onto the next one
			index += increment
			continue
		}
		break
	}
	pv.selectedindex = index
	pv.setTextWithSelectedChild()
}

func (pv *parentView) setTextWithSelectedChild() {
	s := pv.children[pv.selectedindex].String()
	if tv, ok := pv.children[pv.selectedindex].(TextView); ok {
		s = tv.GetText()
	}
	pv.textView.SetText(s)
}

func (pv *parentView) SetSelectedIndexByText(text string) {
	index := -1
	for i, c := range pv.children {
		s := c.String()
		if vt, ok := c.(TextView); ok {
			s = vt.GetText()
		}
		if s == text {
			index = i
			break
		}
	}
	pv.selectedindex = index
}

func padText(s string, width int) string {
	l := len(s)
	i := width - l
	if i <= 0 {
		return s
	}
	return strings.Join([]string{s, strings.Repeat(" ", i)}, "")
}

func NewParentView(label, text string, children ...View) ParentView {
	tv := NewTextView(label, text).(*textView)
	return &parentView{
		textView: *tv,
		children: children,
	}
}

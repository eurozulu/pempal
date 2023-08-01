package ui

import (
	"github.com/nsf/termbox-go"
	"strings"
)

type ParentView interface {
	TextView
	ChildViews() []View
	ChildByLabel(label string) View
	SelectedIndex() int
}

type WindowNotifer interface {
	OnChildUpdate(child View)
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

func (pv parentView) Render(frame ViewFrame) {
	pv.textView.Render(frame)
	if !pv.textView.IsHidden() && len(pv.children) > 0 {
		frame.Println()
		pv.renderChildren(frame)
	}
}

func (pv *parentView) AppendText(ch rune) {
	switch ch {
	case rune(termbox.KeyArrowUp):
		pv.setSelectedIndex(-1)
	case rune(termbox.KeyArrowDown):
		pv.setSelectedIndex(1)
	default:
		if pv.allowInput {
			pv.textView.AppendText(ch)
		}
	}
}

func (pv *parentView) SetText(text string) {
	pv.text = text
	pv.setSelectedByText(pv.text)
}

func (pv parentView) ChildByLabel(label string) View {
	for _, c := range pv.children {
		if c.Label() == label {
			return c
		}
	}
	return nil
}

func (pv parentView) SelectedIndex() int {
	return pv.selectedindex
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
	frame.Print(truncateOrPadValue(child.String(), 25))
}

func (pv *parentView) setSelectedByText(text string) {
	index := -1
	for i, c := range pv.children {
		if c.String() == text {
			index = i
			break
		}
	}
	pv.selectedindex = index
}

func (pv *parentView) setSelectedIndex(relativeIndex int) {
	i := pv.selectedindex
	for {
		i += relativeIndex
		if i < 0 || i >= len(pv.children) {
			return
		}
		if hv, ok := pv.children[i].(HiddenView); ok && hv.IsHidden() {
			continue
		}
		pv.selectedindex = i
		pv.textView.SetText(pv.children[i].String())
		break
	}
}

func truncateOrPadValue(s string, width int) string {
	if len(s) <= width {
		return padText(s, width)
	}
	return s[width-3:] + "..."
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

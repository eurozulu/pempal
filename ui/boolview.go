package ui

import (
	"github.com/nsf/termbox-go"
	"strconv"
)

type BoolView struct {
	parentView
}

func (bv BoolView) Render(frame ViewFrame) {
	bv.textView.Render(frame)
	bv.renderChildrenHorizontally(frame)
}

func (bv BoolView) renderChildrenHorizontally(frame ViewFrame) {
	childFrame := frame.WithRelativeOffset(len(bv.label)+2, 1)
	selected := bv.SelectedIndex()

	for i, v := range bv.children {
		bv.renderChild(childFrame, v, selected == i)
		childFrame.Print("    ")
	}
}

func (bv BoolView) renderChild(frame ViewFrame, child View, selected bool) {
	if selected {
		frame = frame.WithColour(selectedColour)
	}
	frame.Print(padText(child.Label(), 5))
	frame.Print(truncateOrPadValue(child.String(), 5))
}

func (bv *BoolView) AppendText(ch rune) {
	switch ch {
	case rune(termbox.KeyArrowRight), 't', 'T', 'y', 'Y':
		bv.setValue(true)
	case rune(termbox.KeyArrowLeft), 'f', 'F', 'n', 'N':
		bv.setValue(false)
	default:
		// ignore
	}
}

func (bv *BoolView) SetText(text string) {
	bv.text = text
}

func (bv *BoolView) setValue(b bool) {
	bv.SetText(strconv.FormatBool(b))
}

func buildChildViews(labels []string) []View {
	f := strconv.FormatBool(false)
	t := strconv.FormatBool(true)
	fl := f
	tl := t
	if len(labels) > 0 && labels[0] != "" {
		fl = labels[0]
	}
	if len(labels) > 1 && labels[1] != "" {
		tl = labels[1]
	}
	return []View{
		NewLabelView(fl, f),
		NewLabelView(tl, t),
	}
}

func NewBoolView(label string, labels ...string) View {
	pv := NewParentView(label, "", buildChildViews(labels)...).(*parentView)
	return &BoolView{*pv}
}
func NewBoolViewPreSelected(label string, value bool, labels ...string) View {
	bv := NewBoolView(label, labels...).(*BoolView)
	bv.SetText(strconv.FormatBool(value))
	return bv
}

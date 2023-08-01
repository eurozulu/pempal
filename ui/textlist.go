package ui

// TextList is a public interface to the ParentView implementations, to allow external libs to create their own Views.
type TextList struct {
	parentView
}

func NewTextListView(label, text string, children ...View) *TextList {
	pv := NewParentView(label, text, children...).(*parentView)
	return &TextList{*pv}
}

func NewTextList(label, text string, choices ...string) *TextList {
	children := make([]View, len(choices))
	for i, ch := range choices {
		children[i] = NewLabelView("", ch)
	}
	return NewTextListView(label, text, children...)
}

func NewTextListHidden(label, text string, choices ...string) *TextList {
	tl := NewTextList(label, text, choices...)
	tl.hidden = true
	return tl
}

package ui

// ListView is a public implementation of the ParentView implementations, to allow external libs to create their own Views.
type ListView struct {
	parentView
}

func (l *ListView) OnViewOpen() {
	l.SetSelectedIndexByText(l.text)
}

func NewListView(label, text string, children ...View) *ListView {
	pv := NewParentView(label, text, children...).(*parentView)
	return &ListView{*pv}
}

func NewListViewStrings(label, text string, choices ...string) *ListView {
	children := make([]View, len(choices))
	for i, ch := range choices {
		children[i] = NewLabelView("", ch)
	}
	return NewListView(label, text, children...)
}

func NewListViewHidden(label, text string, choices ...string) *ListView {
	tl := NewListViewStrings(label, text, choices...)
	tl.hidden = true
	return tl
}

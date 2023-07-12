package ui

const maxValueWidth = 30

// ListItem represents a single Key/Value pair for display.
// It may optionally specify alternative colours to display, otherwise it defaults to ColourBackground and ColourForeground
type ListItem struct {
	Name             string
	Value            string
	ForegrondColour  ItemColour
	BackgroundColour ItemColour
}

type ItemList []ListItem

func (li ListItem) render(offset ViewOffset, selected bool) {
	bg := ColourBackground.toAttribute()
	if li.BackgroundColour != 0 {
		bg = li.BackgroundColour.toAttribute()
	}
	cleanLine(offset.YOffset, bg)

	if selected {
		bg = ColourBackgroundSelected.toAttribute()
	}

	fg := ColourForeground.toAttribute()
	if li.ForegrondColour != 0 {
		fg = li.ForegrondColour.toAttribute()
	}
	if li.Name != "" {
		tbprint(&offset, fg, bg, li.Name)
		tbprint(&offset, fg, bg, ":  ")
	}
	v := li.Value
	if len(v) > maxValueWidth {
		v = v[:maxValueWidth-3] + "..."
	}
	tbprint(&offset, fg, bg, v)
}

func (sl ItemList) render(offset ViewOffset, selected int) {
	for i, li := range sl {
		li.render(offset, i == selected)
		offset.YOffset++
	}
}

func NewItemListOfValues(values []string) ItemList {
	items := make([]ListItem, len(values))
	for i, v := range values {
		items[i] = ListItem{Value: v}
	}
	return items
}

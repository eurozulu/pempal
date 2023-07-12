package ui

const maxValueWidth = 30

type ListItem struct {
	Name             string
	Value            string
	ForegrondColour  ItemColor
	BackgroundColour ItemColor
}

func (li ListItem) render(offset ViewOffset, selected bool) {
	bg := ColourBackground.ToAttribute()
	if li.BackgroundColour != 0 {
		bg = li.BackgroundColour.ToAttribute()
	}
	if selected {
		bg = ColourBackgroundSelected.ToAttribute()
	}
	cleanLine(offset.YOffset, bg)

	fg := ColourForeground.ToAttribute()
	if li.ForegrondColour != 0 {
		fg = li.ForegrondColour.ToAttribute()
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

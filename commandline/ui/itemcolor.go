package ui

import "github.com/nsf/termbox-go"

const (
	ColourBackground         = ItemColor(termbox.ColorDefault)
	ColourBackgroundSelected = ItemColor(termbox.ColorLightCyan)
	ColourBackgroundEdit     = ItemColor(termbox.ColorLightGreen)
	ColourForeground         = ItemColor(termbox.ColorBlack)
	ColourForegroundOptional = ItemColor(termbox.ColorLightGray)
	ColourForegroundError    = ItemColor(termbox.ColorRed)
)

type ItemColor uint64

func (c ItemColor) ToAttribute() termbox.Attribute {
	return termbox.Attribute(c)
}

package ui

import "github.com/nsf/termbox-go"

const (
	ColourBackground         = ItemColour(termbox.ColorDefault)
	ColourBackgroundSelected = ItemColour(termbox.ColorLightCyan)
	ColourBackgroundEdit     = ItemColour(termbox.ColorLightGreen)
	ColourForeground         = ItemColour(termbox.ColorBlack)
	ColourForegroundOptional = ItemColour(termbox.ColorLightGray)
	ColourForegroundError    = ItemColour(termbox.ColorRed)
)

type ItemColour uint64

func (c ItemColour) toAttribute() termbox.Attribute {
	return termbox.Attribute(c)
}

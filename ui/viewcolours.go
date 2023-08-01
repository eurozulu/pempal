package ui

import "github.com/nsf/termbox-go"

var DefaultColours = ViewColours{
	Foreground: termbox.ColorDefault,
	Background: termbox.ColorDefault,
}
var titleColours = ViewColours{
	Foreground: termbox.ColorBlue,
	Background: termbox.ColorDefault,
}
var textColours = ViewColours{
	Foreground: termbox.ColorBlack,
	Background: termbox.ColorDefault,
}
var choiceColours = ViewColours{
	Background: termbox.ColorLightBlue,
}
var selectedColour = ViewColours{
	Background: termbox.ColorLightMagenta,
}
var ErrorColour = ViewColours{
	Foreground: termbox.ColorRed,
}
var allowEditColour = ViewColours{
	Background: termbox.ColorLightGreen,
}

type ViewColours struct {
	Foreground termbox.Attribute
	Background termbox.Attribute
}

func (vc ViewColours) MergeColours(colours ...ViewColours) ViewColours {
	col := ViewColours{
		Foreground: vc.Foreground,
		Background: vc.Background,
	}
	for _, c := range colours {
		if c.Foreground != 0 {
			col.Foreground = c.Foreground
		}
		if c.Background != 0 {
			col.Background = c.Background
		}
	}
	return col
}

package formselect

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

const (
	colBackground         = termbox.ColorDefault
	colBackgroundSelected = termbox.ColorLightBlue
	colNormal             = termbox.ColorBlack
	colOptional           = termbox.ColorDarkGray
	colMissing            = termbox.ColorMagenta
	colError              = termbox.ColorRed
)

var ERRAborted = fmt.Errorf("aborted")

const minColWidth = 10
const errText = " --------- "
const missingText = errText + "required" + errText

type Form interface {
	Lines() []Line
	SelectLine(selectedIndex int) (int, error)
}

type formOffset struct {
	XOffset, YOffset int
}

type termForm struct {
	lines            []Line
	offset           *formOffset
	firstColumnWidth int
}

func (cf termForm) Lines() []Line {
	return cf.lines
}

func (cf termForm) SelectLine(selectedIndex int) (int, error) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.Clear(colBackground, colBackground)

	lineLen := len(cf.lines)
	index := selectedIndex
	if index < 0 || index >= lineLen {
		index = lineLen - 1
	}
	cf.render(index)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowUp: // Previous line
				if index > 0 {
					index--
				}
			case termbox.KeyArrowDown: // next line
				if index < (lineLen - 1) {
					index++
				}
			case termbox.KeyEnter: // Selected
				return index, nil
			case termbox.KeyEsc: // aborted
				return -1, ERRAborted
			}
		case termbox.EventError:
			return -1, ev.Err

		default:
			// ignore unknowns
			//fmt.Printf("Event: %v", ev)
		}
		cf.render(index)
	}
}

func (cf termForm) render(selectedIndex int) {
	var offset = *cf.offset // take a copy of the form base position
	for i, c := range cf.lines {
		bg := colBackground
		if i == selectedIndex {
			bg = colBackgroundSelected
		}
		fg := colNormal
		if !c.Required() {
			fg = colOptional
		}

		offset.XOffset += tbprint(offset.XOffset, offset.YOffset, fg, bg, c.Name())
		offset.XOffset += tbprint(offset.XOffset, offset.YOffset, colNormal, bg, ":    ")

		value, vCol := cf.valueAndColour(c)
		offset.XOffset += tbprint(offset.XOffset, offset.YOffset, vCol, bg, value)

		offset.YOffset++
		offset.XOffset = 0
	}
}

func (cf termForm) valueAndColour(c Line) (string, termbox.Attribute) {
	fg := colNormal
	s := c.Format()
	// Check if its valid by parsing itself
	if err := c.Parse(s); err != nil {
		s = errText + err.Error() + errText
		fg = colError
	} else if s == "" && c.Required() {
		s = missingText
		fg = colMissing
	}
	return s, fg
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) int {
	var w int
	for _, c := range msg {
		termbox.SetCell(x+w, y, c, fg, bg)
		w += runewidth.RuneWidth(c)
	}
	return w
}

func maxNameWidth(lines []Line) int {
	max := minColWidth
	for _, c := range lines {
		l := len(c.Name())
		if l <= max {
			continue
		}
		max = l
	}
	return max
}

func NewForm(xOffset, yOffset int, lines []Line) Form {
	w := maxNameWidth(lines)
	return &termForm{
		lines: lines,
		offset: &formOffset{
			XOffset: xOffset,
			YOffset: yOffset,
		},
		firstColumnWidth: w,
	}
}

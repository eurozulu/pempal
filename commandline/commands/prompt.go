package commands

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

const colBg = termbox.ColorDefault
const colFg = termbox.ColorBlack
const colSel = termbox.ColorMagenta

type Prompt struct {
	Title            string
	Choice           []string
	DefaultIndex     int
	OffsetX, OffsetY int
	HorizontalList   bool
}

func (p *Prompt) Select() (int, error) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	selectedIndex := p.DefaultIndex
	if selectedIndex < 0 || selectedIndex >= len(p.Choice) {
		selectedIndex = len(p.Choice) - 1
	}

	p.writePrompt()
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEnter:
				return selectedIndex, nil
			case termbox.KeyEsc:
				return -1, fmt.Errorf("aborted")
			case termbox.KeyArrowUp:
				if selectedIndex > 0 {
					selectedIndex--
				}
			case termbox.KeyArrowDown:
				if selectedIndex < (len(p.Choice) - 1) {
					selectedIndex++
				}
			}
		case termbox.EventError:
			return -1, ev.Err

		default:
			fmt.Printf("Event: %v", ev)
		}
		p.writeList(selectedIndex)
	}
}

func (p Prompt) writePrompt() int {
	y := p.OffsetY
	if p.Title != "" {
		tbprint(p.OffsetX, y, colFg, colBg, p.Title)
		y++
	}
	tbprint(p.OffsetX, y, colFg, colBg, "Press ESC to quit")
	y++
	termbox.Flush()
	return y
}

func (p Prompt) writeList(selectedIndex int) {
	y := p.OffsetY + 1
	if p.Title != "" {
		y++
	}
	x := p.OffsetX
	for i, s := range p.Choice {
		bg := colBg
		if i == selectedIndex {
			bg = colSel
		}
		w := tbprint(x, y, colFg, bg, s)
		if !p.HorizontalList {
			y++
		} else {
			x += w
		}
	}
	termbox.Flush()
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) int {
	var w int
	for _, c := range msg {
		termbox.SetCell(x+w, y, c, fg, bg)
		w += runewidth.RuneWidth(c)
	}
	return w
}

func PromptList(prompt string, list []string, def int) (int, error) {
	p := &Prompt{
		Title:          prompt,
		Choice:         list,
		DefaultIndex:   def,
		OffsetX:        0,
		OffsetY:        0,
		HorizontalList: false,
	}
	return p.Select()
}

func PromptYorN(prompt string, def bool) (bool, error) {
	d := 1
	if def {
		d = 0
	}
	p := &Prompt{
		Title:          prompt,
		Choice:         []string{"yes", "no"},
		DefaultIndex:   d,
		OffsetX:        0,
		OffsetY:        0,
		HorizontalList: true,
	}
	i, err := p.Select()
	if err != nil {
		return false, err
	}
	return i == 0, nil
}

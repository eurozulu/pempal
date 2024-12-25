package ui

import (
	"github.com/eurozulu/pempal/logging"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"strings"
)

const escMessage = "Hit enter to select or Esc to cancel"

type SelectList struct {
	Title          string
	Choices        []string
	SelectedIndex  int
	HideEscMessage bool
}

func (sl SelectList) Selected() (string, bool) {
	if sl.SelectedIndex < 0 {
		return "", false
	}
	return sl.Choices[sl.SelectedIndex], true
}

func (sl SelectList) MoveToFirst() {

}

func (sl SelectList) MoveToLast() {

}

func (sl *SelectList) MoveToPrevious() {
	if sl.SelectedIndex > 0 {
		sl.SelectedIndex--
		sl.redraw()
	}
}

func (sl *SelectList) MoveToNext() {
	if sl.SelectedIndex < len(sl.Choices)-1 {
		sl.SelectedIndex++
		sl.redraw()
	}
}

func (sl *SelectList) MoveToRune(r rune) {
	s := string(r)
	for i, choice := range sl.Choices {
		if !strings.EqualFold(choice[:1], s) {
			continue
		}
		// Already selected, move to next
		if sl.SelectedIndex == i {
			continue
		}
		sl.SelectedIndex = i
		return
	}
}

func (sl *SelectList) TextWidth() int {
	var size int
	for _, choice := range sl.Choices {
		if len(choice) > size {
			size = len(choice)
		}
	}
	return size
}

func (sl SelectList) redraw() {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)
	w, h := termbox.Size()
	if h == 0 {
		return
	}
	textWidth := sl.TextWidth()
	midx := (w - textWidth) / 2
	midy := (h / 2) - (len(sl.Choices) / 2)
	bottom := midy + len(sl.Choices)
	if sl.Title != "" {
		tbprint(midx-len(sl.Title), midy-2, coldef, coldef, sl.Title)
		s, _ := sl.Selected()
		tbprint(midx, midy-2, coldef, coldef, s)
	}
	if !sl.HideEscMessage {
		tbprint(midx-(len(escMessage)/2), bottom+1, coldef, coldef, escMessage)
	}

	termbox.SetCell(midx-1, midy-1, '┌', coldef, coldef)
	termbox.SetCell(midx-1, bottom, '└', coldef, coldef)
	termbox.SetCell(midx+textWidth, midy-1, '┐', coldef, coldef)
	termbox.SetCell(midx+textWidth, bottom, '┘', coldef, coldef)
	sl.fill(midx, midy-1, textWidth, 1, termbox.Cell{Ch: '─'})
	sl.fill(midx, bottom, textWidth, 1, termbox.Cell{Ch: '─'})

	for i, choice := range sl.Choices {
		col := coldef
		if i == sl.SelectedIndex {
			col = termbox.ColorCyan
		}
		termbox.SetCell(midx-1, midy+i, '│', coldef, coldef)
		tbprint(midx, midy+i, coldef, col, choice)
		termbox.SetCell(midx+textWidth, midy+i, '│', coldef, coldef)
	}
	termbox.Flush()
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func (sl SelectList) fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

func (sl *SelectList) Show() error {
	if !termbox.IsInit {
		err := termbox.Init()
		if err != nil {
			return err
		}
		defer termbox.Close()
	}
	termbox.SetInputMode(termbox.InputEsc)

	sl.redraw()
mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				sl.SelectedIndex = -1
				break mainloop
			case termbox.KeyEnter:
				break mainloop
			case termbox.KeyArrowUp, termbox.KeyCtrlP:
				sl.MoveToPrevious()
			case termbox.KeyArrowDown, termbox.KeyCtrlN:
				sl.MoveToNext()
			case termbox.KeyHome, termbox.KeyCtrlA:
				sl.MoveToFirst()
			case termbox.KeyEnd, termbox.KeyCtrlE:
				sl.MoveToLast()
			default:
				if ev.Ch != 0 {
					sl.MoveToRune(ev.Ch)
				}
			}
		case termbox.EventError:
			return ev.Err
		}
		sl.redraw()
	}
	return nil
}

func ShowSelectList(prompt string, choices []string) (string, bool) {
	sl := &SelectList{
		Title:   prompt,
		Choices: choices,
	}
	if err := sl.Show(); err != nil {
		logging.Error("ShowSelectList", "Failed to show terminal display. %v", err)
		return "", false
	}
	return sl.Selected()
}

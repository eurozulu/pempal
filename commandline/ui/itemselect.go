package ui

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

type ItemSelect struct {
	Items    []ListItem
	ExitChar rune
}

type ListValues map[string]string

func (sl ItemSelect) Select(offset ViewOffset, selected int) (int, error) {
	for {
		ItemList(sl.Items).render(offset, selected)
		if err := termbox.Flush(); err != nil {
			return -1, err
		}

		ev, err := nextKeyEvent()
		if err != nil {
			return -1, err
		}
		switch ev.Key {
		case termbox.KeyArrowUp: // Previous line
			if selected > 0 {
				selected--
			}
		case termbox.KeyArrowDown: // next line
			if selected < len(sl.Items)-1 {
				selected++
			}
		case termbox.KeyEnter: // Selected last item
			return selected, nil

		case termbox.KeyEsc: // aborted
			return -1, ERRAborted

		default:
			if sl.ExitChar > 0 && ev.Ch == sl.ExitChar {
				return -1, nil
			}
		}
	}
}

func InitUI() (bool, error) {
	isRoot := !termbox.IsInit
	if isRoot {
		if err := termbox.Init(); err != nil {
			return false, err
		}
	}
	return isRoot, nil
}

func Clear() {
	termbox.Clear(ColourBackground.toAttribute(), ColourBackground.toAttribute())
}

func CloseUi() {
	termbox.Close()
}

func PrintF(offset *ViewOffset, fg, bg ItemColour, format string, a ...any) {
	if len(a) > 0 {
		format = fmt.Sprintf(format, a)
	}
	tbprint(offset, fg.toAttribute(), bg.toAttribute(), format)
	termbox.Flush()
}

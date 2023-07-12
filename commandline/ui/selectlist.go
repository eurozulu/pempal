package ui

import "github.com/nsf/termbox-go"

type SelectList struct {
	Names    []string
	ExitChar rune
}

type ListValues map[string]string

func (sl SelectList) Select(offset ViewOffset, selected int, values ListValues) (int, error) {
	isRoot := !termbox.IsInit
	if isRoot {
		if err := termbox.Init(); err != nil {
			return -1, err
		}
		defer termbox.Close()
	}
	list := sl.itemListOfNames(values)

	for {
		if isRoot {
			termbox.Clear(ColourBackground.ToAttribute(), ColourBackground.ToAttribute())
		}
		list.renderList(offset, selected)

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
			if selected < len(list)-1 {
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

func (sl SelectList) itemListOfNames(values ListValues) ItemList {
	items := make([]ListItem, len(sl.Names))
	for i, n := range sl.Names {
		items[i] = ListItem{
			Name:  n,
			Value: values[n],
		}
	}
	return items
}

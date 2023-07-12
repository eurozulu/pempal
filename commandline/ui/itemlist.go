package ui

import (
	"github.com/nsf/termbox-go"
)

type ItemList []ListItem

func (sl ItemList) renderList(offset ViewOffset, selected int) error {
	for i, li := range sl {
		li.render(offset, i == selected)
		offset.YOffset++
	}
	return termbox.Flush()
}

func NewItemListOfValues(values []string) ItemList {
	items := make([]ListItem, len(values))
	for i, v := range values {
		items[i] = ListItem{Value: v}
	}
	return items
}

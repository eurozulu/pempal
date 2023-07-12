package ui

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"strings"
)

type ViewOffset struct {
	XOffset, YOffset int
}

func (vo ViewOffset) OffsetY(yOffset int) ViewOffset {
	return ViewOffset{
		XOffset: vo.XOffset,
		YOffset: vo.YOffset + yOffset,
	}
}

func (vo ViewOffset) OffsetX(xOffset int) ViewOffset {
	return ViewOffset{
		XOffset: vo.XOffset + xOffset,
		YOffset: vo.YOffset,
	}
}

func cleanLine(y int, bg termbox.Attribute) {
	w, h := termbox.Size()
	if y >= h || y < 0 {
		return
	}
	tbprint(&ViewOffset{
		XOffset: 0,
		YOffset: y,
	}, bg, bg, strings.Repeat(" ", w))
}

func tbprint(offset *ViewOffset, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(offset.XOffset, offset.YOffset, c, fg, bg)
		offset.XOffset += runewidth.RuneWidth(c)
	}
}

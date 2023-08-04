package ui

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

type Window interface {
	Show(v View) (View, error)
	Render(v View) error
}

type window struct {
	title  string
	frame  ViewFrame
	isRoot bool
}

func (win window) Render(v View) error {
	if !termbox.IsInit {
		return fmt.Errorf("ui not initalised")
	}
	if win.isRoot {
		win.frame.Clear()
		win.renderTitle()
	} else {
		win.clearArea(childHeight(v))
	}
	v.Render(win.frame)
	return termbox.Flush()
}

func (win window) Show(v View) (View, error) {
	win.isRoot = !termbox.IsInit
	if win.isRoot {
		if err := termbox.Init(); err != nil {
			return nil, err
		}
		defer func() {
			termbox.Close()
		}()
	}
	text := v.String()
	if tv, ok := v.(TextView); ok {
		text = tv.GetText()
	}

	if vo, ok := v.(WindowEventOpen); ok {
		vo.OnViewOpen()
	}

	for {
		if err := win.Render(v); err != nil {
			return nil, err
		}

		if exit, err := win.readInput(v); err != nil {
			if tv, ok := v.(TextView); ok {
				tv.SetText(text)
			}
			return nil, err
		} else if !exit {
			// input processed but no exit signal, rerender and read next char
			continue
		}
		child := selectedChildView(v)

		// if child view accepts input, show that as a View, otherwise return the current view
		if child != nil {
			if _, ok := child.(TextView); ok {
				cv, err := win.showChild(v, child)
				if err == ErrAborted {
					continue
				} // abort back to parent
				if wv, ok := v.(WindowEventClose); ok {
					wv.OnViewClose(cv)
				}
				return cv, err
			}
		}
		// No child or child not a parentview, return current view.
		if wv, ok := v.(WindowEventClose); ok {
			wv.OnViewClose(v)
		}
		return v, nil
	}
}

func (win window) showChild(parent, child View) (View, error) {
	var offset ViewOffset
	if parent != nil {
		offset.X = len(parent.Label()) + 2
		offset.Y = selectedChildIndex(parent)
	}
	win.frame = win.frame.WithRelativeOffset(offset.X, offset.Y)
	nv, err := win.Show(child)
	if err != nil {
		return nil, err
	}
	return nv, nil
}

// readInput reads the next char from the keyboard.
// returns an exit flag or error.
// error (ErrAborted) is returned if ESC pressed
// true exit if Enter is pressed, otherwise returns false, with the char
// having been appened to the given view.
func (win window) readInput(view View) (bool, error) {
	mv, isMutable := view.(TextView)

	ch, err := nextKeyChar()
	if err != nil {
		return true, err
	}
	switch ch {
	case rune(termbox.KeyEsc):
		return true, ErrAborted

	case rune(termbox.KeyEnter):
		return true, nil

	default:
		if isMutable {
			mv.AppendText(ch)
		}
		return false, nil
	}
}

func (win window) clearArea(height int) {
	win.frame.ResetPosition()
	for i := 0; i < height; i++ {
		win.frame.ClearLine()
		win.frame.Println()
	}
	win.frame.ResetPosition()
}

func (win window) renderTitle() {
	if win.title == "" {
		return
	}
	win.frame.WithColour(titleColours).Println(win.title)
}

func childHeight(v View) int {
	c := 1
	if pv, ok := v.(ParentView); ok {
		c += len(pv.ChildViews())
	}
	return c
}

func selectedChildView(v View) View {
	i := selectedChildIndex(v)
	if i < 0 {
		return nil
	}
	return v.(ParentView).ChildViews()[i]
}

func selectedChildIndex(v View) int {
	if v == nil {
		return -1
	}
	pv, ok := v.(ParentView)
	if !ok {
		return -1
	}
	return pv.SelectedIndex()
}

func nextKeyChar() (rune, error) {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			ch := ev.Ch
			if ch == 0 {
				ch = rune(ev.Key)
			}
			return ch, nil
		case termbox.EventError:
			return 0, ev.Err
		default:
			continue
		}
	}
}

func NewWindow(title string, x, y int) Window {
	return &window{
		title: title,
		frame: newViewFrame(ViewOffset{X: x, Y: y}),
	}
}

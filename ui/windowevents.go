package ui

type WindowEventOpen interface {
	OnViewOpen()
}

type WindowEventClose interface {
	OnViewClose(child View)
}

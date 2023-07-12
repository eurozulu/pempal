package ui

type InputType int

const (
	// InputTypeNone only the Options of the field may be selected
	InputTypeNone InputType = iota
	// InputTypeNumbers only the options and any numeric value may be entered
	InputTypeNumbers
	// InputTypeLetters only the options and any string of letters may be entered
	InputTypeLetters
	// InputTypePrintable only the options and any printable string of characters may be entered
	InputTypePrintable
)

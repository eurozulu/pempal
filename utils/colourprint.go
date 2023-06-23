package utils

import (
	"bytes"
	"io"
)

// adapted from: https://golangbyexample.com/print-output-text-color-console/
type Colour string

const (
	ColourRed    = Colour("\033[31m")
	ColourGreen  = Colour("\033[32m")
	ColourYellow = Colour("\033[33m")
	ColourBlue   = Colour("\033[34m")
	ColourPurple = Colour("\033[35m")
	ColourCyan   = Colour("\033[36m")
	ColourWhite  = Colour("\033[37m")

	ColourEnd = Colour("\033[0m")
)

type ColourOut struct {
	Out io.Writer
	Col Colour
}

func (co ColourOut) Write(p []byte) (n int, err error) {
	return co.Out.Write(bytes.Join([][]byte{co.Col.Bytes(), p, ColourEnd.Bytes()}, nil))
}

func (c Colour) Bytes() []byte {
	return []byte(c)
}

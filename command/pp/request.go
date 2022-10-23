package main

import (
	"io"
)

type requestCommand struct {
}

func (r requestCommand) Run(out io.Writer, args ...string) error {
	//TODO implement me
	panic("implement me")
}

func newRequestCommand() Command {
	return &requestCommand{}
}

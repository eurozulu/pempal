package main

import (
	"io"
)

type findCommand struct {
}

func (f findCommand) Run(out io.Writer, args ...string) error {
	//TODO implement me
	panic("implement me")
}

func newFindCommand() Command {
	return &findCommand{}
}

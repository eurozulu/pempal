package main

import "io"

type revokeCommand struct {
}

func (r revokeCommand) Run(out io.Writer, args ...string) error {
	//TODO implement me
	panic("implement me")
}

func newRevokeCommand() Command {
	return &revokeCommand{}
}

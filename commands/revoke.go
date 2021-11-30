package commands

import (
	"context"
	"flag"
	"io"
)

type RevokeCommand struct {
	Key string
}

func (cmd RevokeCommand) Description() string {
	return "creates a new, Certificate revokation list"
}

func (cmd RevokeCommand) Flags(f *flag.FlagSet) {
	f.StringVar(&cmd.Key, "key", "", "specify the key to sign the request")
}

func (cmd RevokeCommand) Run(ctx context.Context, out io.Writer, args ...string) error {
	// TODO:  Implement this
	panic("implement me")
}

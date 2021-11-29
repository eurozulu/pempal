package cmd

import (
	"context"
	"flag"
	"io"
)

type RequestsCommand struct {
}

func (r RequestsCommand) Description() string {
	return "Lists all the Certificate Requests.  Can filter by requests unmatched to any issued certificate"
}

func (r RequestsCommand) Flags(f *flag.FlagSet) {
	panic("implement me")
}

func (r RequestsCommand) Run(ctx context.Context, out io.Writer, args ...string) error {
	panic("implement me")
}

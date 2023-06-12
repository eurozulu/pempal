package commands

import (
	"github.com/eurozulu/pempal/resourceio"
	"io"
)

type FormatCommand struct {
	Format string
}

func (sc FormatCommand) Execute(args []string, out io.Writer) error {
	panic("not done yet")
}

func (sc FormatCommand) resolveResourceFormat() (resourceio.ResourceFormat, error) {
	panic("not done yet")
}

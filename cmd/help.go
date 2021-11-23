package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"sort"
)

type HelpCommand struct {
}

func (h HelpCommand) Flags(f *flag.FlagSet) {
}

func (h HelpCommand) Description() string {
	return "Displays help"
}

// TODO: Tidy up output
func (h HelpCommand) Run(ctx context.Context, out io.Writer, args ...string) error {
	for _, cmd := range sortedCommandKeys(Commands) {
		c := Commands[cmd]
		fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
		c.Flags(fs)
		fs.Usage()
		fmt.Fprintf(out, "%s\n", c.Description())
	}
	return nil
}

func sortedCommandKeys(m map[string]Command) []string {
	keys := make([]string, len(m))
	var i int
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

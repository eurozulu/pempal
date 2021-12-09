package commands

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"
	"text/tabwriter"
)

type HelpCommand struct{}

func (h HelpCommand) Flags(f *flag.FlagSet) {}

func (h HelpCommand) Description() string {
	return "Displays help, what you're looking at!"
}

// TODO: Tidy up output
func (h HelpCommand) Run(ctx context.Context, out io.Writer, args ...string) error {
	if len(args) > 0 {
		return showCommandHelp(args[0])
	}
	return showGeneralHelp()
}

func showGeneralHelp() error {
	fn := path.Base(os.Args[0])
	fmt.Println("\nUsage:")
	fmt.Printf("%s <command> [<optional general flags>] [<optional command flags>] <command arguments>\n", fn)
	fmt.Printf("<command> is required and most commands require at least one argument.\n")
	fmt.Printf("use '%s <command> -help' for detailed help about each command and its flags.\n", fn)
	fmt.Println()
	fmt.Println("<optional general flags> can be used with any command.")
	fs := flag.NewFlagSet("General Flags", flag.ContinueOnError)
	Flags(fs)
	fs.Usage()

	fmt.Println()
	fmt.Println("Available commands.  use the -help flag with any of the commands to get more details about the command.")
	tw := tabwriter.NewWriter(os.Stdout, 8, 4, 4, ' ', 0)
	for _, cmd := range sortedCommandKeys(Commands) {
		c := Commands[cmd]
		// first line of the description
		desc := strings.SplitN(c.Description(), "\n", 2)[0]
		fmt.Fprintf(tw, "%s\t%s\n", cmd, desc)
	}
	tw.Flush()
	return nil
}

// showCommandHelp shows full help about a single command
func showCommandHelp(cmd string) error {
	if acmd, ok := Aliases[cmd]; ok {
		cmd = acmd
	}
	c := Commands[cmd]
	if c == nil {
		return fmt.Errorf("%s is an unknown command", cmd)
	}
	fmt.Printf("\nUsage of %s:\n", cmd)
	aNames := sortedAliasNames(cmd)
	if len(aNames) > 0 {
		fmt.Printf("Shortcut alternatives: %s\n", strings.Join(aNames, ", "))
	}
	fmt.Println(c.Description())
	fs := flag.NewFlagSet(strings.Join([]string{cmd, "flags"}, ""), flag.ContinueOnError)
	c.Flags(fs)
	fs.Usage()
	return nil
}

func sortedAliasNames(cmd string) []string {
	var names []string
	for k, v := range Aliases {
		if v != cmd {
			continue
		}
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}
func sortedCommandKeys(m map[string]Command) []string {
	var keys []string
	for k := range m {
		if k == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

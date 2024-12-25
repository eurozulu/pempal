package commands

import (
	"fmt"
	"strings"
)

type HelpCommand struct{}

func (h HelpCommand) Exec(args ...string) error {
	fmt.Printf("Show help here for: %s\n", strings.Join(args, " "))
	return nil
}

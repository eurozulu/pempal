//go:generate ../../spud generate . -packagename main

package main

import (
	"fmt"
	"github.com/eurozulu/spud/taglibs/subcommander"
	"os"
)

func main() {
	result, err := subcommander.Execute(os.Args[1:]...)
	if result != nil {
		fmt.Println(result)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

}

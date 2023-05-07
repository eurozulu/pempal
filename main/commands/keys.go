package commands

import (
	"fmt"
	"io"
	"pempal/utils"
)

type KeysCommand struct {
}

func (k KeysCommand) Execute(args []string, out io.Writer) error {
	if Keys == nil {
		return fmt.Errorf("no key manager available.  Check configuration for keypath and cert path.")
	}
	colOut := utils.NewColumnOutput(out)
	colOut.ColumnWidths = []int{4}
	colOut.WriteSlice([]string{"Identity", "Names"})
	colOut.Write([]byte("\n"))
	var count int
	for id, users := range Keys.Users() {
		count++
		colOut.Write([]byte(id.String()))
		colOut.Write([]byte(":\n"))
		var ccount int
		for _, user := range users {
			ccount++
			colOut.WriteSlice([]string{"-", user.Certificate().Subject.String()})
			colOut.Write([]byte("\n"))
		}
		fmt.Fprintf(colOut, "-,%d certificates\n", ccount)
		colOut.Write([]byte("\n"))
	}
	fmt.Fprintf(colOut, "%d users\n", count)
	return nil
}

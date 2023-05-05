package commands

import (
	"fmt"
	"io"
	"pempal/utils"
)

type ConfigCommand struct {
	Global bool `flag:global`
}

func (c ConfigCommand) Execute(args []string, out io.Writer) error {
	colOut := utils.NewColumnOutput(out)
	colOut.Delimiter = ":"
	colOut.ColumnWidths = []int{16}
	colOut.ColumnDelimiter = " "
	if _, err := fmt.Fprintln(colOut, "Pempal Configuration:\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(colOut, "Root path: %s\n", Configuration.RootPath); err != nil {
		return err
	}
	if Configuration.RootCertificate != "" {
		if _, err := fmt.Fprintf(colOut, "Root certificate: %s\n", Configuration.RootCertificate); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(colOut, "certificates: %s\n", Configuration.CertPath); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(colOut, "keys: %s\n", Configuration.KeyPath); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(colOut, "requests: %s\n", Configuration.CsrPath); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(colOut, "revoked: %s\n", Configuration.CrlPath); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(colOut, "templates: %s\n", Configuration.TemplatePath); err != nil {
		return err
	}
	return nil
}

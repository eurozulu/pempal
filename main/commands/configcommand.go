package commands

import (
	"fmt"
	"io"
	"pempal/config"
	"pempal/utils"
)

type ConfigCommand struct {
	Global bool `flag:global`
}

func (c ConfigCommand) Execute(args []string, out io.Writer) error {
	cfg := config.NewConfig()

	colOut := utils.NewColumnOutput(out)
	colOut.Delimiter = ":"
	colOut.ColumnWidths = []int{16}
	colOut.ColumnDelimiter = " "
	if _, err := fmt.Fprintln(colOut, "Pempal configuration:\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(colOut, "Root path: %s\n", cfg.RootPath); err != nil {
		return err
	}
	if cfg.RootCertificate != "" {
		if _, err := fmt.Fprintf(colOut, "Root certificate: %s\n", cfg.RootCertificate); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(colOut, "certificates: %s\n", cfg.CertPath); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(colOut, "keys: %s\n", cfg.KeyPath); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(colOut, "requests: %s\n", cfg.CsrPath); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(colOut, "revoked: %s\n", cfg.CrlPath); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(colOut, "templates: %s\n", cfg.TemplatePath); err != nil {
		return err
	}
	return nil
}

package commands

import (
	"fmt"
	"io"
	"pempal/config"
)

type ConfigCommand struct {
	Global bool `flag:global`
}

func (c ConfigCommand) Execute(args []string, out io.Writer) error {
	cfg := config.NewConfig()
	if _, err := fmt.Fprintln(out, "Pempal configuration:\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "Root path: %s\n", cfg.RootPath); err != nil {
		return err
	}
	if cfg.RootCertificate != "" {
		if _, err := fmt.Fprintf(out, "Root certificate: %s\n", cfg.RootCertificate); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(out, "certificates: %s\n", cfg.CertPath); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "keys: %s\n", cfg.KeyPath); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "requests: %s\n", cfg.CsrPath); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "revoked: %s\n", cfg.CrlPath); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "templates: %s\n", cfg.TemplatePath); err != nil {
		return err
	}
	return nil
}

package commands

import (
	"fmt"
	"io"
	"pempal/config"
	"pempal/utils"
)

type ConfigCommand struct {
	RootPath        string `flag:"root-path,rootpath"`
	RootCertificate string `flag:"root-certificate,rootcertificate,rootcert"`
	CertPath        string `flag:"cert-path,certpath"`
	KeyPath         string `flag:"key-path,keypath"`
	CsrPath         string `flag:"csr-path,csrpath"`
	CrlPath         string `flag:"crl-path,crlpath"`
	TemplatePath    string `flag:"template-path,templatepath"`
}

func (cmd ConfigCommand) Execute(args []string, out io.Writer) error {
	colOut := utils.NewColumnOutput(out)
	colOut.Delimiter = ":"
	colOut.ColumnWidths = []int{16}
	colOut.ColumnDelimiter = " "

	cfg, err := GetConfig(args...)
	if err != nil {
		return err
	}
	cmd.applyFlags(cfg.(*config.DefaultConfig))

	if _, err := fmt.Fprintln(colOut, "Pempal Configuration:\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(colOut, "Config path: %s\n", cfg.ConfigLocation()); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(colOut, "Root path: %s\n", cfg.Root()); err != nil {
		return err
	}
	if cfg.RootCertificate() != "" {
		if _, err := fmt.Fprintf(colOut, "Root certificate: %s\n", cfg.RootCertificate()); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(colOut, "certificates: %s\n", cfg.Certificates()); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(colOut, "keys: %s\n", cfg.Keys()); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(colOut, "requests: %s\n", cfg.Requests()); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(colOut, "revoked: %s\n", cfg.Revokations()); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(colOut, "templates: %s\n", cfg.Templates()); err != nil {
		return err
	}
	return nil
}

func (cmd ConfigCommand) applyFlags(cfg *config.DefaultConfig) {
	if cmd.RootPath != "" {
		cfg.RootPath = cmd.RootPath
	}
	if cmd.RootCertificate != "" {
		cfg.RootCertificatePath = cmd.RootCertificate
	}
	if cmd.CertPath != "" {
		cfg.CertPath = cmd.CertPath
	}
	if cmd.KeyPath != "" {
		cfg.KeyPath = cmd.KeyPath
	}
	if cmd.CsrPath != "" {
		cfg.CsrPath = cmd.CsrPath
	}
	if cmd.CrlPath != "" {
		cfg.CrlPath = cmd.CrlPath
	}
	if cmd.TemplatePath != "" {
		cfg.TemplatePath = cmd.TemplatePath
	}
}

func GetConfig(path ...string) (config.Config, error) {
	if len(path) == 0 {
		return Configuration, nil
	}
	return config.NewConfig(path[0])
}

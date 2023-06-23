package commands

import (
	"context"
	"fmt"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/keys"
	"github.com/eurozulu/pempal/resourceio"
	"github.com/eurozulu/pempal/utils"
	"io"
	"strconv"
)

type keysCommand struct {
	Names     bool `flag:"names"`
	Recursive bool `flag:"recursive,r"`

	keys keys.Keys
}

func (kc keysCommand) Execute(args []string, out io.Writer) error {
	if km, err := config.KeyManager(); err != nil {
		return fmt.Errorf("key manager not available %v", err)
	} else {
		kc.keys = km
	}

	if a, err := cleanArguments(args); err != nil {
		return err
	} else {
		args = a
	}

	cols := utils.NewColumnOutput(out)
	if !CommonFlags.Quiet {
		if err := kc.writeKeyHeaders(cols); err != nil {
			return err
		}
	}
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	for loc := range kc.keys.AllKeys(ctx) {
		if err := kc.writeKeyLocation(cols, loc); err != nil {
			return err
		}
	}
	return nil
}

func (kc keysCommand) writeKeyHeaders(out *utils.ColumnOutput) error {
	fields := []string{
		"identity",
		"encrypted",
		"location",
	}
	_, err := out.WriteSlice(fields)
	if err != nil {
		return err
	}
	err = out.WriteString("\n")
	return err
}

func (kc keysCommand) writeCertificateHeaders(out *utils.ColumnOutput) error {
	fields := []string{
		"",
		"serial number",
		"subject",
	}
	_, err := out.WriteSlice(fields)
	if err != nil {
		return err
	}
	err = out.WriteString("\n")
	return err
}

func (kc keysCommand) writeKeyLocation(out *utils.ColumnOutput, loc resourceio.ResourceLocation) error {
	prks := keys.ParseKeyLocation(loc)
	for _, k := range prks {
		fields := []string{
			k.Identity().String(),
			strconv.FormatBool(k.IsEncrypted()),
			loc.Location(),
		}
		if _, err := out.WriteSlice(fields); err != nil {
			return err
		}

		if kc.Names {
			if err := out.WriteString("\n"); err != nil {
				return err
			}

			if !CommonFlags.Quiet {
				if err := kc.writeCertificateHeaders(out); err != nil {
					return err
				}
			}
			if err := kc.writeCertificates(out, k); err != nil {
				return err
			}
		}
		if err := out.WriteString("\n"); err != nil {
			return err
		}
	}
	return nil
}

func (kc keysCommand) writeCertificates(out *utils.ColumnOutput, k keys.Key) error {
	for _, c := range kc.keys.CertificatesById(k.Identity()) {
		fields := []string{
			"",
			c.SerialNumber.String(),
			c.Subject.String(),
		}
		if _, err := out.WriteSlice(fields); err != nil {
			return err
		}
		if err := out.WriteString("\n"); err != nil {
			return err
		}
	}
	return nil
}

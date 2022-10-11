package commands

import (
	"bytes"
	"context"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"pempal/keycache"
	"pempal/keytools"
	"pempal/keytracker"
	"pempal/pemreader/fileformats"
	"pempal/pemresources"
	"pempal/pemwriter"
	"pempal/templates"
	"sort"
)

type RequestCommand struct {
	keyCache *keycache.KeyCache
}

func (cmd RequestCommand) Description() string {
	return "creates a new, signed certificate request"
}

func (cmd RequestCommand) Flags(f *flag.FlagSet) {
}

func (cmd *RequestCommand) Run(ctx context.Context, out io.Writer, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("must provide a template name or path to existing resource to use as the certificate request values.")
	}
	// Build the template for CSR
	var csr = &pemresources.CertificateRequest{}
	t, err := templates.CompoundTemplate(csr, args...)
	if err != nil {
		return err
	}

	der, err := templates.SignRequest(k, cmd.KeyPass, t)
	if err != nil {
		return err
	}
	pemwriter.NewPemFormat(out).Write(&pem.Block{
		Type:  keytools.PEM_CERTIFICATE_REQUEST,
		Bytes: der,
	})
	return nil
}

func (cmd RequestCommand) getUserKey(ctx context.Context, out io.Writer) (*pemresources.PrivateKey, error) {
}

func (cmd *RequestCommand) createNewKey(out io.Writer) (keytracker.Key, error) {
	// use key command to generate new key
	kc := &KeyCommand{
		KeyAlgorithm: cmd.KeyAlgorithm,
		KeyLength:    cmd.KeyLength,
		Quiet:        cmd.Quiet,
		Recursive:    true,
		Public:       true,
		Private:      true,
		Password:     cmd.KeyPass,
		LinkPublic:   false,
	}
	// stream output into buffer and parse back into pem blocks
	buf := bytes.NewBuffer(nil)
	if err := kc.createPrivateKey(buf); err != nil {
		return nil, err
	}
	// copy properties of new key into this command
	cmd.KeyPass = kc.Password
	cmd.KeyLength = kc.KeyLength
	cmd.KeyAlgorithm = kc.KeyAlgorithm

	pf := fileformats.FileFormats["pem"]
	blks, err := pf.Format(buf.Bytes())
	if err != nil {
		return nil, err
	}
	if len(blks) != 2 {
		return nil, fmt.Errorf("unexpected new key output.  Expected private and public pems, found %d blocks", len(blks))
	}
	// write new key to output
	if _, err = io.Copy(out, buf); err != nil {
		return nil, err
	}
	if !keytools.PublicKeyTypes[blks[0].Type] {
		// swap so public is first
		blks = append(blks[1:], blks[0])
	}
	return keytracker.NewKeyWithPublic(blks[1], blks[0]), nil
}

func (cmd *RequestCommand) SetKeys(keys *keycache.KeyCache) {
	cmd.keyCache = keys
}

func fillInTemplate(names []string) templates.Template {
	t := templates.Template{}
	sort.Strings(names)
	for _, k := range names {
		sv := PromptInput(k, "", 1, 0)
		t.SetValue(k, sv)
	}
	return t
}

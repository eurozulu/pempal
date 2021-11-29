package cmd

import (
	"bytes"
	"context"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"pempal/keytools"
	"pempal/keytracker"
	"pempal/pemreader/fileformats"
	"pempal/pemwriter"
	"pempal/templates"
	"pempal/templates/parsers"
	"sort"
	"strings"
)

type RequestCommand struct {
	keyAlgorithm string
	keyLength    int
	keyPass      string

	key     string
	keyPath string
	quiet   bool
}

func (cmd RequestCommand) Description() string {
	return "creates a new, signed certificate request"
}

func (cmd RequestCommand) Flags(f *flag.FlagSet) {
	algos := fmt.Sprintf("The public key algorithm to use when generating new key.\n\tmust be one of %v", keytools.PublicKeyAlgoNames[1:])
	f.StringVar(&cmd.keyAlgorithm, "a", "rsa", algos)
	f.IntVar(&cmd.keyLength, "l", 2048, "The length/curve to use for the key")
	f.StringVar(&cmd.keyPass, "password", "", "password to encrypt new key or decrypt existing key")

	f.StringVar(&cmd.key, "key", "", "specify the key to sign the request")
	f.StringVar(&cmd.keyPath, "keypath", "", fmt.Sprintf("specify a comma delimit list of paths to search for key. Overides %s", ENV_KeyPath))
	f.BoolVar(&cmd.quiet, "q", false, "surpress the confirmation prompt.  Key MUST be supplied in template or flag")
}

func (cmd *RequestCommand) Run(ctx context.Context, out io.Writer, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("must provide a template name or path to existing resource to use as the certificate request values.")
	}
	// Build the template to create cert with
	tb := templates.NewTemplateBuilder()
	if err := tb.Add(args...); err != nil {
		return err
	}
	missing := tb.RequiredNames()
	if len(missing) > 0 {
		tb.AddTemplate(fillInTemplate(missing))
	}
	t, err := tb.Build()
	if err != nil {
		return err
	}

	// establish the key to sign with
	var kp []string
	if cmd.keyPath != "" {
		kp = strings.Split(cmd.keyPath, ":")
	} else {
		// use ENV keypath if not specified
		kp = GetKeyPath(nil)
	}
	if cmd.keyAlgorithm == "" {
		cmd.keyAlgorithm = t.Value(parsers.X509PublicKeyAlgorithm)
	}
	if cmd.key == "" {
		// if not specified, attempt to read from template
		puk, _ := t.PublicKey()
		if puk != nil {
			cmd.key = keytools.PublicKeySha1Hash(puk)
		}
	}
	k, err := cmd.getUserKey(ctx, out, kp...)
	if err != nil {
		return err
	}
	if !cmd.quiet {
		msg := bytes.NewBufferString("\n")
		msg.WriteString(t.String())
		msg.WriteString("\n\nCreate new certificate request:")
		if !PromptConfirm(msg.String(), true) {
			return fmt.Errorf("aborted")
		}
	}

	if k.IsEncrypted() && cmd.keyPass == "" {
		if cmd.quiet {
			return fmt.Errorf("csr sign failed as key %s is encrypted and no password was provided", k.String())
		}
		pwd, err := PromptPassword(fmt.Sprintf("Enter the password for key %s %s", k.String(), k.Location()))
		if err != nil {
			return err
		}
		cmd.keyPass = pwd
	}
	der, err := templates.SignRequest(k, cmd.keyPass, t)
	if err != nil {
		return err
	}
	pemwriter.NewPemFormat(out).Write(&pem.Block{
		Type:  keytools.PEM_CERTIFICATE_REQUEST,
		Bytes: der,
	})
	return nil
}

func (cmd RequestCommand) getUserKey(ctx context.Context, out io.Writer, keyPath ...string) (keytracker.Key, error) {
	var key keytracker.Key
	if cmd.key != "" {
		key = keytracker.KeyTracker{ShowLogs: VerboseFlag, Recursive: true}.FindKey(ctx, cmd.key, keyPath...)
		if key == nil {
			return nil, fmt.Errorf("could not find key %s", cmd.key)
		}
		return key, nil
	}
	// not specified, search for available keys
	keys := Keys(ctx, keyPath, true)
	if cmd.quiet {
		if len(keys) == 1 {
			return keys[0], nil
		}
		return nil, fmt.Errorf("no public key specified in template or flag")
	}

	// Ask user to select or generate a key
	keys = SortKeys(keys)
	names := KeyList(keys)
	names = append(names, "Generate new key")
	var prompt string
	if len(names) == 1 {
		prompt = "No keys found to sign request, create one or zero to abort"
	} else {
		prompt = "Select the key to sign request, create a new one or zero to abort"
	}
	choice := PromptChooseList(prompt, names)
	if choice < 0 {
		return nil, fmt.Errorf("aborted")
	}
	if choice < len(names)-1 {
		return keys[choice], nil
	}
	// request to create new key, return nil without err
	return cmd.createNewKey(out)
}

func (cmd *RequestCommand) createNewKey(out io.Writer) (keytracker.Key, error) {
	// use key command to generate new key
	kc := &KeyCommand{
		keyAlgorithm: cmd.keyAlgorithm,
		keyLength:    cmd.keyLength,
		quiet:        cmd.quiet,
		recursive:    true,
		public:       true,
		private:      true,
		passwd:       cmd.keyPass,
		linkPublic:   false,
	}
	// stream output into buffer and parse back into pem blocks
	buf := bytes.NewBuffer(nil)
	if err := kc.createPrivateKey(buf); err != nil {
		return nil, err
	}
	// copy properties of new key into this command
	cmd.keyPass = kc.passwd
	cmd.keyLength = kc.keyLength
	cmd.keyAlgorithm = kc.keyAlgorithm

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

func fillInTemplate(names []string) templates.Template {
	t := templates.Template{}
	sort.Strings(names)
	for _, k := range names {
		sv := PromptInput(k, "", 1, 0)
		t.SetValue(k, sv)
	}
	return t
}

package cmd

import (
	"context"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"pempal/keytools"
	"pempal/keytracker"
	"pempal/pemwriter"
	"pempal/templates"
	"pempal/templates/parsers"
	"strings"
)

type IssueCommand struct {
	issuer  string
	keyPass string
	keyPath string
	quiet   bool
}

func (cmd *IssueCommand) Description() string {
	lines := []string{"issues a new certificate, based on the given templates/resources"}
	lines = append(lines, "Properties from each resource are merged into one, from left to right, with the right most taking precedence")
	lines = append(lines, "At least one resource should contain the DN of the issuer of the certificare unless specified with the -issuer flag")
	lines = append(lines, "-issuer flag takes precedence over all templates and resource properties.")
	lines = append(lines, "If no issuer is provided, a prompt to select one is presented.")
	return strings.Join(lines, "\n")
}

func (cmd *IssueCommand) Flags(f *flag.FlagSet) {
	f.BoolVar(&cmd.quiet, "q", false, "surpress confirmation prompts")
	f.StringVar(&cmd.issuer, "issuer", "", "set the DN of the issuer for the new certificate. Overrides any template value")
	f.StringVar(&cmd.keyPass, "password", "", "Specify the password for an encrypted issuer key. Will prompt is required and not provided")
	f.StringVar(&cmd.keyPath, "keypath", "", "comma delimited list of directories to search for issuer keys.  Overrides KEYPATH environment variable")
}

func (cmd *IssueCommand) Run(ctx context.Context, out io.Writer, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("must provide a template name or path to existing resource to use as the certificate values.")
	}
	// Build the template to create cert with
	tb := templates.NewTemplateBuilder()
	if err := tb.Add(args...); err != nil {
		return err
	}
	t, err := tb.Build()
	if err != nil {
		return err
	}

	// Establish the signing key
	// If issuer not specified as flag, take it from the template (which might also be empty!)
	var issuer keytracker.Identity
	if cmd.issuer == "" {
		cmd.issuer = t.Value(parsers.X509IssuerDN)
	}
	// check if its self signed
	if cmd.issuer == t.Value(parsers.X509Subject) {
		// create a psudo Identity, with just the key
		issuer, err = keytracker.NewIdentity(nil, nil)
	} else {
		issuer, err = cmd.getIssuer(ctx)
	}
	if err != nil {
		return err
	}

	if issuer == nil {
		// no issuer given in flag or template
	}
	// If encrypted, get the password
	if issuer.Key().IsEncrypted() {
		cmd.keyPass, err = cmd.getKeyPass()
		if err != nil {
			return err
		}
	}
	by, err := templates.IssueCertificate(issuer, cmd.keyPass, t)
	if err != nil {
		return err
	}
	return pemwriter.NewPemFormat(out).Write(&pem.Block{
		Type:  keytools.PEM_CERTIFICATE,
		Bytes: by,
	})
}

// getIssuer attempts to find the issuer certificate, based on the issuer DN name.
// If no Identidy is found and DN (issuer) given, returns error, not found
// If no Identidy is found and no DN (issuer) given, returns nil
// If one found, returns that.
// If more than one, presents a list and prompts to select (or error if quiet flag set)
func (cmd *IssueCommand) getIssuer(ctx context.Context) (keytracker.Identity, error) {
	var kp []string
	if cmd.keyPass != "" {
		kp = strings.Split(cmd.keyPath, ":")
	} else {
		kp = GetKeyPath(nil)
	}

	issuers := issuers(ctx, kp, true, cmd.issuer)
	if len(issuers) == 0 {
		if cmd.issuer != "" {
			return nil, fmt.Errorf("could not locate the issuer %s", cmd.issuer)
		}
		return nil, nil
	}

	if len(issuers) == 1 {
		return issuers[0], nil
	}
	// more than one issuer, present list
	if cmd.quiet {
		return nil, fmt.Errorf("no single issuer was found to sign the certificate.  Be more specific with the issuer flag.")
	}
	// Choose from the list of two or more issuers
	issuers = sortIssuers(issuers)
	names := make([]string, len(issuers))
	for i, id := range issuers {
		names[i] = id.String()
	}
	index := PromptChooseList("Select the issuing certificate for the new certificate or zero to abort:", names)
	if index < 0 {
		return nil, fmt.Errorf("aborted")
	}
	return issuers[index], nil
}

// getKeyPass gets a password for an encrypted private key.
// If cmd.keyPass already set with a flag, that is returned.
// Otherwise user is prompted for input, unless cmd.queit flag set, in which case an error is rasied for the missing password.
func (cmd *IssueCommand) getKeyPass() (string, error) {
	if cmd.keyPass != "" {
		return cmd.keyPass, nil
	}
	if cmd.quiet {
		return "", fmt.Errorf("encrypted issuer key requires password")
	}
	return PromptPassword(fmt.Sprintf("Enter the password for the issuer '%s' private key:", cmd.issuer))
}

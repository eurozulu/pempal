package cmd

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"os"
	"pempal/keytools"
	"pempal/keytracker"
	"pempal/pemwriter"
	"pempal/templates"
	"pempal/templates/parsers"
	"strings"
)

type IssueCommand struct {
	issuer  string
	key     string
	keyPass string
	keyPath string
	quiet   bool
}

func (cmd *IssueCommand) Description() string {
	lines := []string{"<template name or existing certificate/csr> ...[<template name or existing certificate/csr>]"}
	lines = append(lines, "issues a new X509 certificate, based on the given templates/resources")
	lines = append(lines, "Properties from each resource are merged into one, from left to right, with the right most taking precedence")
	lines = append(lines, "At least one resource should contain the DN of the issuer of the certificare unless specified with the -issuer flag")
	lines = append(lines, "-issuer flag takes precedence over all templates and resource properties.")
	lines = append(lines, "If no issuer is provided, a prompt to select one is presented.")
	return strings.Join(lines, "\n")
}

func (cmd *IssueCommand) Flags(f *flag.FlagSet) {
	f.BoolVar(&cmd.quiet, "q", false, "surpress confirmation prompts")
	f.StringVar(&cmd.issuer, "issuer", "", "set the DN of the issuer for the new certificate. Overrides any template value")
	f.StringVar(&cmd.key, "key", "$KEYFILE", "Specify the public key to use. Can be a filename or a sha256 hash")
	f.StringVar(&cmd.keyPass, "password", "", "Specify the password for an encrypted issuer key. Will prompt is required and not provided")
	f.StringVar(&cmd.keyPath, "keypath", "", "comma delimited list of directories to search for issuer keys.  Overrides KEYPATH environment variable")
}

func (cmd *IssueCommand) Run(ctx context.Context, out io.Writer, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("must provide a template name or path to existing resource to use as the certificate values.")
	}
	// Build the template to create cert on
	tb := templates.NewBuilder()
	if err := tb.Add(args...); err != nil {
		return err
	}
	t, err := tb.Build()
	if err != nil {
		return err
	}

	// Establish the signing key
	// If issuer not specified as flag, take it from the template (which might also be empty!)
	if cmd.issuer == "" {
		cmd.issuer = t.Value(parsers.X509IssuerDN)
	}

	issuer, err := cmd.getIssuer(ctx)
	if err != nil {
		return fmt.Errorf("failed to establish a key for issuer %s  %w", cmd.issuer, err)
	}
	// If encrypted, establish the password
	if issuer.Key().IsEncrypted() {
		cmd.keyPass, err = cmd.getKeyPass()
		if err != nil {
			return err
		}
	}
	by, err := templates.GenerateCertificate(issuer, t, cmd.keyPass)
	if err != nil {
		return err
	}
	return pemwriter.NewPemFormat(out).Write(&pem.Block{
		Type:  keytools.PEM_CERTIFICATE,
		Bytes: by,
	})
}

func (cmd *IssueCommand) getIssuerByKey() (keytracker.Identity, error) {
	return nil, fmt.Errorf("key is not yet implemented")
}

func (cmd *IssueCommand) getIssuer(ctx context.Context) (keytracker.Identity, error) {
	if cmd.key != "" {
		return cmd.getIssuerByKey()
	}
	// Search keypath for matching issuer. If flag set, use that, otherwise CWD and any keypath env.
	issuers := collectIssuers(ctx, cmd.issuer, cmd.getKeyPath())
	if len(issuers) == 0 {
		return nil, fmt.Errorf("no %s private keys found to issue a new certificate.  "+
			"Set the $%s to a colon delimited list of directories containing the private keys and CA certificates sign certificates.",
			cmd.issuer, ENV_KeyPath)
	}

	// establish the certificates to sign with
	// Build list of all certificates and index map keeping track of the identity each cert belongs to
	var certs []*x509.Certificate
	var index = map[*x509.Certificate]keytracker.Identity{}
	ids := sortIssuers(issuers)
	for _, id := range ids {
		cs := sortCerts(id.Certificates(0, 0))
		certs = append(certs, cs...)
		for _, c := range cs {
			index[c] = id
		}
	}
	// If just one certificate, return that ID
	if len(certs) == 1 {
		return index[certs[0]], nil
	}
	if cmd.quiet {
		return nil, fmt.Errorf("multiple keys match the issuer '%s'. Must be more specific or use -key to select a unique key", cmd.issuer)
	}
	// propmt user to select a certificate
	names := certificateNames(certs)
	ci := PromptChooseList("Select the issuer to sign the new certificate:", names)
	if ci < 0 {
		return nil, fmt.Errorf("aborted")
	}
	id := index[certs[ci]]
	// return new Identity containing just the relevant certificate
	return keytracker.NewIdentity(id.Key(), []*pem.Block{&pem.Block{
		Type:  keytools.PEM_CERTIFICATE,
		Bytes: certs[ci].Raw,
	}}), nil
}

func (cmd *IssueCommand) getKeyPath() []string {
	if cmd.keyPath != "" {
		return strings.Split(cmd.keyPath, ":")
	}
	return append([]string{os.ExpandEnv("$PWD")}, strings.Split(KeyPath, ":")...)
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
	return PromptPassword(fmt.Sprintf("Enter the password for the '%s' issuer private key:", cmd.issuer))
}

func certificateNames(certs []*x509.Certificate) []string {
	names := make([]string, len(certs))
	for i, c := range certs {
		names[i] = c.Subject.String()
	}
	return names
}

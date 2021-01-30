package main

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"github.com/eurozulu/pempal"
	"github.com/eurozulu/pempal/encoding"
	"github.com/eurozulu/pempal/templates"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// IssueCommand generates new, signed certificates
type IssueCommand struct {
	PublicKey string `flag:"publickey,key,k,puk"`

	// Issuer should be the path to the Issuing CA certificate.
	// When not set, the path is serached for all CA certs and non script execution prompts for a selection.
	// Use a '-' dash, to indicate a self signed certificate, where the source is used for the issuer.
	Issuer string `flag:"issuer,i"`
	// issuerKey should be the
	IssuerKey        string `flag:"issuerkey,ik,prk"`
	IssuerPassphrase string `flag:"issuerpassphrase,ip,pwd"`

	CommonName string `flag:"commonname,cn"`

	Template string `flag:"template, t"`

	OutPath string `flag:"outpath,out,o"`
	Encode  string `flag:"encode,e"`

	Script bool `flag:"script,s"`
}

// Issue generates a newly signed certificate(s).
// issue <source> [<path>[, ...]]
// source is required
// Source may be:
// A Certificate request (CSR)
// An existing Certificate
// A '-' to enter certificate details from the console
// Properties from the source are copied into the new certificate including any existing
// public key and signature.
// The New certificate will be issued by the named issuer flag, using the named issuer key.
// Output is written to stdout, unless -outpath is set.
// A -template flag may specify an existing template.  All properties in the given template will override proeprties in the source.
func (ic IssueCommand) Issue(src string, args ...string) error {
	if src == "" {
		return fmt.Errorf("must specify a source such as a CSR, template or existing certificate")
	}
	if len(args) == 0 {
		args = append(args, ".")
	}
	tc, err := ic.sourceCertificate(src)
	if err != nil {
		return fmt.Errorf("invalid source  %v", err)
	}

	puk, err := ic.publicKey(tc)
	if err != nil {
		return fmt.Errorf("No public key availabe  %v", err)
	}
	tc.PublicKey = puk

	var ca *x509.Certificate
	if ic.Issuer == "" {
		ic.Issuer = tc.Issuer.CommonName
	}
	if ic.Issuer == "-" { // self signed
		ca = tc
	} else {
		ca, err = ic.issuerCertificate(ic.Issuer, args...)
		if err != nil {
			return fmt.Errorf("invalid issuer  %v", err)
		}
	}
	tc.Issuer = ca.Subject

	prk, err := ic.issuerKey(ca, args...)
	if err != nil {
		return fmt.Errorf("invalid issuer  %v", err)
	}

	if ic.CommonName != "" {
		tc.Subject.CommonName = ic.CommonName
	}
	if tc.Subject.CommonName == "" {
		if ic.Script {
			return fmt.Errorf("No common name set")
		}
		cn := PromptInput("Enter the Common Name for the new certificate: ", "")
		if cn == "" {
			return fmt.Errorf("no common name, aborted")
		}
		tc.Subject.CommonName = cn
	}
	if !ic.Script {
		if err := ic.confirmCreate(tc); err != nil {
			return err
		}
	}

	der, err := x509.CreateCertificate(rand.Reader, tc, ca, tc.PublicKey, prk)
	if err != nil {
		return fmt.Errorf("Failed to create certificate  %v", err)
	}

	ct := &templates.CertificateTemplate{}
	if err := ct.UnmarshalBinary(der); err != nil {
		return err
	}
	return ic.writeCertificate(ct)
}

func (ic IssueCommand) sourceCertificate(src string) (*x509.Certificate, error) {
	if src == "-" {
		return &x509.Certificate{}, nil
	}

	tps, err := ic.findResources([]string{"cer", "csr"}, "", "", []string{src})
	if err != nil {
		return nil, err
	}

	t, err := ic.selectSourceList(tps)
	if err != nil {
		return nil, err
	}

	var ct *templates.CertificateTemplate
	switch v := t.(type) {
	case *templates.CertificateTemplate:
		ct = v
	case *templates.CSRTemplate:
		ct = &templates.CertificateTemplate{}
		ct.CopyCSR(v)
	default:
		return nil, fmt.Errorf("%s is an unexpected resource type for source", templates.TemplateType(t))
	}

	by, err := ct.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(by)
}

func (ic IssueCommand) publicKey(c *x509.Certificate) (crypto.PublicKey, error) {
	puk := c.PublicKey
	if ic.PublicKey != "" || puk == nil {
		p := ic.PublicKey
		if p == "" {
			p = "."
		}
		tps, err := ic.findResources([]string{"puk", "cer"}, "", "", []string{p})
		if err != nil {
			return nil, err
		}
		tp, err := ic.selectPublicKey(tps)
		if err != nil {
			return nil, err
		}

		switch v := tp.(type) {
		case *templates.PublicKeyTemplate:
			puk = v.Key()
		case *templates.CertificateTemplate:
			puk = v.Certificate().PublicKey
		default:
			return nil, fmt.Errorf("unextencted key type %s", templates.TemplateType(tp))
		}
	}
	return puk, nil
}

func (ic IssueCommand) issuerCertificate(iname string, args ...string) (*x509.Certificate, error) {
	var q string
	if iname != "" {
		q = fmt.Sprintf("Common Name: %s,", iname)
	}
	q = strings.Join([]string{q, "IsCA: true"}, "")

	tps, err := ic.findResources([]string{"cer"}, "IsCA: true", "", args)
	if err != nil {
		return nil, err
	}
	t, err := ic.selectIssuerList(tps)
	if err != nil {
		return nil, err
	}
	by, err := t.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(by)
}

func (ic IssueCommand) issuerKey(ca *x509.Certificate, args ...string) (crypto.PrivateKey, error) {
	ipuk := ca.PublicKey
	if ipuk == nil {
		return nil, fmt.Errorf("issuer certificate has no public key")
	}

	tps, err := ic.findResources([]string{"prk"}, ic.IssuerKey, "", args)
	if err != nil {
		return nil, err
	}
	for _, kt := range tps {
		pkt, ok := kt.(*templates.PrivateKeyTemplate)
		if !ok {
			continue
		}
		if pkt.Encrypted {
			if ic.IssuerPassphrase == "" {
				if ic.Script {
					return nil, fmt.Errorf("issuer key is encrypted and no passphrase provided")
				}
				ic.IssuerPassphrase, err = PromptPassword(fmt.Sprintf("Enter the private key passphrase for issuer %s:", ca.Subject.CommonName))
				if err != nil {
					return nil, err
				}
			}
			if err := pkt.Decrypt(ic.IssuerPassphrase); err != nil {
				return nil, err
			}
		}
		prk := pkt.Key()
		puk := pempal.PublicKeyFromPrivate(prk)
		if !pempal.ComparePublicKeys(puk, ipuk) {
			continue
		}
		return pkt.Key(), nil
	}
	return nil, fmt.Errorf("private key for CA %s could not be found", ca.Subject.CommonName)
}

func (ic IssueCommand) writeCertificate(ct *templates.CertificateTemplate) error {
	out := os.Stdout
	if ic.OutPath != "" {
		f, err := os.OpenFile(ic.OutPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Println(err)
			}
		}()
		out = f
	}
	if ic.Encode == "" {
		ic.Encode = "pem"
	}
	enc, err := encoding.NewEncoder(ic.Encode, out)
	if err != nil {
		return err
	}
	return enc.Encode([]templates.Template{ct})
}

func (ic IssueCommand) findResources(types []string, query, pwd string, paths []string) ([]templates.Template, error) {
	fc := &FindCommand{
		Recursive:   true,
		Query:       "",
		Insensitive: false,
		Type:        types,
		Password:    pwd,
	}
	ftps, err := fc.FindAllTemplates(paths...)
	if err != nil {
		return nil, err
	}
	tps := make([]templates.Template, len(ftps))
	for i, ft := range ftps {
		tps[i] = ft.Template
	}
	return tps, nil
}

func (ic IssueCommand) selectIssuerList(tps []templates.Template) (templates.Template, error) {
	if len(tps) == 0 {
		return nil, fmt.Errorf("no issuer found. Specify a common name in the -issuer flag")
	}

	if len(tps) == 1 {
		return tps[0], nil
	}
	if ic.Script {
		return nil, fmt.Errorf("multiple issuers found")
	}
	names := make([]string, len(tps))
	for i, t := range tps {
		names[i] = t.String()
	}
	i := PromptChoice("Select the issuer for the new certificate.", names)
	if i < 0 {
		return nil, fmt.Errorf("aborted")
	}
	return tps[i], nil
}

func (ic IssueCommand) selectPublicKey(tps []templates.Template) (templates.Template, error) {
	if len(tps) == 0 {
		return ic.createNewKey()
	}
	if len(tps) == 1 {
		return tps[0], nil
	}
	if ic.Script {
		return nil, fmt.Errorf("multiple public keys found")
	}
	names := make([]string, len(tps))
	for i, t := range tps {
		names[i] = t.String()
	}
	i := PromptChoice("Select the public key for the new certificate.", names)
	if i < 0 {
		return nil, fmt.Errorf("aborted")
	}
	return tps[i], nil
}

func (ic IssueCommand) selectSourceList(tps []templates.Template) (templates.Template, error) {
	if len(tps) == 0 {
		return nil, fmt.Errorf("no templates found")
	}
	if len(tps) == 1 {
		return tps[0], nil
	}
	if ic.Script {
		return nil, fmt.Errorf("multiple sources found")
	}
	names := make([]string, len(tps))
	for i, t := range tps {
		names[i] = t.String()
	}
	i := PromptChoice("Select the source of the new certificate.", names)
	if i < 0 {
		return nil, fmt.Errorf("aborted")
	}
	return tps[i], nil
}

func (ic IssueCommand) confirmCreate(c *x509.Certificate) error {
	ct := templates.NewCertificateTemplate(c)
	for {
		by, err := yaml.Marshal(ct)
		if err != nil {
			return err
		}
		fmt.Println(string(by))
		if PromptConfirm("Create this new certificate", false) {
			return nil
		}
		if err := EditFields(ct); err != nil {
			return err
		}
	}
}

func (ic IssueCommand) createNewKey() (*templates.PublicKeyTemplate, error) {
	if ic.Script {
		return nil, fmt.Errorf("no public keys found")
	}
	var err error
	n := PromptInput("Enter a name for the new key: ", "./certificatekey.pem")
	n, err = filepath.Abs(n)
	if err != nil {
		return nil, err
	}

	pwd, err := PromptCreatePassword("Enter a new password for the key.  (Hit enter for unencrypted key)", 0)
	if err != nil {
		return nil, err
	}
	mkk := MakeKeyCommand{
		NoPublicKey: false,
		Encrypt:     true,
		Script:      false,
		Password:    pwd,
	}
	np := strings.Join([]string{n, "pub"}, ".")
	err = mkk.MakeKey(n, np)
	if err != nil {
		return nil, err
	}
	by, err := ioutil.ReadFile(np)
	if err != nil {
		return nil, err
	}
	tps, err := encoding.ParseTemplates(np, by, pwd)
	if err != nil {
		return nil, err
	}
	if len(tps) < 1 {
		return nil, fmt.Errorf("failed to reload public key %s", np)
	}
	kt, ok := tps[0].(*templates.PublicKeyTemplate)
	if !ok {
		return nil, fmt.Errorf("failed to reload public key %s", np)
	}

	return kt, nil
}

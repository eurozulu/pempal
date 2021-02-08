package main

import (
	"bytes"
	"encoding/pem"
	"fmt"
	"github.com/pempal/pemio"
	"github.com/pempal/pempal"
	"github.com/pempal/templates"
	"log"
	"time"
)

// IssueCommand generates new, signed certificates
type IssueCommand struct {
	Command `yaml:"-" flag:"-"`

	Version               int                                   `yaml:"Version,omitempty" flag:"version"`
	Subject               map[string]interface{}                `yaml:"Subject,omitempty" flag:"subject"`
	PublicKey             *templates.PublicKeyTemplate          `yaml:"PublicKey,omitempty" flag:"publickey"`
	PublicKeyAlgorithm    *templates.PublicKeyAlgorithmTemplate `yaml:"PublicKeyAlgorithm,omitempty" flag:"publickeyalgorithmtemplate,pka"`
	SignatureAlgorithm    *templates.SignatureAlgorithmTemplate `yaml:"SignatureAlgorithm,omitempty" flag:"signaturealgorithm,sa"`
	DNSNames              []string                              `yaml:"DNSNames,omitempty" flag:"dnsnames,dns"`
	EmailAddresses        []string                              `yaml:"EmailAddresses,omitempty" flag:"emailaddresses,emails"`
	IPAddresses           []*templates.IPAddressTemplate        `yaml:"IPAddresses,omitempty" flag:"ipaddresses,ipa,ipas"`
	URIs                  []*templates.URIsTemplate             `yaml:"URIs,omitempty" flag:"uri,uris"`
	SerialNumber          string                                `yaml:"SerialNumber,omitempty" flag:"serialnumber"`
	SubjectKeyId          string                                `yaml:"SubjectKeyId,omitempty" flag:"subjectkeyid"`
	AuthorityKeyId        string                                `yaml:"AuthorityKeyId,omitempty" flag:"authoritykeyid"`
	NotBefore             time.Time                             `yaml:"NotBefore,omitempty" flag:"notbefore"`
	NotAfter              time.Time                             `yaml:"NotAfter,omitempty" flag:"notafter"`
	BasicConstraintsValid bool                                  `yaml:"BasicConstraintsValid,omitempty" flag:"basicconstraintsvalid,bcv"`
	IsCA                  bool                                  `yaml:"IsCA" flag:"isca,ca"`
	MaxPathLen            int                                   `yaml:"MaxPathLen,omitempty" flag:"maxpathlen"`
	MaxPathLenZero        bool                                  `yaml:"MaxPathLenZero" flag:"maxpathlenzero"`
	Signature             string                                `yaml:"Signature,omitempty" flag:"signature"`
	KeyUsage              *templates.KeyUsageTemplate           `yaml:"KeyUsage,omitempty" flag:"keyusage"`
	ExtKeyUsage           []*templates.ExtKeyUsagesTemplate     `yaml:"ExtKeyUsage,omitempty" flag:"extkeyusage"`
	Extensions            []*templates.ExtensionsTemplate       `yaml:"Extensions,omitempty" flag:"extensions"`
	ExtraExtensions       []*templates.ExtensionsTemplate       `yaml:"ExtraExtensions,omitempty" flag:"extraextensions"`
	CRLDistributionPoints string                                `yaml:"CRLDistributionPoints,omitempty" flag:"crldistributionpoints,crlpoints"`

	// Issuer should point to the CA certificate issuing this certificate.
	// Can be a filename, common name or any 'query' value which identifies the certificate uniquely
	// When not set, user asked to select from all available CA's
	// Use '-' dash, to indicate a self signed certificate, where the source is used for the issuer.
	Issuer string `yaml:"Issuer" flag:"issuer,i"`

	// Issuer key passphrase.  When not given and private key encrypted, user prompted for password.
	IssuerPassphrase string `yaml:"-" flag:"issuerpassphrase,ip,pwd"`

	// Script flag surpresses any user prompts, invoking an error if the information is unabailable.
	Script bool `yaml:"-" flag:"script,s"`

	// keyOut contains any key read in from a call to request new csr.
	// Key exists in memory only and written out with the certificate.
	keyOut *templates.PrivateKeyTemplate
}

// Issue will generate a new, signed certificate using the properties from the given templates
func (ic IssueCommand) Issue(tps ...string) error {
	ts, err := NewTemplateFiles(tps)
	if err != nil {
		return err
	}

	// Create a new Cert and apply the templates
	nct := &templates.CertificateTemplate{}
	if err := templates.ApplyTemplates(nct, ts...); err != nil {
		return err
	}
	// Apply this command as a template to assign any matching flag values
	if err := templates.ApplyTemplates(nct, &ic); err != nil {
		return err
	}

	if err := ic.ensureCertSigned(nct); err != nil {
		return fmt.Errorf("certificate is invalid  %v", err)
	}

	// issuer cert
	isc, err := ic.issuerCertificate(nct)
	if err != nil {
		return err
	}
	nct.Issuer = isc.Subject

	isk, err := ic.issuerKey(isc)
	if err != nil {
		return err
	}

	if !ic.Script {
		if err := ConfirmTemplate("create new certificate", nct); err != nil {
			return err
		}
	}

	bl, err := pempal.NewCertificate(nct, isc, isk)
	if err != nil {
		return fmt.Errorf("Failed to sign new certificate  %v", err)
	}

	if err := ic.WriteOutput([]*pem.Block{bl}, 0644); err != nil {
		return err
	}

	if ic.keyOut != nil {
		if err := ic.WriteOutput([]*pem.Block{ic.keyOut.PEMBlock()}, 0600); err != nil {
			return err
		}
	}
	return nil
}

func (ic IssueCommand) String() string {
	return "issues new certificate"
}

// ensureCertSigned ensures the given template has a public key and signature.
// If not present, and not in Script, generates a new request
func (ic IssueCommand) ensureCertSigned(ct *templates.CertificateTemplate) error {
	if ic.PublicKey != nil && ic.Signature != "" {
		return nil
	}
	if ic.Script {
		return fmt.Errorf("No signature / public key found")
	}
	csr, prk, err := ic.requestNewCSR(ct)
	if err != nil {
		return err
	}
	if err := templates.ApplyTemplate(ct, csr); err != nil {
		return err
	}
	ic.keyOut = prk
	return nil
}

func (ic IssueCommand) issuerCertificate(ct *templates.CertificateTemplate) (*templates.CertificateTemplate, error) {
	// is it self signed
	if ic.Issuer == "-" {
		return ct, nil
	}
	// If issuer not stated, take it from the cert template
	is := ic.Issuer
	if is == "" {
		is = ct.Issuer.CommonName
	}
	iscs, err := pempal.CACertificates(is)
	if err != nil {
		return nil, fmt.Errorf("failed to find an issuer CA certificate  %v", err)
	}
	if len(iscs) == 0 {
		return nil, fmt.Errorf("failed to find an issuer CA certificate")
	}
	if ic.Script {
		if len(iscs) == 1 {
			return templates.NewCertificateTemplate(iscs[0].Block)
		}
		return nil, fmt.Errorf("No single CA certificate found.  Found %d certificates", len(iscs))
	}

	i := ChooseTemplate("Select the CA to issue the certificate", iscs, nil)
	if i < 0 {
		return nil, fmt.Errorf("aborted.  No CA certificate")
	}
	return templates.NewCertificateTemplate(iscs[i].Block)
}

func (ic IssueCommand) issuerKey(isc *templates.CertificateTemplate) (*templates.PrivateKeyTemplate, error) {
	q := fmt.Sprintf("PublicKeyFingerprint: %s", templates.PublicKeyFingerprint(isc.PublicKey))
	qrs, err := pempal.FindKeys(q, ic.IssuerPassphrase)
	if err != nil {
		return nil, err
	}
	if len(qrs) == 0 {
		return nil, fmt.Errorf("no issuer key found")
	}
	if len(qrs) == 1 {
		return templates.NewPrivateKeyTemplate(qrs[0].Block), nil
	}
	if ic.Script {
		return nil, fmt.Errorf("no single issuer key found, found %d keys", len(qrs))
	}
	index := ChooseTemplate("Select the key to issue this certificate", qrs, nil)
	if index < 0 {
		return nil, fmt.Errorf("aborted.  No issuer key")
	}
	return templates.NewPrivateKeyTemplate(qrs[index].Block), nil
}

func (ic *IssueCommand) requestNewCSR(t templates.Template) (*templates.RequestTemplate, *templates.PrivateKeyTemplate, error) {
	buf := bytes.NewBuffer(nil)
	kc := &RequestCommand{
		Command: Command{Output: buf},
		KeyOut:  true,
	}
	if err := templates.ApplyTemplate(kc, t); err != nil {
		return nil, nil, err
	}
	if err := kc.Request(); err != nil {
		return nil, nil, err
	}
	bls, err := pemio.ReadPEMs(buf)
	if err != nil {
		return nil, nil, err
	}
	var rt *templates.RequestTemplate
	var prk *templates.PrivateKeyTemplate

	for _, bl := range bls {
		t, err := templates.NewTemplate(bl)
		if err != nil {
			log.Println(err)
			continue
		}
		switch v := t.(type) {
		case *templates.PrivateKeyTemplate:
			prk = v
		case *templates.RequestTemplate:
			rt = v
		default:
			log.Printf("unexpected pem reading back Request %v", v)
		}
	}
	if rt == nil {
		return nil, nil, fmt.Errorf("no request temaplte found back from Request command")
	}
	return rt, prk, nil
}

package encoding

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/templates"
	"golang.org/x/crypto/pkcs12"
	"gopkg.in/yaml.v3"
	"log"
	"path"
	"strings"
)

var ErrorUnknownFormat = fmt.Errorf("unknown format")
var ErrorNotATemplate = fmt.Errorf("not a template")

// ParseTemplates attempts to parse the given bytes into one or more Templates, representing x509 resources.
func ParseTemplates(p string, by []byte, pwd string) ([]templates.Template, error) {
	switch path.Ext(p) {
	case ".yaml":
		return parseYaml(p, by)

	case ".pem":
		return parsePEMs(p, by, pwd)

	case ".p12", ".pfx":
		return parsePKCS12(p, by, pwd)

	case ".der", ".crt", "cer":
		return parseBinary(p, by)
	}

	// Unknown file extension, check if contents look like a PEM
	s := string(by)
	if strings.Contains(s, "----BEGIN ") &&
		strings.Contains(s, "-----END") {
		return parsePEMs(p, by, pwd)
	}
	// Unknown file type and not a pem.  Attempt to parse as a binary der
	return parseBinary(p, by)
}

func parsePEMs(p string, by []byte, pwd string) ([]templates.Template, error) {
	var tpls []templates.Template
	var index int
	data := by
	for {
		bl, r := pem.Decode(data)
		if bl == nil {
			break
		}
		fp := p
		if index > 0 {
			fp = fmt.Sprintf("%s:#%02d", fp, index)
		}
		index++

		tp, err := templates.NewTemplate(fp, bl.Type)
		if err != nil {
			return nil, err
		}

		if strings.Contains(string(by), "ENCRYPTED") {
			if err := tp.UnmarshalBinary(by); err != nil {
				return nil, err
			}
		} else {
			if err := tp.UnmarshalBinary(bl.Bytes); err != nil {
				return nil, err
			}
		}
		if pwd != "" {
			ecPem, ok := tp.(EncryptedPEM)
			if ok {
				err = ecPem.Decrypt(pwd)
				if err != nil {
					log.Println(err)
				}
			}
		}
		tpls = append(tpls, tp)
		data = r
	}
	return tpls, nil
}

// Parse the given byte block as Binary.
// attempts to parse as each type until one is successful
func parseBinary(p string, by []byte) ([]templates.Template, error) {

	// Try to parse as certificates
	if certs, err := x509.ParseCertificates(by); err == nil {
		var tpls []templates.Template
		for i, c := range certs {
			fp := p
			if i > 0 {
				fp = fmt.Sprintf("%s:#%2d", fp, i)
			}
			cr := &templates.CertificateTemplate{FilePath: fp}
			if err := cr.UnmarshalBinary(c.Raw); err != nil {
				return nil, err
			}
			tpls = append(tpls, cr)
		}
		return tpls, nil
	}

	// trey CSR
	csr := &templates.CSRTemplate{FilePath: p}
	if err := csr.UnmarshalBinary(by); err == nil {
		return []templates.Template{csr}, nil
	}

	// Try to parse as a public key
	pubK := &templates.PublicKeyTemplate{FilePath: p}
	if err := pubK.UnmarshalBinary(by); err == nil {
		return []templates.Template{pubK}, nil
	}

	// try as ssh public key
	sspubK := &templates.SSHPublicKeyTemplate{PublicKeyTemplate: *pubK}
	if err := sspubK.UnmarshalBinary(by); err == nil {
		return []templates.Template{sspubK}, nil
	}

	// Try as private key
	prk := &templates.PrivateKeyTemplate{FilePath: p}
	if err := prk.UnmarshalBinary(by); err == nil {
		return []templates.Template{prk}, nil
	}
	// ssh private key
	sshPrk := templates.SSHPrivateKeyTemplate{PrivateKeyTemplate: *prk}
	if err := sshPrk.UnmarshalBinary(by); err == nil {
		return []templates.Template{sspubK}, nil
	}

	// Try as CRL
	crl := &templates.CRLTemplate{FilePath: p}
	if err := crl.UnmarshalBinary(by); err == nil {
		return []templates.Template{crl}, nil
	}

	return nil, ErrorUnknownFormat
}

func parseYaml(p string, by []byte) ([]templates.Template, error) {
	n := readYamlName(by)
	if n == "" {
		return nil, ErrorNotATemplate
	}
	t, err := templates.NewTemplate(p, readYamlName(by))
	if err != nil {
		return nil, fmt.Errorf("failed to read template from '%s'  %v", p, err)
	}
	if err := yaml.Unmarshal(by, t); err != nil {
		return nil, err
	}
	return []templates.Template{t}, nil
}

func readYamlName(by []byte) string {
	i := strings.Index(string(by), "Template:")
	if i < 0 {
		return ""
	}
	s := strings.SplitN(string(by[i:]), "\n", 2)
	if len(s) == 0 {
		return ""
	}
	st := strings.TrimSpace(s[0])
	st = strings.TrimPrefix(st, "\"")
	st = strings.TrimSuffix(st, "\"")
	return st
}

func parsePKCS12(p string, by []byte, pwd string) ([]templates.Template, error) {
	pems, err := pkcs12.ToPEM(by, pwd)
	if err != nil {
		if !strings.Contains(err.Error(), "decryption password incorrect") {
			return nil, err
		}
		return []templates.Template{&templates.PKCS12Template{
			FilePath:    p,
			IsEncrypted: true,
		},
		}, nil
	}

	var pemBytes []byte
	for _, b := range pems {
		pemBytes = append(pemBytes, pem.EncodeToMemory(b)...)
	}
	return parsePEMs(p, pemBytes, pwd)
}

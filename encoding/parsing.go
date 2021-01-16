package encoding

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/templates"
	"gopkg.in/yaml.v3"
	"strings"
)

var ErrorUnknownFormat = fmt.Errorf("unknown format")

func ParseTemplates(p string, by []byte) ([]templates.Template, error) {
	if strings.HasSuffix(p, ".yaml") {
		if readYamlName(by) != "" {
			parseYaml(p, by)
		}
		return nil, ErrorUnknownFormat
	}

	s := string(by)
	if strings.Contains(s, "----BEGIN ") &&
		strings.Contains(s, "-----END") {
		return parsePEMs(p, by)
	}
	return parseBinary(p, by)
}

func parsePEMs(p string, by []byte) ([]templates.Template, error) {
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

		tpPem, ok := tp.(PEMUnmarshaler)
		if !ok {
			return nil, fmt.Errorf("template %s does not support PEM encoded files", bl.Type)
		}
		if err := tpPem.UnmarshalPEM(bl); err != nil {
			return nil, err
		}
		tpls = append(tpls, tp)
		data = r
	}
	return tpls, nil
}

// Parse the given byte block as Binary.
// attempts to parse as each type until one is successful
func parseBinary(p string, by []byte) ([]templates.Template, error) {
	certs, err := x509.ParseCertificates(by)
	if err == nil {
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
	sspubK := &templates.SSHPublicKeyTemplate{FilePath: p}
	err = sspubK.UnmarshalBinary(by)
	if err == nil {
		return []templates.Template{sspubK}, nil
	}
	pubK := &templates.PublicKeyTemplate{FilePath: p}
	err = pubK.UnmarshalBinary(by)
	if err == nil {
		return []templates.Template{pubK}, nil
	}

	csr := &templates.CSRTemplate{FilePath: p}
	err = csr.UnmarshalBinary(by)
	if err == nil {
		return []templates.Template{csr}, nil
	}

	priK := &templates.PrivateKeyTemplate{FilePath: p}
	err = priK.UnmarshalBinary(by)
	if err == nil {
		return []templates.Template{priK}, nil
	}
	crl := &templates.CRLTemplate{FilePath: p}
	err = crl.UnmarshalBinary(by)
	if err == nil {
		return []templates.Template{crl}, nil
	}
	return nil, ErrorUnknownFormat
}

func parseYaml(p string, by []byte) ([]templates.Template, error) {
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
	if strings.HasPrefix(st, "\"") {
		st = st[1:]
	}
	if strings.HasSuffix(st, "\"") {
		st = st[0 : len(st)-1]
	}
	return st
}

package templates

import (
	"crypto/sha256"
	"encoding"
	"fmt"
	"reflect"
	"strings"
)

var AllTemplatTypes = []string{
	"CERTIFICATE",
	"CERTIFICATE REQUEST",
	"X509 CRL",
	"SSH PUBLIC KEY",
	"PUBLIC KEY",
	"OPENSSH PRIVATE KEY",
	"PRIVATE KEY",
}

type Template interface {
	String() string
	Location() string
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

func NewTemplate(p string, tt string) (Template, error) {
	switch strings.ToUpper(tt) {
	case "CERTIFICATE":
		return &CertificateTemplate{FilePath: p}, nil
	case "CERTIFICATE REQUEST":
		return &CSRTemplate{FilePath: p}, nil
	case "PUBLIC KEY":
		return &PublicKeyTemplate{FilePath: p}, nil
	case "SSH PUBLIC KEY":
		pkt := PublicKeyTemplate{FilePath: p}
		return &SSHPublicKeyTemplate{PublicKeyTemplate: pkt}, nil
	case "PRIVATE KEY":
		return &PrivateKeyTemplate{FilePath: p}, nil
	case "OPENSSH PRIVATE KEY":
		return &SSHPrivateKeyTemplate{PrivateKeyTemplate{FilePath: p}}, nil
	case "X509 CRL":
		return &CRLTemplate{FilePath: p}, nil
	default:
		return nil, fmt.Errorf("%s is an unknown template type\n", tt)
	}
}

func TemplateType(t Template) string {
	switch t.(type) {
	case *CertificateTemplate:
		return "CERTIFICATE"
	case *CSRTemplate:
		return "CERTIFICATE REQUEST"
	case *SSHPublicKeyTemplate:
		return "SSH PUBLIC KEY"
	case *PublicKeyTemplate:
		return "PUBLIC KEY"
	case *SSHPrivateKeyTemplate:
		return "OPENSSH PRIVATE KEY"
	case *PrivateKeyTemplate:
		return "PRIVATE KEY"
	case *CRLTemplate:
		return "X509 CRL"
	default:
		return ""
	}
}

func TemplateValues(t Template, flds []string) map[string]interface{} {
	tv := reflect.ValueOf(t).Elem()
	m := map[string]interface{}{}
	for _, fn := range flds {
		v := tv.FieldByName(fn)
		// not a known field
		if !v.IsValid() {
			continue
		}
		m[fn] = fmt.Sprintf("%v", v)
	}
	return m
}

func TemplateFields(t Template) []string {
	return readTypeFields(reflect.TypeOf(t), "")
}

func readTypeFields(tp reflect.Type, p string) []string {
	var flds []string
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}
	for i := 0; i < tp.NumField(); i++ {
		f := tp.Field(i)
		tag, ok := f.Tag.Lookup("yaml")
		if !ok || tag == "-" {
			continue
		}
		n := f.Name
		if p != "" {
			n = strings.Join([]string{p, n}, ".")
		}
		flds = append(flds, n)

		/*
			// If field type is a structure, recurse into that
			if f.Type.Kind() == reflect.Struct ||
				(f.Type.Kind() == reflect.Ptr && f.Type.Elem().Kind() == reflect.Struct){
				flds = append(flds, readTypeFields(f.Type, n)...)
			}*/
	}
	return flds
}

func fingerprint(by []byte) string {
	h := sha256.New()
	_, _ = h.Write(by)
	return string(h.Sum(nil))
}

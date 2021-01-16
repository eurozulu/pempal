package templates

import (
	"bytes"
	"encoding"
	"fmt"
	"reflect"
	"strings"
)


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
		return &SSHPublicKeyTemplate{FilePath: p}, nil
	case "PRIVATE KEY":
		return &PrivateKeyTemplate{FilePath: p}, nil
	case "OPENSSH PRIVATE KEY":
		return &SSHPrivateKeyTemplate{FilePath: p}, nil
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
		return  "CERTIFICATE REQUEST"
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

func TemplateString(t Template) string {
	flds := TemplateFields(t)
	vals := TemplateValues(t, flds)
	buf := bytes.NewBuffer(nil)
	for i, f := range flds {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(f)
		buf.WriteString(fmt.Sprintf("=%v", vals[f]))
	}
	return buf.String()
}

func TemplateValues(t Template, flds []string) map[string] interface{} {
	tv := reflect.ValueOf(t).Elem()
	m := map[string] interface{}{}
	for _, f := range flds {
		m[f] = tv.FieldByName(f).Interface()
	}
	return m
}

func TemplateFields(t Template) []string {
	var flds []string
	tp := reflect.TypeOf(t)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}
	for i := 0; i < tp.NumField(); i++ {
		f := tp.Field(i)
		if _, ok := f.Tag.Lookup("yaml"); !ok {
			continue
		}
		flds = append(flds, f.Name)
	}
	return flds
}
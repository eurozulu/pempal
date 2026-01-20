package templates

import (
	"fmt"
	"reflect"
	"strings"
)

// baseTemplates are the resource type templates, for each of the resources with a corresponding factory.
// Base types are:
// - privatekey
// - certificate
// - csr
// - crl

var baseTemplates = map[string]Template{
	"privatekey":  &PrivateKeyTemplate{},
	"certificate": &CertificateTemplate{},
	"csr":         &CertificateRequestTemplate{},
	"crl":         &RevocationListTemplate{},
}

var baseAliases = map[string]string{
	"p":           "privatekey",
	"prk":         "privatekey",
	"private-key": "privatekey",
	"c":           "certificate",
	"cert":        "certificate",
	"cer":         "certificate",
	"s":           "csr",
	"request":     "csr",
	"r":           "crl",
	"revoke":      "crl",
	"revokation":  "crl",
}

// BaseTemplateNames lists all the names of the base templates.
func BaseTemplateNames() []string {
	var names []string
	for name := range baseTemplates {
		names = append(names, name)
	}
	return names
}

func LookupBaseName(name string) (string, error) {
	name = strings.ToLower(name)
	if _, ok := baseTemplates[name]; ok {
		return name, nil
	}
	if an, ok := baseAliases[name]; ok {
		return an, nil
	}
	return "", fmt.Errorf("template %s is not a known base template", name)
}

func IsBaseTemplate(t Template) bool {
	return TypeOfBaseTemplate(t) != ""
}

func NewBaseTemplate(name string) (Template, error) {
	n, err := LookupBaseName(name)
	if err != nil {
		return nil, err
	}
	p := baseTemplates[n]
	return reflect.New(reflect.TypeOf(p).Elem()).Interface().(Template), nil
}

func TypeOfBaseTemplate(t Template) string {
	ts := reflect.TypeOf(t).String()
	for k, v := range baseTemplates {
		if v == nil {
			continue
		}
		if reflect.TypeOf(v).String() != ts {
			continue
		}
		return k
	}
	return ""
}

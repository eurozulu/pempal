package templates

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"pempal/keytools"
	"pempal/templates/parsers"
	"strings"
)

const KeyDelimiter = "."
const requiredPrefix = "?"
const funcPrefix = "{{"
const funcSuffix = "}}"

// Template represents an x509 resource in plain text (yaml) format
type Template map[string]interface{}

func (t Template) Contains(k string) bool {
	_, ok := t.valueIgnoreCase(k)
	return ok
}

func (t Template) Value(k string) string {
	v, ok := t.valueIgnoreCase(k)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s", v)
}

func (t Template) SetValue(k string, v interface{}) {
	ks := strings.Split(k, ".")
	m := t
	for len(ks) > 1 {
		cm := m.ValueMap(ks[0])
		if cm == nil {
			cm = Template{}
			m[ks[0]] = cm
		}
		m = cm
		ks = ks[1:]
	}
	m[ks[0]] = v
}
func (t Template) ValueMap(k string) map[string]interface{} {
	v, ok := t.valueIgnoreCase(k)
	if !ok {
		return nil
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	return m
}

func (t Template) String() string {
	buf := bytes.NewBuffer(nil)
	if err := yaml.NewEncoder(buf).Encode(&t); err != nil {
		log.Println(err)
		return ""
	}
	return buf.String()
}

func (t Template) RequiredNames() []string {
	return t.requiredNames("")
}

func (t Template) requiredNames(parent string) []string {
	var names []string
	for k, v := range t {
		if parent != "" {
			k = strings.Join([]string{parent, k}, ".")
		}
		switch vt := v.(type) {
		case Template:
			names = append(names, vt.requiredNames(k)...)
		case map[string]interface{}:
			names = append(names, Template(vt).requiredNames(k)...)

		default:
			sv := fmt.Sprintf("%s", v)
			if !strings.HasPrefix(sv, requiredPrefix) {
				continue
			}
			names = append(names, k)
		}
	}
	return names
}

// PublicKey attempts to read the public key from the template
func (t Template) PublicKey() (crypto.PublicKey, x509.PublicKeyAlgorithm) {
	sk := t.Value(parsers.X509PublicKey)
	if sk == "" {
		return nil, 0
	}
	by := stringToBytes(sk)
	pka := keytools.ParsePublicKeyAlgorithm(t.Value(parsers.X509PublicKeyAlgorithm))
	puk, err := keytools.ParsePublicKey(by, pka)
	if err != nil {
		return nil, pka
	}
	return puk, pka
}

func (t Template) funcNames() []string {
	var names []string
	for k := range t {
		v := t.Value(k)
		if !strings.HasPrefix(v, funcPrefix) || !strings.HasSuffix(v, funcSuffix) {
			continue
		}
		names = append(names, k)
	}
	return names
}

func (t Template) valueIgnoreCase(key string) (interface{}, bool) {
	for k, v := range t {
		if strings.EqualFold(k, key) {
			return v, true
		}
	}
	return nil, false
}

package templates

import (
	"fmt"
	"strings"
)

const KeyDelimiter = "."
const requiredPrefix = "?"
const funcPrefix = "${"
const funcSuffix = "}"

// Template represents an x509 resource in plain text (yaml) format
type Template map[string]interface{}

func (t Template) Value(k string) string {
	ks := strings.Split(k, KeyDelimiter)
	li := len(ks) - 1
	parent := t.valueSetIgnoreCase(ks[:li], t)
	if parent == nil {
		return ""
	}
	v := t.valueIgnoreCase(ks[li], parent)
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%s", v)
}

func (t Template) ValueMap(k string) map[string]interface{} {
	return t.valueSetIgnoreCase(strings.Split(k, KeyDelimiter), t)
}

func (t Template) RequiredNames() []string {
	var names []string
	for k := range t {
		if t.Value(k) != requiredPrefix {
			continue
		}
		names = append(names, k)
	}
	return names
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

func (t Template) valueSetIgnoreCase(k []string, m map[string]interface{}) map[string]interface{} {
	if len(k) == 0 {
		return m
	}
	v := t.valueIgnoreCase(k[0], m)
	if v == nil {
		return nil
	}
	vMap, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(k) > 1 {
		vMap = t.valueSetIgnoreCase(k[1:], vMap)
	}
	return vMap

}

func (t Template) valueIgnoreCase(key string, m map[string]interface{}) interface{} {
	for k, v := range m {
		if strings.EqualFold(k, key) {
			return v
		}
	}
	return nil
}

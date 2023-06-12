package model

import "strings"

type FlatMap map[string]*string

func (fm FlatMap) Expand() map[string]interface{} {
	m := map[string]interface{}{}
	for k, v := range fm {
		pm := m
		ks := strings.Split(k, ".")
		last := len(ks) - 1
		if len(ks) > 1 {
			pm = ensureKeyPathPresent(m, ks[:last])
		}
		pm[ks[last]] = v
	}
	return m
}

func ensureKeyPathPresent(m map[string]interface{}, keys []string) map[string]interface{} {
	if len(keys) == 0 {
		return m
	}
	k := keys[0]
	v, ok := m[k]
	if !ok {
		v = map[string]interface{}{}
		m[k] = v
	}
	vm := v.(map[string]interface{})
	return ensureKeyPathPresent(vm, keys[1:])
}

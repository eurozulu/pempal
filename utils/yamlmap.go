package utils

import (
	"fmt"
)

type YamlMap map[string]interface{}

func (y *YamlMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	m := map[interface{}]interface{}{}
	if err := unmarshal(&m); err != nil {
		return err
	}
	*y = cleanMapKeys(m)
	return nil
}

func cleanMapKeys(m map[interface{}]interface{}) map[string]interface{} {
	ym := map[string]interface{}{}
	for k, v := range m {
		switch vt := v.(type) {
		case map[interface{}]interface{}:
			v = cleanMapKeys(vt)
		case map[string]interface{}:
			v = YamlMap(vt)
		default:
			// do nothing
		}
		sk, ok := k.(string)
		if !ok {
			ss, ok := k.(fmt.Stringer)
			if !ok {
				panic("map keys can not be converted to strings")
			}
			sk = ss.String()
		}
		ym[sk] = v
	}
	return ym
}

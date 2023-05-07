package argdecoder

import (
	"fmt"
	"reflect"
)

type mapDecoder struct {
	args []string
}

func (md mapDecoder) Apply(v interface{}) ([]string, error) {
	mapValue := reflect.ValueOf(v)
	if mapValue.Kind() != reflect.Map {
		return nil, fmt.Errorf("can not decode into non map value")
	}
	if mapValue.Type().Key().Kind() != reflect.String {
		return nil, fmt.Errorf("can not decode into maps without string keymanager")
	}

	if mapValue.Type().Elem().Kind() != reflect.Interface {
		return nil, fmt.Errorf("can not decode into maps without interface values")
	}

	params, flags := ParseArgs(md.args)
	m := map[string]interface{}{}
	if len(params) > 0 {
		m[""] = params
	}
	for f, sv := range flags {
		m[f] = sv
	}
	mapValue.Set(reflect.ValueOf(m))
	return nil, nil
}

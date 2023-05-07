package argdecoder

import (
	"fmt"
	"reflect"
	"strings"
)

type stringDecoder struct {
	args []string
}

func (sd stringDecoder) Apply(v interface{}) ([]string, error) {
	stringValue := reflect.ValueOf(v)
	if stringValue.Kind() != reflect.String {
		return nil, fmt.Errorf("can not decode string into %s", stringValue.Kind().String())
	}
	stringValue.Set(reflect.ValueOf(strings.Join(sd.args, " ")))
	return nil, nil
}

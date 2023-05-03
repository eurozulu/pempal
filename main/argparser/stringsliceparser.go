package argparser

import (
	"fmt"
	"reflect"
)

type stringSliceParser struct {
	args []string
}

func (sd stringSliceParser) Apply(v interface{}) ([]string, error) {
	stringValue := reflect.ValueOf(v)
	if stringValue.Kind() != reflect.Slice {
		return nil, fmt.Errorf("can not decode arguments into non string slice")
	}
	if stringValue.Type().Elem().Kind() != reflect.String {
		return nil, fmt.Errorf("can not decode arguments into non string slice")
	}
	stringValue.Set(reflect.ValueOf(sd.args))
	return nil, nil
}

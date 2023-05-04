package argparser

import (
	"fmt"
	"reflect"
	"strings"
)

type ArgumentParser interface {
	// Apply will apply the parser arguments to the given object.
	// Arguments are parsed into flags, args beginning with '-', with the remaining arg value defining the flag name.
	// this name is matched to a public field in the given value either directly by name or via a `flag` tag.
	// If the arguments has another arg following the flag, and not starting with -, this is used as the value for that flag
	// which is assigned to the Field, assuming it can be coersed into the relevant type.
	// Any arguments not matched to fields are returned.
	Apply(v interface{}) ([]string, error)
}

func ApplyArguments(args []string, v interface{}) ([]string, error) {
	vv := reflect.ValueOf(v)
	if vv.IsNil() {
		return nil, fmt.Errorf("can not apply to nil value")
	}
	if vv.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("can not apply to ap non pointer")
	}
	switch vv.Elem().Kind() {
	case reflect.Struct:
		return structParser{args: args}.Apply(v)
	case reflect.Map:
		return mapParser{args: args}.Apply(v)
	case reflect.Slice:
		return stringSliceParser{args: args}.Apply(v)
	case reflect.String:
		return stringParser{args: args}.Apply(v)
	default:
		return nil, fmt.Errorf("%s is not ap supported type", vv.Type().String())
	}
}

func ParseArgs(args []string) (params []string, flags map[string]*string) {
	flags = map[string]*string{}
	for index := 0; index < len(args); index++ {
		if !strings.HasPrefix(args[index], "-") {
			params = append(params, args[index])
			continue
		}
		flag := strings.ToLower(strings.TrimLeft(args[index], "-"))
		var value *string
		if index+1 < len(args) && !strings.HasPrefix(args[index+1], "-") {
			index++
			value = &args[index]
		}
		flags[flag] = value
	}
	return params, flags
}
package argdecoder

import (
	"fmt"
	"pempal/logger"
	"reflect"
	"strconv"
	"strings"
)

const TagKey = "flag"

var boolFailedToDecoder = fmt.Errorf("can not parse as bool")

type structParser struct {
	args []string
}

func (sd structParser) Apply(v interface{}) ([]string, error) {
	structVal := reflect.ValueOf(v)
	if structVal.Kind() != reflect.Ptr ||
		structVal.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("value must be a pointer to a struct")
	}

	params, flags := ParseArgs(sd.args)
	unknownFlags := map[string]*string{}

	// assign flags to Fields
	for f, sv := range flags {
		fld, err := fieldForName(f, structVal.Elem())
		if err != nil {
			logger.Log(logger.Debug, "%v", err)
			unknownFlags[f] = sv
			continue
		}

		v, err := stringAsValue(sv, fld.Type)
		if err != nil {
			// Handel bool as optional values.
			if err != boolFailedToDecoder {
				return nil, err
			}
			// bool value not bool, so treat as param.
			params = append(params, *sv)
			v = reflect.ValueOf(true)

		}
		if err = setFieldValue(structVal, fld, v); err != nil {
			return nil, err
		}
	}
	// assemble what is left over
	result := params
	for k, v := range unknownFlags {
		result = append(result, strings.Join([]string{"-", k}, ""))
		if v != nil {
			result = append(result, *v)
		}
	}
	return result, nil
}

func fieldForName(name string, value reflect.Value) (reflect.StructField, error) {
	t := value.Type()
	for i := 0; i < value.NumField(); i++ {
		fld := t.Field(i)
		if strings.EqualFold(name, fld.Name) {
			return fld, nil
		}
		if isTagName(name, fld.Tag.Get(TagKey)) {
			return fld, nil
		}
	}
	return reflect.StructField{}, fmt.Errorf("%s is not a known field", name)
}

func isTagName(name, tagValue string) bool {
	for _, tn := range strings.Split(tagValue, ",") {
		if strings.EqualFold(strings.TrimSpace(tn), name) {
			return true
		}
	}
	return false
}

func setFieldValue(structVal reflect.Value, fld reflect.StructField, value reflect.Value) error {
	if !fld.IsExported() {
		return fmt.Errorf("field %s is not an exported field", fld.Name)
	}
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != fld.Type.Kind() {
		return fmt.Errorf("field %s could not be assigned with a %s value", fld.Name, value.Kind())
	}
	fldValue := structVal.Elem().FieldByIndex(fld.Index)
	if fldValue.Kind() == reflect.Ptr {
		fldValue = fldValue.Elem()
	}
	fldValue.Set(value)
	return nil
}

func stringAsValue(svalue *string, vtype reflect.Type) (reflect.Value, error) {
	if svalue == nil {
		// No value specified, assign 'zero' value
		// exception is bool fields, which default to true.
		if vtype.Kind() == reflect.Bool {
			return reflect.ValueOf(true), nil
		}
		return reflect.New(vtype), nil
	}
	s := strings.TrimSpace(*svalue)
	switch vtype.Kind() {
	case reflect.String:
		return reflect.ValueOf(s), nil

	case reflect.Bool:
		b := true
		if s != "" {
			bl, err := strconv.ParseBool(s)
			if err != nil {
				return reflect.Value{}, boolFailedToDecoder
			}
			b = bl
		}
		return reflect.ValueOf(b), nil
	case reflect.Uint:
		ui, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("Failed to convert %s into an uint value  %v", s, err)
		}
		return reflect.ValueOf(ui), nil
	case reflect.Int:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("Failed to convert %s into a int value  %v", s, err)
		}
		return reflect.ValueOf(i), nil
	case reflect.Float64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("Failed to convert %s into a floatvalue  %v", s, err)
		}
		return reflect.ValueOf(f), nil
	case reflect.Float32:
		f, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("Failed to convert %s into a floatvalue  %v", s, err)
		}
		return reflect.ValueOf(f), nil

	default:
		return reflect.Value{}, fmt.Errorf("%s type is not supported", vtype.Kind().String())
	}
}

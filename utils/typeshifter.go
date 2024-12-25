package utils

import (
	"encoding"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

const sliceDelimiter = ","

var textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
var urlType = reflect.TypeOf((*url.URL)(nil)).Elem()

func SetValue(s string, v reflect.Value) error {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		return SetValue(s, v.Elem())
	}
	if v.Kind() == reflect.Slice {
		return setSliceValue(strings.Split(s, sliceDelimiter), v)
	}
	if !v.Type().Implements(textUnmarshalerType) {
		return setTextUnmarshalerValue(s, v)
	}
	if v.Type().AssignableTo(urlType) {
		return setURLValue(s, v)
	}
	// none of the specialised types, try as a base type
	ts, err := stringToType(s, v.Type().Kind())
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(ts))
	return nil
}

func StringToType(s string) interface{} {
	if i, err := stringToType(s, reflect.Int64); err == nil {
		return i
	}
	if i, err := stringToType(s, reflect.Int); err == nil {
		return i
	}
	if f, err := stringToType(s, reflect.Float64); err == nil {
		return f
	}
	if b, err := stringToType(s, reflect.Bool); err == nil {
		return b
	}
	return s
}

func stringToType(s string, k reflect.Kind) (interface{}, error) {
	switch k {
	case reflect.String:
		return s, nil
	case reflect.Bool:
		return strconv.ParseBool(s)
	case reflect.Int:
		return strconv.Atoi(s)
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return strconv.ParseInt(s, 64, 64)
	case reflect.Float64:
		return strconv.ParseFloat(s, 64)
	case reflect.Float32:
		return strconv.ParseFloat(s, 32)
	default:
		return nil, fmt.Errorf("%s is an unsupported field type", k)
	}
}

func setSliceValue(ss []string, v reflect.Value) error {
	t := v.Type()
	// TODO Check if value exist and append values
	inst := reflect.MakeSlice(t, len(ss), len(ss))
	for i, s := range ss {
		if err := SetValue(s, inst.Index(i)); err != nil {
			return err
		}
	}
	v.Set(inst)
	return nil
}

func setURLValue(s string, v reflect.Value) error {
	u, err := url.Parse(s)
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(u).Elem())
	return nil
}

func setTextUnmarshalerValue(s string, fld reflect.Value) error {
	fldPtr := fld
	if fld.Type().Kind() != reflect.Ptr {
		fldPtr = fld.Addr()
	}
	tum := fldPtr.Interface().(encoding.TextUnmarshaler)
	return tum.UnmarshalText([]byte(s))
}

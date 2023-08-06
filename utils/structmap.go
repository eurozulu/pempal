package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var ErrUnknownTag = fmt.Errorf("Unknown tag")

type StructMap interface {
	GetValue(tag string) (interface{}, error)
	SetValue(tag string, value interface{}) error
	SetValueString(tag string, value *string) error
}

type structMap struct {
	structVal reflect.Value
	tagType   string
}

func (sm structMap) GetValue(tag string) (interface{}, error) {
	fi := sm.fieldIndexByTag(tag)
	if fi < 0 {
		return nil, ErrUnknownTag
	}
	return sm.structVal.Field(fi).Interface(), nil
}

func (sm structMap) SetValue(tag string, value interface{}) error {
	fi := sm.fieldIndexByTag(tag)
	if fi < 0 {
		return ErrUnknownTag
	}
	return sm.setValue(sm.structVal.Field(fi), value)
}

func (sm structMap) SetValueString(tag string, value *string) error {
	return sm.SetValue(tag, value)
}

func (sm structMap) setValue(field reflect.Value, value interface{}) error {
	switch field.Kind() {
	case reflect.Pointer:
		return sm.setPointerValue(field, value)
	case reflect.String:
		field.SetString(valueAsString(value))
	case reflect.Bool:
		field.SetBool(valueAsBool(value))
	case reflect.Int, reflect.Int64:
		field.SetInt(valueAsInt(value))
	case reflect.Float64, reflect.Float32:
		field.SetFloat(valueAsFloat(value))
	default:
		return fmt.Errorf("field %s type %s not supported", field.Type().Name(), field.Type().String())
	}
	return nil
}

func (sm structMap) setPointerValue(field reflect.Value, value interface{}) error {
	nv := reflect.New(field.Type().Elem())
	field.Set(nv)
	return sm.setValue(field.Elem(), value)
}

func valueAsString(v interface{}) string {
	vo := reflect.ValueOf(v)
	if vo.Type().Kind() == reflect.Pointer {
		if vo.IsNil() {
			return ""
		}
		v = vo.Elem().Interface()
	}
	switch vt := v.(type) {
	case string:
		return vt
	case fmt.Stringer:
		return vt.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func valueAsBool(v interface{}) bool {
	if isNilInterface(v) {
		return true
	}
	b, _ := strconv.ParseBool(valueAsString(v))
	return b
}

func valueAsInt(v interface{}) int64 {
	if v == nil {
		return 0
	}
	i, _ := strconv.ParseInt(valueAsString(v), 10, 64)
	return i
}

func valueAsFloat(v interface{}) float64 {
	if v == nil {
		return 0
	}
	i, _ := strconv.ParseFloat(valueAsString(v), 64)
	return i
}

func isNilInterface(i interface{}) bool {
	if i == nil {
		return true
	}
	vo := reflect.ValueOf(i)
	if vo.Kind() != reflect.Pointer {
		return false
	}
	return vo.IsNil()
}

func (sm structMap) fieldIndexByTag(tag string) int {
	vt := sm.structVal.Type()
	for i := 0; i < vt.NumField(); i++ {
		if sm.tagContainsKey(vt.Field(i).Tag, tag) {
			return i
		}
	}
	return -1
}

func (sm structMap) tagContainsKey(tag reflect.StructTag, key string) bool {
	tagval, ok := tag.Lookup(sm.tagType)
	if !ok || tagval == "-" {
		return false
	}
	for _, tv := range strings.Split(tagval, ",") {
		if strings.EqualFold(tv, key) {
			return true
		}
	}
	return false
}

func NewStructMap(i interface{}, tagname string) (StructMap, error) {
	vo := reflect.ValueOf(i)
	if isNilInterface(vo) || isNilInterface(vo.Elem()) {
		return nil, fmt.Errorf("structure can not be nil")
	}
	if vo.Kind() != reflect.Pointer || vo.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("Invalid value type (%s)!  Must be a pointer to a structure", vo.Type().String())
	}
	return &structMap{
		structVal: vo.Elem(),
		tagType:   tagname,
	}, nil

}

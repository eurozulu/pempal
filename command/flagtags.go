package command

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

const tagName = "flag"

type FlagFields interface {
	Names() []string
	Contains(name string) bool
	IsBool(name string) bool
	SetValue(name string, value string) error
}

type flagFields struct {
	cmd    interface{}
	fields map[string]reflect.StructField
}

func (ff flagFields) Names() []string {
	names := make([]string, len(ff.fields))
	var index int
	for k := range ff.fields {
		names[index] = k
		index++
	}
	return names
}

func (ff flagFields) IsBool(name string) bool {
	fd, ok := ff.fields[name]
	if !ok {
		return false
	}
	return fd.Type.Kind() == reflect.Bool
}

func (ff flagFields) Contains(name string) bool {
	_, ok := ff.fields[name]
	return ok
}

func (ff flagFields) SetValue(name string, value string) error {
	fld, ok := ff.fields[name]
	if !ok {
		return fmt.Errorf("flag %s is not known", name)
	}
	v := reflect.ValueOf(ff.cmd)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	fv := v.FieldByName(fld.Name)

	switch fv.Kind() {
	case reflect.String:
		fv.SetString(value)

	case reflect.Bool:
		// value is optional for bool
		bv := true
		if value != "" {
			var err error
			bv, err = strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("could not read %s value '%s' as a boolean, true or false value  %v", name, value, err)
			}
		}
		fv.SetBool(bv)
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		iv, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("could not read %s value '%s' as an integer number  %v", name, value, err)
		}
		fv.SetInt(iv)

	case reflect.Float64, reflect.Float32:
		flv, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("could not read %s value '%s' as an floating point number  %v", name, value, err)
		}
		fv.SetFloat(flv)

	default:
		return fmt.Errorf("field %s has an unsupported type of %s", name, fv.Type().String())
	}
	return nil
}

func newFlagFields(cmd interface{}) (FlagFields, error) {
	fields := map[string]reflect.StructField{}
	t := reflect.TypeOf(cmd)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("command is not a struct or pointer to one")
	}
	for i := 0; i < t.NumField(); i++ {
		fld := t.Field(i)
		if !fld.IsExported() {
			continue
		}
		tag := fld.Tag.Get(tagName)
		if tag == "" || tag == "-" {
			continue
		}
		names := strings.Split(tag, ",")
		for _, n := range names {
			n = strings.TrimSpace(n)
			// Check if its uniqye
			if df, ok := fields[n]; ok {
				// fatal log as this is not a runtime error but a programming error
				log.Fatalf("Command %s field '%s' has a duplicate flag of '%s' with field %s",
					t.Name(), fld.Name, n, df)
			}
			fields[n] = fld
		}
	}
	return &flagFields{cmd: cmd, fields: fields}, nil
}

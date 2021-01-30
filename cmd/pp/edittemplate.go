package main

import (
	"fmt"
	"github.com/eurozulu/pempal/templates"
	"gopkg.in/yaml.v3"
	"reflect"
	"strconv"
	"time"
)

func EditFields(t interface{}) error {
	by, err := yaml.Marshal(t)
	if err != nil {
		return err
	}
	m := make(map[string]interface{})
	if err := yaml.Unmarshal(by, m); err != nil {
		return err
	}
	flds := templateFieldNames(t)

	for _, fld := range flds {
		v, ok := m[fld]
		if !ok {
			continue
		}

		s, ok := v.(templates.SubjectTemplate)
		if ok {
			if err := EditFields(&s); err != nil {
				return err
			}
			continue
		}
		vs := PromptInput(fld, interfaceToString(v))
		m[fld] = vs
	}

	by, err = yaml.Marshal(m)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(by, t)
}

// templateFieldNames gets the public fields of the template which have a 'yaml' tag
// Ordered in the order they appear in the struct.
func templateFieldNames(t interface{}) []string {
	tp := reflect.TypeOf(t).Elem()
	var names []string
	for i := 0; i < tp.NumField(); i++ {
		fld := tp.Field(i)
		_, ok := fld.Tag.Lookup("yaml")
		if !ok {
			continue
		}
		names = append(names, tp.Field(i).Name)
	}
	return names
}

func interfaceToString(i interface{}) string {
	switch v := i.(type) {
	case string:
		return v

	case int, int64:
		return strconv.FormatInt(i.(int64), 10)

	case bool:
		return strconv.FormatBool(v)

	case time.Time:
		return v.String()

	default :
		return fmt.Sprintf("%v", v)
	}
}

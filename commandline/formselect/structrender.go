package formselect

import "reflect"

func buildFieldMap() map[string]reflect.Type {
	m := map[string]reflect.Type{}
	t := reflect.TypeOf(cf.data)
	c := t.NumField()
	for i := 0; i < c; i++ {
		f := t.Field(i)
		if !f.IsExported() || f.Anonymous {
			continue
		}
		m[f.Name] = f.Type
	}
	return m
}

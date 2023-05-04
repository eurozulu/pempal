package resourceio

import (
	"bytes"
	"fmt"
	"github.com/go-yaml/yaml"
	"io"
	"pempal/model"
	"strings"
)

const defaultFields = "resource-type,identity,subject.common-name,not-after"

type ResourceListPrinter struct {
	Titles     []string
	Fields     []string
	ShowTitles bool
	out        io.Writer
}

func (prn *ResourceListPrinter) Write(location ResourceLocation) error {
	buf := bytes.NewBuffer(nil)
	if prn.ShowTitles {
		prn.writeFields(prn.TitleNames(), prn.out)
		fmt.Fprintf(prn.out, "\n")
		prn.ShowTitles = false
	}
	for _, r := range location.Resources {
		dto, err := model.DTOForResource(r)
		if err != nil {
			return err
		}
		values, err := prn.valuesFromDTO(dto)
		if err != nil {
			return err
		}
		if err = prn.writeFields(values, buf); err != nil {
			return err
		}
		fmt.Fprintf(buf, "\t%s\n", location.Path)
	}
	_, err := prn.out.Write(buf.Bytes())
	return err
}

func (prn ResourceListPrinter) TitleNames() []string {
	names := make([]string, len(prn.Fields)+1)
	for i, n := range prn.Fields {
		if i < len(prn.Titles) && prn.Titles[i] != "" {
			n = prn.Titles[i]
		}
		names[i] = n
	}
	names[len(names)-1] = "location"
	return names
}

func (prn ResourceListPrinter) writeFields(values []string, out io.Writer) error {
	for i, title := range prn.TitleNames() {
		if i >= len(values) {
			break
		}
		v := values[i]
		w := len(title)
		if len(v) > w {
			// truncate value to fit
			v = v[:w+1]

		} else if len(v) < w {
			// padout with spaces
			v = strings.Join([]string{v,
				strings.Repeat(" ", w-len(v))}, "")
		}
		if i > 0 {
			if _, err := fmt.Fprint(out, "\t"); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprint(out, v); err != nil {
			return err
		}
	}
	return nil
}

func (prn ResourceListPrinter) valuesFromDTO(dto model.DTO) ([]string, error) {
	m, err := dtoToMap(dto)
	if err != nil {
		return nil, err
	}
	flds := make([]string, len(prn.Fields))
	for i, f := range prn.Fields {
		v := prn.valueFromMap(f, m)
		if v != nil {
			flds[i] = fmt.Sprint(v)
		}
	}
	return flds, nil
}

func (prn ResourceListPrinter) valueFromMap(key string, m map[string]interface{}) interface{} {
	var val interface{}
	keys := strings.Split(key, ".")
	for i, k := range keys {
		v, ok := m[k]
		if !ok {
			break
		}
		if i == len(keys)-1 {
			val = v
			break
		}
		vm := valueAsMap(v)
		if vm == nil {
			break
		}
		m = vm
	}
	return val
}

func valueAsMap(v interface{}) map[string]interface{} {
	sm, ok := v.(map[string]interface{})
	if ok {
		return sm
	}
	im, ok := v.(map[interface{}]interface{})
	if !ok {
		return nil
	}
	m := map[string]interface{}{}
	for k, val := range im {
		sk, ok := k.(string)
		if !ok {
			//not a string key, abort
			return nil
		}
		m[sk] = val
	}
	return m
}

func dtoToMap(dto model.DTO) (map[string]interface{}, error) {
	data, err := yaml.Marshal(dto)
	if err != nil {
		return nil, err
	}
	m := map[string]interface{}{}
	if err = yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func NewResourceListPrinter(out io.Writer, fields ...string) *ResourceListPrinter {
	if len(fields) == 0 {
		fields = strings.Split(defaultFields, ",")
	}
	return &ResourceListPrinter{
		Fields: fields,
		out:    out,
	}
}

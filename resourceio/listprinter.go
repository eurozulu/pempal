package resourceio

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io"
	"pempal/model"
	"pempal/utils"
	"strings"
)

const defaultFields = "resource-type,identity,subject.common-name,not-after"

type ResourceListPrinter struct {
	Fields []string
	out    *utils.ColumnOutput
}

func (prn *ResourceListPrinter) Write(location ResourceLocation) error {
	for _, r := range location.Resources {
		dto, err := model.DTOForResource(r)
		if err != nil {
			return err
		}
		values, err := prn.valuesFromDTO(dto)
		if err != nil {
			return err
		}
		values = append(values, location.Path)
		if _, err := prn.out.WriteSlice(values); err != nil {
			return err
		}
		fmt.Fprintln(prn.out)
	}
	return nil
}

func (prn ResourceListPrinter) TitleNames() []string {
	return append(prn.Fields, "location")
}
func (prn ResourceListPrinter) WriteTitles() error {
	_, err := prn.out.WriteSlice(prn.TitleNames())
	fmt.Fprintln(prn.out)
	return err
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
		out:    utils.NewColumnOutput(out),
	}
}

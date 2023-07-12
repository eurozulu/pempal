package commands

import (
	"github.com/go-yaml/yaml"
	"strings"
)

type FlagMap map[string]*string

func (fm FlagMap) ApplyAndRemove(v interface{}) error {
	if err := fm.ApplyFlags(v); err != nil {
		return err
	}
	if err := fm.RemoveKeys(v); err != nil {
		return err
	}
	return nil
}

func (fm FlagMap) ApplyFlags(v interface{}) error {
	// Apply yaml (expanded) flags map to the given interface
	by, err := yaml.Marshal(&fm)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(by, v); err != nil {
		return err
	}
	return nil
}

func (fm FlagMap) RemoveKeys(v interface{}) error {
	keys, err := readYamlKeys(v)
	if err != nil {
		return err
	}
	for _, k := range keys {
		if len(fm) == 0 {
			break
		}
		if _, ok := fm[k]; !ok {
			// not found in flags, ignore it
			continue
		}
		delete(fm, k)
	}
	return nil
}

func readYamlKeys(v interface{}) ([]string, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return nil, err
	}
	m := map[string]interface{}{}
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	var names []string
	for k := range m {
		names = append(names, k)
	}
	return names, nil
}

func ParseArgs(args []string) (params []string, flags FlagMap) {
	flags = FlagMap{}
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

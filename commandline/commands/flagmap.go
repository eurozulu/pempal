package commands

import (
	"github.com/eurozulu/pempal/utils"
	"strings"
)

type FlagMap map[string]*string

func (fm *FlagMap) ApplyTo(v interface{}) error {
	sm, err := utils.NewStructMap(v, "yaml")
	if err != nil {
		return err
	}
	var keys []string
	for k, val := range *fm {
		if err := sm.SetValueString(k, val); err != nil {
			if err != utils.ErrUnknownTag {
				return err
			}
			continue
		}
		keys = append(keys, k)
	}
	fm.removeKeys(keys)
	return nil
}

func (fm *FlagMap) removeKeys(keys []string) {
	for _, k := range keys {
		delete(*fm, k)
	}
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

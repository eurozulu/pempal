package command

import (
	"strings"
)

// Arguments represents the command line arguments
type Arguments interface {
	Parameters() []string
	ContainsFlag(name string) bool
	FlagValue(name string) string
	Args() []string
	Flags() map[string]string
}

type arguments struct {
	args []string
}

func (a arguments) Args() []string {
	return a.args
}
func (a arguments) Flags() map[string]string {
	m := map[string]string{}
	for i := 0; i < len(a.args); i++ {
		arg := a.args[i]
		if strings.HasPrefix(arg, "-") {
			k := strings.TrimLeft(arg, "-")
			var v string
			if a.HasValue(k) {
				v = a.FlagValue(k)
				i++
			}
			m[k] = v
		}
	}
	return m
}

func (a arguments) Parameters() []string {
	var parms []string
	var index int
	for ; index < len(a.args); index++ {
		arg := a.args[index]
		if strings.HasPrefix(arg, "-") {
			if a.HasValue(strings.TrimLeft(arg, "-")) {
				index++
			}
			continue
		}
		parms = append(parms, a.args[index])
	}
	return parms
}

func (a arguments) ContainsFlag(name string) bool {
	return a.FlagIndex(name) >= 0
}

func (a arguments) FlagValue(name string) string {
	i := a.FlagIndex(name)
	if i < 0 || i+1 >= len(a.args) {
		return ""
	}
	return a.args[i+1]
}

func (a arguments) HasValue(name string) bool {
	i := a.FlagIndex(name)
	return i >= 0 && i+1 < len(a.args)
}

func (a arguments) FlagIndex(name string) int {
	var index int
	for _, arg := range a.args {
		if strings.HasPrefix(arg, "-") {
			if strings.TrimLeft(arg, "-") == name {
				return index
			}
		}
		index++
	}
	return -1
}

func (a *arguments) removeFlag(name string) {
	i := a.FlagIndex(name)
	if i < 0 {
		return
	}
	var last []string
	if i+1 < len(a.args) {
		last = a.args[i+1:]
	}
	a.args = append(a.args[:i], last...)
}

func (a *arguments) removeFlagAndValue(name string) {
	iStart := a.FlagIndex(name)
	if iStart < 0 {
		return
	}
	iEnd := iStart + 1
	var last []string
	if iEnd+1 < len(a.args) {
		last = a.args[iEnd+1:]
	}
	a.args = append(a.args[:iStart], last...)
}

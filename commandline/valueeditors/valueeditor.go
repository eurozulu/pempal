package valueeditors

import (
	"github.com/eurozulu/pempal/commandline/ui"
	"strconv"
	"strings"
)

// ValueEditor represents a single named property which may be edited.
type ValueEditor interface {
	Name() string
	Edit(offset ui.ViewOffset, value string) (string, error)
}

type StringEditor struct {
	// PropertyName of the name of the property to edit
	PropertyName string
	// Selection of values which the property may contain
	Choice []string
	// DefaultChoice of the choices if no value is present
	DefaultChoice int
	// AllowInput flag indicates if user input values are allowed. False = only Choice strings allowed
	AllowedInput bool
}

type NumberEditor struct {
	// PropertyName of the name of the property to edit
	PropertyName string
	// Selection of values which the property may contain
	Choice []int
	// DefaultChoice the index of the choices if no value is present
	DefaultChoice int
	// AllowInput flag indicates if user input values are allowed. False = only Choice numbers allowed
	// If Choices is empty, this is forced to true.  i.e. if no choices, input is always allowed
	AllowedInput bool
}

type BoolEditor struct {
	// PropertyName of the name of the property to edit
	PropertyName string

	// DefaultChoice the default value
	DefaultChoice bool
}

type StringSliceEditor struct {
	// PropertyName of the name of the property to edit
	PropertyName string
	// Selection of values which the property may contain
	Choice []string
	// AllowInput flag indicates if user input values are allowed. False = only Choice strings allowed
	AllowedInput bool
}

func (se StringEditor) Name() string {
	return se.PropertyName
}
func (se StringEditor) Edit(offset ui.ViewOffset, value string) (string, error) {
	it := ui.InputTypeNone
	if se.AllowedInput || len(se.Choice) == 0 {
		it = ui.InputTypePrintable
	}
	ev := ui.EditItem{
		Name:      se.Name(),
		ValueType: it,
		Options:   se.Choice,
	}
	return ev.Edit(offset, value)
}

func (ne NumberEditor) Name() string {
	return ne.PropertyName
}

func (ne NumberEditor) Edit(offset ui.ViewOffset, value string) (string, error) {
	it := ui.InputTypeNone
	if ne.AllowedInput || len(ne.Choice) == 0 {
		it = ui.InputTypeNumbers
	}
	ev := &ui.EditItem{
		Name:      ne.Name(),
		ValueType: it,
		Options:   ne.choiceAsStrings(),
	}
	return ev.Edit(offset, value)
}

func (ne NumberEditor) choiceAsStrings() []string {
	items := make([]string, len(ne.Choice))
	for i, iv := range ne.Choice {
		items[i] = strconv.FormatInt(int64(iv), 10)
	}
	return items
}

func (be BoolEditor) Name() string {
	return be.PropertyName
}

func (be BoolEditor) Edit(offset ui.ViewOffset, value string) (string, error) {
	v := be.DefaultChoice
	if value != "" {
		b, err := parseBool(value)
		if err != nil {
			return "", err
		}
		v = b
	}
	choice := []string{"yes", "no"}
	if !v {
		choice = []string{"no", "yes"}
	}
	ev := ui.EditItem{
		Name:      be.Name(),
		ValueType: ui.InputTypeNone,
		Options:   choice,
	}
	return ev.Edit(offset, value)
}

func (sle StringSliceEditor) Name() string {
	return sle.PropertyName
}

func (sle StringSliceEditor) Edit(offset ui.ViewOffset, value string) (string, error) {
	it := ui.InputTypeNone
	if sle.AllowedInput || len(sle.Choice) == 0 {
		it = ui.InputTypePrintable
	}
	ev := ui.EditItem{
		Name:      sle.Name(),
		ValueType: it,
		Options:   sle.Choice,
	}
	v, err := ev.Edit(offset, value)
	if err != nil {
		return "", err
	}
	return v, nil
}

func parseBool(s string) (bool, error) {
	if strings.EqualFold(s, "yes") || strings.EqualFold(s, "true") || strings.EqualFold(s, "1") {
		return true, nil
	}
	if strings.EqualFold(s, "no") || strings.EqualFold(s, "false") || strings.EqualFold(s, "0") {
		return false, nil
	}
	return strconv.ParseBool(s)
}

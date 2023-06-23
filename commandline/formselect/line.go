package formselect

import (
	"reflect"
	"strconv"
)

type Line interface {
	Name() string
	Format() string
	Parse(s string) error
	Required() bool
}

type StringLine struct {
	name     string
	value    string
	required bool
}

func (sc StringLine) Name() string {
	return sc.name
}
func (sc StringLine) Required() bool {
	return sc.required
}
func (sc StringLine) Format() string {
	return sc.value
}
func (sc *StringLine) Parse(s string) error {
	sc.value = s
	return nil
}

type BoolLine struct {
	name     string
	value    bool
	required bool
}

func (sc BoolLine) Name() string {
	return sc.name
}
func (sc BoolLine) Required() bool {
	return sc.required
}
func (sc BoolLine) Format() string {
	return strconv.FormatBool(sc.value)
}
func (sc *BoolLine) Parse(s string) error {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	sc.value = b
	return nil
}

type IntLine struct {
	name     string
	value    int64
	required bool
}

func (sc IntLine) Name() string {
	return sc.name
}
func (sc IntLine) Required() bool {
	return sc.required
}
func (sc IntLine) Format() string {
	return strconv.FormatInt(sc.value, 10)
}
func (sc *IntLine) Parse(s string) error {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	sc.value = i
	return nil
}

func lineByType(name string, v interface{}, required bool) (Line, error) {
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.String:
		return &StringLine{
			name:     name,
			value:    v.(string),
			required: required,
		}, nil
	case reflect.Int, reflect.Int64, reflect.Int32:
		return &IntLine{
			name:     name,
			value:    v.(int64),
			required: required,
		}, nil
	case reflect.Bool:
		return &BoolLine{
			name:     name,
			value:    v.(bool),
			required: required,
		}, nil
	default:

	}
}
func NewLine(v interface{}) (Line, error) {

}

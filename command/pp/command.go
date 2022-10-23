package main

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

const FlagToken = "-"

type Command interface {
	Run(out io.Writer, args ...string) error
}

type NewCommandFnc func() Command

func ApplyFlags(cmd interface{}, args ...string) ([]string, error) {
	// get the names of fields tagged as flags
	flagNames := commandFlagNames(cmd)

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, FlagToken) {
			continue
		}
		arg = strings.TrimLeft(arg, FlagToken)
		fieldName := findFlag(arg, flagNames)
		if fieldName == "" {
			// arg has no matching field name, ignore it.
			continue
		}
		// Remove matched argument from the list & establish value
		args = removeElement(args, i)
		var val string
		if isBoolField(cmd, fieldName) {
			val = "true"
			// Check if following arg is a bool value
			if i < len(args) && !strings.HasPrefix(args[i], FlagToken) {
				_, err := strconv.ParseBool(args[i])
				if err == nil {
					val = args[i]
					args = removeElement(args, i)
				}
			}
		} else {
			// non bool field, Must have value
			if i >= len(args) || strings.HasPrefix(args[i], FlagToken) {
				return nil, fmt.Errorf("%s is missing a value", arg)
			}
			val = args[i]
			args = removeElement(args, i)
		}
		if err := applyFlag(cmd, fieldName, val); err != nil {
			return nil, err
		}
	}
	return args, nil
}

func applyFlag(cmd interface{}, name string, value string) error {
	vCmd := reflect.ValueOf(cmd)
	if vCmd.IsNil() {
		return fmt.Errorf("unknown command")
	}
	vF := vCmd.FieldByName(name)
	if vF.IsZero() {
		return fmt.Errorf("%s is an unknown argument", name)
	}
	switch vF.Kind() {
	case reflect.String:
		vF.SetString(value)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("%s is not a valud boolean value")
		}
		vF.SetBool(b)
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("%s could not be read an a number  %v", name, err)
		}
		vF.SetInt(i)

	default:
		return fmt.Errorf("%s is an unsupported type")
	}
	return nil
}

func commandFlagNames(cmd interface{}) []string {
	var names []string
	t := reflect.TypeOf(cmd)
	count := t.Elem().NumField()
	for i := 0; i < count; i++ {
		f := t.Field(i)
		if f.Anonymous || !unicode.IsUpper(rune(f.Name[0])) {
			continue
		}
		name := readTagValue(string(f.Tag))
		if name == "" {
			continue
		}
		names = append(names, name)
	}
	return names
}

func readTagValue(tag string) string {
	i := strings.Index(tag, "flag:")
	if i < 0 {
		return ""
	}
	is := strings.Split(tag[i+len("flag:"):], ",")
	return strings.TrimSpace(is[0])
}

func isBoolField(i interface{}, name string) bool {
	return reflect.ValueOf(i).FieldByName(name).Kind() == reflect.Bool
}

func removeElement(s []string, index int) []string {
	var end []string
	if index+1 < len(s) {
		end = s[index+1:]
	}
	return append(s[:index], end...)
}

func findFlag(arg string, names []string) string {
	k := strings.ToLower(arg)
	for _, n := range names {
		if k == strings.ToLower(n) {
			return n
		}
	}
	return ""
}

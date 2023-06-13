package utils

import (
	"github.com/eurozulu/pempal/logger"
	"reflect"
	"strconv"
	"strings"
)

type FlatMap map[string]*string

func (fm FlatMap) Expand() map[string]interface{} {
	m := map[string]interface{}{}
	for k, v := range fm {
		pm := m
		ks := strings.Split(k, ".")
		last := len(ks) - 1
		if len(ks) > 1 {
			pm = ensureKeyPathPresent(m, ks[:last])
		}
		pm[ks[last]] = stringToValue(v)
	}
	return m
}

func (fm FlatMap) MarshalYAML() (interface{}, error) {
	return fm.Expand(), nil
}

func (fm *FlatMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	m := map[string]interface{}{}
	if err := unmarshal(&m); err != nil {
		return err
	}
	fm.importIntoKey("", m)
	return nil
}

func (fm *FlatMap) mergeIntoKey(key string, m FlatMap) {
	for k, v := range m {
		nk := strings.Join([]string{key, k}, ".")
		(*fm)[nk] = v
	}
}

func (fm *FlatMap) importIntoKey(key string, m map[string]interface{}) {
	expandKey := key != ""
	for k, v := range m {
		if expandKey {
			k = strings.Join([]string{key, k}, ".")
		}
		switch vt := v.(type) {
		case *string:
			(*fm)[k] = vt
		case string:
			(*fm)[k] = &vt
		case bool:
			s := strconv.FormatBool(vt)
			(*fm)[k] = &s
		case int:
			s := strconv.Itoa(vt)
			(*fm)[k] = &s
		case int8, int16, int32, int64:
			s := strconv.FormatInt(v.(int64), 10)
			(*fm)[k] = &s
		case uint, uint8, uint16, uint32, uint64:
			s := strconv.FormatUint(v.(uint64), 10)
			(*fm)[k] = &s
		case float32, float64:
			s := strconv.FormatFloat(v.(float64), 'f', -1, 64)
			(*fm)[k] = &s

		case map[string]*string:
			fm.mergeIntoKey(k, FlatMap(vt))
		case map[string]interface{}:
			fm.importIntoKey(k, vt)
		case map[interface{}]interface{}:
			fm.importIntoKey(k, flipInterfaceKeyMap(vt))

		default:
			logger.Warning("Unable to format key %s value of type %s, into a string value", k, reflect.TypeOf(v).String())
			(*fm)[k] = nil
		}
	}
}

func ensureKeyPathPresent(m map[string]interface{}, keys []string) map[string]interface{} {
	if len(keys) == 0 {
		return m
	}
	k := keys[0]
	v, ok := m[k]
	if !ok {
		v = map[string]interface{}{}
		m[k] = v
	}
	vm := v.(map[string]interface{})
	return ensureKeyPathPresent(vm, keys[1:])
}

func stringToValue(s *string) interface{} {
	if s == nil {
		return nil
	}
	b, err := strconv.ParseBool(*s)
	if err == nil {
		return &b
	}
	i, err := strconv.ParseInt(*s, 10, 64)
	if err == nil {
		return &i
	}
	f, err := strconv.ParseFloat(*s, 64)
	if err == nil {
		return &f
	}
	if strings.HasPrefix(*s, "[") && strings.HasSuffix(*s, "]") {
		ss := strings.Split(strings.Trim(*s, "[]"), ",")
		return &ss
	}
	return s
}

func flipInterfaceKeyMap(m map[interface{}]interface{}) map[string]interface{} {
	nm := map[string]interface{}{}
	for k, v := range m {
		sk, ok := k.(string)
		if !ok {
			sp, ok := k.(*string)
			if !ok {
				continue
			}
			sk = *sp
		}
		vm, ok := v.(map[interface{}]interface{})
		if ok {
			v = flipInterfaceKeyMap(vm)
		}
		nm[sk] = v
	}
	return nm
}

func NewFlatMap(m map[string]interface{}) FlatMap {
	fm := FlatMap{}
	fm.importIntoKey("", m)
	return fm
}

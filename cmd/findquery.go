package cmd

import (
	"fmt"
	"regexp"
	"strings"
)

type FindQuery interface {
	Match(data map[string]interface{}) bool
	Values(data map[string]interface{}) []string
}

type findQuery struct {
	names   []string
	queries map[string]*regexp.Regexp
}

func (lq findQuery) Match(data map[string]interface{}) bool {
	for k, v := range lq.queries {
		if v == nil {
			continue
		}

		// If query has condition, that property must match.
		dv, ok := data[k]
		dvs := fmt.Sprintf("%v", dv)
		if !ok {
			return false
		}
		if !v.MatchString(dvs) {
			return false
		}
	}
	return true
}

func (lq findQuery) Values(data map[string]interface{}) []string {
	vals := make([]string, len(lq.names))
	var found = false
	for i, name := range lq.names {
		k := keyWithoutCase(name, data)
		if k == "" {
			continue
		}
		v, ok := data[k]
		if !ok {
			continue
		}
		found = true
		vals[i] = fmt.Sprintf("%s", v)
	}
	if !found {
		return nil
	}
	return vals
}

func (lq findQuery) IsEmpty() bool {
	return len(lq.names) == 0
}

func keyWithoutCase(name string, m map[string]interface{}) string {
	for k := range m {
		if strings.EqualFold(k, name) {
			return k
		}
	}
	return ""
}
func indexOf(s string, ss []string) int {
	for i, sz := range ss {
		if s == sz {
			return i
		}
	}
	return -1
}

func ParseQuery(query string) (FindQuery, error) {
	var names []string
	queries := map[string]*regexp.Regexp{}
	for _, q := range strings.Split(query, ",") {
		qq := strings.SplitN(q, "=", 2)
		_, ok := queries[qq[0]]
		if !ok {
			names = append(names, qq[0])
		}
		var v *regexp.Regexp
		if len(qq) > 1 {
			r, err := regexp.Compile(qq[1])
			if err != nil {
				return nil, fmt.Errorf("Failed to evaluate expression %s  %v", q, err)
			}
			v = r
		}
		queries[qq[0]] = v
	}
	return &findQuery{
		names:   names,
		queries: queries,
	}, nil
}

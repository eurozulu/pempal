package commands

import (
	"fmt"
	"pempal/templates"
	"regexp"
	"strings"
)

type FindQuery interface {
	Match(data templates.Template) bool
	Values(data templates.Template) []string
}

type findQuery struct {
	orderedNames []string
	queries      map[string]*regexp.Regexp
}

func (lq findQuery) Match(data templates.Template) bool {
	for k, v := range lq.queries {
		// No expression, accept
		if v == nil {
			continue
		}
		// If query has condition, that property must match.
		if !data.Contains(k) {
			return false
		}
		dv := data.Value(k)
		if !v.MatchString(dv) {
			return false
		}
	}
	return true
}

func (lq findQuery) Values(data templates.Template) []string {
	vals := make([]string, len(lq.orderedNames))
	for i, name := range lq.orderedNames {
		vals[i] = data.Value(name)
	}
	return vals
}

func (lq findQuery) IsEmpty() bool {
	return len(lq.orderedNames) == 0
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
		name := strings.TrimSpace(qq[0])
		_, ok := queries[name]
		if !ok {
			names = append(names, name)
		}
		var v *regexp.Regexp
		if len(qq) > 1 {
			value := strings.TrimSpace(qq[1])
			// Convert short cut wildcrd to a regex equiv
			if value == "*" {
				value = ".*"
			}
			r, err := regexp.Compile(value)
			if err != nil {
				return nil, fmt.Errorf("Failed to evaluate expression %s  %v", q, err)
			}
			v = r
		}
		queries[name] = v
	}
	return &findQuery{
		orderedNames: names,
		queries:      queries,
	}, nil
}

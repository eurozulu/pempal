package templates

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"pempal/fileformats"
	"sort"
	"strings"
)

const RequiredKey = "?"
const funcStartKey = "{{"
const funcEndKey = "}}"
const AppendKeyPrefix = "+"

type TemplateBuilder struct {
	Templates []Template
}

func (tb TemplateBuilder) Build() (Template, error) {
	base, req := tb.mergedMaps()
	missing := falseKeys(req)
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required properties: %s", strings.Join(missing, ", "))
	}
	pt, ok := base["pem_type"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'pem_type'")
	}
	by, err := yaml.Marshal(base)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(pt, fileformats.PEM_TEMPLATE) {
		pt = strings.Join([]string{pt, fileformats.PEM_TEMPLATE}, "")
	}
	return BlockToTemplate(&pem.Block{
		Type:  pt,
		Bytes: by,
	})
}

func (tb TemplateBuilder) MissingNames() []string {
	return falseKeys(tb.RequiredNames())
}

func (tb TemplateBuilder) RequiredNames() map[string]bool {
	_, r := tb.mergedMaps()
	return r
}

func (tb TemplateBuilder) mergedMaps() (base map[string]interface{}, required map[string]bool) {
	base = map[string]interface{}{}
	required = map[string]bool{}
	ms, err := templatesToMaps(tb.Templates)
	if err != nil {
		log.Println(err)
		return base, required
	}
	for _, m := range ms {
		mergeMap(m, base, "", required)
	}
	return base, required
}

func mergeMap(m map[string]interface{}, base map[string]interface{}, keyheader string, required map[string]bool) {
	for k, v := range m {
		switch vt := v.(type) {
		case string:
			if vt == RequiredKey {
				required[k] = required[k]
				continue
			}
			// empty string doesn't overwrite
			if vt == "" {
				if _, exists := base[k]; exists {
					continue
				}
			}
			// If this is a required field, mark it as satisfied
			if _, ok := required[k]; ok {
				required[k] = true
			}
			base[k] = vt

		case map[string]interface{}:
			existMap, ok := base[k].(map[string]interface{})
			if !ok {
				existMap = map[string]interface{}{}
				base[k] = existMap
			}
			mergeMap(vt, existMap, strings.Join([]string{keyheader, k}, "."), required)

		default:
			base[k] = v
		}
	}
}

func templatesToMaps(ts []Template) ([]map[string]interface{}, error) {
	ms := make([]map[string]interface{}, len(ts))
	for i, t := range ts {
		by, err := yaml.Marshal(t)
		if err != nil {
			return nil, err
		}
		m := map[string]interface{}{}
		if err = json.Unmarshal(by, &m); err != nil {
			return nil, err
		}
		ms[i] = m
	}
	return ms, nil
}

func falseKeys(m map[string]bool) []string {
	var names []string
	for k, v := range m {
		if v {
			continue
		}
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

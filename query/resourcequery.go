package query

import (
	"context"
	"fmt"
	"github.com/eurozulu/pempal/logging"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/utils"
	"strings"
)

type ResourceQuery struct {

	// Fields lists the field names to output.
	// If none are specified, the default fields are used.
	// Default fields are defined by the Types set.
	// If no types are set, basic, generic fields are used.
	// When set, fields relevant to the type(s) are used.
	// A Field name may be preceeded with a '+' sign to add that field to the default fields.
	// If any given field name has a plus, the default fields will be shown.
	// If all the given fields have no plus, then only the given fields are shown.
	Fields []string `flag:"fields,f"`

	// Types lists the resource types to be included in the query
	// If none are set, all resource types are queried.
	// Default Fields are defined by the Types set.
	// These can be adjusted using the Fields
	Types []model.ResourceType `flag:"type,t"`

	// NonRecursive prevents the query scanning subdirectories of any given path.
	// By default, all subdirectories are scanned.
	NonRecursive bool `flag:"nonrecursive,nr"`
	Limit        int  `flag:"limit,l"`
	Conditions   []Condition
}

func (rq ResourceQuery) QueryAll(path ...string) []ResourceProperties {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	var found []ResourceProperties
	for p := range rq.Query(ctx, path...) {
		found = append(found, p...)
	}
	return found
}

func (rq ResourceQuery) Query(ctx context.Context, path ...string) <-chan []ResourceProperties {
	ch := make(chan []ResourceProperties)
	go func(ch chan<- []ResourceProperties, path []string) {
		defer close(ch)
		scan := resources.NewPemScan(rq.Types...)
		scan.NonRecursive = scan.NonRecursive
		colNames := rq.ColumnNames()
		count := 0
		for pemRes := range scan.ScanPath(ctx, path...) {
			props, err := parseResourceProperties(pemRes)
			if err != nil {
				logging.Warning("ResourceQuery", "failed to parse properties in %s. %v", pemRes.Path, err.Error())
				continue
			}
			props = rq.filterPropertiesByConditions(props)
			if len(props) == 0 {
				continue
			}
			props = rq.filterPropertyNames(props, colNames)

			select {
			case <-ctx.Done():
				return
			case ch <- props:
				count++
			}

			if rq.Limit > 0 && count >= rq.Limit {
				return
			}
		}
	}(ch, path)
	return ch
}

func (rq ResourceQuery) ColumnNames() []string {
	var names []string
	if rq.usingDefaultFields() {
		names = append(names, DefaultColumnNames(rq.Types...)...)
	}
	if len(rq.Fields) > 0 {
		cleanFields := make([]string, len(rq.Fields))
		for i, field := range rq.Fields {
			if strings.HasPrefix(field, "+") {
				field = strings.TrimPrefix(field, "+")
			}
			cleanFields[i] = field
		}
		names = append(cleanFields, names...)
	}
	return names
}

func (rq ResourceQuery) filterPropertyNames(props []ResourceProperties, names []string) []ResourceProperties {
	var isWildCard bool
	names, isWildCard = containsWildCardField(names)
	result := make([]ResourceProperties, len(props))
	for i, prop := range props {
		if isWildCard {
			// add to existing
			result[i] = prop
		} else {
			result[i] = ResourceProperties{}
		}
		// This ensures all names exist, even if null
		copyKeys(result[i], prop, names)
	}
	return result
}

func (rq ResourceQuery) filterPropertiesByConditions(props []ResourceProperties) []ResourceProperties {
	if len(rq.Conditions) == 0 {
		return props
	}
	var found []ResourceProperties
	for _, prop := range props {
		if !rq.compareWithConditions(prop) {
			continue
		}
		found = append(found, prop)
	}
	return found
}

func (rq ResourceQuery) compareWithConditions(props ResourceProperties) bool {
	for _, cond := range rq.Conditions {
		if !cond.IsTrue(props) {
			logging.Trace("ResourceQuery", "condition %s rejected properties from %s", cond, props["filename"])
			return false
		}
	}
	return true
}

func containsWildCardField(names []string) ([]string, bool) {
	for i, name := range names {
		if name != "*" {
			continue
		}
		var end []string
		if i+1 < len(names) {
			end = names[i+1:]
		}
		return append(names[:i], end...), true
	}
	return names, false
}

// usingDefaultFields checks if the predefined column names are being used.
// predefined names are used when either no -fields names are set or none of the names set are preceeded with the '+'.
func (rq ResourceQuery) usingDefaultFields() bool {
	if len(rq.Fields) == 0 {
		return true
	}
	for _, field := range rq.Fields {
		if strings.HasPrefix(field, "+") {
			return true
		}
	}
	return false
}

func copyKeys(dst, src ResourceProperties, keys []string) {
	for _, key := range keys {
		dst[key] = src[key]
	}
}

func (rq ResourceQuery) formatRow(cols []utils.Column, values map[string]interface{}) []string {
	var found []string
	for _, col := range cols {
		v, ok := values[col.Name]
		if !ok {
			v = ""
		}
		found = append(found, fmt.Sprintf("%v", v))
	}
	return found
}

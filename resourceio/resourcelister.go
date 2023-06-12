package resourceio

import (
	"bytes"
	"fmt"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/utils"
	"github.com/go-yaml/yaml"
	"strings"
)

const LocationTitle = "location"

// ResourceLister parses Resources from each Location into a slice of pre-defined property values
type ResourceLister interface {
	// Limits the listing to a specific resource type.  When empty, all resources are listed.
	ResourceTypes() []model.ResourceType

	// Names the fields of each resource to be listed.  If a resource type does not contain a named field it will show as empty.
	Fields() []string

	// List all the preset Fields from the given locations channel
	List(locs <-chan ResourceLocation) <-chan []string
}

type resourceLister struct {
	resourceTypes []model.ResourceType
	fields        []string
	yamlFormatter ResourceFormatter
	query         Query
}

func (rl resourceLister) Fields() []string {
	return rl.fields
}

func (rl resourceLister) ResourceTypes() []model.ResourceType {
	return rl.resourceTypes
}

func (rl resourceLister) List(locs <-chan ResourceLocation) <-chan []string {
	ch := make(chan []string)
	go func(locs <-chan ResourceLocation) {
		defer close(ch)
		for loc := range locs {
			err := rl.parseLocation(loc, ch)
			if err != nil {
				logger.Error("failed to parse location %s  %v", loc.Location(), err)
				break
			}
		}
	}(locs)
	return ch
}

func (rl resourceLister) parseLocation(loc ResourceLocation, out chan<- []string) error {
	for _, res := range loc.Resources(rl.resourceTypes...) {
		m, err := rl.resourceToMap(res)
		if err != nil {
			return err
		}
		m[LocationTitle] = loc.Location()
		if rl.query != nil {
			ok, err := rl.matchResource(m, false)
			if err != nil {
				return err
			}
			if !ok {
				continue
			}
		}
		out <- readDotKeyValues(rl.fields, m)
	}
	return nil
}

func (rl resourceLister) matchResource(m map[string]interface{}, all bool) (bool, error) {
	qm, err := rl.query.Match(m)
	if err != nil {
		return false, err
	}
	for _, b := range qm {
		if b && !all {
			return true, nil
		}
		if !b && all {
			return false, nil
		}
	}
	return all, nil
}

// resourceToMap flips a PEM Resource into a map of properties.
// The resource is marshalled into yaml and unmarshalled back into a map.
func (rl resourceLister) resourceToMap(res model.Resource) (map[string]interface{}, error) {
	yamlData, err := rl.yamlFormatter.FormatResources(res)
	if err != nil {
		return nil, err
	}
	yamlDec := yaml.NewDecoder(bytes.NewBuffer(yamlData))
	m := utils.YamlMap{} // map[string]interface{}{}
	if err := yamlDec.Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}

func readDotKeyValues(keys []string, m map[string]interface{}) []string {
	values := make([]string, len(keys))
	for i, key := range keys {
		v, ok := readValueFromDotKey(key, m)
		if !ok {
			continue
		}
		values[i] = fmt.Sprintf("%v", v)
	}
	return values
}

func readValueFromDotKey(dotKey string, m map[string]interface{}) (interface{}, bool) {
	ks := strings.Split(dotKey, ".")
	v, ok := m[ks[0]]
	if !ok {
		return nil, false
	}
	if len(ks) == 1 {
		// no dots found, return value found (if any)
		return v, ok
	}
	// key has dotted subname: ensure parent value is a map
	vm, ok := v.(map[string]interface{})
	if !ok {
		return nil, false
	}
	return readValueFromDotKey(strings.Join(ks[1:], "."), vm)
}

func NewResourceLister(fields []string, query Query, resourceTypes ...model.ResourceType) ResourceLister {
	if len(fields) == 0 {
		fields = []string{LocationTitle}
	}
	return &resourceLister{
		resourceTypes: resourceTypes,
		fields:        fields,
		yamlFormatter: NewResourceFormatter(FormatYAML),
		query:         query,
	}
}

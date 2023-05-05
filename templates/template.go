package templates

import (
	"bytes"
	"fmt"
	"github.com/go-yaml/yaml"
	"strings"
	gotemplate "text/template"
)

type Template interface {
	// Raw returns the original, raw bytes of the template, prior to any parsing and attribution
	Raw() []byte

	// String returns a string of the template bytes after attribution (tags applied)
	fmt.Stringer

	// Tags returns any #tags found in this template
	Tags() Tags

	Apply(in interface{}) error
}

type yamlTemplate struct {
	raw    []byte
	tags   Tags
	parsed []byte
}

func (t yamlTemplate) Raw() []byte {
	return t.raw
}

func (t yamlTemplate) String() string {
	return string(t.parsed)
}

func (t yamlTemplate) Tags() Tags {
	return t.tags
}

// Apply applies this template to the given object
func (t *yamlTemplate) Apply(out interface{}) error {
	return yaml.Unmarshal(t.parsed, out)
}

func mergeTags(tags Tags, templates []Template) Tags {
	for _, temp := range templates {
		tags = append(tags, temp.Tags()...)
	}
	// Ensure no repeated tags found
	unique := map[string]bool{}
	var utags Tags
	for _, tg := range tags {
		k := tg.String()
		if unique[k] {
			continue
		}
		unique[k] = true
		utags = append(utags, tg)
	}
	return utags
}

func mergeTemplatesToYaml(templates []Template) ([]byte, error) {
	m := map[string]interface{}{}
	for _, et := range templates {
		if err := et.Apply(&m); err != nil {
			return nil, err
		}
	}
	return yaml.Marshal(&m)
}

func executeGoTemplate(text string, data map[string]interface{}) ([]byte, error) {
	gt, err := gotemplate.New("template-manager").Parse(text)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	if err = gt.Execute(buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func containsGoTemplates(text string) bool {
	i := strings.Index(text, "{{")
	if i < 0 {
		return false
	}
	return strings.Index(text[i+2:], "}}") >= 0
}

func newYamlTemplate(raw []byte, extends []Template, imports map[string]interface{}) (Template, error) {
	tags, base := parseTags(raw)
	tags = mergeTags(tags, extends)
	t := &yamlTemplate{
		raw:    raw,
		parsed: base,
		tags:   tags,
	}
	parsed, err := mergeTemplatesToYaml(append(extends, t))
	if err != nil {
		return nil, fmt.Errorf("failed to merge extended templates as yaml  %v", err)
	}
	t.parsed = parsed

	if containsGoTemplates(string(parsed)) {
		// Apply this template to imports so available in data
		if err := t.Apply(imports); err != nil {
			return nil, err
		}
		parsed, err := executeGoTemplate(string(parsed), imports)
		if err != nil {
			return nil, fmt.Errorf("failed to execute templates  %v", err)
		}
		t.parsed = parsed
	}

	return t, nil
}

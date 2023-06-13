package templates

import (
	"bytes"
	"fmt"
	"github.com/go-yaml/yaml"
	"strings"
)

// ParseInlineTemplate parses given, comma delimited lines of text into a template.
// names may use dot notation to indicate bud properties (e.g. subject.common-name)
// line may be preceeded with #tags
func ParseInlineTemplate(s string) (Template, error) {
	var tags []string
	var readTags bool

	lines := strings.Split(s, ",")
	m := map[string]interface{}{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		if !readTags && strings.HasPrefix(line, TAG_TOKEN) {
			tags = append(tags, line)
			continue
		}
		readTags = true
		ls := strings.SplitN(line, ":", 2)
		key := strings.TrimSpace(ls[0])
		var val string
		if len(ls) > 1 {
			val = strings.TrimSpace(ls[1])
		}
		if err := addDotKeyValue(key, val, m); err != nil {
			return nil, err
		}
	}
	buf := bytes.NewBuffer(nil)
	for _, tag := range tags {
		buf.WriteString(tag)
		buf.WriteRune('\n')
	}

	if err := yaml.NewEncoder(buf).Encode(&m); err != nil {
		return nil, err
	}
	return NewTemplate(buf.Bytes())
}

func addDotKeyValue(key string, v interface{}, m map[string]interface{}) error {
	i := strings.IndexRune(key, '.')
	if i < 0 {
		m[key] = v
		return nil
	}
	iv, ok := m[key[:i]]
	if !ok {
		iv = map[string]interface{}{}
		m[key[:i]] = iv
	}
	im, ok := iv.(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to apply property %s, as already holds an incompatible type")
	}
	if i+1 >= len(key) {
		return fmt.Errorf("key '%s' is invalid. no name found after dot", key)
	}
	return addDotKeyValue(key[i+1:], v, im)
}

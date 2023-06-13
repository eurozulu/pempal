package templates

import (
	"bytes"
	"github.com/eurozulu/pempal/utils"
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
	m := utils.FlatMap{}
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
		key := ls[0]
		var val *string
		if len(ls) > 1 {
			val = &ls[1]
		}
		m[key] = val
	}
	buf := bytes.NewBuffer(nil)
	for _, tag := range tags {
		buf.WriteString(tag)
		buf.WriteRune('\n')
	}
	em := m.Expand()
	if err := yaml.NewEncoder(buf).Encode(&em); err != nil {
		return nil, err
	}
	return NewTemplate(buf.Bytes())
}

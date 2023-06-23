package templates

import (
	"bytes"
	"github.com/eurozulu/pempal/utils"
	"github.com/go-yaml/yaml"
	"strings"
)

func ParseInlineTemplate(s string) (Template, error) {
	var tags []string
	var tagsAllRead bool

	lines := strings.Split(s, ",")
	m := utils.FlatMap{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if !tagsAllRead && strings.HasPrefix(line, TAG_TOKEN) {
			tags = append(tags, line)
			continue
		}
		tagsAllRead = true
		ls := strings.SplitN(line, ":", 2)
		key := strings.TrimSpace(ls[0])
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
	if err := yaml.NewEncoder(buf).Encode(&m); err != nil {
		return nil, err
	}
	return NewTemplate(buf.Bytes())
}

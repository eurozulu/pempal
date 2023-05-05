package templates

import (
	"bytes"
	"strings"
)

const TAG_TOKEN = "#"
const TAG_EXTENDS = "extends"
const TAG_IMPORTS = "imports"

type Tags []Tag

type Tag struct {
	Name  string
	Value string
}

func (tag Tag) String() string {
	return strings.Join([]string{TAG_TOKEN, tag.Name, " ", tag.Value}, "")
}

func (tag Tag) ParseAsImport() (name string, alias string) {
	names := strings.SplitN(tag.Value, " ", 2)
	alias = strings.TrimSpace(names[0])
	if len(names) > 1 {
		name = strings.TrimSpace(names[1])
	} else {
		name = alias
	}
	return name, alias
}

func (ts Tags) TagsByName(name string) []Tag {
	var tags []Tag
	name = strings.TrimSpace(name)
	for _, tg := range ts {
		if tg.Name != name {
			continue
		}
		tags = append(tags, tg)
	}
	return tags
}

func tagValues(tags []Tag) []string {
	ss := make([]string, len(tags))
	for i, t := range tags {
		ss[i] = strings.TrimSpace(t.Value)
	}
	return ss
}

func parseTag(s string) Tag {
	ss := strings.SplitN(s, " ", 2)
	n := strings.TrimSpace(strings.TrimLeft(ss[0], TAG_TOKEN)) // Trim off the tag token
	var v string
	if len(ss) > 1 {
		v = strings.TrimSpace(ss[1])
	}
	return Tag{
		Name:  n,
		Value: v,
	}
}

func parseTags(data []byte) (Tags, []byte) {
	var tags Tags
	// search first lines for tokens
	for len(data) > 0 {
		eol := bytes.Index(data, []byte{'\n'})
		if eol < 0 {
			eol = len(data)
		}
		line := strings.TrimSpace(string(data[0:eol]))
		if len(line) > 0 {
			// first line NOT starting with a token stops search
			if !strings.HasPrefix(line, TAG_TOKEN) {
				break
			}
			tags = append(tags, parseTag(line))
		}
		data = bytes.TrimLeft(data[eol:], "\n")
	}
	return tags, data
}

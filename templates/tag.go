package templates

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

const TAG_TOKEN = "#"
const TAG_EXTENDS = "extends"
const TAG_IMPORTS = "imports"

type Tags []Tag

type Tag interface {
	Name() string
	Value() string
	fmt.Stringer
}

type uniqueNames map[string]bool

// tag represents a single key/value pair which optionally appears at the start of a template.
// Every tag starts with the tag token '#' followed directly by the tag name.
// Following the name a space delimits the value to assign to the tag.
// e.g.. #extend mytemplate
// Valid tags are:
// #extends
// #imports
type tag struct {
	name  string
	value string
}

func (t tag) Name() string {
	return t.name
}

func (t tag) Value() string {
	return t.value
}

func (t tag) String() string {
	return strings.Join([]string{TAG_TOKEN, t.name, " ", t.value}, "")
}

func (ts Tags) ContainsTag(name string) bool {
	return ts.tagIndex(name) >= 0
}

// TagByName gets the tag for the given name
// If the tag name is unknown, returns an empty tag
func (ts Tags) TagByName(name string) Tag {
	i := ts.tagIndex(name)
	if i < 0 {
		return nil
	}
	return ts[i]
}

func (ts Tags) tagIndex(name string) int {
	name = strings.TrimSpace(name)
	for i, tg := range ts {
		if tg.Name() == name {
			return i
		}
	}
	return -1
}

func parseTag(s string) Tag {
	ss := strings.SplitN(s, " ", 2)
	n := strings.TrimSpace(strings.TrimLeft(ss[0], TAG_TOKEN)) // Trim off the tag token
	var v string
	if len(ss) > 1 {
		v = strings.TrimSpace(ss[1])
	}
	return &tag{
		name:  n,
		value: v,
	}
}

func parseImport(value string) (name string, alias string) {
	names := strings.SplitN(value, " ", 2)
	name = strings.TrimSpace(names[0])
	if len(names) > 1 {
		a, err := strconv.Unquote(strings.TrimSpace(names[1]))
		if err != nil {
			a = strings.TrimSpace(names[1])
		}
		alias = a
	} else {
		alias = name
	}
	return name, alias
}

func parseTags(data []byte) (Tags, []byte, error) {
	var tags Tags
	names := uniqueNames{}
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
			tag := parseTag(line)
			if names[tag.Name()] {
				return nil, nil, fmt.Errorf("duplicate tag name '%s'. tag names must be unique.")
			}
			names[tag.Name()] = true
			tags = append(tags, parseTag(line))
		}
		data = bytes.TrimLeft(data[eol:], "\n")
	}
	return tags, data, nil
}

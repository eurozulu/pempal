package templates

import (
	"bytes"
	"strconv"
	"strings"
)

const TAG_TOKEN = "#"
const TAG_EXTENDS = "extends"
const TAG_IMPORTS = "imports"

type Tags []Tag

// Tag represents a single key/value pair which optionally appears at the start of a template.
// Every tag starts with the tag token '#' followed directly by the tag name.
// Following the name a space delimits the value to assign to the tag.
// e.g.. #extend mytemplate
// Valid tags are:
// #extends
// #imports
type Tag struct {
	Name  string
	Value string
}

func (tag Tag) String() string {
	return strings.Join([]string{TAG_TOKEN, tag.Name, " ", tag.Value}, "")
}

// ParseAsImport will attempt to split the Tag value into a name and alias.
// imports may be expressed as a simple template name or be given an alternative name, using a space delited name:
// e.g. #imports default-certificate defcert
// This will import the "default-certificate" and label is as 'defcert'
func (tag Tag) ParseAsImport() (name string, alias string) {
	names := strings.SplitN(tag.Value, " ", 2)
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

// TagsByName gets all the tags for the given name
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

// tagValues gets all the values if the given tags
func tagValues(tags Tags) []string {
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

package tools

import (
	"encoding/base64"
	"fmt"
	"strings"
	"unicode"
)

func StringsWithSuffix(s []string, suffix string) []string {
	var found []string
	for _, sz := range s {
		if len(sz) < len(suffix) {
			continue
		}
		if !strings.EqualFold(sz[len(sz)-len(suffix):], suffix) {
			continue
		}
		found = append(found, sz)
	}
	return found
}

func StringsWithPrefix(s []string, prefix string) []string {
	var found []string
	for _, sz := range s {
		if len(sz) < len(prefix) {
			continue
		}
		if !strings.EqualFold(sz[:len(prefix)], prefix) {
			continue
		}
		found = append(found, sz)
	}
	return found
}

func CleanQuotes(s string) string {
	s = strings.TrimLeft(s, "\"'")
	return strings.TrimRight(s, "\"'")
}

func SingleQuote(s string) string {
	if s == "" {
		return "''"
	}
	if s[0] != '\'' {
		s = string(append([]rune{'\''}, []rune(s)...))
	}
	if s[len(s)-1] != '\'' {
		s = string(append([]rune(s), '\''))
	}
	return s
}

func StringerToString[T fmt.Stringer](s ...T) []string {
	ss := make([]string, len(s))
	for i, sz := range s {
		ss[i] = sz.String()
	}
	return ss
}

func ToTitle(s string) string {
	lastSpace := true
	rz := make([]rune, len(s))
	for i, r := range []rune(s) {
		if unicode.IsSpace(r) {
			lastSpace = true
			rz[i] = r
			continue
		}
		if lastSpace && !unicode.IsUpper(r) {
			r = unicode.ToUpper(r)
		} else if !lastSpace && unicode.IsUpper(r) {
			r = unicode.ToLower(r)
		}
		rz[i] = r
		lastSpace = false
	}
	return string(rz)
}

func DeCodeBase64(bytes []byte) []byte {
	l := base64.StdEncoding.DecodedLen(len(bytes))
	decoded := make([]byte, l)
	_, err := base64.StdEncoding.Decode(decoded, bytes)
	if err == nil {
		return decoded
	}
	return bytes
}

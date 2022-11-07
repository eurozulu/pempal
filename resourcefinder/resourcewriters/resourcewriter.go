package resourcewriters

import (
	"bytes"
	"encoding/pem"
	"io"
	"pempal/resources"
	"strings"
)

const longestTypeString = 20

type ResourceWriter interface {
	Write(r resources.Resource) error
}

type pemWriter struct {
	out io.Writer
}

func (rw pemWriter) Write(r resources.Resource) error {
	return pem.Encode(rw.out, r.Pem())
}

type lineWriter struct {
	out io.Writer
}

func (l lineWriter) Write(r resources.Resource) error {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(padString(r.Type().String(), longestTypeString, " "))
	buf.WriteString(r.Location())
	buf.WriteRune('\n')
	_, err := l.out.Write(buf.Bytes())
	return err
}

func padString(s string, size int, pad string) string {
	pl := size - len(s)
	if pl < 0 {
		// crop string if longer than pad size
		return s[:size]
	}
	return strings.Join([]string{s, strings.Repeat(pad, pl)}, "")
}

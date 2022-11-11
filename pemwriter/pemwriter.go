package pemwriter

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"pempal/pemtypes"
	"strings"
)

type PemWriter interface {
	Write(resource pemtypes.PemResource) error
}
type PemWriterType int

const (
	Unknown PemWriterType = iota
	Text
	Pem
	Der
	Yaml
)

var PemWriterNames = [...]string{
	Text: "text",
	Pem:  "pem",
	Der:  "der",
	Yaml: "yaml",
}

type pemWriter struct {
	out io.Writer
}

func (p PemWriterType) String() string {
	return PemWriterNames[p]
}
func ParseWriterType(s string) PemWriterType {
	i := typeIndex(s)
	if i < 0 {
		return Unknown
	}
	return PemWriterType(i)
}

func (pw pemWriter) Write(resource pemtypes.PemResource) error {
	text, err := resource.MarshalText()
	if err != nil {
		return err
	}
	_, err = pw.out.Write(text)
	return err
}

type derWriter struct {
	out io.Writer
}

func (dw derWriter) Write(resource pemtypes.PemResource) error {
	data, err := resource.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = dw.out.Write(data)
	return err
}

type yamlWriter struct {
	encoder *yaml.Encoder
}

func (yw yamlWriter) Write(resource pemtypes.PemResource) error {
	return yw.encoder.Encode(resource)
}

type stringWriter struct {
	out io.Writer
}

func (sw stringWriter) Write(resource pemtypes.PemResource) error {
	_, err := fmt.Fprintf(sw.out, "%s", resource.String())
	return err
}

func typeIndex(s string) int {
	s = strings.ToLower(s)
	for i, ss := range PemWriterNames {
		if ss == s {
			return i
		}
	}
	return -1
}

func NewPemWriter(writerType PemWriterType, out io.Writer) PemWriter {
	switch writerType {
	case Text:
		return stringWriter{out: out}
	case Pem:
		return pemWriter{out: out}
	case Der:
		return derWriter{out: out}
	case Yaml:
		return yamlWriter{encoder: yaml.NewEncoder(out)}
	default:
		return stringWriter{out: out}
	}
}

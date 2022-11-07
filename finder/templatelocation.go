package finder

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"pempal/templates"
)

type templateLocation struct {
	path  string
	temps []templates.Template
}

func (r templateLocation) Path() string {
	return r.path
}

func (r templateLocation) MarshalText() (text []byte, err error) {
	buf := bytes.NewBuffer(nil)
	if _, err := fmt.Fprintln(buf, r.Path()); err != nil {
		return nil, err
	}
	if err := yaml.NewEncoder(buf).Encode(r.temps); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

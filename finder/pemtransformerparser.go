package finder

import (
	"bytes"
	"encoding/pem"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"pempal/encoders"
	"pempal/pemtypes"
	"pempal/templates"
)

type pemTransformerParser struct {
	pParser *pemParser
	filter  PemQuery
}

type PemQuery *string

func (p pemTransformerParser) MatchPath(path string) bool {
	return p.pParser.MatchPath(path)
}

func (p pemTransformerParser) FilterLocation(rl Location) Location {
	tl, ok := rl.(*templateLocation)
	if !ok {
		return nil
	}
	if p.filter == nil {
		return tl
	}
	var temps []templates.Template
	f := []byte(*p.filter)
	for _, t := range tl.temps {
		txt, err := marshalTemplateText(t)
		if err != nil {
			log.Println(err)
			continue
		}
		if !bytes.Contains(txt, f) {
			continue
		}
		temps = append(temps, t)
	}
	if len(temps) == 0 {
		return nil
	}
	return &templateLocation{
		path:  rl.Path(),
		temps: temps,
	}
}

func (p pemTransformerParser) Parse(path string, data []byte) (Location, error) {
	pl, err := p.pParser.Parse(path, data)
	if err != nil {
		return nil, err
	}
	pl = p.pParser.FilterLocation(pl)
	if pl == nil {
		return nil, fmt.Errorf("no pems found in %s", path)
	}
	ppl := pl.(*pemLocation)
	temps := p.parseTemplates(ppl.pems)
	if len(temps) == 0 {
		return nil, nil
	}
	return &templateLocation{
		path:  path,
		temps: temps,
	}, nil
}
func (p pemTransformerParser) parseTemplates(blocks []*pem.Block) []templates.Template {
	var temps []templates.Template
	for _, b := range blocks {
		pt := pemtypes.ParsePEMType(b.Type)
		e, err := encoders.NewEncoder(pt)
		if err != nil {
			log.Println(err)
			continue
		}
		t, err := e.Encode(b)
		if err != nil {
			log.Println(err)
			continue
		}
		temps = append(temps, t)
	}
	return temps
}

func marshalTemplateText(t templates.Template) (text []byte, err error) {
	buf := bytes.NewBuffer(nil)
	if err := yaml.NewEncoder(buf).Encode(t); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func newPemTransformerParser(q PemQuery, pemType ...pemtypes.PEMType) *pemTransformerParser {
	if len(pemType) == 0 {
		pemType = []pemtypes.PEMType{pemtypes.Unknown}
	}
	pp := newPemParser(pemType...)
	return &pemTransformerParser{
		filter:  q,
		pParser: pp}
}

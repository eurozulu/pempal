package templates

import (
	"context"
	"encoding/pem"
	"fmt"
	"pempal/pemreader"
	"strings"
)

const AppendKeyPrefix = "+"

type TemplateBuilder interface {
	Add(p string) error
	Templates() []Template
	Build() (Template, error)
}

type builder struct {
	temps []Template
	pr    *pemreader.PemScanner
}

func (tb builder) Templates() []Template {
	return tb.temps
}

func (tb builder) Build() (Template, error) {
	mt := mergeTemplates(tb.temps...)
	missing := mt.RequiredNames()
	if len(missing) > 0 {
		return nil, fmt.Errorf("The following properties are required: %v", missing)
	}

	funcs := mt.funcNames()
	if len(funcs) > 0 {
		return nil, fmt.Errorf("template functions not yet supported")
	}
	return mt, nil
}

func (tb *builder) Add(names ...string) error {
	var err error
	for _, p := range names {
		if !strings.HasPrefix(p, FileTag) {
			err = tb.addPemResource(p)
		} else {
			err = tb.addTemplate(p)
		}
		if err != nil {
			return fmt.Errorf("failed to open %s  %v", p, err)
		}
	}
	return nil
}

func (tb *builder) addTemplate(p string) error {
	t, err := FindTemplate(p)
	if err != nil {
		return err
	}
	tb.temps = append(tb.temps, t)
	return nil
}

func (tb builder) addPemResource(p string) error {
	pb, err := tb.findPemResource(p)
	if err != nil {
		return err
	}
	t, err := ParseBlock(pb)
	if err != nil {
		return err
	}
	tb.temps = append(tb.temps, t)
	return nil
}

func (tb builder) findPemResource(p string) (*pem.Block, error) {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	pemIn := tb.pr.Find(ctx, p)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case pb, ok := <-pemIn:
			if !ok {
				return nil, fmt.Errorf("%s not found", p)
			}
			return pb, nil
		}
	}
}

func mergeTemplates(ms ...Template) Template {
	nt := Template{}
	for _, m := range ms {
		for k, v := range nt {
			vs := m.Value(k)
			if vs == requiredPrefix {
				// only overwrite with required symbol when no value present
				_, ok := m[k]
				if ok {
					continue
				}
			}
			m[k] = v
		}
	}
	return nt
}

func NewTemplateBuilder() *builder {
	return &builder{
		temps: nil,
		pr: &pemreader.PemScanner{
			AddLocationHeader: true,
			Recursive:         false,
		},
	}
}

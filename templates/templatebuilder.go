package templates

import (
	"bytes"
	"fmt"
	"github.com/eurozulu/pempal/logging"
	"gopkg.in/yaml.v2"
)

var ErrTemplateBuildEmpty = fmt.Errorf("no templates set, nothing to build")

type TemplateBuilder interface {
	// AddTemplate adds the named template (and any parent templates it extends) to the builder.
	// Templates are added in order of heritage, with the named, leaf template being last.
	// Any template in the resulting chain which already exists in the builder is ignored.
	AddTemplate(name ...string) error

	AddNewTemplate(name string, data []byte) error

	// BaseTemplate returns the first BaseTemplate in the stack.
	BaseTemplate() Template

	Build() (Template, error)
}

type templateBuilder struct {
	tempLib    TemplateLib
	buildStack []Template
}

func (tb *templateBuilder) Templates() []Template {
	return tb.buildStack
}

func (tb *templateBuilder) ClearTemplates() {
	tb.buildStack = nil
}

func (tb *templateBuilder) BaseTemplate() Template {
	for _, t := range tb.buildStack {
		if bt := tb.tempLib.BaseTemplate(t.Name()); bt == nil {
			continue
		}
		return t
	}
	return nil
}

func (tb *templateBuilder) AddTemplate(name ...string) error {
	for _, n := range name {
		tps, err := tb.tempLib.GetTemplates(n)
		if err != nil {
			return err
		}
		if err := tb.addTemplates(tps); err != nil {
			return err
		}
	}
	return nil
}

func (tb *templateBuilder) AddNewTemplate(name string, data []byte) error {
	return tb.addTemplates([]Template{
		&template{
			name: name,
			data: data,
		},
	})
}

func (tb *templateBuilder) Build() (Template, error) {
	if len(tb.buildStack) == 0 {
		return nil, ErrTemplateBuildEmpty
	}
	target := tb.BaseTemplate()
	if target == nil {
		return nil, fmt.Errorf("template %q is not a known base template", target.Name())
	}
	for _, t := range tb.buildStack {
		if t == target {
			continue
		}
		if err := ApplyTemplateToTarget(t, target); err != nil {
			return nil, err
		}
	}
	return target, nil
}

func (tb *templateBuilder) addTemplates(tps []Template) error {
	for _, t := range tps {
		// Skip if that Name already set to build
		if tb.containsName(t.Name()) {
			logging.Debug("TemplateBuilder", "Ignoring adding template %s as already exists", t.Name())
			continue
		}
		logging.Debug("TemplateBuilder", "adding template %s", t.Name())
		tb.buildStack = append(tb.buildStack, t)
	}
	return nil
}

func (tb *templateBuilder) containsName(name string) bool {
	for _, t := range tb.buildStack {
		if t.Name() == name {
			return true
		}
	}
	return false
}

func ApplyTemplateToTarget(t Template, target interface{}) error {
	data := bytes.NewBuffer(nil)
	err := yaml.NewEncoder(data).Encode(t)
	if err != nil {
		return fmt.Errorf("template %s is invalid. %v", t.Name(), err)
	}
	if err = yaml.Unmarshal(data.Bytes(), target); err != nil {
		return fmt.Errorf("failed to apply template %q  %v", t.Name(), err)
	}
	return nil
}

func NewTemplateBuilder(templatePath ...string) TemplateBuilder {
	return &templateBuilder{
		tempLib:    NewTemplateLib(templatePath...),
		buildStack: nil,
	}
}

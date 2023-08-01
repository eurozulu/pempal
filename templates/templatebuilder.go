package templates

type TemplateBuilder interface {
	Clear()
	Add(t ...Template)
	Build() Template
	Templates() []Template
}

type templateBuilder struct {
	temps []Template
}

func (tb *templateBuilder) Clear() {
	tb.temps = nil
}

func (tb *templateBuilder) Add(t ...Template) {
	if len(t) > 0 {
		tb.temps = append(tb.temps, t...)
	}
}

func (tb templateBuilder) Build() Template {
	temp := Template{}
	for _, t := range tb.temps {
		mergeTemplate(temp, t)
	}
	return temp
}

func (tb *templateBuilder) Templates() []Template {
	return tb.temps
}

func NewTemplateBuilder(t ...Template) TemplateBuilder {
	return &templateBuilder{temps: t}
}

func mergeTemplate(dst, src Template) {
	for k, v := range src {
		if v == "" {
			// don't overwrite with empty strings.
			if _, exists := dst[k]; exists {
				continue
			}
		}
		dst[k] = v
	}
}

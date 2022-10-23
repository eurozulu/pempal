package templates

type TemplateParser interface {
	Parse(by []byte) (Template, error)
}

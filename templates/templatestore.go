package templates

type TemplateStore interface {
	Names(s ...string) []string

	// SaveTemplate adds a new template to the store under the given name.
	// returns error if the name already exists.
	SaveTemplate(name string, t Template) error

	// DeleteTemplate removes a named template from the store.
	//  returns error if the name is not known
	DeleteTemplate(name string) error
}

func NewTemplateStore(rootpath string) (TemplateStore, error) {
	tm, err := NewTemplateManager(rootpath)
	if err != nil {
		return nil, err
	}
	return tm.(TemplateStore), nil
}

package prompts

import (
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/ui"
)

type RequestTypeView struct {
	ui.TextList
}

func (rv RequestTypeView) Show() (resources.ResourceType, error) {
	win := ui.NewWindow("Select the type of resource", 0, 0)
	v, err := win.Show(&rv)
	if err != nil {
		return 0, err
	}
	return resources.ParseResourceType(v.String()), nil
}

func typeNames(types []resources.ResourceType) []string {
	ss := make([]string, len(types))
	for i, t := range types {
		ss[i] = t.String()
	}
	return ss
}

func NewRequestTypeView(label, text string, types ...resources.ResourceType) *RequestTypeView {
	if len(types) == 0 {
		types = resources.ResourceTypes[2:]
	}
	return &RequestTypeView{*ui.NewTextList(label, text, typeNames(types)...)}
}

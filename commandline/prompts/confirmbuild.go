package prompts

import (
	"github.com/eurozulu/pempal/builders"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/templates"
	"github.com/eurozulu/pempal/ui"
	"strings"
)

var applyLabel = ui.NewLabelView("Apply", "")

type ConfirmBuild interface {
	Confirm(template templates.Template) (templates.Template, error)
}

type confirmBuild struct {
	buildType resources.ResourceType
}

func (cb confirmBuild) Confirm(template templates.Template) (templates.Template, error) {
	builder, err := builders.NewBuilder(cb.buildType)
	if err != nil {
		return nil, err
	}

	tb := templates.NewTemplateBuilder(template)
	window := ui.NewWindow("", 0, 0)
	for {
		t := tb.Build()
		errs := builder.Validate(t)
		hasErrors := len(errs) > 0
		view := cb.buildChildView(t)

		if hasErrors {
			// Set text to first error to select it.
			firstErr := errorsIntoChildViews(view, errs)
			view.SetText(firstErr)

		} else {
			setChildViewsHidden(view, false)
			// Add append label to end of list
			view = ui.NewParentView(view.Label(), view.String(), append(view.ChildViews(), applyLabel)...)
		}
		sv, err := window.Show(view)
		if err != nil {
			return nil, err
		}
		if selectedView(sv) == applyLabel {
			break
		}
		tb.Add(map[string]string{sv.Label(): sv.String()})
	}
	// strip original template to return just what's changed
	var updates []templates.Template
	if len(tb.Templates()) > 1 {
		updates = tb.Templates()[1:]
	} else {
		return nil, err
	}
	return templates.NewTemplateBuilder(updates...).Build(), nil
}

func selectedView(v ui.View) ui.View {
	if pv, ok := v.(ui.ParentView); ok {
		if pv.SelectedIndex() >= 0 && pv.SelectedIndex() < len(pv.ChildViews()) {
			return pv.ChildViews()[pv.SelectedIndex()]
		}
	}
	return nil
}

func (cb confirmBuild) buildChildView(t templates.Template) ui.ParentView {
	pv, _ := createResourceTypeView(cb.buildType)
	valuesIntoChildViews(pv, t)
	return pv
}

func valuesIntoChildViews(view ui.ParentView, values map[string]string) {
	for _, cv := range view.ChildViews() {
		if tv, ok := cv.(ui.TextView); ok {
			tv.SetText(values[cv.Label()])
		}
	}
}

func valuesFromChildViews(view ui.ParentView) map[string]string {
	values := map[string]string{}
	for _, cv := range view.ChildViews() {
		s := cv.String()
		if s == "" {
			continue
		}
		values[cv.Label()] = s
	}
	return values
}

func setChildViewsHidden(view ui.ParentView, hidden bool) {
	for _, cv := range view.ChildViews() {
		if hv, ok := cv.(ui.HiddenView); ok {
			hv.SetHidden(hidden)
		}
	}
}

func errorsIntoChildViews(view ui.ParentView, errs []error) string {
	var first string
	for _, cv := range view.ChildViews() {
		errIndex := errorIndexOfName(cv.Label(), errs)
		if errIndex < 0 {
			continue
		}
		if tv, ok := cv.(ui.TextView); ok {
			tv.SetText(errs[errIndex].Error())
			tv.SetColours(ui.ErrorColour)
			if first == "" {
				first = tv.String()
			}
		}
	}
	return first
}

func errorIndexOfName(label string, errs []error) int {
	for i, err := range errs {
		if strings.HasPrefix(err.Error(), label) {
			return i
		}
	}
	return -1
}

func NewConfirmBuild(rt resources.ResourceType) (ConfirmBuild, error) {
	if err := supportseResourceTypeView(rt); err != nil {
		return nil, err
	}
	return &confirmBuild{buildType: rt}, nil
}

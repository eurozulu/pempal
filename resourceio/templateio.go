package resourceio

import (
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/templates"
)

func ResourceLocationToTemplates(loc ResourceLocation, resourceType ...resources.ResourceType) ([]templates.Template, error) {
	var temps []templates.Template
	for _, r := range loc.Resources(resourceType...) {
		t, err := ResourceToTemplate(r)
		if err != nil {
			return nil, err
		}
		temps = append(temps, t)
	}
	return temps, nil
}

func ResourceToTemplate(r resources.Resource) (templates.Template, error) {
	dto, err := resources.NewResourceDTO(r)
	if err != nil {
		return nil, err
	}
	return resources.DTOToTemplate(dto, false)
}

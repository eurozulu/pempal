package resourceio

import (
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/templates"
)

func LoadTemplatesFromFile(path string) ([]templates.Template, error) {
	locs, err := ParseLocation(path)
	if err != nil {
		return nil, err
	}
	var temps []templates.Template
	for _, res := range locs.Resources() {
		dto, err := resources.NewResourceDTO(res)
		if err != nil {
			return nil, err
		}
		t, err := resources.DTOToTemplate(dto)
		if err != nil {
			return nil, err
		}
		temps = append(temps, t)
	}
	return temps, nil
}

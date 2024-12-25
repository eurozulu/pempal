package query

import (
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/templates"
)

type ResourceProperties map[string]interface{}

func parseResourceProperties(pr resources.PemResource) ([]ResourceProperties, error) {
	result := make([]ResourceProperties, len(pr.Content))
	for i, blk := range pr.Content {
		t, err := templates.TemplateOfPem(blk)
		if err != nil {
			return nil, err
		}
		m, err := templates.TemplateAsMap(t)
		if err != nil {
			return nil, err
		}
		m["type"] = model.ParseResourceTypeFromPEMType(blk.Type).String()
		m["filename"] = pr.Path
		result[i] = m
	}
	return result, nil
}

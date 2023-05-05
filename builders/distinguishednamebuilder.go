package builders

import (
	"crypto/x509/pkix"
	"fmt"
	"pempal/model"
	"pempal/templates"
)

type DistinguishedNameBuilder struct {
	dto model.DistinguishedNameDTO
}

func (db DistinguishedNameBuilder) ApplyTemplate(tp ...templates.Template) error {
	//TODO implement me
	panic("implement me")
}

func (db DistinguishedNameBuilder) Validate() []error {
	var errs []error
	for k := range db.RequiredValues() {
		errs = append(errs, fmt.Errorf("%s missing", k))
	}
	return errs
}

func (db DistinguishedNameBuilder) RequiredValues() map[string]interface{} {
	m := map[string]interface{}{}
	if db.dto.CommonName == "" {
		m["common-name"] = ""
	}
	return m
}

func (db DistinguishedNameBuilder) BuildName() (pkix.Name, error) {
	if errs := db.Validate(); len(errs) > 0 {
		return pkix.Name{}, fmt.Errorf("%s", collectErrorList(errs, ", "))
	}
	return db.dto.ToName(), nil
}

func newDistinguishedNameBuilder(dto *model.DistinguishedNameDTO) *DistinguishedNameBuilder {
	return &DistinguishedNameBuilder{dto: *dto}
}
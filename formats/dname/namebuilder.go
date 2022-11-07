package dname

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"gopkg.in/yaml.v3"
	"pempal/resources"
	"pempal/templates"
)

type nameBuilder struct {
	nameTemp *templates.NameTemplate
	location string
}

func (fm *nameBuilder) SetLocation(l string) {
	fm.location = l
}

func (fm *nameBuilder) AddTemplate(ts ...templates.Template) error {
	for _, t := range ts {
		by, err := yaml.Marshal(t)
		if err != nil {
			return err
		}
		if err = yaml.Unmarshal(by, fm.nameTemp); err != nil {
			return err
		}
	}
	return nil
}

func (fm nameBuilder) Template() templates.Template {
	return fm.nameTemp
}

func (fm nameBuilder) Build() (resources.Resources, error) {
	var n pkix.Name
	fm.ApplyTemplate(&n)
	by, err := asn1.Marshal(&n)
	if err != nil {
		return nil, err
	}
	return resources.Resources{resources.NewResource(fm.location, &pem.Block{
		Type:  resources.Name.String(),
		Bytes: by,
	})}, nil
}

func (fm nameBuilder) ApplyTemplate(name *pkix.Name) {
	if fm.nameTemp.CommonName != "" {
		name.CommonName = fm.nameTemp.CommonName
	}
	if fm.nameTemp.SerialNumber != "" {
		name.SerialNumber = fm.nameTemp.SerialNumber
	}
	if len(fm.nameTemp.Organization) > 0 {
		name.Organization = fm.nameTemp.Organization
	}
	if len(fm.nameTemp.OrganizationalUnit) > 0 {
		name.OrganizationalUnit = fm.nameTemp.OrganizationalUnit
	}
	if len(fm.nameTemp.Country) > 0 {
		name.Country = fm.nameTemp.Country
	}
	if len(fm.nameTemp.Locality) > 0 {
		name.Locality = fm.nameTemp.Locality
	}
	if len(fm.nameTemp.Province) > 0 {
		name.Province = fm.nameTemp.Province
	}
	if len(fm.nameTemp.StreetAddress) > 0 {
		name.StreetAddress = fm.nameTemp.StreetAddress
	}
	if len(fm.nameTemp.PostalCode) > 0 {
		name.PostalCode = fm.nameTemp.PostalCode
	}
}

func NewNameBuilder(template ...templates.Template) *nameBuilder {
	nb := &nameBuilder{}
	if len(template) > 0 {
		nb.AddTemplate(template...)
	}
	return nb
}

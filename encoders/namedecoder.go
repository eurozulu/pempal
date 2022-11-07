package encoders

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"pempal/pemtypes"
	"pempal/templates"
)

type nameDecoder struct {
}

func (nd nameDecoder) Decode(t templates.Template) (*pem.Block, error) {
	nt, ok := t.(*templates.NameTemplate)
	if !ok {
		return nil, fmt.Errorf("template is not for a Name")
	}
	var n pkix.Name
	nd.ApplyTemplate(nt, &n)
	by, err := asn1.Marshal(&n)
	if err != nil {
		return nil, err
	}
	return &pem.Block{
		Type:  pemtypes.Name.String(),
		Bytes: by,
	}, nil
}

func (nd nameDecoder) ApplyTemplate(t *templates.NameTemplate, name *pkix.Name) {
	if t.CommonName != "" {
		name.CommonName = t.CommonName
	}
	if t.SerialNumber != "" {
		name.SerialNumber = t.SerialNumber
	}
	if len(t.Organization) > 0 {
		name.Organization = t.Organization
	}
	if len(t.OrganizationalUnit) > 0 {
		name.OrganizationalUnit = t.OrganizationalUnit
	}
	if len(t.Country) > 0 {
		name.Country = t.Country
	}
	if len(t.Locality) > 0 {
		name.Locality = t.Locality
	}
	if len(t.Province) > 0 {
		name.Province = t.Province
	}
	if len(t.StreetAddress) > 0 {
		name.StreetAddress = t.StreetAddress
	}
	if len(t.PostalCode) > 0 {
		name.PostalCode = t.PostalCode
	}
}

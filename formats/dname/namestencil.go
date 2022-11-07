package dname

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"pempal/resources"
	"pempal/templates"
)

type nameStencil struct {
}

func (st nameStencil) MakeTemplate(r resources.Resource) (templates.Template, error) {
	blk := r.Pem()
	if blk == nil {
		return nil, nil
	}
	var n pkix.Name
	if _, err := asn1.Unmarshal(blk.Bytes, &n); err != nil {
		return nil, err
	}
	t := &templates.NameTemplate{}
	st.copyToTemplate(t, n)
	return t, nil
}

func (st nameStencil) copyToTemplate(t *templates.NameTemplate, name pkix.Name) {
	t.CommonName = name.CommonName
	t.SerialNumber = name.SerialNumber
	t.Organization = name.Organization
	t.OrganizationalUnit = name.OrganizationalUnit
	t.Country = name.Country
	t.Locality = name.Locality
	t.Province = name.Province
	t.StreetAddress = name.StreetAddress
	t.PostalCode = name.PostalCode
}

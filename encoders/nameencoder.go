package encoders

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"pempal/templates"
)

type nameEncoder struct {
}

func (ne nameEncoder) Encode(p *pem.Block) (templates.Template, error) {
	var n pkix.Name
	if _, err := asn1.Unmarshal(p.Bytes, &n); err != nil {
		return nil, err
	}
	t := &templates.NameTemplate{}
	ne.ApplyPem(&n, t)
	return t, nil
}

func (ne nameEncoder) ApplyPem(name *pkix.Name, t *templates.NameTemplate) {
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

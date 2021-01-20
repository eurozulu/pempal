package templates

import (
	"crypto/x509/pkix"
)

type SubjectTemplate struct {
	CommonName         string   `yaml:"CommonName,omitempty"`
	SerialNumber       string   `yaml:"SerialNumber,omitempty"`
	Organization       []string `yaml:"Organization,omitempty"`
	OrganizationalUnit []string `yaml:"OrganizationalUnit,omitempty"`
	StreetAddress      []string `yaml:"StreetAddress,omitempty"`
	Locality           []string `yaml:"Locality,omitempty"`
	Province           []string `yaml:"Province,omitempty"`
	Country            []string `yaml:"Country,omitempty"`
	PostalCode         []string `yaml:"PostalCode,omitempty"`

	// Names contains all parsed attributes. When parsing distinguished names,
	// this can be used to extract non-standard attributes that are not parsed
	// by this package. When marshaling to RDNSequences, the Names field is
	// ignored, see ExtraNames.
	Names []AttributeTypeAndValue `yaml:"Names,omitempty"`

	// ExtraNames contains attributes to be copied, raw, into any marshaled
	// distinguished names. Values override any attributes with the same OID.
	// The ExtraNames field is not populated when parsing, see Names.
	ExtraNames []AttributeTypeAndValue `yaml:"ExtraNames,omitempty"`
}

func (st SubjectTemplate) String() string {
	return st.CommonName
}

func (st SubjectTemplate) Subject() pkix.Name {
	var n pkix.Name
	st.Apply(&n)
	return n
}

func (st SubjectTemplate) Apply(n *pkix.Name) {
	if len(st.Organization) > 0 {
		n.Organization = append(n.Organization, st.Organization...)
	}
	if len(st.OrganizationalUnit) > 0 {
		n.OrganizationalUnit = append(n.OrganizationalUnit, st.OrganizationalUnit...)
	}
	if len(st.StreetAddress) > 0 {
		n.StreetAddress = append(n.StreetAddress, st.StreetAddress...)
	}
	if len(st.Locality) > 0 {
		n.Locality = append(n.Locality, st.Locality...)
	}
	if len(st.Province) > 0 {
		n.Province = append(n.Province, st.Province...)
	}
	if len(st.Country) > 0 {
		n.Country = append(n.Country, st.Country...)
	}
	if len(st.PostalCode) > 0 {
		n.PostalCode = append(n.PostalCode, st.PostalCode...)
	}
	if len(st.Names) > 0 {
		n.Names = append(n.Names, AttributeTypeAndValueReslice(st.Names)...)
	}
	if len(st.ExtraNames) > 0 {
		n.ExtraNames = append(n.ExtraNames, AttributeTypeAndValueReslice(st.ExtraNames)...)
	}
	if st.SerialNumber != "" {
		n.SerialNumber = st.SerialNumber
	}
}

func NewSubjectTemplate(subject pkix.Name) SubjectTemplate {
	return SubjectTemplate{
		CommonName:         subject.CommonName,
		SerialNumber:       subject.SerialNumber,
		Organization:       subject.Organization,
		OrganizationalUnit: subject.OrganizationalUnit,
		StreetAddress:      subject.StreetAddress,
		Locality:           subject.Locality,
		Province:           subject.Province,
		Country:            subject.Country,
		PostalCode:         subject.PostalCode,
		Names:              AttributeTypeAndValueSlice(subject.Names),
		ExtraNames:         AttributeTypeAndValueSlice(subject.ExtraNames),
	}
}

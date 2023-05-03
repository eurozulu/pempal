package model

import "crypto/x509/pkix"

type DistinguishedNameDTO struct {
	CommonName         string   `yaml:"common-name"`
	Country            []string `yaml:"country,omitempty"`
	Organization       []string `yaml:"organization,omitempty"`
	OrganizationalUnit []string `yaml:"organizational-unit,omitempty"`
	Locality           []string `yaml:"locality,omitempty"`
	Province           []string `yaml:"province,omitempty"`
	StreetAddress      []string `yaml:"street-address,omitempty"`
	PostalCode         []string `yaml:"postal-code,omitempty"`
	SerialNumber       string   `yaml:"serial.txt-number,omitempty"`
}

func (d DistinguishedNameDTO) ToName() pkix.Name {
	return pkix.Name{
		Country:            d.Country,
		Organization:       d.Organization,
		OrganizationalUnit: d.OrganizationalUnit,
		Locality:           d.Locality,
		Province:           d.Province,
		StreetAddress:      d.StreetAddress,
		PostalCode:         d.PostalCode,
		SerialNumber:       d.SerialNumber,
		CommonName:         d.CommonName,
	}
}

func NewDistinguishedNameDTO(n pkix.Name) *DistinguishedNameDTO {
	return &DistinguishedNameDTO{
		Country:            n.Country,
		Organization:       n.Organization,
		OrganizationalUnit: n.OrganizationalUnit,
		Locality:           n.Locality,
		Province:           n.Province,
		StreetAddress:      n.StreetAddress,
		PostalCode:         n.PostalCode,
		SerialNumber:       n.SerialNumber,
		CommonName:         n.CommonName,
	}
}

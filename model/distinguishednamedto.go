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

func (dn DistinguishedNameDTO) ToName() pkix.Name {
	return pkix.Name{
		Country:            dn.Country,
		Organization:       dn.Organization,
		OrganizationalUnit: dn.OrganizationalUnit,
		Locality:           dn.Locality,
		Province:           dn.Province,
		StreetAddress:      dn.StreetAddress,
		PostalCode:         dn.PostalCode,
		SerialNumber:       dn.SerialNumber,
		CommonName:         dn.CommonName,
	}
}

func newDistinguishedNameDTO(n pkix.Name) *DistinguishedNameDTO {
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

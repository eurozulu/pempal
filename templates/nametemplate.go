package templates

import "pempal/resources"

type NameTemplate struct {
	SerialNumber       string   `yaml:"serial-number"`
	CommonName         string   `yaml:"common-name"`
	Organization       []string `yaml:"organization"`
	OrganizationalUnit []string `yaml:"organizational-unit"`
	Country            []string `yaml:"country"`
	Locality           []string `yaml:"locality"`
	Province           []string `yaml:"province"`
	StreetAddress      []string `yaml:"street-address"`
	PostalCode         []string `yaml:"postal-code"`
}

func (nt NameTemplate) IsEmpty() bool {
	return nt.SerialNumber == "" &&
		nt.CommonName == "" &&
		len(nt.Organization) == 0 &&
		len(nt.OrganizationalUnit) == 0 &&
		len(nt.Country) == 0 &&
		len(nt.Locality) == 0 &&
		len(nt.Province) == 0 &&
		len(nt.StreetAddress) == 0 &&
		len(nt.PostalCode) == 0
}

func (nt NameTemplate) Type() resources.ResourceType {
	return resources.Name
}

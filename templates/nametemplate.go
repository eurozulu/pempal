package templates

type NameTemplate struct {
	SerialNumber       string   `yaml:"serial-number,omitempty"`
	CommonName         string   `yaml:"common-name,omitempty"`
	Organization       []string `yaml:"organization,omitempty"`
	OrganizationalUnit []string `yaml:"organizational-unit,omitempty"`
	Country            []string `yaml:"country,omitempty"`
	Locality           []string `yaml:"locality,omitempty"`
	Province           []string `yaml:"province,omitempty"`
	StreetAddress      []string `yaml:"street-address,omitempty"`
	PostalCode         []string `yaml:"postal-code,omitempty"`
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

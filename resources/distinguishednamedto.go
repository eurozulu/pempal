package resources

import (
	"crypto/x509/pkix"
	"fmt"
	"strings"
)

type DistinguishedNameDTO struct {
	CommonName         string   `yaml:"common-name"`
	SerialNumber       string   `yaml:"serial-number,omitempty"`
	Country            []string `yaml:"country,omitempty"`
	Organization       []string `yaml:"organization,omitempty"`
	OrganizationalUnit []string `yaml:"organizational-unit,omitempty"`
	Locality           []string `yaml:"locality,omitempty"`
	Province           []string `yaml:"province,omitempty"`
	StreetAddress      []string `yaml:"street-address,omitempty"`
	PostalCode         []string `yaml:"postal-code,omitempty"`
}

func (dn DistinguishedNameDTO) ToName() pkix.Name {
	var name pkix.Name
	name.CommonName = dn.CommonName
	name.SerialNumber = dn.SerialNumber
	name.Country = dn.Country
	name.Organization = dn.Organization
	name.OrganizationalUnit = dn.OrganizationalUnit
	name.Locality = dn.Locality
	name.Province = dn.Province
	name.StreetAddress = dn.StreetAddress
	name.PostalCode = dn.PostalCode
	return name
}

func (dn *DistinguishedNameDTO) UnmarshalBinary(data []byte) error {
	for _, n := range strings.Split(string(data), ",") {
		ns := strings.SplitN(n, "=", 2)
		if len(ns) != 2 {
			return fmt.Errorf("invalid '%s', no value", n)
		}

		if err := dn.addValue(ns[0], parseValueToSlice(ns[1])...); err != nil {
			return err
		}
	}
	return nil
}

func (dn DistinguishedNameDTO) MarshalBinary() (data []byte, err error) {
	return []byte(dn.ToName().ToRDNSequence().String()), nil
}

func (dn DistinguishedNameDTO) String() string {
	return dn.ToName().String()
}

func (dn *DistinguishedNameDTO) addValue(key string, value ...string) error {
	switch strings.ToUpper(key) {
	case "CN":
		dn.CommonName = value[0]
	case "O":
		dn.Organization = value
	case "OU":
		dn.OrganizationalUnit = value
	case "C":
		dn.Country = value
	case "L":
		dn.Locality = value
	case "ST":
		dn.Province = value
	case "STREET":
		dn.StreetAddress = value
	case "POSTALCODE":
		dn.PostalCode = value
	case "SERIALNUMBER":
		dn.SerialNumber = value[0]
	default:
		return fmt.Errorf("%s is not a known DN key", key)
	}
	return nil
}

func parseValueToSlice(value string) []string {
	var found []string
	for _, v := range strings.Split(value, "+") {
		if strings.Contains(v, "=") {
			v = strings.SplitN(v, "=", 2)[1]
		}
		found = append(found, v)
	}
	return found
}

func ParseDistinguishedName(dn string) (*DistinguishedNameDTO, error) {
	dto := &DistinguishedNameDTO{}
	if err := dto.UnmarshalBinary([]byte(dn)); err != nil {
		return nil, err
	}
	return dto, nil
}

package model

import (
	"crypto/x509/pkix"
	"fmt"
	"github.com/eurozulu/pempal/tools"

	"strings"
)

type DistinguishedName pkix.Name

func (dn DistinguishedName) String() string {
	return pkix.Name(dn).String()
}

func (dn DistinguishedName) IsEmpty() bool {
	return dn.String() == ""
}

func (dn DistinguishedName) Matches(o DistinguishedName) bool {
	if o.CommonName != "" && !strings.Contains(dn.CommonName, o.CommonName) {
		return false
	}
	if o.SerialNumber != "" && !strings.Contains(dn.SerialNumber, o.SerialNumber) {
		return false
	}
	if len(o.Country) > 0 && !tools.HasSameElements(dn.Country, o.Country) {
		return false
	}
	if len(o.Organization) > 0 && !tools.HasSameElements(dn.Organization, o.Organization) {
		return false
	}
	if len(o.OrganizationalUnit) > 0 && !tools.HasSameElements(dn.OrganizationalUnit, o.OrganizationalUnit) {
		return false
	}
	if len(o.Locality) > 0 && !tools.HasSameElements(dn.Locality, o.Locality) {
		return false
	}
	if len(o.Province) > 0 && !tools.HasSameElements(dn.Province, o.Province) {
		return false
	}
	if len(o.StreetAddress) > 0 && !tools.HasSameElements(dn.StreetAddress, o.StreetAddress) {
		return false
	}
	if len(o.PostalCode) > 0 && !tools.HasSameElements(dn.PostalCode, o.PostalCode) {
		return false
	}
	return true
}

func (dn DistinguishedName) Equals(o DistinguishedName) bool {
	if dn.CommonName != o.CommonName {
		return false
	}
	if dn.SerialNumber != o.SerialNumber {
		return false
	}
	if !tools.HasSameElements(dn.Country, o.Country) {
		return false
	}
	if !tools.HasSameElements(dn.Organization, o.Organization) {
		return false
	}
	if !tools.HasSameElements(dn.OrganizationalUnit, o.OrganizationalUnit) {
		return false
	}
	if !tools.HasSameElements(dn.Locality, o.Locality) {
		return false
	}
	if !tools.HasSameElements(dn.Province, o.Province) {
		return false
	}
	if !tools.HasSameElements(dn.StreetAddress, o.StreetAddress) {
		return false
	}
	if !tools.HasSameElements(dn.PostalCode, o.PostalCode) {
		return false
	}
	//if len(dn.Names) != len(o.Names) {
	//	return false
	//}
	if len(dn.ExtraNames) != len(o.ExtraNames) {
		return false
	}
	return true
}

func (dn *DistinguishedName) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}
	d, err := ParseDistinguishedName(string(text))
	if err != nil {
		return err
	}
	dn.Merge(*d)
	return nil
}

func (dn DistinguishedName) MarshalText() (text []byte, err error) {
	return []byte(dn.String()), nil
}

func (dn DistinguishedName) ToName() pkix.Name {
	return pkix.Name(dn)
}

// Merge replaces any field with the given non empty field.
// If the given Names field is not empty it will be added to this name.
// For single value fields, the value is replaced. For slices, the slices are merged with the unique values of existing and given.
func (dn *DistinguishedName) Merge(o DistinguishedName) {
	if o.SerialNumber != "" {
		dn.SerialNumber = o.SerialNumber
	}
	if o.CommonName != "" {
		dn.CommonName = o.CommonName
	}
	if len(o.Country) > 0 {
		dn.Country = tools.AppendUnique(dn.Country, o.Country...)
	}
	if len(o.Organization) > 0 {
		dn.Organization = tools.AppendUnique(dn.Organization, o.Organization...)
	}
	if len(o.OrganizationalUnit) > 0 {
		dn.OrganizationalUnit = tools.AppendUnique(dn.OrganizationalUnit, o.OrganizationalUnit...)
	}
	if len(o.Locality) > 0 {
		dn.Locality = tools.AppendUnique(dn.Locality, o.Locality...)
	}
	if len(o.Province) > 0 {
		dn.Province = tools.AppendUnique(dn.Province, o.Province...)
	}
	if len(o.StreetAddress) > 0 {
		dn.StreetAddress = tools.AppendUnique(dn.StreetAddress, o.StreetAddress...)
	}
	if len(o.PostalCode) > 0 {
		dn.PostalCode = tools.AppendUnique(dn.PostalCode, o.PostalCode...)
	}
}

func (dn *DistinguishedName) addValue(key string, value ...string) error {
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
	// DC not supported

	default:
		return fmt.Errorf("%s is not a known DN key", key)
	}
	return nil
}

func parseRDNSValueToSlice(value string) []string {
	var found []string
	for _, v := range strings.Split(value, "+") {
		v = strings.Trim(v, "\"")
		if strings.Contains(v, "=") {
			v = strings.SplitN(v, "=", 2)[1]
		}
		found = append(found, v)
	}
	return found
}

// ParseName parses the given name as a DN.  If the given string is not a valid DN,
// The string is used as the Common Name of a DN.
func ParseName(name string) (*DistinguishedName, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if !strings.Contains(name, "=") {
		name = "CN=" + name
	}
	return ParseDistinguishedName(name)
}

// ParseDistinguishedName parses the given atring as a valid DN.
func ParseDistinguishedName(s string) (*DistinguishedName, error) {
	dn := &DistinguishedName{}
	for _, n := range strings.Split(s, ",") {
		ns := strings.SplitN(n, "=", 2)
		if len(ns) != 2 {
			return nil, fmt.Errorf("invalid DN '%s' has no '='", n)
		}

		if err := dn.addValue(ns[0], parseRDNSValueToSlice(ns[1])...); err != nil {
			return nil, err
		}
	}
	return dn, nil
}

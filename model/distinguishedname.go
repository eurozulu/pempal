package model

import (
	"crypto/x509/pkix"
	"fmt"
	"strings"
)

type DistinguishedName pkix.Name

func (dn DistinguishedName) String() string {
	return pkix.Name(dn).String()
}

func (dn DistinguishedName) Equals(o DistinguishedName) bool {
	if dn.CommonName != o.CommonName {
		return false
	}
	if dn.SerialNumber != o.SerialNumber {
		return false
	}
	if !equalsSlice(dn.Country, o.Country) {
		return false
	}
	if !equalsSlice(dn.Organization, o.Organization) {
		return false
	}

	if !equalsSlice(dn.OrganizationalUnit, o.OrganizationalUnit) {
		return false
	}
	if !equalsSlice(dn.Locality, o.Locality) {
		return false
	}
	if !equalsSlice(dn.Province, o.Province) {
		return false
	}
	if !equalsSlice(dn.StreetAddress, o.StreetAddress) {
		return false
	}
	if !equalsSlice(dn.PostalCode, o.PostalCode) {
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
	dn.merge(*d)
	return nil
}

func (dn DistinguishedName) MarshalText() (text []byte, err error) {
	return []byte(dn.String()), nil
}

func (dn DistinguishedName) ToName() pkix.Name {
	return pkix.Name(dn)
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

func (dn *DistinguishedName) merge(o DistinguishedName) {
	if o.SerialNumber != "" {
		dn.SerialNumber = o.SerialNumber
	}
	if o.CommonName != "" {
		dn.CommonName = o.CommonName
	}
	if len(o.Country) > 0 {
		dn.Country = mergeSlice(dn.Country, o.Country)
	}
	if len(o.Organization) > 0 {
		dn.Organization = mergeSlice(dn.Organization, o.Organization)
	}
	if len(o.OrganizationalUnit) > 0 {
		dn.OrganizationalUnit = mergeSlice(dn.OrganizationalUnit, o.OrganizationalUnit)
	}
	if len(o.Locality) > 0 {
		dn.Locality = mergeSlice(dn.Locality, o.Locality)
	}
	if len(o.Province) > 0 {
		dn.Province = mergeSlice(dn.Province, o.Province)
	}
	if len(o.StreetAddress) > 0 {
		dn.StreetAddress = mergeSlice(dn.StreetAddress, o.StreetAddress)
	}
	if len(o.PostalCode) > 0 {
		dn.PostalCode = mergeSlice(dn.PostalCode, o.PostalCode)
	}
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

func equalsSlice(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for _, s := range s1 {
		if !containsString(s, s2) {
			return false
		}
	}
	return true
}

func containsString(s string, ss []string) bool {
	for _, sz := range ss {
		if sz == s {
			return true
		}
	}
	return false
}

func mergeSlice(s1, s2 []string) []string {
	ss := make([]string, len(s1))
	copy(ss, s1)
	for _, sz := range s2 {
		if containsString(sz, ss) {
			continue
		}
		ss = append(ss, sz)
	}
	return ss
}

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

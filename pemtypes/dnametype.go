package pemtypes

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"gopkg.in/yaml.v3"
	"pempal/templates"
)

type dnameType struct {
	dname pkix.Name
}

func (nt dnameType) String() string {
	return fmt.Sprintf("%s\t%s", Name.String(), nt.dname.String())
}

func (nt dnameType) MarshalBinary() (data []byte, err error) {
	return asn1.Marshal(&nt.dname)
}

func (nt *dnameType) UnmarshalBinary(data []byte) error {
	_, err := asn1.Unmarshal(data, &nt.dname)
	return err
}

func (nt dnameType) MarshalText() (text []byte, err error) {
	data, err := nt.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  Name.String(),
		Bytes: data,
	}), nil
}

func (nt *dnameType) UnmarshalText(text []byte) error {
	blocks := ReadPEMBlocks(text, Name)
	if len(blocks) == 0 {
		return fmt.Errorf("no name pem found")
	}
	return nt.UnmarshalBinary(blocks[0].Bytes)
}

func (nt dnameType) MarshalYAML() (interface{}, error) {
	t := templates.NameTemplate{}
	nt.applyToTemplate(&t)
	return yaml.Marshal(&t)
}

func (nt *dnameType) UnmarshalYAML(value *yaml.Node) error {
	t := templates.NameTemplate{}
	if err := value.Decode(&t); err != nil {
		return err
	}
	nt.applyTemplate(t)
	return nil
}

func (nt dnameType) applyToTemplate(t *templates.NameTemplate) {
	t.CommonName = nt.dname.CommonName
	t.SerialNumber = nt.dname.SerialNumber
	t.Organization = nt.dname.Organization
	t.OrganizationalUnit = nt.dname.OrganizationalUnit
	t.Country = nt.dname.Country
	t.Locality = nt.dname.Locality
	t.Province = nt.dname.Province
	t.StreetAddress = nt.dname.StreetAddress
	t.PostalCode = nt.dname.PostalCode
}

func (nt *dnameType) applyTemplate(t templates.NameTemplate) {
	if t.CommonName != "" {
		nt.dname.CommonName = t.CommonName
	}
	if t.SerialNumber != "" {
		nt.dname.SerialNumber = t.SerialNumber
	}
	if len(t.Organization) > 0 {
		nt.dname.Organization = t.Organization
	}
	if len(t.OrganizationalUnit) > 0 {
		nt.dname.OrganizationalUnit = t.OrganizationalUnit
	}
	if len(t.Country) > 0 {
		nt.dname.Country = t.Country
	}
	if len(t.Locality) > 0 {
		nt.dname.Locality = t.Locality
	}
	if len(t.Province) > 0 {
		nt.dname.Province = t.Province
	}
	if len(t.StreetAddress) > 0 {
		nt.dname.StreetAddress = t.StreetAddress
	}
	if len(t.PostalCode) > 0 {
		nt.dname.PostalCode = t.PostalCode
	}
}

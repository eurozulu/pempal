package pemtypes

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"gopkg.in/yaml.v3"
	"pempal/pemproperties"
	"pempal/templates"
	"strings"
)

type requestType struct {
	request x509.CertificateRequest
}

func (rt requestType) String() string {
	return fmt.Sprintf("%s\t%s", Request.String(), rt.request.Subject.String())
}

func (rt requestType) MarshalBinary() (data []byte, err error) {
	return rt.request.Raw, err
}

func (rt requestType) UnmarshalBinary(data []byte) error {
	csr, err := x509.ParseCertificateRequest(data)
	if err != nil {
		return err
	}
	rt.request = *csr
	return nil
}

func (rt requestType) MarshalText() (text []byte, err error) {
	b := &pem.Block{
		Type:  Request.String(),
		Bytes: rt.request.Raw,
	}
	return pem.EncodeToMemory(b), nil
}

func (rt requestType) UnmarshalText(text []byte) error {
	blocks := ReadPEMBlocks(text, Request)
	if len(blocks) == 0 {
		return fmt.Errorf("no request pem found")
	}
	return rt.UnmarshalBinary(blocks[0].Bytes)
}

func (rt requestType) MarshalYAML() (interface{}, error) {
	t := templates.CSRTemplate{}
	rt.applyToTemplate(&t)
	return yaml.Marshal(&t)
}

func (rt requestType) UnmarshalYAML(value *yaml.Node) error {
	t := templates.CSRTemplate{}
	if err := value.Decode(&t); err != nil {
		return err
	}
	rt.applyTemplate(t)
	return nil
}

func (rt *requestType) applyTemplate(t templates.CSRTemplate) {
	csr := rt.request
	if t.Version != 0 {
		csr.Version = t.Version
	}
	if t.SignatureAlgorithm != "" {
		csr.SignatureAlgorithm = pemproperties.SignatureAlgorithmProperty{}.Parse(t.SignatureAlgorithm)
	}
	if t.PublicKeyAlgorithm != "" {
		csr.PublicKeyAlgorithm = pemproperties.PublicKeyAlgorithmProperty{}.Parse(t.PublicKeyAlgorithm)
	}
	if t.Subject != nil {
		nt := &dnameType{}
		nt.applyTemplate(*t.Subject)
		csr.Subject = nt.dname
	}

	//TODO:
	//Attributes         []pkix.AttributeTypeAndValueSET
	//Extensions         []pkix.Extension
	//ExtraExtensions    []pkix.Extension
	if len(t.DNSNames) > 0 {
		csr.DNSNames = t.DNSNames
	}
	if len(t.EmailAddresses) > 0 {
		csr.EmailAddresses = t.EmailAddresses
	}
	if len(t.IPAddresses) > 0 {
		csr.IPAddresses = pemproperties.IPAddressListProperty{}.Parse(strings.Join(t.IPAddresses, ","))
	}
	if len(t.URIs) > 0 {
		csr.URIs = pemproperties.URIListProperty{}.Parse(strings.Join(t.URIs, ","))
	}
}

func (rt requestType) applyToTemplate(t *templates.CSRTemplate) {
	csr := rt.request
	t.Version = csr.Version
	t.SignatureAlgorithm = csr.SignatureAlgorithm.String()
	t.PublicKeyAlgorithm = csr.PublicKeyAlgorithm.String()

	t.Subject = &templates.NameTemplate{}
	nt := &dnameType{dname: csr.Subject}
	nt.applyToTemplate(t.Subject)

	if len(csr.DNSNames) > 0 {
		t.DNSNames = csr.DNSNames
	}
	if len(csr.EmailAddresses) > 0 {
		t.EmailAddresses = csr.EmailAddresses
	}
	if len(csr.IPAddresses) > 0 {
		t.IPAddresses = strings.Split(pemproperties.IPAddressListProperty{}.String(csr.IPAddresses), ",")
	}
	if len(csr.URIs) > 0 {
		t.URIs = strings.Split(pemproperties.URIListProperty{}.String(csr.URIs), ",")
	}
}

package model

import (
	"fmt"
	"strings"
)

type ResourceType int

const (
	UnknownResourceType ResourceType = iota
	PublicKey
	PrivateKey
	CertificateRequest
	Certificate
	RevokationList
)

var resourceTypeNames = []string{
	"UnknownResourceType",
	"PublicKey",
	"PrivateKey",
	"CertificateRequest",
	"Certificate",
	"RevokationList",
}
var resourceTypePEMNames = []string{
	"",
	"PUBLIC KEY",
	"PRIVATE KEY",
	"CERTIFICATE REQUEST",
	"CERTIFICATE",
	"X509 CRL",
}

var aliases = map[string]string{
	"puk":  "PublicKey",
	"prk":  "PrivateKey",
	"cert": "Certificate",
	"csr":  "CertificateRequest",
	"crl":  "RevokationList",
}

func (rt ResourceType) String() string {
	i := int(rt)
	if i < 0 || i >= len(resourceTypeNames) {
		i = 0
	}
	return resourceTypeNames[i]
}

func (rt ResourceType) PEMString() string {
	i := int(rt)
	if i < 0 || i >= len(resourceTypePEMNames) {
		i = 0
	}
	return resourceTypePEMNames[i]
}

func (rt ResourceType) MarshalText() (text []byte, err error) {
	return []byte(rt.String()), nil
}

func (rt *ResourceType) UnmarshalText(text []byte) error {
	s := string(text)
	if real, ok := aliases[strings.ToLower(s)]; ok {
		s = real
	}
	for i, n := range resourceTypeNames {
		if strings.EqualFold(n, s) {
			*rt = ResourceType(i)
			return nil
		}
	}
	return fmt.Errorf("unknown resource type %q", s)
}

func ContainsResourceType(types []ResourceType, typ ResourceType) bool {
	for _, t := range types {
		if t == typ {
			return true
		}
	}
	return false
}

func ParseResourceType(s string) (ResourceType, error) {
	rt := UnknownResourceType
	if err := rt.UnmarshalText([]byte(s)); err != nil {
		return UnknownResourceType, err
	}
	return rt, nil
}

func ParseResourceTypeFromPEMType(pemType string) ResourceType {
	n := strings.ToUpper(pemType)
	for i, pn := range resourceTypePEMNames {
		if strings.Contains(pn, " KEY") {
			// keys do a contains search
			if strings.Contains(n, pn) {
				return ResourceType(i)
			}
			continue
		}
		if strings.EqualFold(n, pn) {
			return ResourceType(i)
		}
	}
	return UnknownResourceType
}

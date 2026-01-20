package model

import (
	"fmt"
	"strings"
)

type ResourceType int

const (
	ResourceTypeUnknown ResourceType = iota
	ResourceTypePublicKey
	ResourceTypePrivateKey
	ResourceTypeCertificateRequest
	ResourceTypeCertificate
	ResourceTypeRevokationList
)

var resourceTypeNames = []string{
	"UNKNOWN",
	"PUBLIC KEY",
	"PRIVATE KEY",
	"CERTIFICATE REQUEST",
	"CERTIFICATE",
	"X509 CRL",
}

var aliases = map[string]string{
	"puk":     "PUBLIC KEY",
	"prk":     "PRIVATE KEY",
	"key":     "PRIVATE KEY",
	"cert":    "CERTIFICATE",
	"cer":     "CERTIFICATE",
	"csr":     "CERTIFICATE REQUEST",
	"request": "CERTIFICATE REQUEST",
	"req":     "CERTIFICATE REQUEST",
	"crl":     "X509 CRL",
}

func (rt ResourceType) String() string {
	i := int(rt)
	if i < 0 || i >= len(resourceTypeNames) {
		i = 0
	}
	return resourceTypeNames[i]
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

func ParseResourceType(s string) ResourceType {
	rt := ResourceTypeUnknown
	if err := rt.UnmarshalText([]byte(s)); err != nil {
		return ResourceTypeUnknown
	}
	return rt
}

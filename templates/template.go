package templates

import "pempal/pemtypes"

type Template interface {
}

type BlankTemplate map[string]interface{}

func NewTemplate(pt pemtypes.PEMType) Template {
	switch pt {
	case pemtypes.Certificate:
		return &CertificateTemplate{}
	case pemtypes.Request:
		return &CSRTemplate{}
	case pemtypes.PrivateKey:
		return &KeyTemplate{}
	case pemtypes.PrivateKeyEncrypted:
		return &KeyTemplate{}
	case pemtypes.RevocationList:
		return &CRLTemplate{}
	case pemtypes.Name:
		return &NameTemplate{}
	default:
		return BlankTemplate{}
	}
}

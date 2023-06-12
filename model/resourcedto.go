package model

import (
	"fmt"
)

type ResourceDTO interface {
	UnmarshalPEM(data []byte) error
}

// DTOForResource attempts to marshal the given PEMResource into a DTO for its type.
func DTOForResource(r Resource) (ResourceDTO, error) {
	dto := NewDTOForResourceType(r.ResourceType())
	if dto == nil {
		return nil, fmt.Errorf("'%s' is an unsupported resource type", r.ResourceType().String())
	}
	if err := dto.UnmarshalPEM([]byte(r.String())); err != nil {
		return nil, err
	}
	return dto, nil
}

// NewDTOForResourceType gets an empty ResourceDTO specific to the given type
func NewDTOForResourceType(t ResourceType) ResourceDTO {
	switch t {
	case Certificate:
		return &CertificateDTO{}
	case CertificateRequest:
		return &CertificateRequestDTO{}
	case PublicKey:
		return &PublicKeyDTO{}
	case PrivateKey:
		return &PrivateKeyDTO{}
	case RevokationList:
		return &RevocationListDTO{}
	default:
		return nil
	}
}

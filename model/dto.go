package model

import (
	"encoding"
)

const ResourceTypeName = "resource-type"

type DTO interface {
	encoding.BinaryUnmarshaler
}

// DTOForResource attempts to marshal the given PEMResource into a DTO for its type.
func DTOForResource(r PEMResource) (DTO, error) {
	der, err := r.MarshalBinary()
	if err != nil {
		return nil, err
	}
	dto := DTOForResourceType(r.ResourceType())
	if err = dto.UnmarshalBinary(der); err != nil {
		return nil, err
	}
	return dto, nil
}

// DTOForResourceType gets an empty DTO specific to the given type
func DTOForResourceType(t ResourceType) DTO {
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

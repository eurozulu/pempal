package resources

import (
	"bytes"
	"encoding"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/templates"
	"github.com/go-yaml/yaml"
)

// ResourceDTO is the intermediary object between Resources and Templates
// ResourceDTOs are specific to a Resource Type
// Every DTO must also marshal/unmarshal yaml
type ResourceDTO interface {
	encoding.BinaryUnmarshaler
	encoding.BinaryMarshaler
	fmt.Stringer
}

func DTOToResource(dto ResourceDTO) (Resource, error) {
	blk, _ := pem.Decode([]byte(dto.String()))
	if blk == nil {
		return nil, fmt.Errorf("insufficient data to convert resourcedto to resource")
	}
	return NewResource(blk), nil
}

func ApplyTemplateToDTO(dto ResourceDTO, t templates.Template) error {
	if err := transferViaYaml(dto, t); err != nil {
		return err
	}
	return nil
}

func DTOToTemplate(dto ResourceDTO) (templates.Template, error) {
	m := map[string]string{}
	if err := transferViaYaml(&m, dto); err != nil {
		return nil, err
	}
	return m, nil
}

func transferViaYaml(target, src interface{}) error {
	buf := bytes.NewBuffer(nil)
	if err := yaml.NewEncoder(buf).Encode(src); err != nil {
		return err
	}
	if err := yaml.NewDecoder(buf).Decode(target); err != nil {
		return err
	}
	return nil
}

func NewResourceDTO(r Resource) (ResourceDTO, error) {
	dto, err := NewResourceDTOByType(r.ResourceType())
	if err != nil {
		return nil, err
	}
	if err = dto.UnmarshalBinary(r.Bytes()); err != nil {
		return nil, err
	}
	return dto, nil
}

func NewResourceDTOByType(resourceType ResourceType) (ResourceDTO, error) {
	switch resourceType {
	case PublicKey:
		return &publicKeyDTO{}, nil

	case PrivateKey:
		return &PrivateKeyDTO{}, nil

	case Certificate:
		return &CertificateDTO{}, nil

	default:
		return nil, fmt.Errorf("no resourcedto available for resource type %s", resourceType.String())
	}
}

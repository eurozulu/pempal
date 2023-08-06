package resources

import (
	"bytes"
	"encoding"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/templates"
	"github.com/go-yaml/yaml"
	"log"
	"strings"
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

func DTOToTemplate(dto ResourceDTO, includeAll bool) (templates.Template, error) {
	m := map[string]interface{}{}

	if !includeAll {
		if err := transferViaYaml(&m, dto); err != nil {
			return nil, err
		}
	} else {
		if err := transferViaJson(&m, dto); err != nil {
			return nil, err
		}
	}
	sm := map[string]string{}
	for k, v := range m {
		sm[k] = fmt.Sprintf("%v", v)
	}
	return sm, nil
}

func TemplateTypes(t templates.Template) ([]ResourceType, error) {
	var found []ResourceType
	for _, rt := range resourceTypes[1:] {
		// Ensure template contains all of the keys in the resource template
		if !containsRequiredKeys(rt, t) {
			continue
		}
		if containsAllKeys(rt, t) {
			return []ResourceType{rt}, nil
		}
		found = append(found, rt)
	}
	if len(found) == 0 {
		return nil, fmt.Errorf("template type could not be determined")
	}
	return found, nil
}

func ResourceTemplateByType(rt ResourceType, includeAll bool) (templates.Template, error) {
	dto, err := NewResourceDTOByType(rt)
	if err != nil {
		return nil, err
	}
	return DTOToTemplate(dto, includeAll)
}

func containsRequiredKeys(rt ResourceType, t map[string]string) bool {
	rtt, err := ResourceTemplateByType(rt, false)
	if err != nil {
		log.Fatalf("Unexpected resource type error ", err)
	}
	for k := range rtt {
		if _, ok := t[k]; !ok {
			return false
		}
	}
	return true
}

func containsAllKeys(rt ResourceType, t map[string]string) bool {
	rtt, err := ResourceTemplateByType(rt, true)
	if err != nil {
		log.Fatalf("Unexpected resource type error ", err)
	}
	for k := range t {
		if strings.EqualFold(k, rt.String()) {
			continue
		} // ignore the type property
		if _, ok := rtt[k]; !ok {
			return false
		}
	}
	return true
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
func transferViaJson(target, src interface{}) error {
	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(src); err != nil {
		return err
	}
	if err := json.NewDecoder(buf).Decode(target); err != nil {
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
		return &PublicKeyDTO{}, nil
	case PrivateKey:
		return &PrivateKeyDTO{}, nil
	case Certificate:
		return &CertificateDTO{}, nil
	case CertificateRequest:
		return &CertificateRequestDTO{}, nil
	case RevocationList:
		return &RevocationListDTO{}, nil

	default:
		return nil, fmt.Errorf("no resourcedto available for resource type %s", resourceType.String())
	}
}

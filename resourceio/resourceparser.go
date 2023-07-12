package resourceio

import (
	"bytes"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/resources"
	"os"
)

const pemBegin = "-----BEGIN "
const pemEnd = "-----END "

func ParseLocation(path string) (ResourceLocation, error) {
	logger.Debug("reading location %s", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	res, err := ParseResources(data)
	if err != nil {
		return nil, err
	}
	logger.Debug("found %d resources in location %s", len(res), path)
	return NewResourceLocation(path, res), nil
}

func ParseResources(data []byte) ([]resources.Resource, error) {
	if containPem(data) {
		return parseResourcesAsPEM(data)
	}
	// try as der for each type
	for _, rt := range resources.ResourceTypes {
		r, err := parseResourcesAsDER(rt, data)
		if err != nil {
			continue
		}
		return []resources.Resource{r}, nil
	}
	return nil, fmt.Errorf("failed to parse as a pem or known der type")
}

func parseResourcesAsDER(rt resources.ResourceType, data []byte) (resources.Resource, error) {
	dto, err := resources.NewResourceDTOByType(rt)
	if err != nil {
		return nil, fmt.Errorf("unexpected resourcedto creation error %v", err)
	}
	if err = dto.UnmarshalBinary(data); err != nil {
		return nil, err
	}
	return resources.DTOToResource(dto)
}

func parseResourcesAsPEM(data []byte) ([]resources.Resource, error) {
	var found []resources.Resource
	for len(data) > 0 {
		blk, rest := pem.Decode(data)
		if blk == nil {
			break
		}
		r, err := parseResourcesAsDER(resources.ParsePEMType(blk.Type), blk.Bytes)
		if err != nil {
			return nil, err
		}
		found = append(found, r)
		data = rest
	}
	if len(found) == 0 {
		return nil, fmt.Errorf("no pem encoded resources found")
	}
	return found, nil
}

func containPem(data []byte) bool {
	pe := []byte(pemBegin)
	i := bytes.Index(data, pe)
	return i >= 0 && bytes.Index(data[i+len(pe):], []byte(pemEnd)) > 0
}

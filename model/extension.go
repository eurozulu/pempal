package model

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
	"strconv"
	"strings"
)

type Extension pkix.Extension

func (e Extension) MarshalYAML() (interface{}, error) {
	dto := &extensionDTO{
		ID:       e.Id.String(),
		Value:    EncodeAsBase64(e.Value),
		Critical: e.Critical,
	}
	return dto, nil
}

func (e *Extension) UnmarshalYAML(unmarshal func(interface{}) error) error {
	dto := &extensionDTO{}
	if err := unmarshal(dto); err != nil {
		return err
	}
	id, err := dto.ParsedID()
	if err != nil {
		return err
	}
	val, err := DecodeAsBase64(dto.Value)
	if err != nil {
		return err
	}

	e.Id = id
	e.Value = val
	e.Critical = dto.Critical
	return nil
}

type extensionDTO struct {
	ID       string `yaml:"id"`
	Value    string `yaml:",omitempty"`
	Critical bool   `yaml:"critical,omitempty"`
}

func (dto extensionDTO) ParsedID() (asn1.ObjectIdentifier, error) {
	ids := strings.Split(dto.ID, ".")
	oid := make(asn1.ObjectIdentifier, len(ids))
	for i, s := range ids {
		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid id %s. %v", dto.ID, err)
		}
		oid[i] = v
	}
	return oid, nil
}

func ExtensionsToModel(exts []pkix.Extension) []Extension {
	result := make([]Extension, len(exts))
	for i, ext := range exts {
		result[i] = Extension(ext)
	}
	return result
}
func ModelToExtensions(exts []Extension) []pkix.Extension {
	result := make([]pkix.Extension, len(exts))
	for i, ext := range exts {
		result[i] = pkix.Extension(ext)
	}
	return result
}

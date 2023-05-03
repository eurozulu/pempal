package resourceio

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io"
	"pempal/logger"
	"pempal/model"
	"strings"
)

type ResourceFormat int

const (
	PEM ResourceFormat = iota
	DER
	YAML
	LIST
)

type ResourcePrinter interface {
	Write(location ResourceLocation) error
}

type resourcePemPrinter struct {
	out io.Writer
}

func (prn resourcePemPrinter) Write(location ResourceLocation) error {
	for _, r := range location.Resources {
		by, err := r.MarshalPEM()
		if err != nil {
			return err
		}
		if _, err = prn.out.Write(by); err != nil {
			return err
		}
	}
	return nil
}

type resourceYamlPrinter struct {
	out io.Writer
}

func (prn resourceYamlPrinter) Write(location ResourceLocation) error {
	for i, r := range location.Resources {
		if i > 0 {

		}
		der, err := r.MarshalBinary()
		if err != nil {
			return err
		}
		dto := model.DTOForResourceType(r.ResourceType())
		if err = dto.UnmarshalBinary(der); err != nil {
			return err
		}

		if err = yaml.NewEncoder(prn.out).Encode(dto); err != nil {
			return err
		}
	}
	return nil
}

type resourceDerPrinter struct {
	out io.Writer
}

func (prn resourceDerPrinter) Write(location ResourceLocation) error {
	if len(location.Resources) == 0 {
		return nil
	}
	if len(location.Resources) > 1 {
		return fmt.Errorf("%s contains multiple resources, only a single resource can be printed out in DER")
	}
	der, err := location.Resources[0].MarshalBinary()
	if err != nil {
		return err
	}
	if _, err = prn.out.Write(der); err != nil {
		return err
	}
	return nil
}

var ResourceFormatNames = []string{"pem", "der", "yaml", "list"}

func ParseResourceFormat(s string) (ResourceFormat, error) {
	for i, rn := range ResourceFormatNames {
		if strings.EqualFold(rn, s) {
			return ResourceFormat(i), nil
		}
	}
	return -1, fmt.Errorf("%s is not a known format", s)
}

func NewResourcePrinter(out io.Writer, format ResourceFormat) ResourcePrinter {
	switch format {
	case PEM:
		return &resourcePemPrinter{out: out}
	case DER:
		return &resourceDerPrinter{out: out}
	case YAML:
		return &resourceYamlPrinter{out: out}
	case LIST:
		return NewResourceListPrinter(out)
	default:
		logger.Log(logger.Error, "Invalid ResourceFormat %v", format)
		return nil
	}
}

package builders

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/identity"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/templates"
	"github.com/eurozulu/pempal/utils"
	"strconv"
)

const (
	property_cert_serial_number = "serial-number"
	property_cert_subject_name  = "subject"
	property_cert_issuer_name   = "issuer"
	property_cert_common_name   = "common-name"
	property_cert_public_key    = "public-key"
)

type certificateBuilder struct {
	knownIssuers identity.Issuers
	temps        []templates.Template
}

func (c *certificateBuilder) AddTemplate(t ...templates.Template) {
	c.temps = append(c.temps, t...)
}

func (c certificateBuilder) Validate() utils.CompoundErrors {
	return c.validate(c.BuildTemplate())
}

func (c certificateBuilder) BuildTemplate() templates.Template {
	return templates.MergeTemplates(c.temps...)
}

func (c certificateBuilder) Build() (resources.Resource, error) {
	t := c.BuildTemplate()
	if errs := c.validate(t); len(errs) > 0 {
		return nil, errs
	}

	cert, err := c.buildTemplateCertificate(t)
	if err != nil {
		return nil, err
	}

	issuer, err := c.knownIssuers.IssuerByName(cert.Issuer.String())
	if err != nil {
		return nil, err
	}

	der, err := x509.CreateCertificate(rand.Reader, cert, issuer.Certificate().Certificate(), cert.PublicKey, issuer.Key().PrivateKey())
	if err != nil {
		return nil, err
	}
	return resources.NewResource(&pem.Block{
		Type:  resources.Certificate.PEMString(),
		Bytes: der,
	}), nil
}

func (c certificateBuilder) validate(t templates.Template) utils.CompoundErrors {
	var errs utils.CompoundErrors

	if vp := t[property_cert_serial_number]; vp == "" {
		errs = append(errs, fmt.Errorf("%s missing", property_cert_serial_number))
	} else if _, err := strconv.ParseInt(vp, 10, 64); err != nil {
		errs = append(errs, fmt.Errorf("%s invalid", property_cert_serial_number))
	}

	if vp := t[property_cert_subject_name]; vp == "" {
		errs = append(errs, fmt.Errorf("%s missing", property_cert_subject_name))
	} else if dto, err := resources.ParseDistinguishedName(vp); err != nil {
		errs = append(errs, fmt.Errorf("%s invalid %v", property_cert_subject_name, err))
	} else {
		if dto.CommonName == "" {
			errs = append(errs, fmt.Errorf("%s.%s missing", property_cert_subject_name, property_cert_common_name))
		}
	}
	if vp := t[property_cert_issuer_name]; vp == "" {
		errs = append(errs, fmt.Errorf("%s missing", property_cert_issuer_name))
	} else if dto, err := resources.ParseDistinguishedName(vp); err != nil {
		errs = append(errs, fmt.Errorf("%s invalid %v", property_cert_issuer_name, err))
	} else {
		if dto.CommonName == "" {
			errs = append(errs, fmt.Errorf("%s.%s missing", property_cert_issuer_name, property_cert_common_name))
		}
	}

	if vp := t[property_cert_public_key]; vp == "" {
		errs = append(errs, fmt.Errorf("%s missing", property_cert_public_key))
	} else if _, err := c.resolveKeyId(vp); err != nil {
		errs = append(errs, fmt.Errorf("%s %v", property_cert_public_key, err))
	}

	return errs
}

func (c certificateBuilder) buildTemplateCertificate(t templates.Template) (*x509.Certificate, error) {
	dto := &resources.CertificateDTO{}
	err := resources.ApplyTemplateToDTO(dto, t)
	if err != nil {
		return nil, err
	}
	return dto.ToCertificate(), nil
}

func (c certificateBuilder) resolveKeyId(s string) (identity.Identity, error) {
	if identity.IsIdentity(s) {
		return identity.Identity(s), nil
	}
	k, err := c.knownIssuers.Keys().KeyByIdentity(s)
	if err != nil {
		return "", err
	}
	return k.Identity(), err
}

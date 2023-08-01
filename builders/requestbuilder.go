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
)

const (
	Property_signature           = "signature"
	Property_signature_algorithm = "signature-algorithm"
	Property_public_key          = "public-key"
	Property_subject             = "subject"
	Property_common_name         = "common-name"
)

type requestBuilder struct {
	knownKeys identity.Keys
}

func (c requestBuilder) Validate(t templates.Template) utils.CompoundErrors {
	var errs utils.CompoundErrors
	if sub := t[Property_subject]; sub == "" {
		errs = append(errs, fmt.Errorf("%s missing", Property_subject))
	} else if dto, err := resources.ParseDistinguishedName(sub); err != nil {
		errs = append(errs, fmt.Errorf("%s invalid %v", Property_subject, err))
	} else {
		if dto.CommonName == "" {
			errs = append(errs, fmt.Errorf("%s.%s missing", Property_subject, Property_common_name))
		}
	}

	if ks := t[Property_public_key]; ks == "" {
		errs = append(errs, fmt.Errorf("%s missing", Property_public_key))
	} else if !identity.IsIdentity(ks) {
		errs = append(errs, fmt.Errorf("%s is invaliud", Property_public_key))
	}

	if sa := t[Property_signature_algorithm]; sa == "" {
		errs = append(errs, fmt.Errorf("%s missing", Property_signature_algorithm))
	} else if utils.ParseSignatureAlgorithm(sa) == x509.UnknownSignatureAlgorithm {
		errs = append(errs, fmt.Errorf("%s is invalid", Property_signature_algorithm))
	}
	return errs
}

func (c requestBuilder) Build(t templates.Template) (resources.Resource, error) {
	if errs := c.Validate(t); len(errs) > 0 {
		return nil, errs
	}
	req, err := c.buildTemplateRequest(t)
	if err != nil {
		return nil, err
	}
	id := identity.Identity(t[Property_public_key])
	signer, err := c.knownKeys.KeyByIdentity(id.String())
	if err != nil {
		return nil, err
	}

	der, err := x509.CreateCertificateRequest(rand.Reader, req, signer.PrivateKey())
	if err != nil {
		return nil, err
	}
	return resources.NewResource(&pem.Block{
		Type:  resources.Certificate.PEMString(),
		Bytes: der,
	}), nil
}

func (c requestBuilder) buildTemplateRequest(t templates.Template) (*x509.CertificateRequest, error) {
	dto := &resources.CertificateRequestDTO{}
	err := resources.ApplyTemplateToDTO(dto, t)
	if err != nil {
		return nil, err
	}
	return dto.ToCertificateRequest()
}

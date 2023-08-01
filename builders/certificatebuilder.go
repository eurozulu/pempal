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
	Property_version                 = "version"
	Property_serial_number           = "serial-number"
	Property_issuer                  = "issuer"
	Property_not_before              = "not-before"
	Property_not_after               = "not-after"
	Property_is_ca                   = "is-ca"
	Property_basic_constraints_valid = "basic-constraints-valid"
	Property_max_path_len            = "max-path-len"
	Property_max_path_len_zero       = "max-path-len-zero"
	Property_key_usage               = "max-path-len-zero"
	Property_extended_key_usage      = "extended-key-usage"
)

type certificateBuilder struct {
	knownIssuers identity.Users
}

func (c certificateBuilder) Validate(t templates.Template) utils.CompoundErrors {
	var errs utils.CompoundErrors
	if sn := t[Property_serial_number]; sn == "" {
		errs = append(errs, fmt.Errorf("%s missing", Property_serial_number))
	} else if _, err := strconv.ParseInt(sn, 10, 64); err != nil {
		errs = append(errs, fmt.Errorf("%s invalid", Property_serial_number))
	}

	if ver := t[Property_version]; ver == "" {
		errs = append(errs, fmt.Errorf("%s missing", Property_version))
	} else if _, err := strconv.ParseInt(ver, 10, 64); err != nil {
		errs = append(errs, fmt.Errorf("%s invalid", Property_version))
	}

	if sub := t[Property_subject]; sub == "" {
		errs = append(errs, fmt.Errorf("%s missing", Property_subject))
	} else if dto, err := resources.ParseDistinguishedName(sub); err != nil {
		errs = append(errs, fmt.Errorf("%s invalid %v", Property_subject, err))
	} else {
		if dto.CommonName == "" {
			errs = append(errs, fmt.Errorf("%s.%s missing", Property_subject, Property_common_name))
		}
	}
	if iss := t[Property_issuer]; iss == "" {
		errs = append(errs, fmt.Errorf("%s missing", Property_issuer))
	} else if dto, err := resources.ParseDistinguishedName(iss); err != nil {
		errs = append(errs, fmt.Errorf("%s invalid %v", Property_issuer, err))
	} else {
		if dto.CommonName == "" {
			errs = append(errs, fmt.Errorf("%s.%s missing", Property_issuer, Property_common_name))
		}
	}

	if puks := t[Property_public_key]; puks == "" {
		errs = append(errs, fmt.Errorf("%s missing", Property_public_key))
	} else if _, err := c.resolveKeyId(puks); err != nil {
		errs = append(errs, fmt.Errorf("%s %v", Property_public_key, err))
	}

	return errs
}

func (c certificateBuilder) Build(t templates.Template) (resources.Resource, error) {
	if errs := c.Validate(t); len(errs) > 0 {
		return nil, errs
	}

	cert, err := c.buildTemplateCertificate(t)
	if err != nil {
		return nil, err
	}

	issuer, err := c.knownIssuers.UserByName(cert.Issuer.String())
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

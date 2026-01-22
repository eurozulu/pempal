package factories

import (
	"crypto/rand"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/repositories"
	"github.com/eurozulu/pempal/templates"
	"time"
)

type CertificateFactory struct{}

func (cf CertificateFactory) Make(ct *templates.CertificateTemplate) ([]model.PemResource, error) {
	var newKey *model.PrivateKey
	if ct.PublicKey == nil {
		prk, err := CreateDefaultKey()
		if err != nil {
			return nil, err
		}
		newKey = prk
		ct.PublicKey = newKey.Public()
	}

	if err := ValidateCertificateTemplate(ct); err != nil {
		return nil, err
	}
	cert := &model.Certificate{}
	ct.ApplyTo(cert)

	var issuer *model.Issuer
	if ct.SelfSigned {
		// If using an existing key, go find it
		if newKey == nil {
			prk, err := repositories.Keys(config.KeyPath()).ByPublicKey(ct.PublicKey)
			if err != nil {
				return nil, fmt.Errorf("no private key found for certificate %s. %v", ct.Subject, err)
			}
			newKey = prk
		}
		issuer = model.NewIssuer(cert, newKey)
	} else {
		// Using existing issuer
		is, err := resolveIssuer(ct.Issuer)
		if err != nil {
			return nil, err
		}
		issuer = is
	}
	der, err := x509.CreateCertificate(
		rand.Reader,
		(*x509.Certificate)(cert),
		(*x509.Certificate)(issuer.Certificate()),
		ct.PublicKey.Public(),
		issuer.PrivateKey().Private(),
	)
	if err != nil {
		return nil, err
	}
	c := &model.Certificate{}
	// marshal DER into pem
	if err := c.UnmarshalBinary(der); err != nil {
		return nil, err
	}
	result := []model.PemResource{c}
	if newKey != nil {
		result = append(result, newKey)
	}
	return result, nil
}

func ValidateCertificateTemplate(ct *templates.CertificateTemplate) error {
	if ct.Subject.IsEmpty() {
		return fmt.Errorf("Invalid certificate template. Certificate subject is empty")
	}
	if ct.Subject.CommonName == "" {
		return fmt.Errorf("Invalid certificate template. Certificate subject: common name is empty")
	}
	if ct.SelfSigned {
		if !ct.Issuer.IsEmpty() && !ct.Issuer.Equals(ct.Subject) {
			return fmt.Errorf("mismatched issuer for self signed certificate %q. For self-signed, leave issuer blank or set to same as subject.", ct.Subject)
		}
		ct.Issuer = ct.Subject
	}
	if ct.Issuer.IsEmpty() {
		return fmt.Errorf("Certificate issuer is empty. use self-signed flag if certificate is self signed")
	}
	//if ct.SerialNumber == nil {
	//	return fmt.Errorf("Serial number is missing")
	//}
	puk := ct.PublicKey
	if puk == nil {
		return fmt.Errorf("Public key is missing")
	}
	if puk.PublicKeyAlgorithm() != ct.PublicKeyAlgorithm {
		ct.PublicKeyAlgorithm = puk.PublicKeyAlgorithm()
	}

	if time.Time(ct.NotAfter).Before(time.Time(ct.NotBefore)) {
		return errors.New("not_after is before not_before")
	}
	return nil
}

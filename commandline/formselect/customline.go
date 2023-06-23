package formselect

import (
	"crypto/x509"
	"fmt"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/utils"
)

var errInvalidAlgorithm = fmt.Errorf("invalid public-key-algorithm")

type PublicKeyAlgorithmLine struct {
	name     string
	value    x509.PublicKeyAlgorithm
	required bool
}

func (cl PublicKeyAlgorithmLine) Name() string {
	return cl.name
}

func (cl PublicKeyAlgorithmLine) Required() bool {
	return cl.required
}

func (cl PublicKeyAlgorithmLine) Format() string {
	if cl.value == x509.UnknownPublicKeyAlgorithm {
		return errText + errInvalidAlgorithm.Error() + errText
	}
	return cl.value.String()
}

func (cl *PublicKeyAlgorithmLine) Parse(s string) error {
	pka := utils.ParsePublicKeyAlgorithm(s)
	if cl.value == x509.UnknownPublicKeyAlgorithm {
		return errInvalidAlgorithm
	}
	var c model.CertificateDTO
	cl.value = pka
	return nil
}

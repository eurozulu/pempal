package factories

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/templates"
	"github.com/eurozulu/pempal/utils"
	"github.com/eurozulu/pempal/validation"
	"os"
	"path/filepath"
	"strings"
)

type csrfactory struct {
	keyrepo resources.Keys
}

func (cf csrfactory) Build(template templates.Template) ([]byte, error) {
	vld := validation.ValidatorForTemplate("certificaterequest")
	if err := vld.Validate(template); err != nil {
		return nil, err
	}
	t := template.(*templates.CSRTemplate)
	// public key for the certificate
	puk, err := cf.publicKeyForCSR(t)
	if err != nil {
		return nil, err
	}
	id, err := model.NewKeyIdFromKey(puk)
	if err != nil {
		return nil, err
	}

	prk, err := cf.keyrepo.PrivateKeyFromID(id)
	if err != nil {
		return nil, validation.ValidationError{
			PropertyName: "public-key",
			Message:      err.Error(),
		}
	}

	der, err := x509.CreateCertificateRequest(rand.Reader, t.ToCSR(), prk)
	if err != nil {
		return nil, err
	}
	// write result to file
	data := utils.DERToPem(der, model.CertificateRequest)
	return data, nil
}

func (cf csrfactory) publicKeyForCSR(template *templates.CSRTemplate) (crypto.PublicKey, error) {
	if template.PublicKey.PublicKey != nil {
		return template.PublicKey.PublicKey, nil
	}
	if template.ID != nil {
		return cf.keyrepo.PublicKeyFromID(template.ID)
	}
	return nil, validation.ValidationError{
		PropertyName: "public-key",
		Message:      "invalid CSR no 'public-key' or 'id' found",
	}
}

func saveCSR(path string, pemdata []byte) (string, error) {
	blk, _ := pem.Decode(pemdata)
	cert, err := x509.ParseCertificateRequest(blk.Bytes)
	if err != nil {
		return "", err
	}
	fp := strings.Join([]string{strings.ReplaceAll(cert.Subject.String(), "/", "_"), ".pem"}, "")
	return fp, os.WriteFile(filepath.Join(path, fp), pemdata, 0644)
}

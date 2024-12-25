package factories

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/logging"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/templates"
	"github.com/eurozulu/pempal/utils"
	"github.com/eurozulu/pempal/validation"
	"os"
	"path/filepath"
	"strings"
)

type CRLFactory struct {
	certrepo resources.Certificates
	keyrepo  resources.Keys
}

func (cf CRLFactory) Build(template templates.Template) ([]byte, error) {
	vld := &validation.CRLValidator{}
	if err := vld.Validate(template); err != nil {
		return nil, err
	}
	t := template.(*templates.CRLTemplate)

	// find issuer key and certificate to sign with
	issuerCert, err := cf.issuerCertificate(t.Issuer.ToName())
	if err != nil {
		return nil, err
	}
	logging.Debug("CRLFactory", "creating new CRL issued by %s", issuerCert.Subject.String())

	issuerId, err := model.NewKeyIdFromKey(issuerCert.PublicKey)
	if err != nil {
		return nil, err
	}
	logging.Debug("CRLFactory", "issued ID: %s", issuerId.String())

	prk, err := cf.keyrepo.PrivateKeyFromID(issuerId)
	if err != nil {
		return nil, err
	}
	logging.Debug("CRLFactory", "located private key for ", issuerId.String())
	signer, ok := prk.(crypto.Signer)
	if !ok {
		return nil, fmt.Errorf("unknown key type %T.  Expected Signer key", prk)
	}
	der, err := x509.CreateRevocationList(rand.Reader, t.ToCRL(), issuerCert, signer)
	if err != nil {
		return nil, err
	}
	return utils.DERToPem(der, model.RevokationList), nil
}

func (cf CRLFactory) issuerCertificate(name pkix.Name) (*x509.Certificate, error) {
	// find issuer key and certificate to sign with
	certs := cf.certrepo.CertificatesBySubject(name)
	if len(certs) == 0 {
		return nil, fmt.Errorf("issuer %q no known", name.String())
	}
	if len(certs) > 1 {
		return nil, fmt.Errorf("issuer %q has more than one certificate", name.String())
	}
	return certs[0], nil
}

func saveCRL(path string, pemdata []byte) (string, error) {
	blk, _ := pem.Decode(pemdata)
	crl, err := x509.ParseRevocationList(blk.Bytes)
	if err != nil {
		return "", err
	}
	fp := strings.Join([]string{
		strings.ReplaceAll(crl.Issuer.String(), "/", "_"),
		fmt.Sprintf("_%40d", crl.Number.Int64()),
		".pem"}, "")
	return fp, os.WriteFile(filepath.Join(path, fp), pemdata, 0644)
}

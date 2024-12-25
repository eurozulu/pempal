package factories

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
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

type certificateFactory struct {
	certrepo resources.Certificates
	keyrepo  resources.Keys
}

func (cf certificateFactory) Build(template templates.Template) ([]byte, error) {
	cv := validation.ValidatorForTemplate("certificate")
	if err := cv.Validate(template); err != nil {
		return nil, err
	}
	t := template.(*templates.CertificateTemplate)
	certKey := t.PublicKey.PublicKey
	// issuer key and certificate to sign with
	issuerCert, err := cf.issuerCertificate(t)
	if err != nil {
		return nil, err
	}
	issuerKey, err := cf.issuerKey(issuerCert)
	if err != nil {
		return nil, err
	}

	der, err := x509.CreateCertificate(rand.Reader, t.ToCertificate(), issuerCert, certKey, issuerKey)
	if err != nil {
		return nil, err
	}
	// write result to file
	data := utils.DERToPem(der, model.Certificate)
	return data, nil
}

func (cf certificateFactory) issuerCertificate(template *templates.CertificateTemplate) (*x509.Certificate, error) {
	if template.SelfSigned {
		return template.ToCertificate(), nil
	}
	certs := cf.certrepo.CertificatesBySubject(template.Issuer.ToName())
	if len(certs) == 0 {
		return nil, &validation.ValidationError{
			PropertyName: "issuer",
			Message:      fmt.Sprintf("%s not known", template.Issuer.String()),
		}
	}
	if len(certs) > 1 {
		return nil, &validation.ValidationError{
			PropertyName: "issuer",
			Message:      fmt.Sprintf("%s is ambiguous. It matches %d certificates", template.Issuer.String(), len(certs)),
		}
	}
	return certs[0], nil
}

func (cf certificateFactory) issuerKey(issuerCert *x509.Certificate) (crypto.PrivateKey, error) {
	issuerId, err := model.NewKeyIdFromKey(issuerCert.PublicKey)
	if err != nil {
		return nil, err
	}
	logging.Debug("CertificateFactory", "issuer %s id %s", issuerCert.Subject, issuerId)

	// locate private key for the given issuerCert
	prk, err := cf.keyrepo.PrivateKeyFromID(issuerId)
	if err != nil {
		return nil, &validation.ValidationError{
			PropertyName: "issuer",
			Message:      fmt.Sprintf("could not locate private key for %q. %v", issuerCert.Subject.String(), err.Error()),
		}
	}
	return prk, nil
}

func saveCertificate(path string, pemdata []byte) (string, error) {
	blk, _ := pem.Decode(pemdata)
	cert, err := x509.ParseCertificate(blk.Bytes)
	if err != nil {
		return "", err
	}
	fn := strings.ReplaceAll(cert.Subject.String(), "/", "_")
	fn = strings.ReplaceAll(fn, " ", "_")
	fp := uniqueName(path, fn+".pem")
	if err := utils.EnsureDirExists(path); err != nil {
		return "", err
	}
	return fp, os.WriteFile(fp, pemdata, 0644)
}

func uniqueName(path, name string) string {
	fn := name
	ext := filepath.Ext(name)
	name = name[:len(name)-len(ext)]
	count := 0
	for utils.FileExists(filepath.Join(path, fn)) {
		fn = fmt.Sprintf("%s_%04d%s", name, count, ext)
		count++
	}
	return filepath.Join(path, fn)
}

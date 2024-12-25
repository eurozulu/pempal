package validation

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/resources"
	"math/big"
	"sort"
)

type SerialNumberFactory struct {
}

func (snf SerialNumberFactory) NextSerialNumberFor(issuer pkix.Name) *big.Int {
	cert := getLastCertificateOfIssuer(issuer)
	if cert == nil {
		return big.NewInt(1)
	}
	return cert.SerialNumber.Add(cert.SerialNumber, big.NewInt(1))
}

func getLastCertificateOfIssuer(issuer pkix.Name) *x509.Certificate {
	certPath := resources.NewCertificates(config.Config.CertPath)
	certs := certPath.CertificatesBySubject(issuer)
	if len(certs) == 0 {
		return nil
	}
	sort.Slice(certs, func(i, j int) bool {
		return certs[i].SerialNumber.Int64() < certs[j].SerialNumber.Int64()
	})
	return certs[0]
}

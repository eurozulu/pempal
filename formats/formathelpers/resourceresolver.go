package formathelpers

import (
	"crypto"
	"crypto/x509"
	"pempal/templates"
)

func ResolveKeys(puk string) (crypto.PrivateKey, crypto.PublicKey, error) {

}

func ResolvePublicKey(puk string) (crypto.PublicKey, error) {

}

func ResolveIssuer(dn templates.NameTemplate, locations []string) (*x509.Certificate, error) {

}

func MatchPrivateKeyFromPublic(puk crypto.PublicKey) (crypto.PrivateKey, error) {

}

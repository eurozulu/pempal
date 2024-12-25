package factories

import (
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/templates"
)

type Factory interface {
	Build(t templates.Template) ([]byte, error)
}

func FactoryForType(name string) (Factory, error) {
	switch name {
	case "privatekey":
		return &keyFactory{}, nil
	case "certificate":
		return &certificateFactory{
			certrepo: resources.NewCertificates(config.Config.CertPath),
			keyrepo:  resources.NewKeys(config.Config.CertPath),
		}, nil
	case "request":
		return &csrfactory{
			keyrepo: resources.NewKeys(config.Config.CertPath),
		}, nil
	case "revoke":
		return &CRLFactory{
			certrepo: resources.NewCertificates(config.Config.CertPath),
			keyrepo:  resources.NewKeys(config.Config.CertPath),
		}, nil
	default:
		return nil, fmt.Errorf("unknown template type: %s", name)
	}
}

func SaveResource(pemdata []byte) (string, error) {
	blk, _ := pem.Decode(pemdata)
	rt := model.ParseResourceTypeFromPEMType(blk.Type)
	path, err := config.Config.ResourcePath(rt)
	if err != nil {
		return "", err
	}

	switch rt {
	case model.PublicKey:
		return savePublicKey(path, pemdata)
	case model.PrivateKey:
		return savePrivateKey(path, pemdata)
	case model.Certificate:
		return saveCertificate(path, pemdata)
	case model.CertificateRequest:
		return saveCSR(path, pemdata)
	case model.RevokationList:
		return saveCRL(path, pemdata)
	default:
		return "", fmt.Errorf("unknown resource type: %s", rt)
	}
}

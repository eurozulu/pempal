package model

import (
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/logging"
)

type PemFile struct {
	Path   string
	Blocks []*pem.Block
}

func (f PemFile) Resources() []PemResource {
	resz := make([]PemResource, 0, len(f.Blocks))
	for _, blk := range f.Blocks {
		var res PemResource
		var err error
		switch ParseResourceType(blk.Type) {
		case ResourceTypeCertificate:
			res, err = NewCertificateFromPem(blk)
		case ResourceTypePrivateKey:
			res, err = NewPrivateKeyFromPem(blk)
		case ResourceTypePublicKey:
			res, err = NewPublicKeyFromPem(blk)
		case ResourceTypeCertificateRequest:
			res, err = NewCertificateRequestFromPem(blk)
		case ResourceTypeRevokationList:
			res, err = NewRevocationListFromPem(blk)
		default:
			err = fmt.Errorf("unsupported resource type %v ignored", blk.Type)
		}
		if err != nil {
			logging.Warning("failed to read pem %s %v", f.Path, err)
			continue
		}
		resz = append(resz, res)
	}
	return resz
}

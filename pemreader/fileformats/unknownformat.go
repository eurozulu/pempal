package fileformats

import (
	"encoding/pem"
	"fmt"
)

type unknownFormat struct {
}

func (u unknownFormat) Format(by []byte) ([]*pem.Block, error) {
	blks, err := pemFileFormat{}.Format(by)
	if err == nil {
		return blks, nil
	}
	blks, err = derKeyFormat{}.Format(by)
	if err == nil {
		return blks, nil
	}
	blks, err = derCertificateFormat{}.Format(by)
	if err == nil {
		return blks, nil
	}
	blks, err = derCertificateRequestFormat{}.Format(by)
	if err == nil {
		return blks, nil
	}
	return nil, fmt.Errorf("unknown file format")
}

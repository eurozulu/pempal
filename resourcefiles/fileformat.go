package resourcefiles

import (
	"encoding/pem"
	"github.com/eurozulu/pempal/logging"
)

type FileFormat interface {
	Format(data []byte) ([]*pem.Block, error)
}

var knownPemFormats = []FileFormat{
	&PemFileFormat{},
	&Derfileformat{},
}

func FormatAsPems(data []byte) ([]*pem.Block, error) {
	var blocks []*pem.Block
	for _, format := range knownPemFormats {
		pemz, err := format.Format(data)
		if err != nil {
			logging.Warning("Failed to parse pem: %s", err)
			continue
		}
		if len(pemz) == 0 {
			continue
		}
		blocks = append(blocks, pemz...)
		break
	}
	return blocks, nil
}

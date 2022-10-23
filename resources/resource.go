package resources

import (
	"encoding"
	"encoding/pem"
)

type Resource interface {
	ResourceType() string
	Location() string
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

type PemMarshaler interface {
	MarshalPEM() (data []byte, err error)
}

type PemUnmarshaler interface {
	UnmarshalPem(data []byte) error
}

func UnmarshalResources(data []byte) ([]Resource, error) {
	var resources []Resource
	resources, _ = parsePemBlocks(data)
	if len(resources) > 0 {
		return resources, nil
	}
	r, err := parseDerData(data)
	if err != nil {
		return nil, err
	}
	return []Resource{r}, err
}

func parsePemBlocks(data []byte) ([]Resource, []byte) {
	var blocks []Resource
	var blk *pem.Block
	d := data[:]
	for len(d) > 0 {
		blk, d = pem.Decode(d)
		if blk == nil {
			break
		}
		blocks = append(blocks, &pemResource{block: blk})
	}
	return blocks, d
}

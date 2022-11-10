package finder

import (
	"bytes"
	"encoding"
	"encoding/pem"
	"fmt"
)

type pemLocation struct {
	path string
	pems []*pem.Block
}

func (r pemLocation) Path() string {
	return r.path
}

func (r pemLocation) MarshalText() (text []byte, err error) {
	m := map[string]int{}
	for _, p := range r.pems {
		m[p.Type] += 1
	}
	buf := bytes.NewBuffer(nil)
	if _, err := fmt.Fprintln(buf, r.Path()); err != nil {
		return nil, err
	}
	for k, v := range m {
		if _, err := fmt.Fprintf(buf, "\t%s\t%d\n", k, v); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (r pemLocation) MarshalBinary() (data []byte, err error) {
	buf := bytes.NewBuffer(nil)
	for _, p := range r.pems {
		if err := pem.Encode(buf, p); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func ReadLocationPems(l Location) ([]*pem.Block, error) {
	data, err := l.(encoding.BinaryMarshaler).MarshalBinary()
	if err != nil {
		return nil, err
	}
	var blocks []*pem.Block
	var blk *pem.Block
	for len(data) > 0 {
		blk, data = pem.Decode(data)
		if blk == nil {
			break
		}
		blocks = append(blocks, blk)
	}
	return blocks, nil
}

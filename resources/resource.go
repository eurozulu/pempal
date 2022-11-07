package resources

import (
	"bytes"
	"encoding/pem"
	"fmt"
	"strings"
)

const locationToken = "#@"

// Resource represents a single x509 resource.
type Resource interface {
	// The resource as a pem block
	Pem() *pem.Block
	// Type represents the type of resource, CERTIFICATE, CERTIFICATE REQUEST, PRIVATE KEY etc
	Type() ResourceType
	// Location is the optional location of the resource, if it has a persistent location
	Location() string
}

type Resources []Resource

type pemResource struct {
	location string
	blk      *pem.Block
}

func (r pemResource) Pem() *pem.Block {
	return &pem.Block{
		Type:    r.blk.Type,
		Headers: r.blk.Headers,
		Bytes:   r.blk.Bytes,
	}
}

func (r pemResource) Type() ResourceType {
	if r.blk == nil {
		return Unknown
	}
	return ParseResourceType(r.blk.Type)
}

func (r pemResource) Location() string {
	return r.location
}

func (r pemResource) MarshalText() (text []byte, err error) {
	buf := bytes.NewBuffer(nil)
	if r.location != "" {
		buf.WriteString(locationToken)
		buf.WriteString(r.location)
		buf.WriteByte('\n')
	}
	if r.blk != nil {
		if _, err := buf.Write(pem.EncodeToMemory(r.blk)); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), err
}

func (r *pemResource) UnmarshalText(text []byte) error {
	var leadingLine string
	// optional leading line starting with '@' indicates a location present on first line
	if bytes.HasPrefix(text, []byte{'@'}) {
		leadingLine, text = readLine(text)
	}
	b, _ := pem.Decode(text)
	if b == nil {
		return fmt.Errorf("failed to read as PEM")
	}
	r.location = strings.TrimLeft(leadingLine, "@")
	r.blk = b
	return nil
}

func (r pemResource) MarshalBinary() (data []byte, err error) {
	if r.blk == nil {
		return nil, nil
	}
	return r.blk.Bytes, err
}

func readLine(text []byte) (string, []byte) {
	i := bytes.IndexByte(text, '\n')
	if i < 0 {
		return string(text), nil
	}
	l := string(text[:i])
	var remain []byte
	if i+1 < len(text) {
		remain = text[i+1:]
	}
	return l, remain
}

func NewResource(location string, p *pem.Block) Resource {
	return &pemResource{
		location: location,
		blk:      p,
	}
}

package model

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

type RevocationList x509.RevocationList

type RevocationListEntry x509.RevocationListEntry

func (r RevocationList) ResourceType() ResourceType {
	return ResourceTypeRevokationList
}

func (r RevocationList) String() string {
	return fmt.Sprintf("%s @%v", r.Issuer.String(), r.ThisUpdate)
}

func (r RevocationList) Fingerprint() Fingerprint {
	return NewFingerPrint(r.Raw)
}

func (r RevocationList) MarshalBinary() (data []byte, err error) {
	return r.Raw, nil
}

func (r *RevocationList) UnmarshalBinary(data []byte) error {
	crl, err := x509.ParseRevocationList(data)
	if err != nil {
		return err
	}
	*r = RevocationList(*crl)
	return nil
}

func (r RevocationList) MarshalText() (text []byte, err error) {
	der, err := r.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  ResourceTypeRevokationList.String(),
		Bytes: der,
	}), nil
}

func (r RevocationList) UnmarshalText(text []byte) error {
	blk, _ := pem.Decode(text)
	if blk == nil {
		return fmt.Errorf("no pem found")
	}
	if ParseResourceType(blk.Type) != ResourceTypeRevokationList {
		return fmt.Errorf("pem not a %s, found %s", ResourceTypeRevokationList, blk.Type)
	}
	return r.UnmarshalBinary(blk.Bytes)
}

func NewRevocationListFromPem(blk *pem.Block) (*RevocationList, error) {
	crl := &RevocationList{}
	if err := crl.UnmarshalText(pem.EncodeToMemory(blk)); err != nil {
		return nil, err
	}
	return crl, nil
}

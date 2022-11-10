package pemtypes

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"gopkg.in/yaml.v3"
	"pempal/templates"
)

type revokeListType struct {
	crl x509.RevocationList
}

func (rt revokeListType) String() string {
	return fmt.Sprintf("%s\t%s", RevocationList.String(), rt.crl.Issuer.String())
}

func (rt revokeListType) MarshalBinary() (data []byte, err error) {
	return rt.crl.Raw, nil
}

func (rt revokeListType) UnmarshalBinary(data []byte) error {
	cl, err := x509.ParseRevocationList(data)
	if err != nil {
		return err
	}
	rt.crl = *cl
	return nil
}

func (rt revokeListType) MarshalText() (text []byte, err error) {
	der, err := rt.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  RevocationList.String(),
		Bytes: der,
	}), nil
}

func (rt revokeListType) UnmarshalText(text []byte) error {
	blocks := ReadPEMBlocks(text, RevocationList)
	if len(blocks) == 0 {
		return fmt.Errorf("no revokation list pem found")
	}
	return rt.UnmarshalBinary(blocks[0].Bytes)
}

func (rt revokeListType) MarshalYAML() (interface{}, error) {
	t := templates.CRLTemplate{}
	rt.applyToTemplate(&t)
	return yaml.Marshal(&t)
}

// UnmarshalYAML attempts to read the given value as a YAML encoded certificate
func (rt *revokeListType) UnmarshalYAML(value *yaml.Node) error {
	t := templates.CRLTemplate{}
	if err := value.Decode(&t); err != nil {
		return err
	}
	rt.applyTemplate(t)
	return nil
}

func (rt revokeListType) applyToTemplate(t *templates.CRLTemplate) {
	if rt.crl.Number != nil {
		t.Number = rt.crl.Number.Int64()
	}

}
func (rt *revokeListType) applyTemplate(t templates.CRLTemplate) {

}

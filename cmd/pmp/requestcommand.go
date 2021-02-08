package main

import (
	"bytes"
	"encoding/pem"
	"fmt"
	"github.com/pempal/pemio"
	"github.com/pempal/pempal"
	"github.com/pempal/templates"
	"strings"
)

type RequestCommand struct {
	Command
	Version            int                                   `yaml:"Version,omitempty" flag:"version"`
	Subject            map[string]interface{}                `yaml:"Subject,omitempty" flag:"subject"`
	PublicKey          *templates.PublicKeyTemplate          `yaml:"PublicKey,omitempty" flag:"publickey"`
	PublicKeyAlgorithm *templates.PublicKeyAlgorithmTemplate `yaml:"PublicKeyAlgorithm,omitempty" flag:"publickeyalgorithmtemplate,pka"`
	SignatureAlgorithm *templates.SignatureAlgorithmTemplate `yaml:"SignatureAlgorithm,omitempty" flag:"signaturealgorithm,sa"`
	DNSNames           []string                              `yaml:"DNSNames,omitempty" flag:"dnsnames,dns"`
	EmailAddresses     []string                              `yaml:"EmailAddresses,omitempty" flag:"emailaddresses,emails"`
	IPAddresses        []*templates.IPAddressTemplate        `yaml:"IPAddresses,omitempty" flag:"ipaddresses,ipa,ipas"`
	URIs               []*templates.URIsTemplate             `yaml:"URIs,omitempty" flag:"uri,uris"`

	Key string `yaml:"-" flag:"key,k"`

	// KeyOut, when true will output the private key after the request is written.
	// IF a new key is created during a request, Keyout if forced to true to prevent the new key being lost
	// If key is encrypted, the encrypted form is output.
	KeyOut   bool   `yaml:"-" flag:"keyout,ko"`
	Password string `yaml:"-" flag:"password,p"`

	Script bool `yaml:"-" flag:"script,s"`
}

func (rc RequestCommand) String() string {
	return "create a new certificate signing request"
}

// Request is the command line entry point.
func (rc RequestCommand) Request(tps ...string) error {
	ts, err := NewTemplateFiles(tps)
	if err != nil {
		return err
	}

	// Create a new request and apply the templates
	nrt := &templates.RequestTemplate{}
	if err := templates.ApplyTemplates(nrt, ts...); err != nil {
		return err
	}
	// Apply this command as a template to assign any matching flag values
	if err := templates.ApplyTemplates(nrt, &rc); err != nil {
		return err
	}

	// locate the private key to sign request
	prk, err := rc.findRequestKey()
	if err != nil {
		return err
	}
	puk, err := prk.PublicKey()
	if err != nil {
		return fmt.Errorf("failed to read public key %v", err)
	}
	nrt.PublicKey = puk
	if err := ConfirmTemplate("sign new request", nrt); err != nil {
		return err
	}

	bl, err := pempal.NewRequest(nrt, prk)
	if err != nil {
		return err
	}

	if err := rc.WriteOutput([]*pem.Block{bl}, 0644); err != nil {
		return err
	}
	if rc.KeyOut {
		rc.Truncate = false
		if err := rc.WriteOutput([]*pem.Block{prk.PEMBlock()}, 0600); err != nil {
			return err
		}
	}
	return nil
}

func (rc *RequestCommand) findRequestKey() (*templates.PrivateKeyTemplate, error) {
	if rc.Key == "-" {
		return rc.keyFromInput()
	}
	if rc.Key == "+" {
		rc.KeyOut = true // Force key out as to prevent new being lost
		return rc.requestNewKey()
	}

	// Search for private keys in current dir or rc.Key (key may be filename OR a query)
	qrs, err := pempal.FindKeys(rc.Key, rc.Password)
	if err != nil {
		return nil, err
	}

	if len(qrs) == 0 {
		if !pempal.PromptConfirm("No private key found, would you like to generate a new key?", true) {
			return nil, fmt.Errorf("aborted, no key found")
		}
		return rc.requestNewKey()
	}

	i := ChooseTemplate("for the private key to sign new request", qrs, []string{"Create new key"})
	// was 'option' New Key choosen
	if i == len(qrs) {
		rc.KeyOut = true // Force key out as to prevent new being lost
		return rc.requestNewKey()
	}
	return templates.NewPrivateKeyTemplate(qrs[i].Block), nil
}

func (rc RequestCommand) requestNewKey() (*templates.PrivateKeyTemplate, error) {
	buf := bytes.NewBuffer(nil)
	kc := KeyCommand{
		Command:            Command{Output: buf},
		PublicKeyAlgorithm: rc.PublicKeyAlgorithm.String(),
		Password:           rc.Password,
	}

	if err := kc.Key(); err != nil {
		return nil, err
	}
	bls, err := pemio.ReadPEMs(buf)
	if err != nil {
		return nil, err
	}
	return templates.NewPrivateKeyTemplate(bls[0]), nil
}

func (rc RequestCommand) keyFromInput() (*templates.PrivateKeyTemplate, error) {
	bls, err := rc.ReadInput()
	if err != nil {
		return nil, err
	}
	for _, bl := range bls {
		if !strings.Contains(bl.Type, "PRIVATE KEY") {
			continue
		}
		return templates.NewPrivateKeyTemplate(bl), nil
	}
	return nil, fmt.Errorf("no private key found in %s", rc.Key)
}

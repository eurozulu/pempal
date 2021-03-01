package main

import (
	"encoding/pem"
	"github.com/pempal/pempal"
	"github.com/pempal/templates"
)

type KeyCommand struct {
	// PublicKeyAlgorithm defines the type of key.
	PublicKeyAlgorithm string `flag:"publickeyalgorithm,pka" yaml:"PublicKeyAlgorithm"`
	// PublicKeyKeySize defines the complexity/length of the key
	PublicKeyKeySize string `flag:"keysize,size,s" yaml:"PublicKeyKeySize"`

	IsEncrypted bool `flag:"isencrypted,encrypted,e" yaml:"IsEncrypted"`

	// Password defines the password to set on the new key
	// When provided, implies 	IsEncrypted = true
	Password string `flag:"password,p" yaml:"Password"`

	// PEMCipher defines the cipher to encrypt the new key. (If password empty, this is ignored)
	// Valid values are: "PEMCipherDES","PEMCipher3DES","PEMCipherAES128","PEMCipherAES192","PEMCipherAES256"
	PEMCipher string `flag:"PEMCipherTemplate,c" yaml:"PEMCipher"`

	Script bool `flag:"script,s" yaml:"-"`
}

func (kc KeyCommand) String() string {
	return "key command generates a new key."
}

func (kc KeyCommand) Key(tps ...string) error {
	ts, err := NewTemplateFiles(tps)
	if err != nil {
		return err
	}

	// Stack any given templates into new key template, then apply 'this' command so its flags are applied.
	nkt := templates.NewNewKeyTemplate()
	if err := templates.ApplyTemplates(nkt, ts...); err != nil {
		return err
	}
	if err := templates.ApplyTemplates(nkt, &kc); err != nil {
		return err
	}

	if !kc.Script {
		if err := ConfirmTemplate("Create new key:", nkt); err != nil {
			return err
		}
	}

	bl, err := pempal.NewKey(nkt)
	if err != nil {
		return err
	}
	return writePemsToOutput([]*pem.Block{bl}, 0600)
}

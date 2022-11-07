package templates

import "pempal/resources"

type KeyTemplate struct {
	KeyType     string `yaml:"key-type"`
	Size        string `yaml:"size"`
	IsEncrypted bool   `yaml:"is-encrypted"`
	PublicKey   string `yaml:"public-key,omitempty"`
}

func (k KeyTemplate) Type() resources.ResourceType {
	return resources.PrivateKey
}

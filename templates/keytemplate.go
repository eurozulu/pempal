package templates

type KeyTemplate struct {
	PublicKey   *PublicKeyTemplate `yaml:"public-key,omitempty"`
	IsEncrypted bool               `yaml:"is-encrypted"`
}

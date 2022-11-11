package templates

type PrivateKeyTemplate struct {
	PublicKey   *PublicKeyTemplate `yaml:"public-key,omitempty"`
	IsEncrypted bool               `yaml:"is-encrypted"`
}

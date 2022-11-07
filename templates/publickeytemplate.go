package templates

type PublicKeyTemplate struct {
	KeyType   string `yaml:"key-type"`
	Size      string `yaml:"size"`
	PublicKey string `yaml:"public-key,omitempty"`
}

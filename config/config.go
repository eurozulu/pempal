package config

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"os"
	"path/filepath"
	"pempal/utils"
	"strings"
)

const (
	ENV_CA_ROOT             = "CA_ROOT"
	ENV_CA_ROOT_CERTIFICATE = "CA_ROOT_CERT"
	ENV_CA_CERTS            = "CA_CERTS"
	ENV_CA_KEYS             = "CA_KEYS"
	ENV_CA_REQUESTS         = "CA_REQUESTS"
	ENV_CA_REVOKELISTS      = "CA_REVOKELISTS"
	ENV_CA_TEMPLATES        = "CA_TEMPLATES"
	ENV_CA_CONFIG           = "CA_CONFIG"
)

const defaultConfigPath = "$HOME/.pempal/.config"

type Config interface {
	Root() string
	RootCertificate() string
	Certificates() string
	Keys() string
	Requests() string
	Revokations() string
	Templates() string
	ConfigLocation() string
}

type DefaultConfig struct {
	RootPath            string `yaml:"root-path"`
	RootCertificatePath string `yaml:"root-certificate,omitempty"`
	CertPath            string `yaml:"cert-path"`
	KeyPath             string `yaml:"key-path"`
	CsrPath             string `yaml:"csr-path"`
	CrlPath             string `yaml:"crl-path"`
	TemplatePath        string `yaml:"template-path"`
	configLocation      string `yaml:"-"`
}

func (cfg DefaultConfig) Root() string {
	return cfg.RootPath
}

func (cfg DefaultConfig) RootCertificate() string {
	return cfg.resolvePath(cfg.RootCertificatePath)
}

func (cfg DefaultConfig) Certificates() string {
	return cfg.resolvePath(cfg.CertPath)
}

func (cfg DefaultConfig) Keys() string {
	return cfg.resolvePath(cfg.KeyPath)
}

func (cfg DefaultConfig) Requests() string {
	return cfg.resolvePath(cfg.CsrPath)
}

func (cfg DefaultConfig) Revokations() string {
	return cfg.resolvePath(cfg.CrlPath)
}

func (cfg DefaultConfig) Templates() string {
	return cfg.resolvePath(cfg.TemplatePath)
}

func (cfg DefaultConfig) ConfigLocation() string {
	return cfg.configLocation
}

func (cfg DefaultConfig) resolvePath(path string) string {
	path = os.ExpandEnv(path)
	if path == "" || filepath.IsLocal(path) {
		path = filepath.Join(os.ExpandEnv(cfg.RootPath), path)
	}
	return path
}

func applyENVValues(cfg *DefaultConfig) {
	cfg.RootPath = envOrDefault(ENV_CA_ROOT, cfg.RootPath)
	cfg.RootCertificatePath = envOrDefault(ENV_CA_ROOT_CERTIFICATE, cfg.RootCertificatePath)
	cfg.CertPath = envOrDefault(ENV_CA_CERTS, cfg.CertPath)
	cfg.KeyPath = envOrDefault(ENV_CA_KEYS, cfg.KeyPath)
	cfg.CsrPath = envOrDefault(ENV_CA_REQUESTS, cfg.CsrPath)
	cfg.CrlPath = envOrDefault(ENV_CA_REVOKELISTS, cfg.CrlPath)
	cfg.TemplatePath = envOrDefault(ENV_CA_TEMPLATES, cfg.TemplatePath)
}

func envOrDefault(name string, def string) string {
	s, ok := os.LookupEnv(name)
	if !ok {
		return def
	}
	return s
}

func LoadConfig(name string, cfg Config) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = yaml.NewDecoder(f).Decode(cfg); err != nil {
		return err
	}
	return nil
}

func SaveConfig(name string, cfg Config) error {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewEncoder(f).Encode(&cfg)
}

func resolveConfigPath(path string) (string, error) {
	cfgPath := envOrDefault(ENV_CA_CONFIG, "")
	if path != "" {
		cfgPath = path
	}
	if cfgPath == "" && utils.FileExists(os.ExpandEnv(defaultConfigPath)) {
		cfgPath = defaultConfigPath
	}
	if cfgPath != "" {
		cfgPath = os.ExpandEnv(cfgPath)
		if filepath.IsLocal(cfgPath) {
			cfgPath = filepath.Join(os.ExpandEnv("$PWD"), cfgPath)
		}
		if utils.FileExists(cfgPath) {
			return cfgPath, nil
		}
		if !strings.HasSuffix(strings.ToLower(cfgPath), ".config") {
			cfgPath = filepath.Join(cfgPath, ".config")
			if utils.FileExists(cfgPath) {
				return cfgPath, nil
			}
		}
		return "", fmt.Errorf("config file %s not found", cfgPath)
	}
	return cfgPath, nil
}

func NewConfig(path string) (Config, error) {
	cfgPath, err := resolveConfigPath(path)
	if err != nil {
		return nil, err
	}
	cfg := &DefaultConfig{configLocation: cfgPath}
	if cfgPath != "" {
		if err := LoadConfig(cfgPath, cfg); err != nil {
			return nil, err
		}
	}
	applyENVValues(cfg)
	if cfg.RootPath == "" || cfg.RootPath == "." {
		root := "$PWD"
		if cfgPath != "" {
			root = filepath.Dir(cfgPath)
		}
		cfg.RootPath = os.ExpandEnv(root)
	}
	return cfg, nil
}

package config

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"os"
	"path/filepath"
	"pempal/logger"
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

type Config struct {
	RootPath        string `yaml:"root-path"`
	RootCertificate string `yaml:"root-certificate,omitempty"`
	CertPath        string `yaml:"cert-path"`
	KeyPath         string `yaml:"key-path"`
	CsrPath         string `yaml:"csr-path"`
	CrlPath         string `yaml:"crl-path"`
	TemplatePath    string `yaml:"template-path"`
	ConfigPath      string `yaml:"-"`
}

var defaultConfig = Config{
	RootPath:        "$PWD",
	RootCertificate: "root-certificate.pem",
	CertPath:        "",
	KeyPath:         "",
	CsrPath:         "",
	CrlPath:         "",
	TemplatePath:    "templates",
	ConfigPath:      "$HOME/.pempal/.config",
}

func (cfg Config) ResolveWithRootPath(p string) string {
	return filepath.Join(cfg.RootPath, p)
}

func resolvePaths(cfg *Config) {
	cfg.RootPath = resolvePath("", cfg.RootPath)
	cfg.RootCertificate = resolvePath(cfg.RootPath, cfg.RootCertificate)
	cfg.CertPath = resolvePath(cfg.RootPath, cfg.CertPath)
	cfg.KeyPath = resolvePath(cfg.RootPath, cfg.KeyPath)
	cfg.CsrPath = resolvePath(cfg.RootPath, cfg.CsrPath)
	cfg.CrlPath = resolvePath(cfg.RootPath, cfg.CrlPath)
	cfg.TemplatePath = resolvePath(cfg.RootPath, cfg.TemplatePath)
}

func resolvePath(base, path string) string {
	path = os.ExpandEnv(path)
	if path == "" || filepath.IsLocal(path) {
		path = filepath.Join(base, path)
	}
	return path
}

func applyENVValues(cfg *Config) {
	cfg.RootPath = envOrDefault(ENV_CA_ROOT, cfg.RootPath)
	cfg.RootCertificate = envOrDefault(ENV_CA_ROOT_CERTIFICATE, cfg.RootCertificate)
	cfg.CertPath = envOrDefault(ENV_CA_CERTS, cfg.CertPath)
	cfg.KeyPath = envOrDefault(ENV_CA_KEYS, cfg.KeyPath)
	cfg.CsrPath = envOrDefault(ENV_CA_REQUESTS, cfg.CsrPath)
	cfg.CrlPath = envOrDefault(ENV_CA_REVOKELISTS, cfg.CrlPath)
	cfg.TemplatePath = envOrDefault(ENV_CA_TEMPLATES, cfg.TemplatePath)
	cfg.ConfigPath = envOrDefault(ENV_CA_CONFIG, cfg.ConfigPath)
}

func envOrDefault(name string, def string) string {
	s, ok := os.LookupEnv(name)
	if !ok {
		return def
	}
	return s
}

func applyGlobalValues(cfg *Config) error {
	p := resolvePath("", cfg.ConfigPath)
	if p == "" {
		logger.Log(logger.Debug, "no config path set")
		return nil
	}
	logger.Log(logger.Debug, "trying config at %s", p)
	if err := LoadConfig(p, cfg); err != nil {
		if p == resolvePath(cfg.RootPath, defaultConfig.ConfigPath) {
			logger.Log(logger.Debug, "default config %s not found, using default config", p)
			return nil
		}
		return fmt.Errorf("config path %s could not be found ", p)
	}
	logger.Log(logger.Warning, "applied config from %s", p)
	return nil
}

func LoadConfig(name string, cfg *Config) error {
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

func NewConfig(path string) Config {
	cfg := &Config{}
	*cfg = defaultConfig
	if path != "" {
		cfg.ConfigPath = path
	}
	if err := applyGlobalValues(cfg); err != nil {
		logger.Log(logger.Error, "Config Error: %v\n", err)
	}
	applyENVValues(cfg)
	resolvePaths(cfg)
	return *cfg
}

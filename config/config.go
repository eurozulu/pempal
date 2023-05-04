package config

import (
	"os"
	"path/filepath"
)

const (
	ENV_CA_ROOT             = "CA_ROOT"
	ENV_CA_ROOT_CERTIFICATE = "CA_ROOT_CERT"
	ENV_CA_CERTS            = "CA_CERTS"
	ENV_CA_KEYS             = "CA_KEYS"
	ENV_CA_REQUESTS         = "CA_REQUESTS"
	ENV_CA_REVOKELISTS      = "CA_REVOKELISTS"
	ENV_CA_TEMPLATES        = "CA_TEMPLATES"
)

type Config struct {
	RootPath        string `yaml:"root-path"`
	RootCertificate string `yaml:"root-certificate,omitempty"`
	CertPath        string `yaml:"cert-path"`
	KeyPath         string `yaml:"key-path"`
	CsrPath         string `yaml:"csr-path"`
	CrlPath         string `yaml:"crl-path"`
	TemplatePath    string `yaml:"template-path"`
}

var defaultConfig = Config{
	RootPath:        "$PWD",
	RootCertificate: "root-certificate.pem",
	CertPath:        "certs",
	KeyPath:         "private",
	CsrPath:         "requests",
	CrlPath:         "revoked",
	TemplatePath:    "templates",
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
	if filepath.IsLocal(path) {
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
}

func envOrDefault(name string, def string) string {
	s, ok := os.LookupEnv(name)
	if !ok {
		return def
	}
	return s
}

func NewConfig() Config {
	cfg := &Config{}
	*cfg = defaultConfig
	applyENVValues(cfg)
	resolvePaths(cfg)
	return *cfg
}

package config

import (
	"github.com/go-yaml/yaml"
	"os"
	"path/filepath"
	"pempal/logger"
	"pempal/utils"
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
	RootPath:        "",
	RootCertificate: "root-certificate.pem",
	CertPath:        "",
	KeyPath:         "",
	CsrPath:         "",
	CrlPath:         "",
	TemplatePath:    "",
	ConfigPath:      "$HOME/.pempal/.config",
}

func resolvePaths(cfg *Config) {
	cfg.RootPath = resolvePath("$PWD", cfg.RootPath)
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
		path = filepath.Join(os.ExpandEnv(base), path)
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

func applyConfigFileValues(cfg *Config) error {
	logger.Log(logger.Debug, "trying config at %s", cfg.ConfigPath)
	if err := LoadConfig(cfg.ConfigPath, cfg); err != nil {
		return err
	}
	logger.Log(logger.Warning, "applied config from %s", cfg.ConfigPath)
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

func resolveConfigPath(path string) string {
	p := resolvePath("${PWD}", path)
	logger.Log(logger.Debug, "trying config at %s", p)
	if utils.FileExists(p) {
		logger.Log(logger.Debug, "found config at %s", p)
		return p
	}
	logger.Log(logger.Debug, "no config found at %s", p)
	return ""
}

func NewConfig(path string) Config {
	cfg := &Config{}
	*cfg = defaultConfig
	if path != "" {
		cfg.ConfigPath = path
	}
	cfg.ConfigPath = resolveConfigPath(cfg.ConfigPath)
	if cfg.ConfigPath != "" {
		if err := applyConfigFileValues(cfg); err != nil {
			logger.Log(logger.Error, "Config Error: %v\n", err)
		}
	}
	applyENVValues(cfg)
	resolvePaths(cfg)
	return *cfg
}

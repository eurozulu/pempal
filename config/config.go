package config

import (
	"fmt"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/utils"
	"github.com/go-yaml/yaml"
	"io"
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

var ConfigNames = [...]string{
	"root-path",
	"root-certificate",
	"cert-path",
	"key-path",
	"csr-path",
	"crl-path",
	"template-path",
}

type defaultConfig struct {
	RootPath            string `yaml:"root-path"`
	RootCertificatePath string `yaml:"root-certificate,omitempty"`
	CertPath            string `yaml:"cert-path,omitempty"`
	KeyPath             string `yaml:"key-path,omitempty"`
	CsrPath             string `yaml:"csr-path,omitempty"`
	CrlPath             string `yaml:"crl-path,omitempty"`
	TemplatePath        string `yaml:"template-path,omitempty"`
	configLocation      string `yaml:"-"`
}

func (cfg defaultConfig) Root() string {
	return cfg.RootPath
}

func (cfg defaultConfig) RootCertificate() string {
	return cfg.resolvePath(cfg.RootCertificatePath)
}

func (cfg defaultConfig) Certificates() string {
	return cfg.resolvePath(cfg.CertPath)
}

func (cfg defaultConfig) Keys() string {
	return cfg.resolvePath(cfg.KeyPath)
}

func (cfg defaultConfig) Requests() string {
	return cfg.resolvePath(cfg.CsrPath)
}

func (cfg defaultConfig) Revokations() string {
	return cfg.resolvePath(cfg.CrlPath)
}

func (cfg defaultConfig) Templates() string {
	return cfg.resolvePath(cfg.TemplatePath)
}

func (cfg defaultConfig) ConfigLocation() string {
	return cfg.configLocation
}

func (cfg defaultConfig) resolvePath(path string) string {
	path = os.ExpandEnv(path)
	if path == "" || filepath.IsLocal(path) {
		path = filepath.Join(os.ExpandEnv(cfg.RootPath), path)
	}
	return path
}

func applyENVValues(cfg *defaultConfig) {
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

func loadConfig(path string) (*defaultConfig, error) {
	cfg := &defaultConfig{}
	if path != "" {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer func(out io.WriteCloser) {
			if err := out.Close(); err != nil {
				logger.Error("Failed to close %s  %v", path, err)
			}
		}(f)
		if err = yaml.NewDecoder(f).Decode(cfg); err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
	}
	return cfg, nil
}

func SaveConfig(cfg Config) error {
	if cfg.ConfigLocation() == "" {
		return fmt.Errorf("configuration has no location path set")
	}
	f, err := os.OpenFile(cfg.ConfigLocation(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0640)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewEncoder(f).Encode(&cfg)
}

func resolveConfigPath(path string) (string, error) {
	// If no path given, check current directory for .config
	if path == "" {
		p := filepath.Join(os.ExpandEnv("$PWD"), ".config")
		if utils.FileExists(p) {
			path = p
			logger.Debug("using config file found in current directory: %s", path)
		}
	}
	// not in current dir, check if ENV var is set
	if path == "" {
		envPath, ok := os.LookupEnv(ENV_CA_CONFIG)
		if ok {
			if !utils.FileExists(os.ExpandEnv(envPath)) {
				logger.Warning("config path %s in %s could not be found", envPath, ENV_CA_CONFIG)
			} else {
				path = envPath
				logger.Debug("using config file found in environment var %s: %s", ENV_CA_CONFIG, path)
			}
		}
	}
	// not in ENV, check default location
	if path == "" {
		p := os.ExpandEnv(defaultConfigPath)
		if utils.FileExists(p) {
			path = p
			logger.Debug("using config file found in default path : %s", path)
		}
	}

	if path != "" && !utils.FileExists(path) {
		return "", fmt.Errorf("config file %s could not be found", path)
	}
	return path, nil
}

func NewConfig(path string) (Config, error) {
	cfgPath, err := resolveConfigPath(path)
	if err != nil {
		return nil, err
	}
	cfg, err := loadConfig(cfgPath)
	if err != nil {
		return nil, err
	}

	applyENVValues(cfg)
	if cfg.RootPath == "" || cfg.RootPath == "." {
		// default to PWD as root
		root := "$PWD"
		// if a loaded config file, use containing dir of that file as root
		if cfgPath != "" {
			root = filepath.Dir(cfgPath)
		}
		cfg.RootPath = os.ExpandEnv(root)
	}
	return cfg, nil
}

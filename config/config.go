package config

import (
	"fmt"
	"github.com/eurozulu/pempal/logging"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/utils"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"time"
)

var Config = newConfig()

var DurationYear = time.Hour * 24 * 365
var DurationMonth = DurationYear / 12

const configFileName = ".config"

const (
	ENVPPPATH = "PP_PATH"
	ENVPPROOT = "PP_ROOT"
)

type Configuration struct {
	CertPath     []string `yaml:"cert-path,omitempty"`
	Root         string   `yaml:"root,omitempty"`
	Keys         string   `yaml:"keys,omitempty"`
	Certificates string   `yaml:"certificates,omitempty"`
	Requests     string   `yaml:"requests,omitempty"`
	Revoked      string   `yaml:"revoked,omitempty"`
	Templates    []string `yaml:"templates,omitempty"`
	Archive      string   `yaml:"archive,omitempty"`
}

func (cfg Configuration) ResourcePath(rt model.ResourceType) (string, error) {
	switch rt {
	case model.PublicKey:
		return cfg.Keys, nil
	case model.PrivateKey:
		return cfg.Keys, nil
	case model.Certificate:
		return cfg.Certificates, nil
	case model.CertificateRequest:
		return cfg.Requests, nil
	case model.RevokationList:
		return cfg.Revoked, nil
	default:
		return "", fmt.Errorf("unknown resource type: %s", rt)
	}
}

func defaultConfiguration() *Configuration {
	return &Configuration{
		CertPath:     []string{"."},
		Root:         ".",
		Keys:         "./private",
		Certificates: "./certs",
		Requests:     "./csrs",
		Revoked:      "./crls",
		Templates:    []string{"./templates"},
		Archive:      "./archive",
	}
}

func applyFileConfiguration(path string, config *Configuration) error {
	if !utils.FileExists(path) {
		logging.Debug("config", "no config file found at %s", path)
		return nil
	}
	logging.Debug("config", "found config file at %s", path)
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return err
	}
	return yaml.NewDecoder(f).Decode(config)
}

func applyEnv(config *Configuration) {
	if s, ok := os.LookupEnv(ENVPPPATH); ok {
		logging.Debug("config", "found env %s setting of %q", ENVPPPATH, s)
		config.CertPath = filepath.SplitList(s)
	}
	if s, ok := os.LookupEnv(ENVPPROOT); ok {
		logging.Debug("config", "found env %s setting of %q", ENVPPROOT, s)
		config.Root = s
	}
}

func resolveConfigPaths(cfg *Configuration) error {
	cfg.CertPath = resolvePaths("", cfg.CertPath)
	cfg.Root = resolvePath("", cfg.Root)

	cfg.Keys = resolvePath(cfg.Root, cfg.Keys)
	cfg.Certificates = resolvePath(cfg.Root, cfg.Certificates)
	cfg.Requests = resolvePath(cfg.Root, cfg.Requests)
	cfg.Revoked = resolvePath(cfg.Root, cfg.Revoked)
	cfg.Templates = resolvePaths(cfg.Root, cfg.Templates)
	cfg.Archive = resolvePath(cfg.Root, cfg.Archive)
	return nil
}

func resolvePath(root, path string) string {
	path = os.ExpandEnv(path)
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, path)
}

func resolvePaths(root string, paths []string) []string {
	var resolved []string
	for _, path := range paths {
		path = filepath.Join(root, path)
		path = os.ExpandEnv(path)
		if !filepath.IsAbs(path) {
			p, err := filepath.Abs(path)
			if err != nil {
				logging.Warning("config", "failed to resolve path %s: %s", path, err)
				continue
			}
			resolved = append(resolved, p)
		}
	}
	return resolved
}

func GlobalFilePath() string {
	return filepath.Join(os.ExpandEnv("$HOME"), configFileName)
}

func (cfg Configuration) LocalFilePath() string {
	return filepath.Join(cfg.Root, configFileName)
}

func newConfig() *Configuration {
	cfg := defaultConfiguration()
	applyEnv(cfg)

	// apply 'global' config from home

	if err := applyFileConfiguration(GlobalFilePath(), cfg); err != nil {
		logging.Error("Configuration", "Failed to read global config file.  %s", err)
	}
	if err := applyFileConfiguration(cfg.LocalFilePath(), cfg); err != nil {
		logging.Error("Configuration", "Failed to read local config file.  %s", err)
	}

	resolveConfigPaths(cfg)
	return cfg
}

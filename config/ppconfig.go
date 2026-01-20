package config

import (
	"fmt"
	"github.com/eurozulu/pempal/tools"
	"gopkg.in/yaml.v2"

	"log"
	"os"
	"path/filepath"
	"strings"
)

const PP_CONFIG_NAME = ".ppconfig"

const ENV_PP_ROOT_PATH = "PP_ROOT"
const ENV_PP_SEARCH_PATH = "PP_PATH"

var DefaultPPConfig = NewPPConfig()
var rootPath string = "."
var searchPath string = "."

type PPConfig struct {
	KeyPath            string   `yaml:"key-path"`
	CertPath           string   `yaml:"cert-path"`
	CSRPath            string   `yaml:"csr-path"`
	CRLPath            string   `yaml:"crl-path"`
	TemplatePath       string   `yaml:"template-path,omitempty"`
	DefaultKeyTemplate string   `yaml:"default-key-template,omitempty"`
	FileExt            []string `yaml:"file-extensions"`
}

func init() {
	var err error
	if err = setEnvironmentVars(); err != nil {
		log.Fatalf("Error setting environment variables: %v", err)
	}
	// ensure path entries are unique and root path is first in search path
	searchPath = strings.Join(tools.AppendUnique([]string{RootPath()},
		strings.Split(searchPath, string(os.PathListSeparator))...), string(os.PathListSeparator))

	if err = LoadConfig(ConfigPath(), &DefaultPPConfig); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
}

func setEnvironmentVars() error {
	if s, ok := os.LookupEnv(ENV_PP_ROOT_PATH); ok {
		s = os.ExpandEnv(s)
		if !tools.IsDirExists(s) {
			return fmt.Errorf("Root path %s, found in %s could not be found", s, ENV_PP_ROOT_PATH)
		}
		rootPath = s
	}
	if s, ok := os.LookupEnv(ENV_PP_SEARCH_PATH); ok {
		searchPath = s
	}
	return nil
}

func configAsMap(cfg *PPConfig) (map[string]interface{}, error) {
	m := map[string]interface{}{}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func NewPPConfig() *PPConfig {
	return &PPConfig{
		KeyPath:            "./private",
		CertPath:           "./certs",
		CSRPath:            "./requests",
		CRLPath:            "./revoked",
		TemplatePath:       "./templates",
		DefaultKeyTemplate: "key",
		FileExt: []string{
			"", ".pem",
			".crt", ".cert", ".cer",
			".key", ".pub", ".prk", ".puk", ".rsa",
			".x509",
			".csr", ".request",
			".crl", ".revoke",
		},
	}
}

func LoadConfig(path string, cfg interface{}) error {
	if path == "" {
		return nil
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = yaml.NewDecoder(f).Decode(cfg); err != nil {
		return fmt.Errorf("failed to load config file %s: %s", path, err)
	}
	return nil
}

func SaveConfig(cfg interface{}) error {
	current, err := configAsMap(DefaultPPConfig)
	if err != nil {
		return err
	}
	def, err := configAsMap(NewPPConfig())
	if err != nil {
		return err
	}
	for k, v := range current {
		if v == def[k] {
			delete(def, k)
			continue
		}
		def[k] = v
	}
	data, err := yaml.Marshal(def)
	if err != nil {
		return err
	}
	path := filepath.Join(RootPath(), PP_CONFIG_NAME)
	return os.WriteFile(path, data, 0644)
}

func ConfigPath() string {
	path := filepath.Join(RootPath(), PP_CONFIG_NAME)
	if tools.IsFileExists(path) {
		return path
	}
	return ""
}

func RootPath() string {
	return rootPath
}
func SearchPath() string {
	return searchPath
}
func TemplatePath() string {
	return filepath.Join(RootPath(), DefaultPPConfig.TemplatePath)
}
func KeyPath() string {
	return filepath.Join(RootPath(), DefaultPPConfig.KeyPath)
}
func CertificatePath() string {
	return filepath.Join(RootPath(), DefaultPPConfig.CertPath)
}
func CSRPath() string {
	return filepath.Join(RootPath(), DefaultPPConfig.CSRPath)
}
func CRLPath() string {
	return filepath.Join(RootPath(), DefaultPPConfig.CRLPath)
}
func FileExtensions() []string {
	return DefaultPPConfig.FileExt
}
func DefaultKeyTemplateName() string {
	return DefaultPPConfig.DefaultKeyTemplate
}

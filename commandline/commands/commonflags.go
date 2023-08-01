package commands

import (
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/utils"
	"os"
	"path/filepath"
	"strings"
)

const (
	ENV_PP_HOME      = "PP_HOME"
	ENV_PP_CERTS     = "PP_CERT_PATH"
	ENV_PP_KEYS      = "PP_KEY_PATH"
	ENV_PP_CSRS      = "PP_CSR_PATH"
	ENV_PP_CRLS      = "PP_CRL_PATH"
	ENV_PP_TEMPLATES = "PP_TEMPLATES_PATH"
)

const defaultHomePath = "$HOME/.pempal"
const defaultCertPath = "$PWD:" + defaultHomePath + "/certs:/etc/ssl/certs"
const defaultKeyPath = "$PWD:" + defaultHomePath + "/private:$HOME/.ssh"
const defaultTemplatePath = defaultHomePath + "/templates"

var CommonFlags = &DefaultCommonFlags{}

type DefaultCommonFlags struct {
	CertPath     string `yaml:"cert-path"`
	KeyPath      string `yaml:"key-path"`
	CsrPath      string `yaml:"csr-path,omitempty"`
	CrlPath      string `yaml:"crl-path,omitempty"`
	TemplatePath string `yaml:"template-path"`
	HomePath     string `yaml:"home-path,omitempty"`
	Quiet        bool   `yaml:"q"`
	Verbose      bool   `yaml:"v"`
	Debug        bool   `yaml:"debug"`
	Help         bool   `yaml:"help"`
	Output       string `yaml:"out"`
	ForceOut     bool   `yaml:"force"`
}

func init() {
	CommonFlags.HomePath = envOrDefault(ENV_PP_HOME, defaultHomePath)
	CommonFlags.CertPath = envOrDefault(ENV_PP_CERTS, defaultCertPath)
	CommonFlags.KeyPath = envOrDefault(ENV_PP_KEYS, defaultKeyPath)
	CommonFlags.CsrPath = envOrDefault(ENV_PP_CSRS, "")
	CommonFlags.CrlPath = envOrDefault(ENV_PP_CRLS, "")
	CommonFlags.TemplatePath = envOrDefault(ENV_PP_TEMPLATES, defaultTemplatePath)
}

func ResolvePath(p string) []string {
	if p == "" {
		p = "$PWD"
	}
	var found []string
	for _, path := range strings.Split(p, ":") {
		path = os.ExpandEnv(path)
		if filepath.IsLocal(path) {
			path = filepath.Join(os.ExpandEnv("$PWD"), path)
		}
		if !utils.DirectoryExists(path) && !utils.FileExists(path) {
			logger.Debug("ignoring path entry %s as could not be found", path)
			continue
		}
		found = append(found, path)
	}
	return found
}

func envOrDefault(name string, def string) string {
	s, ok := os.LookupEnv(name)
	if !ok {
		return def
	}
	return s
}

package config

import (
	"github.com/eurozulu/pempal/keys"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/templates"
	"strings"
)

var ConfigPath string

func CurrentConfig() (Config, error) {
	return NewConfig(ConfigPath)
}

func TemplateStore() (templates.TemplateStore, error) {
	cfg, err := CurrentConfig()
	if err != nil {
		return nil, err
	}
	return templates.NewTemplateStore(cfg.Templates(), model.DefaultResourceTemplates)
}

func KeyManager() (keys.Keys, error) {
	cfg, err := CurrentConfig()
	if err != nil {
		return nil, err
	}
	return keys.NewKeys(strings.Split(cfg.Keys(), ":"), strings.Split(cfg.Certificates(), ":")), nil
}

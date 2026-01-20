package config

import (
	"fmt"
	"os"
	"testing"
)

func TestPPConfigRootPath(t *testing.T) {
	if err := setEnv(ENV_PP_ROOT_PATH, ""); err != nil {
		t.Errorf("Error setting environment variables: %v", err)
	}
	expected := "."
	if RootPath() != expected {
		t.Errorf("RootPath() = %s, expected %s", RootPath(), expected)
	}

	expected = os.ExpandEnv("${HOME}")
	if err := setEnv(ENV_PP_ROOT_PATH, "$HOME"); err != nil {
		t.Errorf("Error setting environment variables: %v", err)
	}
	if RootPath() != expected {
		t.Errorf("RootPath() = %s, expected %s", RootPath(), expected)
	}

	os.Setenv(ENV_PP_ROOT_PATH, "/doesntexist")
	if err := setEnvironmentVars(); err == nil {
		t.Errorf("Expected error setting environment variable %s to non existing directory: %v", ENV_PP_ROOT_PATH, err)
	}
}

func TestConfigSearchPath(t *testing.T) {
	if err := setSearchPathEnv(""); err != nil {
		t.Errorf("Error setting search path: %v", err)
	}
	expectPath := os.ExpandEnv("${PWD}")
	if SearchPath() != expectPath {
		t.Errorf("SearchPath() = %s, expected %s", SearchPath(), expectPath)
	}

	if err := setSearchPathEnv("${HOME}/.ssh"); err != nil {
		t.Errorf("Error setting search path: %v", err)
	}
	expectPath = os.ExpandEnv("${PWD}:${HOME}/.ssh")
	if SearchPath() != expectPath {
		t.Errorf("SearchPath() = %s, expected %s", SearchPath(), expectPath)
	}

}

func setEnv(name, value string) error {
	if value == "" {
		os.Unsetenv(name)
	} else {
		os.Setenv(name, value)
	}
	if err := setEnvironmentVars(); err != nil {
		return fmt.Errorf("Error setting environment variables: %v", err)
	}
	return nil
}

func setSearchPathEnv(path string) error {
	if err := setEnv(ENV_PP_SEARCH_PATH, path); err != nil {
		return err
	}
	if err := cleanSearchPath(); err != nil {
		return fmt.Errorf("Error cleaning search path: %s", err)
	}
	return nil
}

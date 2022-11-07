package filepaths

import (
	"os"
	"strings"
)

const envKeyPath = "PP_KEY_PATH"
const envCertPath = "PP_CERTIFICATE_PATH"

const defaultKeyPath = "/etc/ssl/private"
const defaultCertPath = "/etc/ssl/certs"

func KeyPath() string {
	// read ENV
	kp := os.ExpandEnv(os.Getenv(envKeyPath))
	if kp == "" {
		// nopthing in env, see if default ixists
		if DirExists(defaultKeyPath) {
			kp = defaultKeyPath
		}
	}
	// Always start with the working dir
	wd, err := os.Getwd()
	if err == nil {
		kp = strings.Join([]string{wd, kp}, ":")
	}
	return kp
}

func CertPath() string {
	// read ENV
	cp := os.ExpandEnv(os.Getenv(envCertPath))
	if cp == "" {
		// nopthing in env, see if default exists
		if DirExists(defaultCertPath) {
			cp = defaultCertPath
		}
	}
	// Always start with the working dir
	wd, err := os.Getwd()
	if err == nil {
		cp = strings.Join([]string{wd, cp}, ":")
	}
	return cp
}

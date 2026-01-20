package factories

import (
	"fmt"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/tools"
	"os"
	"path/filepath"
	"strings"
)

func PathForResource(rt model.ResourceType) string {
	switch rt {
	case model.ResourceTypeCertificate:
		return config.CertificatePath()
	case model.ResourceTypePrivateKey:
		return config.KeyPath()
	case model.ResourceTypeCertificateRequest:
		return config.CSRPath()
	case model.ResourceTypeRevokationList:
		return config.CRLPath()
	default:
		return config.RootPath()

	}
}

// PublicFingerPrint gets the Fingerprint of the given resource.
// If the resource is a private key, its public key fingerprint is returned
func PublicFingerPrint(res model.PemResource) model.Fingerprint {
	if res.ResourceType() == model.ResourceTypePrivateKey {
		return res.(*model.PrivateKey).Public().Fingerprint()
	}
	return res.Fingerprint()
}

func uniqueFileName(path string) string {
	if !tools.IsPathExists(path) {
		return path
	}
	var count int
	p := filepath.Dir(path)
	name := filepath.Base(path)
	ext := filepath.Ext(name)
	name = strings.TrimSuffix(name, ext)
	for tools.IsFileExists(path) {
		path = filepath.Join(filepath.Dir(p), fmt.Sprintf("%s_%3d%s", name, count, ext))
		count++
	}
	return path
}

func SaveResource(res ...model.PemResource) error {
	for _, r := range res {
		name := strings.Join([]string{PublicFingerPrint(r).String(), ".pem"}, "")
		path := PathForResource(r.ResourceType())
		perm := os.FileMode(0644)
		if r.ResourceType() == model.ResourceTypePrivateKey {
			if err := tools.EnsureSecurePath(path); err != nil {
				return err
			}
			perm = 0600
		} else {
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}
		}
		path = uniqueFileName(filepath.Join(path, name))
		data, err := r.MarshalText()
		if err != nil {
			return err
		}
		if err := os.WriteFile(path, data, perm); err != nil {
			return err
		}
	}
	return nil
}

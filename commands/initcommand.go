package commands

import (
	"fmt"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/logging"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/repositories"
	"github.com/eurozulu/pempal/tools"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

const rootCACert = "rootca"
const interCACert = "intermediateca"
const defaultNameName = "default-name"
const rootIssuerName = "root-issuer"
const defaultIssuerName = "default-issuer"

const userName = "usercertificate"

// InitCommand initalises a new PKI in the current pempal root path.
// A Root CA, Inter CA and User certificate of the current users name are created
// @Command(init)
type InitCommand struct {
}

// InitPKI requires an 'org' name to define the organisation
// @Action
func (cmd InitCommand) InitPKI(baseName string, args []string) error {
	if err := cmd.InitTemplates(baseName); err != nil {
		return err
	}

	if err := cmd.InitCertificates(); err != nil {
		return err
	}
	return nil
}

// InitTemplates initalises new name templates in the current pempal root/templates directory.
// 'default-name', 'root-issuer' and 'default-issuer' are all created containing their respective names
// based on the given base name.
// Given base name is designed as the bases for all other names, usually contaning minimal information such as Organisation.
// Base name must NOT contain a common-name.
// @Action(templates,t)
func (cmd InitCommand) InitTemplates(baseName string) error {
	// write default-name template
	name, err := parseBaseName(baseName)
	if err != nil {
		return err
	}
	if templateExists(defaultNameName) {
		logging.Info("Template %s already exists. Skipping.", defaultNameName)
	} else {
		if err := writeTemplate(defaultNameName, []byte("subject: "+name.String())); err != nil {
			return err
		}
		logging.Info("created Template %s with name %s", defaultNameName, name.String())
	}

	// write rootissuer template
	if templateExists(rootIssuerName) {
		logging.Info("Template %s already exists. Skipping.", rootIssuerName)
	} else {
		subject, err := readTemplateProperty(rootCACert, "subject")
		if err != nil {
			return err
		}
		rootDN, err := model.ParseDistinguishedName(subject)
		if err != nil {
			return err
		}
		rootDN.Merge(*name)
		if err := writeTemplate(rootIssuerName, []byte("issuer: "+rootDN.String())); err != nil {
			return err
		}
		logging.Info("created Template %s with name %s", rootIssuerName, rootDN.String())
	}

	// write default issuer template
	if templateExists(defaultIssuerName) {
		logging.Info("Template %s already exists. Skipping.", defaultIssuerName)
	} else {
		subject, err := readTemplateProperty(interCACert, "subject")
		if err != nil {
			return err
		}
		issuerDN, err := model.ParseDistinguishedName(subject)
		if err != nil {
			return err
		}
		issuerDN.Merge(*name)
		if err := writeTemplate(defaultIssuerName, []byte("issuer: "+issuerDN.String())); err != nil {
			return err
		}
		logging.Info("created Template %s with name %s", defaultIssuerName, issuerDN.String())
	}
	return nil
}

// InitCertificates generates the basic certificates for a new PKI.
// Using the existing name templates, it generates:
// Root CA - self signed
// Inter CA - signed by root
// User cert - signed by InterCA.
func (cmd InitCommand) InitCertificates() error {
	rootName, err := readTemplateProperty(rootIssuerName, "issuer")
	if err != nil {
		return err
	}
	if certificateExists(rootName) {
		logging.Info("Certificate %s already exists. Skipping.", rootName)
	} else {
		_, err := makeCertificate(rootCACert, rootName)
		if err != nil {
			return err
		}
		logging.Info("created root CA certificate %s")
	}

	issuerName, err := readTemplateProperty(defaultIssuerName, "issuer")
	if err != nil {
		return err
	}
	if certificateExists(issuerName) {
		logging.Info("Certificate %s already exists. Skipping.", issuerName)
	} else {
		_, err = makeCertificate(interCACert, issuerName)
		if err != nil {
			return err
		}
		logging.Info("created Intermediate CA certificate %s", issuerName)
	}
	uName, err := model.ParseDistinguishedName("CN=" + os.ExpandEnv("${USER}"))
	if err != nil {
		return err
	}
	if certificateExists(uName.String()) {
		logging.Info("Certificate %s already exists. Skipping.", uName.String())
	} else {
		_, err = makeCertificate(userName, uName.String())
		if err != nil {
			return err
		}
		logging.Info("created user certificate %s", uName.String())
	}
	return nil
}

func parseBaseName(name string) (*model.DistinguishedName, error) {
	if name == "" {
		return nil, fmt.Errorf("an organisation name must be provided")
	}
	if !strings.Contains(name, "=") {
		name = strings.Join([]string{"O=", name}, "")
	}
	dn, err := model.ParseDistinguishedName(name)
	if err != nil {
		return nil, err
	}
	if dn.CommonName != "" {
		return nil, fmt.Errorf("Base name can not contain a common name")
	}
	return dn, nil
}

func templateExists(name string) bool {
	tpath := filepath.Join(config.TemplatePath(), name)
	tpath = strings.Join([]string{tpath, "yml"}, ".")
	return tools.IsFileExists(tpath)
}

func writeTemplate(name string, data []byte) error {
	path := templatePathFromName(name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return nil
}

func readTemplateProperty(templateName, propertyName string) (string, error) {
	m, err := readTemplate(templateName)
	if err != nil {
		return "", err
	}
	return m[propertyName], nil
}

func readTemplate(templateName string) (map[string]string, error) {
	t, err := repositories.Templates(config.TemplatePath()).ByName(templateName)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	if err = yaml.Unmarshal([]byte(t[0].String()), &m); err != nil {
		return nil, err
	}
	return m, nil
}

func templatePathFromName(name string) string {
	filename := filepath.Join(config.TemplatePath(), name)
	return strings.Join([]string{filename, "yml"}, ".")
}

func certificateExists(name string) bool {
	c, err := readCertificate(name)
	if err != nil {
		return false
	}
	return c != nil
}
func readCertificate(name string) (*model.Certificate, error) {
	dn, err := model.ParseDistinguishedName(name)
	if err != nil {
		return nil, err
	}
	return repositories.Certificates(config.SearchPath()).ByName(*dn)
}

func makeCertificate(templateName string, subject string) (string, error) {
	args := []string{templateName, "-subject", subject}
	dn, err := MakeCommand{Save: true}.Create(args...)
	if err != nil {
		return "", err
	}
	return dn, nil
}

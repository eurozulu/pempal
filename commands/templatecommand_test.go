package commands

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/templates"
	"os"
	"testing"
)

func TestTemplateCommand_Exec(t *testing.T) {
	os.Setenv(config.ENVPKIROOTPATH, "../testing")
	buf := bytes.NewBuffer(nil)
	tc := &TemplateCommand{Output: buf}

	if err := tc.Exec("rootcakey"); err != nil {
		t.Fatal(err)
	}
	fmt.Println(buf.String())
	var keyTemplate templates.PrivateKeyTemplate
	if err := json.Unmarshal(buf.Bytes(), &keyTemplate); err != nil {
		t.Fatal(err)
	}
	if x509.PublicKeyAlgorithm(keyTemplate.PublicKeyAlgorithm) != x509.Ed25519 {
		t.Errorf("unexpected PublicKeyAlgorithm, expected %s, found %s", x509.Ed25519, keyTemplate.PublicKeyAlgorithm)
	}
	if keyTemplate.KeySize != 4096 {
		t.Errorf("unexpected keysize, expected %d, found %d", 4096, keyTemplate.KeySize)
	}

}

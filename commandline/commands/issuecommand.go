package commands

import (
	"fmt"
	"github.com/eurozulu/pempal/builders"
	"github.com/eurozulu/pempal/commandline/valueeditors"
	"github.com/eurozulu/pempal/identity"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/utils"
	"io"
	"strings"
)

var certificatePropertyEditors = []valueeditors.ValueEditor{
	valueeditors.NumberEditor{PropertyName: "version"},
	valueeditors.NumberEditor{PropertyName: "serial-number"},
	valueeditors.StringEditor{
		PropertyName:  "signature-algorithm",
		Choice:        utils.SignatureAlgorithmNames(),
		DefaultChoice: 4,
	},
	valueeditors.BoolEditor{
		PropertyName:  "is-ca",
		DefaultChoice: false,
	},
	valueeditors.DistinguishedNameEditor{
		PropertyName: "subject",
	},
	valueeditors.DistinguishedNameEditor{
		PropertyName: "issuer",
	},
}

// issueCommand manages x509 certificate
// when no arguments are given lists all the known certificates of the identity with known private identity
// when arguments are given, issue will attempt to generate a new certificate.
// To succesfully create a certificate the command requires these properties:
// serial-number: non zero value
// subject.common-name: Non empty string
// public-key: A key ID, name or PEM
// issuer: A unique DN of a known issuer certificate. (Self signed must contain the certificate Subject)
//
// These may be provided as flags or within named templates.
// as flags, the values should be quoted strings.
// e.g. issue -public-key "identity/rootkey.pem" -subject.common-name "My Root Certificate"

// A template is simply a pre-formed representation of the same thing stored under a unique name.
// myroottemplate=`
//  public-key: "identity/rootkey.pem"
//  subject.common-name: "My Root Certificate"
// e.g. issue myroottemplate

// flags and template names may be combined. Flag values will always take precedences and be applied as the last values in the chain.
// e.g. issue myroottemplate -subject.common-name: "Some other name"
// will issue a new certificate, using the same key ("identity/rootkey.pem") but with a new name "Some other name"
//
// An argument may be a named template or a named resource
// When an argument points to a known resource, by specifying a partial file path which uniqely identifies that resource,
// the resource is transformed into a template containing the values of that resource.

type issueCommand struct {
	flagValues resources.CertificateDTO
}

func (cmd issueCommand) Execute(args []string, out io.Writer) error {
	keyz := identity.NewIssuers(strings.Split(CommonFlags.KeyPath, ":"), strings.Split(CommonFlags.CertPath, ":"))
	builder, err := builders.NewSigningBuilder(resources.Certificate, keyz)

	// check if any named templates or file resources given in args
	temps, err := argumentsToTemplates(args)
	if err != nil {
		return err
	}
	builder.AddTemplate(temps...)

	// get flag properties as template
	t, err := resources.DTOToTemplate(&cmd.flagValues)
	if err != nil {
		return err
	}
	builder.AddTemplate(t)

	if err = confirmBuild("Issue new certificate", certificatePropertyEditors, builder); err != nil {
		return err
	}

	cr, err := builder.Build()
	if err != nil {
		return err
	}

	// flip new cert resource into a PEM string
	dto, err := resources.NewResourceDTO(cr)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, dto.String())
	return err
}

func (cmd issueCommand) MarshalYAML() (interface{}, error) {
	return &cmd.flagValues, nil
}

func (cmd issueCommand) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(&cmd.flagValues)
}

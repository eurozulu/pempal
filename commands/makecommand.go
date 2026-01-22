package commands

import (
	"bytes"
	"fmt"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/factories"
	"github.com/eurozulu/pempal/repositories"
	"github.com/eurozulu/pempal/templates"
	"github.com/eurozulu/pempal/tools"
	"strings"
)

// MakeCommand creates new resources using the given template names.
// Templates are merged into one ao the resulting template is a base template.
// The relevant resource for that base template is then generated.
// @Command(make)
type MakeCommand struct {
	// Key flag is optional, when given should be the fingerprint (or unique oartial fingerprint) of the key to use to identify the new object.
	// Only applies to Certificates and CSRs, public-key property.  Ignored for other resources, (keys and CRLs).
	// When not set, a new key is generated using the 'key' template.
	// @Flag(key)
	Key string

	// Save when set will save the new resource into the PKI repository.
	// @Flag(save, s)
	Save bool
}

// Create generates a new resource witht he resulting template from merging the given named templates
// requires one or more known template names, which are merged into a single base template.
// returns either the PEM encoded resource or, when Save is set, the fingerprint of the new resource
// @Action
func (cmd MakeCommand) Create(args ...string) (string, error) {
	argFlags, argz, err := ArgFlagsToTemplate(args)
	if err != nil {
		return "", err
	}
	temps, err := templateRepo.ExpandedByName(argz...)
	if err != nil {
		return "", err
	}

	if argFlags.String() != "" {
		temps = append(temps, argFlags)
	}
	if cmd.Key != "" {
		puk, err := resolveKey(cmd.Key)
		if err != nil {
			return "", err
		}
		t, _, err := ArgFlagsToTemplate([]string{"-public-key", tools.SingleQuote(string(puk))})
		if err != nil {
			return "", err
		}
		temps = append(temps, t)
	}
	t, err := templates.MergeTemplates(temps)
	if err != nil {
		return "", err
	}
	resz, err := factories.Make(t)
	if err != nil {
		return "", err
	}

	if cmd.Save {
		if err := factories.SaveResource(resz...); err != nil {
			return "", err
		}
		return strings.Join(tools.StringerToString(resz...), "\n"), nil
	}
	buf := bytes.NewBuffer(nil)
	for _, res := range resz {
		data, err := res.MarshalText()
		if err != nil {
			return "", err
		}
		buf.Write(data)
		buf.WriteString("\n")
	}
	return buf.String(), nil
}

func resolveKey(fingerprint string) ([]byte, error) {
	keyz := repositories.Keys(config.KeyPath()).MatchByAnyFingerPrint(fingerprint)
	if len(keyz) == 0 {
		return nil, fmt.Errorf("%q key not found", fingerprint)
	}
	if len(keyz) > 1 {
		fpz := strings.Join(tools.StringerToString(keyz...), ", ")
		return nil, fmt.Errorf("%q key matches multiple keys: %s", fingerprint, fpz)
	}
	return keyz[0].Public().MarshalText()
}

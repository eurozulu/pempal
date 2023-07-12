package valueeditors

import (
	"github.com/eurozulu/pempal/commandline/ui"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/utils"
)

var dnEditors = []ValueEditor{
	&StringEditor{PropertyName: "common-name"},
	&StringEditor{PropertyName: "serial-number"},
	&StringSliceEditor{PropertyName: "organization"},
	&StringSliceEditor{PropertyName: "organizational-unit"},
	&StringSliceEditor{PropertyName: "locality"},
	&StringSliceEditor{PropertyName: "street-address"},
	&StringSliceEditor{PropertyName: "province"},
	&StringSliceEditor{PropertyName: "postal-code"},
	&StringSliceEditor{PropertyName: "country"},
}

type DistinguishedNameEditor struct {
	// PropertyName of the name of the property to edit
	PropertyName string
	ReadOnly     bool
}

func (de DistinguishedNameEditor) Name() string {
	return de.PropertyName
}

func (de DistinguishedNameEditor) Edit(offset ui.ViewOffset, value string) (string, error) {
	offset.XOffset += len(de.PropertyName) // move to the right of title name
	var edits []ValueEditor
	if !de.ReadOnly { //TODO: Fix with a view only Editor
		edits = dnEditors
	}

	dnvalues, err := parseRDNSToMap(value)
	if err != nil {
		return "", err
	}

	form := EditorList(edits)
	result, err := form.Show(offset, dnvalues)
	if err != nil {
		return "", err
	}
	return parseMapToRDNSString(utils.MergeMap(dnvalues, result))
}

func parseRDNSToMap(rdns string) (map[string]string, error) {
	// convert value into template via dnDTO
	if rdns == "" {
		return nil, nil
	}
	// Unmarshall RDNSequence string into DN-dto
	dto := &resources.DistinguishedNameDTO{}
	if err := dto.UnmarshalBinary([]byte(rdns)); err != nil {
		return nil, err
	}
	return resources.DTOToTemplate(dto)
}

func parseMapToRDNSString(m map[string]string) (string, error) {
	dto := &resources.DistinguishedNameDTO{}
	if err := resources.ApplyTemplateToDTO(dto, m); err != nil {
		return "", err
	}
	return dto.ToName().String(), nil
}

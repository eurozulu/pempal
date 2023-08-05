package prompts

import (
	"fmt"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/ui"
)

type DNView struct {
	ui.ListView
}

func buildDNChildViews() []ui.View {
	return []ui.View{
		ui.NewTextView("common-name", ""),
		ui.NewTextView("serial-number", ""),
		ui.NewTextView("organization", ""),
		ui.NewTextView("organizational-unit", ""),
		ui.NewTextView("locality", ""),
		ui.NewTextView("street-address", ""),
		ui.NewTextView("province", ""),
		ui.NewTextView("postal-code", ""),
		ui.NewTextView("country", ""),
	}
}

func (dnv *DNView) OnViewOpen() {
	dnv.setChildValues(dnv.GetText())
	dnv.SetSelectedIndexByText("") // select first empty view
}

func (dnv *DNView) OnViewClose(child ui.View) ui.View {
	dnv.SetText(dnv.getChildValues())
	return dnv
}

func (dnv *DNView) setChildValues(rdns string) {
	values, err := parseRDNSToMap(rdns)
	if err != nil {
		dnv.SetText(err.Error())
	}
	if values == nil {
		values = map[string]string{}
	}
	for _, cv := range dnv.ChildViews() {
		if tv, ok := cv.(ui.TextView); ok {
			tv.SetText(values[tv.Label()])
		}
	}
}

func (dnv *DNView) getChildValues() string {
	s, err := formatMapToRDNSString(dnv.getChildValuesMap())
	if err != nil {
		return err.Error()
	}
	return s
}

func (dnv *DNView) getChildValuesMap() map[string]string {
	m := map[string]string{}
	for _, cv := range dnv.ChildViews() {
		s := cv.String()
		if tv, ok := cv.(ui.TextView); ok {
			s = tv.GetText()
		}
		m[cv.Label()] = s
	}
	return m
}

func parseRDNSToMap(rdns string) (map[string]string, error) {
	// convert value into template via dnDTO
	if rdns == "" {
		return nil, fmt.Errorf("missing common-name")
	}
	// Unmarshall RDNSequence string into DN-dto
	dto := &resources.DistinguishedNameDTO{}
	if err := dto.UnmarshalBinary([]byte(rdns)); err != nil {
		return nil, err
	}
	return resources.DTOToTemplate(dto, false)
}

func formatMapToRDNSString(m map[string]string) (string, error) {
	dto := &resources.DistinguishedNameDTO{}
	if err := resources.ApplyTemplateToDTO(dto, m); err != nil {
		return "", err
	}
	return dto.ToName().String(), nil
}

func NewDNView(label, rdn string) ui.ParentView {
	return &DNView{*ui.NewListView(label, rdn, buildDNChildViews()...)}
}

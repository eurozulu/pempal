package prompts

import (
	"context"
	"github.com/eurozulu/pempal/identity"
	"github.com/eurozulu/pempal/ui"
)

const defaultLabelWidth = 36

type publicKeyView struct {
	ui.ListView
	keys identity.Keys
}

func (bv *publicKeyView) String() string {
	text := bv.GetText()
	if identity.IsIdentity(text) {
		return identity.Identity(text).String()
	}
	return text
}

func (bv *publicKeyView) OnViewOpen() {
	bv.buildKeyChoice()
	bv.LabelWidth = defaultLabelWidth
	if id := bv.textAsIdentity(); id != "" {
		bv.setSelectedIndexByLabel(id.String())
	}
}

func (bv *publicKeyView) OnViewClose(child ui.View) ui.View {
	selected := bv.SelectedIndex()
	if selected < 0 {
		return child
	}
	sv := bv.ChildViews()[selected]
	if sv.Label() == createKeyLabel {

	}
	bv.setTextFromId(bv.ChildViews()[selected].Label())
	return child
}

func (bv publicKeyView) textAsIdentity() identity.Identity {
	if identity.IsIdentity(bv.GetText()) {
		return identity.Identity(bv.GetText())
	}
	return ""
}

func (bv *publicKeyView) setTextFromId(id string) {
	k, err := bv.keys.KeyByIdentity(id)
	if err != nil {
		bv.SetText(err.Error())
		bv.SetColours(ui.ErrorColour)
		return
	}
	bv.SetText(k.String())
}

const createKeyLabel = "Create New Key"

func (bv *publicKeyView) buildKeyChoice() {
	kez := bv.listAllKeys()
	childs := make([]ui.View, len(kez)+1)
	for i, k := range kez {
		childs[i] = ui.NewLabelView(k.Identity().String(), k.Location())
	}
	childs[len(kez)] = ui.NewLabelView(createKeyLabel, "")
	ui.MutableParentView(&bv.ListView).SetChildViews(childs)
}

func (bv *publicKeyView) listAllKeys() []identity.Key {
	if bv.keys == nil {
		return nil
	}
	var found []identity.Key
	for key := range bv.keys.AllKeys(context.Background()) {
		found = append(found, key)
	}
	return found
}

func (bv *publicKeyView) setSelectedIndexByLabel(label string) {
	selected := -1
	for i, c := range bv.ChildViews() {
		if c.Label() == label {
			selected = i
			break
		}
	}
	bv.SetSelectedIndex(selected)
}

func NewPublicKeyView(label, text string, keys identity.Keys) ui.ParentView {
	return &publicKeyView{
		ListView: *ui.NewListViewStrings(label, "Create New Key", "Select Existing Key"),
		keys:     keys,
	}
}

package prompts

import "github.com/eurozulu/pempal/ui"

type publicKeyView struct {
	ui.BoolView
}

func NewPublicKeyView(label, text string) ui.View {
	return ui.NewBoolView(label, "Create New Key", "Select Existing Key")
}

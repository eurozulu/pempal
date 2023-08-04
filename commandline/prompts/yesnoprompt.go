package prompts

import (
	"github.com/eurozulu/pempal/ui"
	"strconv"
)

type YesNoPrompt struct {
	ui.ListView
}

func NewYesNoPrompt(label string, value bool) *YesNoPrompt {
	tv := ui.NewListViewStrings(label, strconv.FormatBool(value), "No", "Yes")
	return &YesNoPrompt{ListView: *tv}
}

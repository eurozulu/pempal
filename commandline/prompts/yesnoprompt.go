package prompts

import (
	"github.com/eurozulu/pempal/ui"
	"strconv"
)

type YesNoPrompt struct {
	ui.TextList
}

func NewYesNoPrompt(label string, value bool) *YesNoPrompt {
	tv := ui.NewTextList(label, strconv.FormatBool(value), "No", "Yes")
	return &YesNoPrompt{TextList: *tv}
}

package valueeditors

import (
	"github.com/eurozulu/pempal/commandline/ui"
	"github.com/eurozulu/pempal/logger"
	"strings"
)

type EditorList struct {
	Editors          []ValueEditor
	ForegrondColour  ui.ItemColour
	BackgroundColour ui.ItemColour
	ExitChar         rune
}

func (le EditorList) Show(offset ui.ViewOffset, values map[string]string, errs []error) (map[string]string, error) {
	var clearScreen bool
	if isRoot, err := ui.InitUI(); err != nil {
		return nil, err
	} else if isRoot {
		clearScreen = true
		defer ui.CloseUi()
	}

	list := ui.ItemSelect{ExitChar: le.ExitChar}
	list.Items = le.buildItemList(values, errs)
	var selected int

	deltaValues := map[string]string{}

	for {
		if clearScreen {
			ui.Clear()
		}
		// If no errors, Add the "Confirm key" prompt
		if len(errs) == 0 {
			list.ExitChar = 'Y'
			writePrompt(offset.OffsetY(len(list.Items)+1), "Hit 'Y' (Captial Letter) to confirm")
		}

		index, err := list.Select(offset, selected)
		if err != nil {
			return nil, err
		}
		if index < 0 {
			// Exit char hit
			return nil, nil
		}
		selected = index
		name := list.Items[selected].Name
		editor := le.editorByName(name)
		if editor == nil {
			logger.Debug("Ignoring edit of %s as no editor of that name found", name)
			continue
		}
		v, err := editor.Edit(offset.OffsetY(selected), values[name])
		if err != nil {
			if err == ui.ERRAborted {
				// abort value editor, return to main list
				continue
			}
			return nil, err
		}
		// no change?
		if v == values[name] {
			continue
		}

		deltaValues[name] = v
		// no exit char = exit on first change
		if le.ExitChar == 0 {
			return deltaValues, nil
		}
	}
}

func writePrompt(offset ui.ViewOffset, prompt string) {
	ui.PrintF(&offset, ui.ColourForeground, ui.ColourBackground, prompt)
}

func (le EditorList) editorByName(name string) ValueEditor {
	for _, ed := range le.Editors {
		if ed.Name() == name {
			return ed
		}
	}
	return nil
}

func (le EditorList) editorNames() []string {
	names := make([]string, len(le.Editors))
	for i, ve := range le.Editors {
		names[i] = ve.Name()
	}
	return names
}

func (le EditorList) buildItemList(values ui.ListValues, errs []error) []ui.ListItem {
	var items []ui.ListItem
	for _, editName := range le.editorNames() {
		v, ok := values[editName]
		errIndex := errorIndexByName(editName, errs)

		// if no value to edit or no error with editable name, ignore it.
		if !ok && errIndex < 0 {
			continue
		}
		li := ui.ListItem{
			Name:             editName,
			Value:            v,
			ForegrondColour:  le.ForegrondColour,
			BackgroundColour: le.BackgroundColour,
		}
		if errIndex >= 0 {
			// replace name with the value, within the same error message.
			li.Value = strings.Replace(errs[errIndex].Error(), editName, v, -1)
			li.ForegrondColour = ui.ColourForegroundError
		}
		items = append(items, li)
	}
	return items
}

func errorIndexByName(name string, errs []error) int {
	name = strings.ToLower(name)
	for i, err := range errs {
		if strings.Contains(strings.ToLower(err.Error()), name) {
			return i
		}
	}
	return -1
}

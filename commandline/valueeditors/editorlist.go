package valueeditors

import (
	"github.com/eurozulu/pempal/commandline/ui"
	"strings"
)

type EditorList []ValueEditor

func (le EditorList) Show(offset ui.ViewOffset, listValues map[string]string, errs ...error) (map[string]string, error) {
	isRoot, err := ui.InitUI()
	if err != nil {
		return nil, err
	}
	if isRoot {
		defer ui.CloseUi()
	}

	deltaValues := ui.ListValues{}
	var selected int
	names := le.editorNames()
	list := ui.ValueSelect{Names: names, ExitChar: 'Y'}
	values := ui.ListValues{}
	copyMapValues(values, listValues)

	for {
		if isRoot {
			ui.Clear()
		}

		sl, err := list.Select(offset, selected, values)
		if err != nil {
			return nil, err
		}
		if sl < 0 {
			// Exit char hit
			return deltaValues, nil
		}
		selected = sl
		name := names[selected]
		edit := le.editorByName(name)
		if edit == nil {
			continue
		}
		os := offset
		os.YOffset += selected
		v, err := edit.Edit(os, values[name])
		if err != nil {
			if err == ui.ERRAborted {
				// abort value edit, return to form
				continue
			}
			return nil, err
		}
		deltaValues[name] = v
		values[name] = v
	}
}

func (le EditorList) editorByName(name string) ValueEditor {
	for _, ed := range le {
		if ed.Name() == name {
			return ed
		}
	}
	return nil
}

func (le EditorList) editorNames() []string {
	names := make([]string, len(le))
	for i, ve := range le {
		names[i] = ve.Name()
	}
	return names
}

func copyMapValues(dst ui.ListValues, src ui.ListValues) {
	for k, v := range src {
		dst[k] = v
	}
}

func errorByName(name string, errs []error) error {
	name = strings.ToLower(name)
	for _, err := range errs {
		if strings.Contains(strings.ToLower(err.Error()), name) {
			return err
		}
	}
	return nil
}

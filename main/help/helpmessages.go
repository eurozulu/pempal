package help

import (
	"bytes"
	"pempal/utils"
	"sort"
)

var helpCommands = map[string]string{
	"find":      "locates resources based on search criteria set in its flags.",
	"make":      "creates new resources based on the given template(s)",
	"config":    "shows or sets the application configuration",
	"templates": "manages the templates, displaying names and add/remove of custom templates",
	"template":  "displays a template based on the given template name(s). Multiple names are merged into a single template.",
}

func HelpCommands() string {
	return mapToOrderedList(helpCommands)
}

func mapToOrderedList(m map[string]string) string {
	keys := make([]string, len(m))
	var index int
	for key := range m {
		keys[index] = key
		index++
	}
	sort.Strings(keys)
	buf := bytes.NewBuffer(nil)
	out := utils.NewColumnOutput(buf)
	for _, key := range keys {
		out.WriteSlice([]string{key, m[key]})
		out.Write([]byte("\n"))
	}
	return buf.String()
}

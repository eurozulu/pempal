package help

import (
	"bytes"
	"pempal/utils"
	"sort"
)

var helpCommands = map[string]string{
	"find":      "Locates resources based on search criteria set in its flags.",
	"make":      "Creates new resources based on the given template(s)",
	"config":    "Shows or sets the application configuration",
	"templates": "Manages the templates, displaying names and add/remove of custom templates",
	"template":  "Displays a template based on the given template name(s). Multiple names are merged into a single template.",
	"keys":      "Lists the identities of all private keys and the certificates they've signed..",
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

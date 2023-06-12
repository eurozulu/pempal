package help

import "bytes"

// HelpPage represents a single page of help information
type HelpPage struct {
	Title       string
	Aliases     []string
	Description string
	Format      string
	Flags       []HelpPageFlag
}

// HelpPageFlag represents the help contents for a single flag
type HelpPageFlag struct {
	Name        string
	Description string
	Aliases     []string
}

func (h HelpPage) String() string {
	buf := bytes.NewBufferString("\n")
	buf.WriteString(h.Format)
	buf.WriteString("\n")
	buf.WriteString(h.Title)
	buf.WriteString("\n\n")
	if len(h.Aliases) > 0 {
		buf.WriteString("Alias names: ")
		for i, a := range h.Aliases {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(a)
		}
		buf.WriteRune('\n')
	}
	buf.WriteString(h.Description)
	buf.WriteString("\n\nFlags:\n")
	if len(h.Flags) == 0 {
		buf.WriteString("no flags\n")
	}
	for _, flag := range h.Flags {
		buf.WriteString(flag.String())
		buf.WriteString("\n")
	}
	return buf.String()
}

func (h HelpPageFlag) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteRune('-')
	buf.WriteString(h.Name)
	if len(h.Aliases) > 0 {
		buf.WriteString("\t(")
		for i, a := range h.Aliases {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteRune('-')
			buf.WriteString(a)
		}
		buf.WriteString(")")
	}
	buf.WriteRune('\n')
	buf.WriteString(h.Description)
	buf.WriteRune('\n')
	return buf.String()
}

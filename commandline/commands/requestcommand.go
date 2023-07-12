package commands

import "io"

// requestCommand manages the Certificate requests (CSRs)
// When no arguments are given, lists all the known CSRs with matching private identity.
// When one or more arguments or flags are given attempts to generate a new CSR based on the named templates and flags.
// Editors required to build a new CSR are:
// subject.common-name  A non empty string
// public-key The key id, name or PEM to sign the new CSR
// These may be provided as flags or from within named templates.
type requestCommand struct {
}

func (r requestCommand) Execute(args []string, out io.Writer) error {
	//TODO implement me
	panic("implement me")
}

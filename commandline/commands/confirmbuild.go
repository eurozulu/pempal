package commands

import (
	"bufio"
	"bytes"
	"github.com/eurozulu/pempal/builders"
	"github.com/eurozulu/pempal/commandline/ui"
	"github.com/eurozulu/pempal/commandline/valueeditors"
	"github.com/eurozulu/pempal/templates"
	"io"
	"os"
)

func confirmBuild(prompt string, editors []valueeditors.ValueEditor, builder builders.Builder) error {
	offset := ui.ViewOffset{}
	tedit := &valueeditors.EditorList{Editors: editors}
	for {
		errs := builder.Validate()
		if CommonFlags.Quiet {
			// when quiet, either report errors or continue without confirmation when no errors
			return errs
		}

		t, err := tedit.Show(offset, builder.BuildTemplate(), errs)
		if err != nil {
			return err
		}
		if len(t) > 0 {
			// New property values, add to builder and revalidate
			builder.AddTemplate(templates.Template(t))
			continue
		}
		// Nothing new added, build confirmed
		return nil
	}
}

func readStdIn() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		buf.Write(scan.Bytes())
	}
	if scan.Err() != nil && scan.Err() != io.EOF {
		return nil, scan.Err()
	}
	return buf.Bytes(), nil
}

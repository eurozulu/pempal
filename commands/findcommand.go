package commands

import (
	"bytes"
	"context"
	"fmt"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/resourcefiles"
	"github.com/eurozulu/pempal/resourceformat"
	"github.com/eurozulu/pempal/tools"
	"io"
	"path/filepath"
	"strings"
)

// FindCommand finds the resources in the given path(s)
// @Command(find)
type FindCommand struct {

	// Counts when set displays a list of totals of the resources found
	// @Flag(count,c)
	Counts bool

	// Where specifies a simple expression to filter the results
	// @Flag(where, w)
	Where string
}

// ViewResources lists any pem reslurces found in the given path(s)
// path is a required file path to a directory or file.
// optionally additional paths may be given, seperated with a space.
// @Action
func (cmd FindCommand) ViewResources(path string, paths ...string) (string, error) {
	paths = append([]string{path}, paths...)
	path = strings.Join(paths, string(filepath.ListSeparator))
	format, err := resourceformat.NewResourceFormat("list")
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(nil)
	filter := cmd.buildFilter()
	counts := map[string]int{}
	total := 0
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	pemz := resourcefiles.PemFiles(path)
	for pemFile := range pemz.Find(ctx, filter) {
		if cmd.Counts {
			total++
			addToTotals(pemFile, counts)
		}

		if err := format.Format(buf, pemFile); err != nil {
			return "", err
		}
	}
	if cmd.Counts {
		if err := writeTotals(counts, total, buf); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}

func (cmd FindCommand) buildFilter() resourcefiles.PemFileFilter {
	return nil
}

func addToTotals(pf *model.PemFile, counts map[string]int) {
	for _, blk := range pf.Blocks {
		counts[tools.ToTitle(model.ParseResourceType(blk.Type).String())]++
	}
}

func writeTotals(counts map[string]int, total int, out io.Writer) error {
	out.Write([]byte("\n"))
	for k, v := range counts {
		if v > 1 {
			k = strings.Join([]string{k, "s"}, "")
		}
		if _, err := fmt.Fprintf(out, "%d %s\n", v, k); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(out, "Total: %d \n", total); err != nil {
		return err
	}
	return nil
}

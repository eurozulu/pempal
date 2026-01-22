package commands

import (
	"bytes"
	"context"
	"github.com/eurozulu/pempal/resourcefiles"
	"github.com/eurozulu/pempal/resourceformat"
	"path/filepath"
	"strings"
)

// ViewCommand display the existing resources in various formats.
// @Command(view)
type ViewCommand struct {

	// Format specifies the output format of the resource(s)
	// Valid formats are:
	// list	The default, lists basic details of each resource
	// pem	Outputs resource(s) as Pem encoded blocks
	// yaml Outputs a yaml document(s) of the properties in each resource
	// json Outputs a json document(s) of the properties in each resource
	// @Flag(format,f)
	Format string
}

// ViewResources lists any pem reslurces found in the given path(s)
// path is a required file path to a directory or file.
// optionally additional paths may be given, seperated with a space.
// @Action
func (cmd ViewCommand) ViewResources(path string, paths ...string) (string, error) {
	paths = append([]string{path}, paths...)
	path = strings.Join(paths, string(filepath.ListSeparator))
	format, err := resourceformat.NewResourceFormat(cmd.Format)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(nil)
	filter := cmd.buildFilter()
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	pemz := resourcefiles.PemFiles(path)
	for pemFile := range pemz.Find(ctx, filter) {
		if err := format.Format(buf, pemFile); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}

func (cmd ViewCommand) buildFilter() resourcefiles.PemFileFilter {
	return nil
}

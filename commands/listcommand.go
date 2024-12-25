package commands

import (
	"context"
	"fmt"
	"github.com/eurozulu/pempal/query"
	"github.com/eurozulu/pempal/utils"
	"io"
	"os"
	"sort"
)

// ListCommand lists the PEM resources in the given path(s)
// Each file containing one or more pems is show with the details of those pem resources.
// The properties show are dependent on the types of pems being listed.
// If no type is stated, the generic type and ID are show.
// Specifying one or more types with the -type flag will show properties more relevant to those types.
// Additional fields can be show using the -fields flag and a comma delimited list of field names to show.
// If the field name begins with a '+' that field is added to the regular fields for the resource type.
// if no given -fields have the preceeding plus, only the given fields are show.
type ListCommand struct {
	Output    io.Writer
	ListQuery query.ResourceQuery `flag:"+"`
}

func (cmd ListCommand) Exec(args ...string) error {
	if len(args) == 0 {
		args = []string{"."}
	}

	out := &utils.ColumnWriter{
		Columns:      query.ColumnsByName(cmd.ListQuery.ColumnNames()),
		ColumnSpacer: " ",
		Out:          cmd.Output,
	}
	if out.Out == nil {
		out.Out = os.Stdout
	}

	if !CommonFlags.Quiet {
		if err := cmd.WriteHeaderNames(out); err != nil {
			return err
		}
	}

	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	for res := range cmd.ListQuery.Query(ctx, args...) {
		if err := cmd.ListResources(res, out); err != nil {
			return err
		}
	}
	return nil
}

func (cmd ListCommand) WriteHeaderNames(out *utils.ColumnWriter) error {
	if _, err := out.WriteStrings(out.ColumnNames()); err != nil {
		return err
	}
	return nil
}

func (cmd ListCommand) ListResources(res []query.ResourceProperties, out *utils.ColumnWriter) error {
	sort.Slice(res, func(i, j int) bool {
		f1 := res[i]["filename"].(string)
		f2 := res[j]["filename"].(string)
		return f1 < f2
	})
	for _, prop := range res {
		vals := orderPropertiesByName(prop, out.Columns)
		if _, err := out.WriteStrings(vals); err != nil {
			return err
		}
	}
	return nil
}

func orderPropertiesByName(prop query.ResourceProperties, cols []utils.Column) []string {
	sz := make([]string, len(prop))
	for i, c := range cols {
		v, ok := prop[c.Name]
		if !ok {
			continue
		}
		sz[i] = fmt.Sprintf("%v", v)
	}
	return sz
}

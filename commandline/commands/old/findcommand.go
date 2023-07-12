package old

import (
	"context"
	"fmt"
	"github.com/eurozulu/pempal/commandline/commands"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/resourceio"
	"github.com/eurozulu/pempal/utils"
	"io"
	"strconv"
	"strings"
)

var defaultFieldsForTypes = map[model.ResourceType]string{
	model.Unknown:            "resource-type,identity,subject.common-name,location",
	model.Certificate:        "identity,version[7],serial-number,subject.common-name[30],issuer.common-name[30],not-before,not-after,location",
	model.CertificateRequest: "identity,version[7],subject.common-name,location",
	model.PrivateKey:         "identity,public-key-algorithm,key-param,is-encrypted,location",
	model.PublicKey:          "identity,public-key-algorithm,location",
	model.RevokationList:     "identity,issuer.common-name,number,this-update,next-update,location",
}

type FindCommand struct {
	ResourceType string   `flag:"resource-type,type"`
	Query        string   `flag:"query,qy"`
	Recursive    bool     `flag:"recursive,r"`
	Fields       []string `flag:"fields,fd"`
	//SortField    string   `yaml:"sort,sort-field"`

	ResourceTypes []model.ResourceType `flag:"-"`
	columnWidths  []int
	query         resourceio.Query
}

func (fc FindCommand) Execute(args []string, out io.Writer) error {
	if a, err := cleanArguments(args); err != nil {
		return err
	} else {
		args = a
	}

	if len(args) == 0 {
		return fmt.Errorf("Find requires at least one path to a file or directory.")
	}

	if err := fc.resolveResourceTypes(); err != nil {
		return err
	}
	fc.resolveFields()
	if err := fc.resolveQuery(); err != nil {
		return err
	}

	listOut := utils.NewColumnOutput(out, fc.columnWidths...)
	if !commands.CommonFlags.Quiet {
		if err := fc.writeTitles(listOut); err != nil {
			return err
		}
		listOut.WriteString("\n")
	}
	var count int

	for locFields := range fc.ListResources(context.Background(), args...) {
		if _, err := listOut.WriteSlice(locFields); err != nil {
			return err
		}
		if err := listOut.WriteString("\n"); err != nil {
			return err
		}
		count++
	}
	if !commands.CommonFlags.Quiet {
		fmt.Fprintf(out, "found %d resources\n", count)
	}
	return nil
}

func (fc FindCommand) ListResources(ctx context.Context, paths ...string) <-chan []string {
	locations := resourceio.NewResourceScanner(fc.Recursive).Scan(ctx, paths...)
	return resourceio.NewResourceLister(fc.Fields, fc.query, fc.ResourceTypes...).List(locations)
}

func (fc FindCommand) writeTitles(out *utils.ColumnOutput) error {
	titles := make([]string, len(fc.Fields))
	for i, fld := range fc.Fields {
		titles[i] = strings.Title(fld)
	}
	_, err := out.WriteSlice(titles)
	return err
}

func (fc *FindCommand) resolveQuery() error {
	if fc.Query == "" {
		return nil
	}
	var err error
	fc.query, err = resourceio.ParseQuery(fc.Query)
	return err

}

func (fc *FindCommand) resolveFields() {
	// Add/replace with flag given fields
	var fields []string
	var hardField bool // hard field indicates field name without '+' precurser. When present default fields are excluded.

	for _, f := range fc.Fields {
		if strings.HasPrefix(f, "+") {
			fields = append(fields, f[1:])
			continue
		}
		hardField = true
		fields = append(fields, f)
	}
	if !hardField {
		// no hard fields, include default fields for type (or generic)
		var rtk model.ResourceType
		if len(fc.ResourceTypes) == 1 {
			rtk = fc.ResourceTypes[0]
		}
		dfs := strings.Split(defaultFieldsForTypes[rtk], ",")
		fields = append(dfs, fields...)
	}
	fc.Fields = fields

	// Check for column widths
	fc.columnWidths = make([]int, len(fc.Fields))
	for i, fld := range fc.Fields {
		if !strings.HasSuffix(fld, "]") {
			continue
		}
		pos := strings.IndexRune(fld, '[')
		if pos < 0 {
			continue
		}
		si := strings.TrimSpace(fld[pos+1 : len(fld)-1])
		colWidth, err := strconv.Atoi(si)
		if err != nil {
			logger.Error("Invalid column width '%s' for field %s", si, fld[:pos])
		}
		fc.Fields[i] = fld[:pos]
		fc.columnWidths[i] = colWidth
	}
}

func (fc *FindCommand) resolveResourceTypes() error {
	fc.ResourceTypes = nil
	if fc.ResourceType == "" {
		return nil
	}
	for _, r := range strings.Split(fc.ResourceType, ",") {
		rt := model.ParseResourceType(r)
		if rt == model.Unknown {
			return fmt.Errorf("%s is not a known resource type", r)
		}
		fc.ResourceTypes = append(fc.ResourceTypes, rt)
	}
	return nil
}

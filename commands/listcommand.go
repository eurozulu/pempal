package commands

import (
	"context"
	"fmt"
	"io"
	"pempal/resourcefinder"
	"pempal/resources"
	"strings"
)

type ListCommand struct {
	Recursive bool `flag:"recursive"`
	recursive bool `flag:"r"`

	ResourceTypes string `flag:"type"`
	resourceTypes string `flag:"t"`

	TypeCertificates bool `flag:"certificates"`
	typeCertificates bool `flag:"certs"`

	TypeKeys bool `flag:"private-keys"`
	typeKeys bool `flag:"keys"`

	TypeRequests bool `flag:"requests"`
	typeRequests bool `flag:"csrs"`
	typeRequest  bool `flag:"csr"`

	TypeRevokeLists bool `flag:"revoke-lists"`
	typeRevokeLists bool `flag:"crls"`
	typeRevokeList  bool `flag:"crl"`

	// Specifies the output format
	Format string `flag:"format"`
	format string `flag:"f"`
}

func (cmd ListCommand) Run(ctx context.Context, args Arguments, out io.Writer) error {
	locs := args.Parameters()
	if len(locs) == 0 {
		return fmt.Errorf("no location specified.  specify at least one file or path to search")
	}
	types, err := cmd.getTypes()
	if err != nil {
		return err
	}

	query := cmd.makeQuery(args)

	scanner := resourcefinder.NewResourceScanner(types...)
	scanner.Find(ctx)
	return nil
}

func (cmd ListCommand) makeQuery(args Arguments) resourcefinder.Query {
	q := resourcefinder.Query{}
	for _, k := range args.FlagNames() {
		q[k] = resourcefinder.QueryValue{Value: args.FlagValue(k)}
	}
	return q
}

func (cmd ListCommand) getTypes() ([]resources.ResourceType, error) {
	// use map of types to ensure a unique set
	var rts map[resources.ResourceType]bool
	if cmd.TypeCertificates || cmd.typeCertificates {
		rts[resources.Certificate] = true
	}
	if cmd.TypeKeys || cmd.typeKeys {
		rts[resources.Key] = true
	}
	if cmd.TypeRequests || cmd.typeRequests || cmd.typeRequest {
		rts[resources.Request] = true
	}
	if cmd.TypeRevokeLists || cmd.typeRevokeLists || cmd.typeRevokeList {
		rts[resources.RevocationList] = true
	}

	ts := strings.Join([]string{cmd.ResourceTypes, cmd.resourceTypes}, ",")
	for _, s := range strings.Split(ts, ",") {
		rt := resources.ParseResourceType(strings.TrimSpace(s))
		if rt == resources.Unknown {
			return nil, fmt.Errorf("-type '%s' is unknown", s)
		}
		rts[rt] = true
	}
	var rs []resources.ResourceType
	for rt := range rts {
		rs = append(rs, rt)
	}
	return rs, nil
}

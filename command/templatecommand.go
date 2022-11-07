package command

import (
	"context"
	"encoding"
	"fmt"
	"io"
	"pempal/finder"
)

type templateCommand struct {
	VerboseOutput   bool   `flag:"verbose,v"`
	RecursiveSearch bool   `flag:"recursive,r"`
	PemTypes        string `flag:"type,t"`
}

func (tc templateCommand) Run(ctx context.Context, args Arguments, out io.Writer) error {
	locs := args.Parameters()
	if len(locs) == 0 {
		return fmt.Errorf("no location to search given")
	}

	pts, err := parsePemTypes(tc.PemTypes)
	if err != nil {
		return err
	}
	// Scan available pem finder encoding them into templates
	res := finder.NewTransformerFinder(nil, tc.RecursiveSearch, tc.VerboseOutput, pts...)
	for l := range res.Find(ctx, locs...) {
		data, err := l.(encoding.TextMarshaler).MarshalText()
		if err != nil {
			return err
		}
		_, err = out.Write(data)
		if err != nil {
			return err
		}
	}
	return nil
}

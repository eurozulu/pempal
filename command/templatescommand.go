package command

import (
	"context"
	"encoding"
	"fmt"
	"io"
	"os"
	"pempal/finder"
)

type templatesCommand struct {
	VerboseOutput   bool `flag:"verbose,v"`
	RecursiveSearch bool `flag:"recursive,r"`
	Show            bool `flag:"show"`
}

func (tc templatesCommand) Run(ctx context.Context, args Arguments, out io.Writer) error {
	locs := args.Parameters()
	if len(locs) == 0 {
		return fmt.Errorf("no location to search given")
	}

	res := finder.NewTemplateFinder(tc.RecursiveSearch, tc.VerboseOutput)
	for tl := range res.Find(ctx, locs...) {
		if tc.Show {
			data, err := tl.(encoding.TextMarshaler).MarshalText()
			if err != nil {
				return err
			}
			_, err = out.Write(data)
			if err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintln(os.Stdout, tl.Path()); err != nil {
				return err
			}
		}
	}
	return nil
}

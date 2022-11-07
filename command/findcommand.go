package command

import (
	"context"
	"encoding"
	"fmt"
	"io"
	"pempal/finder"
	"pempal/pemtypes"
	"strings"
)

type findCommand struct {
	VerboseOutput   bool   `flag:"verbose,v"`
	RecursiveSearch bool   `flag:"recursive,r"`
	OutputFormat    string `flag:"format,f"`

	Query    string `flag:"query,q"`
	PemTypes string `flag:"type,t"`
}

type outputFunc func(l finder.Location) ([]byte, error)

func (fc findCommand) Run(ctx context.Context, args Arguments, out io.Writer) error {
	locs := args.Parameters()
	if len(locs) == 0 {
		return fmt.Errorf("no location to search given")
	}

	pts, err := parsePemTypes(fc.PemTypes)
	if err != nil {
		return err
	}

	write, err := fc.outputFormat()
	if err != nil {
		return err
	}
	res := finder.NewTransformerFinder(&fc.Query, fc.RecursiveSearch, fc.VerboseOutput, pts...)

	for l := range res.Find(ctx, locs...) {
		data, err := write(l)
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

func (fc findCommand) outputFormat() (outputFunc, error) {
	of := strings.ToLower(fc.OutputFormat)
	if of == "" || of == "text" {
		return fc.writeAsText, nil
	}
	if of == "pem" {
		return fc.writeAsPem, nil
	}
	return nil, fmt.Errorf("%s is not a know output format. use 'pem' or 'text'", of)
}

func (fc findCommand) writeAsPem(l finder.Location) ([]byte, error) {
	if br, ok := l.(encoding.BinaryMarshaler); !ok {
		return nil, fmt.Errorf("%s is not a binary resource", l.Path())
	} else {
		return br.MarshalBinary()
	}
}
func (fc findCommand) writeAsText(l finder.Location) ([]byte, error) {
	if tr, ok := l.(encoding.TextMarshaler); !ok {
		return nil, fmt.Errorf("%s is not a text resource", l.Path())
	} else {
		return tr.MarshalText()
	}
}

func parsePemTypes(types string) ([]pemtypes.PEMType, error) {
	if types == "" {
		return nil, nil
	}
	pts := strings.Split(types, ",")
	var ts []pemtypes.PEMType
	for _, p := range pts {
		pt := pemtypes.ParsePEMType(strings.ToUpper(strings.TrimSpace(p)))
		if pt == pemtypes.Unknown {
			return nil, fmt.Errorf("%s is not a known pem type", p)
		}
		ts = append(ts, pt)
	}
	return ts, nil
}

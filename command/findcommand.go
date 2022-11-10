package command

import (
	"bytes"
	"context"
	"encoding"
	"fmt"
	"io"
	"log"
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

	//indexr indexer.Indexer
}

type outputFunc func(l finder.PemLocation) ([]byte, error)

func (fc findCommand) Run(ctx context.Context, args Arguments, out io.Writer) error {
	prms := args.Parameters()
	// Add any ENV paths, filtered by types being sought
	if finder.PPPath != "" {
		pathParms, err := finder.CleanPPPath()
		if err != nil {
			log.Printf("Warning.  %s contains an invalid entry: %v", finder.ENV_PATH, err)
		}
		prms = append(prms, pathParms...)
	}

	if len(prms) == 0 {
		return fmt.Errorf("no location to search given or no $%s set", finder.ENV_PATH)
	}
	locs, plusNames := splitPlusParams(prms)
	if len(plusNames) > 0 {
		pts, err := parsePemTypes(strings.Join(plusNames, ","))
		if err != nil {
			return fmt.Errorf("Resource indexing failed, +%v", err)
		}
		fc.startIndex(pts)
	}

	pts, err := parsePemTypes(fc.PemTypes)
	if err != nil {
		return err
	}

	outputWrite, err := fc.outputFormat()
	if err != nil {
		return err
	}

	finder := finder.NewFinder(nil, fc.RecursiveSearch, fc.VerboseOutput, pts...)
	foundCh, err := finder.Find(ctx, locs...)

	if err != nil {
		return fmt.Errorf("invalid parameter  %v", err)
	}

	for pl := range foundCh {
		// format the location into output
		data, err := outputWrite(pl)
		if err != nil {
			return err
		}

		// write formatted output to final output
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

func (fc findCommand) writeAsPem(pl finder.PemLocation) ([]byte, error) {
	if br, ok := pl.(encoding.BinaryMarshaler); !ok {
		return nil, fmt.Errorf("%s is not a binary resource", pl.Location())
	} else {
		return br.MarshalBinary()
	}
}
func (fc findCommand) writeAsText(pl finder.PemLocation) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(pl.Location())
	buf.WriteRune('\n')
	for _, r := range pl.Resources() {
		buf.WriteRune('\t')
		buf.WriteString(r.String())
		buf.WriteRune('\n')
	}
	return buf.Bytes(), nil
}

func (fc findCommand) startIndex(types []pemtypes.PEMType) error {
	log.Println("Index not yet functioning")
	log.Printf("types: %v", types)
	log.Printf(" will NOT be indexed")
	return nil
}

func splitPlusParams(prms []string) (nonplus, plus []string) {
	for _, p := range prms {
		if strings.HasPrefix(p, "+") {
			plus = append(plus, strings.TrimLeft(p, "+"))
			continue
		}
		nonplus = append(nonplus, p)
	}
	return nonplus, plus
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

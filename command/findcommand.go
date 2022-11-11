package command

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"pempal/finder"
	"pempal/pemtypes"
	"pempal/pemwriter"
	"strings"
)

type findCommand struct {
	VerboseOutput   bool   `flag:"verbose,v"`
	RecursiveSearch bool   `flag:"recursive,r"`
	OutputFormat    string `flag:"format,f"`
	Query           string `flag:"query,q"`
	PemTypes        string `flag:"type,t"`

	//indexr indexer.Indexer
}

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

	var outformat pemwriter.PemWriterType
	if fc.OutputFormat != "" {
		outformat = pemwriter.ParseWriterType(fc.OutputFormat)
		if outformat == pemwriter.Unknown {
			return fmt.Errorf("'%s' is not a known output format. Use one of %v", fc.OutputFormat, pemwriter.PemWriterNames)
		}
	}
	pemOut := pemwriter.NewPemWriter(outformat, out)

	finder := finder.NewFinder(nil, fc.RecursiveSearch, fc.VerboseOutput, pts...)
	foundCh, err := finder.Find(ctx, locs...)

	if err != nil {
		return fmt.Errorf("invalid parameter  %v", err)
	}

	var lc, rc int
	for pl := range foundCh {
		lc++
		fmt.Fprintln(os.Stdout, pl.Location())
		for _, r := range pl.Resources() {
			rc++
			if err := pemOut.Write(r); err != nil {
				return err
			}
			fmt.Fprintln(os.Stdout)
		}
	}
	fmt.Fprintf(os.Stdout, "\nfound %d resources in %d locations\n", rc, lc)
	return nil
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

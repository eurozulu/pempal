package main

import (
	"context"
	"github.com/eurozulu/pempal/encoding"
	"github.com/eurozulu/pempal/filescan"
	"log"
	"os"
)

type ViewCommand struct {
	Verbose bool   `flag:"verbose,v"`
	Encode  string `flag:"encode,e"`

	OutPath string `flag:"outpath,out"`
}

// View displays the given pathname(s)
// path can point to certificate, key, csr, crl or a container, a directory, pkcs#7, pkcs#12.
func (sc ViewCommand) View(args ...string) error {
	if len(args) == 0 {
		log.Fatalln("must provide at least one path")
	}
	if sc.Encode == "" {
		sc.Encode = "yaml"
	}
	out := os.Stdout
	if sc.OutPath != "" {
		f, err := os.OpenFile(sc.OutPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Println(err)
			}
		}()
		out = f
	}
	ec, err := encoding.NewEncoder(sc.Encode)
	if err != nil {
		return err
	}

	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()

	ds := filescan.DirectoryScanner{}
	ch := ds.ScanDirectories(ctx, args)
	for {
		select {
		case <-ctx.Done():
			return nil
		case tps, ok := <-ch:
			if !ok {
				return nil
			}
			if err := ec.Encode(out, tps); err != nil {
				return err
			}
		}
	}
	return nil
}

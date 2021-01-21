package main

import (
	"bufio"
	"context"
	"github.com/eurozulu/pempal/encoding"
	"github.com/eurozulu/pempal/filescan"
	"log"
	"os"
	"strings"
)

type ViewCommand struct {
	Verbose     bool `flag:"verbose,v"`
	VeryVerbose bool `flag:"vv"`
	Recursive   bool `flag:"recursive,r"`

	Encode string `flag:"encode,e"`

	OutPath string `flag:"outpath,out"`

	Password string `flag:"pass,p"`
}

// View displays the given pathname(s)
// path can point to certificate, key, csr, crl or a container, a directory, pkcs#7, pkcs#12.
func (sc ViewCommand) View(args ...string) error {
	if len(args) == 0 {
		log.Fatalln("must provide at least one path to view")
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
	ec, err := encoding.NewEncoder(sc.Encode, out)
	if err != nil {
		return err
	}

	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()

	ds := filescan.DirectoryScanner{
		Recursive:   sc.Recursive,
		PrintErrors: sc.VeryVerbose,
		Password:    sc.Password,
	}

	if args[0] == "-" {
		args = append(args[1:], scanInput()...)
	}

	ch := ds.ScanDirectories(ctx, args)
	for {
		select {
		case <-ctx.Done():
			return nil
		case tps, ok := <-ch:
			if !ok {
				return nil
			}
			if err := ec.Encode(tps); err != nil {
				return err
			}
		}
	}
}

func scanInput() []string {
	var lines []string
	scn := bufio.NewScanner(os.Stdin)
	for scn.Scan() {
		l := strings.Fields(scn.Text())
		if len(l) == 0 || l[len(l)-1] == "" {
			continue
		}
		lines = append(lines, l[len(l)-1])
	}
	return lines
}

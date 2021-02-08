package main

import (
	"encoding/pem"
	"fmt"
	"github.com/pempal/pemio"
)

type ViewCommand struct {
	Command
	Password string `flag:"password,pwd,p"`
	Verbose  bool   `flag:"verbose,v"`
}

func (vc ViewCommand) ViewItems(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("provide at least on file path to view")
	}
	if vc.Encode == "" {
		vc.Encode = "yaml"
	}

	for _, p := range args {
		ps := pemio.PEMScanner{
			FilePath:  p,
			Recursive: false,
			Verbose:   vc.Verbose,
		}
		pfs, err := ps.ScanPath()
		if err != nil {
			return fmt.Errorf("failed to read %s  %v", p, err)
		}

		if err := vc.writePemFiles(pfs); err != nil {
			return fmt.Errorf("failed to format pem from %s  %v", p, err)
		}
	}
	return nil
}

func (vc ViewCommand) writePemFiles(pfs []*pemio.PEMFile) error {
	var bls []*pem.Block
	for _, pf := range pfs {
		bls = append(bls, pf.Blocks...)
	}
	return vc.WriteOutput(bls, 0600)
}

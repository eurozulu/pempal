package main

import (
	"fmt"
	"github.com/pempal/pemio"
)

type ViewCommand struct {
	Password string `flag:"password,pwd,p"`
	Verbose  bool   `flag:"verbose,v"`
}

func (vc ViewCommand) ViewItems(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("provide at least on file path to view")
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

		if err := writePemFilesToOutput(pfs, 0640); err != nil {
			return fmt.Errorf("failed to format pem from %s  %v", p, err)
		}
	}
	return nil
}

package main

import (
	"encoding/pem"
	"github.com/pempal/pemio"
	"os"
)

// Out sets an output filename. Defaults to the standard output
var Out string

// Encode sets the output format. valid values are 'pem', 'der' or 'p12'
var Encode string

func writePemFilesToOutput(pems []*pemio.PEMFile, perm os.FileMode) error {
	var bls []*pem.Block
	for _, pf := range pems {
		bls = append(bls, pf.Blocks...)
	}
	return writePemsToOutput(bls, perm)
}

func writePemsToOutput(pems []*pem.Block, perm os.FileMode) error {
	if Encode == "" {
		Encode = "yaml"
	}
	if Out == "" {
		return pemio.WritePEMs(os.Stdout, pems, Encode)
	} else {
		return pemio.WritePEMsFile(Out, pems, Encode, perm, false)
	}
}

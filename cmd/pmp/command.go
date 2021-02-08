package main

import (
	"encoding/pem"
	"github.com/pempal/pemio"
	"io"
	"os"
)

const defaultEncode = "pem"

type Command struct {
	Input  io.Reader `flag:"-" yaml:"-"`
	Output io.Writer `flag:"-" yaml:"-"`

	// Out sets an output filename. Defaults to the standard output
	Out string `flag:"out,o" yaml:"-"`

	// Encode sets the output format. valid values are 'pem', 'der' or 'p12'
	Encode string `flag:"encode,e" yaml:"-"`

	// Truncate, when true, will truncate any file being written to, wiping any previous contents.
	// Ignored when not writing to file
	Truncate bool `flag:"encode,e" yaml:"-"`
}

func (c Command) ReadInput() ([]*pem.Block, error) {
	if c.Input == nil {
		c.Input = os.Stdin
	}
	return pemio.ReadPEMs(c.Input)
}

func (c Command) WriteOutput(bls []*pem.Block, perm os.FileMode) error {
	if c.Output == nil {
		c.Output = os.Stdout
	}
	if c.Encode == "" {
		c.Encode = defaultEncode
	}
	if c.Out != "" {
		return pemio.WritePEMsFile(c.Out, bls, c.Encode, perm, c.Truncate)
	}
	return pemio.WritePEMs(c.Output, bls, c.Encode)
}

package pemio

import (
	"log"
	"os"
	"path"

	"encoding/pem"
	"io/ioutil"
	"path/filepath"
)

// PEMScanner will scan a file or directory for resources which can be read into PEMS.
// Reads DER, PEM and p12 encoded files.
type PEMScanner struct {
	FilePath  string
	Recursive bool
	Verbose   bool
}
type PEMFile struct {
	Blocks   []*pem.Block
	FilePath string
}

func (ps PEMScanner) ScanPath() ([]*PEMFile, error) {
	p, err := filepath.Abs(ps.FilePath)
	if err != nil {
		return nil, err
	}

	return ps.scanPath(p)
}

func (ps PEMScanner) scanPath(p string) ([]*PEMFile, error) {
	fi, err := os.Stat(p)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return ps.readPEMFile(p)
	}
	fis, err := ioutil.ReadDir(p)
	if err != nil {
		return nil, err
	}

	var pfs []*PEMFile
	for _, fi := range fis {
		if !ps.Recursive && fi.IsDir() {
			continue
		}
		fp := path.Join(p, fi.Name())
		spfs, err := ps.scanPath(fp)
		if err != nil {
			if ps.Verbose {
				log.Printf("%s  %v", fp, err)
			}
			continue
		}
		if len(spfs) > 0 {
			pfs = append(pfs, spfs...)
		}
	}
	return pfs, nil
}

func (ps PEMScanner) readPEMFile(p string) ([]*PEMFile, error) {
	bls, err := ReadPEMsFile(p)
	if err != nil {
		return nil, err
	}
	return []*PEMFile{{
		Blocks:   bls,
		FilePath: p,
	}}, nil
}

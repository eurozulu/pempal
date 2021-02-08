package pemio

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"encoding/pem"
	"io/ioutil"
)

// ReadPEMs reads any resources which can be parsed into PEM files, from the given reader.
func ReadPEMs(in io.Reader) ([]*pem.Block, error) {
	by, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}

	pps := AllPemParsers()
	for _, pr := range pps {
		pbs, err := pr.ParsePems(by)
		if err != nil {
			if !strings.Contains(err.Error(), "unknown format") &&
				!strings.Contains(err.Error(), "pkcs12: error reading P12 data: asn1: structure error:") {
				return nil, err
			}
			continue
		}
		return pbs, nil
	}
	return nil, fmt.Errorf("no resources found")

}

// ReadPEMsFile reads any resources which can be parsed into PEM files, from the file at the given path.
func ReadPEMsFile(p string) ([]*pem.Block, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()
	return ReadPEMs(f)
}

func ReadPEMsFiles(ps []string) ([]*pem.Block, error) {
	var bls []*pem.Block
	for _, p := range ps {
		blss, err := ReadPEMsFile(p)
		if err != nil {
			return nil, err
		}
		bls = append(bls, blss...)
	}
	return bls, nil
}

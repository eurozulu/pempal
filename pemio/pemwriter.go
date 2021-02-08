package pemio

import (
	"fmt"
	"io"
	"log"
	"os"

	"encoding/pem"
	"path/filepath"
)

func WritePEMsFile(p string, bls []*pem.Block, encode string, perm os.FileMode, trunc bool) error {
	p, err := filepath.Abs(p)
	if err != nil {
		return err
	}

	flag := os.O_CREATE | os.O_WRONLY
	if trunc {
		flag |= os.O_TRUNC
	} else {
		flag |= os.O_APPEND
	}
	f, err := os.OpenFile(p, flag, perm)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()
	return WritePEMs(f, bls, encode)
}

func WritePEMs(out io.Writer, bls []*pem.Block, encode string) error {
	pm := NewPEMMarshaler(encode)
	if pm == nil {
		return fmt.Errorf("%s is not a known encoding", encode)
	}
	return pm.Marshal(out, bls)
}

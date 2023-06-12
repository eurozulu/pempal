package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/eurozulu/pempal/logger"
	"io"
	"os"
)

func promptYorN(msg string, defalt bool) bool {
	fmt.Print(msg)
	choice := []byte("Ny")
	if defalt {
		choice = []byte("nY")
	}
	fmt.Printf(" [%s]", choice)
	selection := bytes.ToLower(choice)
	buf := bufio.NewReader(os.Stdin)
	for {
		b, err := buf.ReadByte()
		if err != nil {
			if err != io.EOF {
				logger.Error("Failed to read standard in  %v", err)
			}
			return false
		}
		if i := bytes.IndexByte(selection, b); i >= 0 {
			return selection[i] == 'y'
		}
		if b == '\n' {
			return defalt
		}
	}
}

func readStdIn() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		buf.Write(scan.Bytes())
	}
	if scan.Err() != nil && scan.Err() != io.EOF {
		return nil, scan.Err()
	}
	return buf.Bytes(), nil
}

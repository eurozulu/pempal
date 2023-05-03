package model

import (
	"math/big"
	"os"
	"path/filepath"
	"strconv"
)

const serialNumberFilename = "serialnumber.txt"

type SerialNumber uint64

func (s SerialNumber) ToBigInt() *big.Int {
	return big.NewInt(int64(s))
}

func (s SerialNumber) MarshalText() (text []byte, err error) {
	return []byte(strconv.FormatUint(uint64(s), 16)), nil
}

func (s *SerialNumber) UnmarshalText(text []byte) error {
	i, err := strconv.ParseUint(string(text), 16, 64)
	if err != nil {
		return err
	}
	*s = SerialNumber(i)
	return nil
}

func (s SerialNumber) Save(dirpath string) error {
	by, err := s.MarshalText()
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dirpath, serialNumberFilename), by, 0640)
}

func (s *SerialNumber) Load(dirpath string) error {
	by, err := os.ReadFile(filepath.Join(dirpath, serialNumberFilename))
	if err != nil {
		return err
	}
	return s.UnmarshalText(by)
}

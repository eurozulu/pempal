package model

import (
	"math/big"
	"strconv"
)

type SerialNumber big.Int

func (s SerialNumber) String() string {
	i := big.Int(s)
	return strconv.FormatInt((&i).Int64(), 10)
}

func (s *SerialNumber) MarshalText() (text []byte, err error) {
	return []byte(s.String()), nil
}

func (s *SerialNumber) UnmarshalText(text []byte) error {
	n, err := strconv.ParseInt(string(text), 10, 64)
	if err != nil {
		return err
	}
	i := big.NewInt(n)
	*s = SerialNumber(*i)
	return nil
}

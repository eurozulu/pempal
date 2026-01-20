package model

import (
	"crypto/elliptic"
	"fmt"
	"strings"
)

type EllipticCurve int

const (
	UnknownCurve EllipticCurve = iota
	P224
	P256
	P384
	P521
)

var CurveNames = []string{
	"UnknownCurve",
	"P224",
	"P256",
	"P384",
	"P521",
}

func (c EllipticCurve) String() string {
	i := int(c)
	if i < 0 || i >= len(CurveNames) {
		i = 0
	}
	return CurveNames[i]
}

func (c EllipticCurve) ToCurve() elliptic.Curve {
	switch c {
	case P224:
		return elliptic.P224()
	case P256:
		return elliptic.P256()
	case P384:
		return elliptic.P384()
	case P521:
		return elliptic.P521()
	default:
		return nil
	}
}

func (c EllipticCurve) MarshalText() (text []byte, err error) {
	return []byte(c.String()), nil
}

func (c *EllipticCurve) UnmarshalText(text []byte) error {
	cu, err := ParseCurve(string(text))
	if err != nil {
		return err
	}
	*c = cu
	return nil
}

func ParseCurve(s string) (EllipticCurve, error) {
	for i, n := range CurveNames {
		if strings.EqualFold(n, s) {
			return EllipticCurve(i), nil
		}
	}
	return UnknownCurve, fmt.Errorf("unknown curve %s", s)
}

func NewCurve(c elliptic.Curve) (EllipticCurve, error) {
	return ParseCurve(c.Params().Name)
}

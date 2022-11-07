package encoders

import (
	"crypto"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"log"
	"net"
	"net/url"
	"strconv"
	"time"
)

func MarshalURIs(urls []*url.URL) []string {
	var ss []string
	for _, u := range urls {
		ss = append(ss, u.String())
	}
	return ss
}

func MarshalIPAddresses(ips []net.IP) []string {
	var ss []string
	for _, ip := range ips {
		ss = append(ss, ip.String())
	}
	return ss
}

func MarshalKeyUsage(k x509.KeyUsage) string {
	if k < 1 || int(k) >= len(keyUsageNames) {
		k = 0
	}
	return keyUsageNames[k]
}

func MarshalExtKeyUsage(kus []x509.ExtKeyUsage) []string {
	var usage []string
	for _, k := range kus {
		if k < 0 || int(k) >= len(ExtKeyUsageNames) {
			// Unknown ExtKeyUsage, ignore it
			continue
		}
		usage = append(usage, ExtKeyUsageNames[k])
	}
	return usage
}

func MarshalTime(t time.Time) string {
	return t.Format(TimeFormat)
}

func MarshalSizeFromKey(prk crypto.PublicKey) string {
	if prk == nil {
		return ""
	}
	switch v := prk.(type) {
	case rsa.PublicKey:
		return strconv.Itoa(v.Size())
	case *rsa.PublicKey:
		return strconv.Itoa(v.Size())
	case rsa.PrivateKey:
		return strconv.Itoa(v.Size())
	case *rsa.PrivateKey:
		return strconv.Itoa(v.Size())
	case *ecdsa.PublicKey:
		return MarshalCurve(v.Curve)
	case *ecdsa.PrivateKey:
		return MarshalCurve(v.Curve)
	case *ed25519.PrivateKey, *dsa.PrivateKey:
		return ""
	default:
		return ""
	}
}

func MarshalCurve(c elliptic.Curve) string {
	switch c {
	case elliptic.P224():
		return "p224"
	case elliptic.P256():
		return "p256"
	case elliptic.P384():
		return "p384"
	case elliptic.P521():
		return "p521"
	default:
		return ""
	}
}

func MarshalPublicKey(puk crypto.PublicKey) string {
	rKey, ok := puk.(rsa.PublicKey)
	if ok {
		return hex.EncodeToString(x509.MarshalPKCS1PublicKey(&rKey))
	}
	by, err := x509.MarshalPKIXPublicKey(puk)
	if err != nil {
		log.Println(err)
		return ""
	}
	return hex.EncodeToString(by)
}

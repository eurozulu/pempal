package model

import (
	"crypto/x509/pkix"
	"encoding/hex"
	"time"
)

type RevokedCertificateDTO struct {
	SerialNumber   SerialNumber   `yaml:"serial.txt-number"`
	RevocationTime time.Time      `yaml:"revocation-time"`
	Extensions     []ExtensionDTO `yaml:"extensions"`
}

type ExtensionDTO struct {
	Id       []int  `yaml:"id"`
	Critical bool   `yaml:"critical"`
	Value    string `yaml:"value"`
}

func newRevokedCertificateDTO(cert pkix.RevokedCertificate) RevokedCertificateDTO {
	var sn SerialNumber
	if cert.SerialNumber != nil {
		sn = SerialNumber(cert.SerialNumber.Uint64())
	}

	return RevokedCertificateDTO{
		SerialNumber:   sn,
		RevocationTime: cert.RevocationTime,
		Extensions:     newExtentionsDTOs(cert.Extensions),
	}
}

func newExtentionsDTO(extension pkix.Extension) ExtensionDTO {
	return ExtensionDTO{
		Id:       extension.Id,
		Critical: extension.Critical,
		Value:    hex.EncodeToString(extension.Value),
	}
}

func newExtentionsDTOs(extensions []pkix.Extension) []ExtensionDTO {
	exts := make([]ExtensionDTO, len(extensions))
	for i, e := range extensions {
		exts[i] = newExtentionsDTO(e)
	}
	return exts
}

func newRevokedCertificateDTOs(certs []pkix.RevokedCertificate) []RevokedCertificateDTO {
	cs := make([]RevokedCertificateDTO, len(certs))
	for i, c := range certs {
		cs[i] = newRevokedCertificateDTO(c)
	}
	return cs
}

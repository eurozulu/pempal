package templates

var builtInTemplates = map[string]Template{
	"selfsigned": {
		"Description": "Issues a certificate signed by the certificate owner.",
		"NotBefore":   "{{ (time.ANSIC) }}",
		"Subject": Template{
			"CommonName": "?",
		},
		"IssuedBy":        "{{.Subject }}",
		"PublicKeyAlgo":   "?",
		"PublicKeyLength": "?",
		"Certificates":    "KeyUsageCertSign, KeyUsageKeyAgreement, KeyUsageDigitalSignature, KeyUsageKeyEncipherment, KeyUsageDataEncipherment",
	},
	"client": {
		"Subject": Template{
			"Organisation": "?",
			"CommonName":   "?",
		},
		"PublicKeyAlgo":   "?",
		"PublicKeyLength": "?",
		"IssuedBy":        "?",
		"Certificates":    "KeyUsageKeyAgreement, KeyUsageDigitalSignature, KeyUsageKeyEncipherment, KeyUsageDataEncipherment",
	},
	"server": {
		"Subject": Template{
			"Organisation":   "?",
			"CommonName":     "?",
			"EmailAddresses": "?",
		},
		"PublicKeyAlgo":   "?",
		"PublicKeyLength": "?",
		"IssuedBy":        "{{.Subject}}",
		"Certificates":    "KeyUsageKeyAgreement, KeyUsageDigitalSignature, KeyUsageKeyEncipherment, KeyUsageDataEncipherment",
	},
	"caroot": {
		"Subject": Template{
			"Organisation":   "?",
			"CommonName":     "?",
			"EmailAddresses": "?",
		},
		"PublicKeyAlgo":   "?",
		"PublicKeyLength": "?",
		"IssuedBy":        "{{.Subject}}",
		"IsCA":            "true",
		"Certificates":    "KeyUsageCertSign, KeyUsageCRLSign, KeyUsageKeyAgreement, KeyUsageDigitalSignature, KeyUsageKeyEncipherment, KeyUsageDataEncipherment",
	},
	"ca": {
		"Subject": Template{
			"Organisation": "{{.IssuedBy.Organisation}}",
			"CommonName":   "?",
		},
		"PublicKeyAlgo":   "?",
		"PublicKeyLength": "?",
		"IssuedBy":        "?",
		"IsCA":            "true",
		"Certificates":    "KeyUsageCertSign, KeyUsageKeyAgreement, KeyUsageDigitalSignature, KeyUsageKeyEncipherment, KeyUsageDataEncipherment",
	},
}

package templates

var builtInTemplates = map[string]Template{
	"selfsigned": {
		"Subject.CommonName": "?",
		"IssuedBy":           "${.Subject}",
		"PublicKeyAlgo":      "?RSA",
		"PublicKeyLength":    "?2048",
		"Certificates":       "KeyUsageCertSign, KeyUsageKeyAgreement, KeyUsageDigitalSignature, KeyUsageKeyEncipherment, KeyUsageDataEncipherment",
	},
	"client": {
		"Subject.Organisation": "?",
		"Subject.CommonName":   "?",
		"PublicKeyAlgo":        "?RSA",
		"PublicKeyLength":      "?2048",
		"IssuedBy":             "?",
		"Certificates":         "KeyUsageKeyAgreement, KeyUsageDigitalSignature, KeyUsageKeyEncipherment, KeyUsageDataEncipherment",
	},
	"server": {
		"Subject.Organisation": "?",
		"Subject.CommonName":   "?",
		"PublicKeyAlgo":        "?RSA",
		"PublicKeyLength":      "?2048",
		"IssuedBy":             "?",
		"Certificates":         "KeyUsageKeyAgreement, KeyUsageDigitalSignature, KeyUsageKeyEncipherment, KeyUsageDataEncipherment",
	},
	"rootca": {
		"Subject.Organisation": "?",
		"Subject.CommonName":   "?",
		"PublicKeyAlgo":        "?RSA",
		"PublicKeyLength":      "?2048",
		"IssuedBy":             "${.Subject}",
		"Certificates":         "KeyUsageCertSign, KeyUsageCRLSign, KeyUsageKeyAgreement, KeyUsageDigitalSignature, KeyUsageKeyEncipherment, KeyUsageDataEncipherment",
	},
	"ca": {
		"Subject.Organisation": "${.IssuedBy.Organisation}",
		"Subject.CommonName":   "?",
		"PublicKeyAlgo":        "?RSA",
		"PublicKeyLength":      "?2048",
		"IssuedBy":             "?",
		"Certificates":         "KeyUsageCertSign, KeyUsageKeyAgreement, KeyUsageDigitalSignature, KeyUsageKeyEncipherment, KeyUsageDataEncipherment",
	},
}

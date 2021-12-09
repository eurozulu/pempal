package templates

import "pempal/pemresources"

var baseTemplates = map[string]Template{
	"certificate": &pemresources.Certificate{},
	"request":     &pemresources.CertificateRequest{},
	"privatekey":  &pemresources.PrivateKey{},
	"publickey":   &pemresources.PublicKey{},
	"rovokelist":  nil,
}

var builtInTemplates = map[string]interface{}{
	"privatekey": `{
		"pem_type":             "PRIVATE KEY",
		"public_key_algorithm": "?",
		"public_key_length":    "?",
		"private_key":          "?",
	}`,
	"rsakey": `{
		"wxtends":              "privatekey",
		"pem_type":             "RSA PRIVATE KEY",
		"public_key_algorithm": "rsa",
		"public_key_length":    2048,
		"is_encrypted":         true,
	}`,
	"eckey": `{
		"extends":              "rsakey",
		"pem_type":             "EC PRIVATE KEY",
		"public_key_algorithm": "ecdsa",
		"public_key_length":    348,
	}`,

	"request": `{
		"description":          "Represents a certificate request with the base properties.",
		"pem_type":             "CERTIFICATE REQUEST",
		"public_key":           "?",
		"public_key_algorithm": "?",
		"public_key_length":    "?",

		"subject":             "?",
		"signature_algorithm": "?",
	}`,

	"certificate": `{
		"description":        "Issues a certificate with the bare properties required for a certificate.",
		"pem_type":            "CERTIFICATE",
		"public_key":          "?",
		"public_key_algorithm": "?",
		"public_key_length":    "?",

		"subject":             "?",
		"issuer":              "?",
		"is_ca":               false,
		"max_path_len_zero":   true,
		"max_path_len":        0,
		"not_before":          "{{ (time.ANSIC) }}",
		"not_after":           "{{ (time.ANSIC.Add(time.Year) }}",
		"signature_algorithm": "?",
	}`,
	"selfsigned": `{
		"extends":     "certificate",
		"description": "Issues a certificate signed by the certificate owner.",
		"subject":     {
			"common_name": "?",
		},
		"issuer":      "{{.Subject }}",
		"key_usage":   "KeyUsageKeyAgreement, KeyUsageDigitalSignature, KeyUsageKeyEncipherment, KeyUsageDataEncipherment",
	}`,
	"caroot": `{
		"extends":     "selfsigned",
		"description": "Issues a certificate self signed root CA certificate",
		"subject": {
			"organisation":    "?",
			"common_name":     "-",
			"email_addresses": "?",
		},
		"is_ca":             true,
		"max_path_len_zero": false,
		"max_path_len":      "1",
		"key_usage":         "KeyUsageCertSign, KeyUsageCRLSign, KeyUsageKeyAgreement, KeyUsageDigitalSignature, KeyUsageKeyEncipherment, KeyUsageDataEncipherment",
	}`,
	"ca": `{
		"extends": "certificate",
		"subject": {
			"organisation":    "{{.issuer.organisation}}",
			"email_addresses": "{{.issuer.email_addresses}}",
			"common_name":     "-",
		},
		"key_usage": "KeyUsageCertSign, KeyUsageKeyAgreement, KeyUsageDigitalSignature, KeyUsageKeyEncipherment, KeyUsageDataEncipherment",
	}`,

	"server": `{
		"extends": "certificate",
		"subject": {
			"organisation": "?",
			"common_name":   "-",
		},
		"key_usage": "KeyUsageKeyAgreement, KeyUsageDigitalSignature, KeyUsageKeyEncipherment, KeyUsageDataEncipherment",
	}`,
	"client": `{
		"extends": "server",
		"subject": {
			"email_addresses": "?",
		},
	}`,
}

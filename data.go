package certsurfer

type HeartbeatData struct {
	Timestamp float64 `json:"timestamp"`
}

type DomainData struct {
	Data []string `json:"data"`
}

type CertificateData struct {
	Data struct {
		CertIndex int    `json:"cert_index"`
		CertLink  string `json:"cert_link"`
		LeafCert  struct {
			AllDomains []string `json:"all_domains"`
			Extensions struct {
				AuthorityInfoAccess           string `json:"authorityInfoAccess"`
				AuthorityKeyIdentifier        string `json:"authorityKeyIdentifier"`
				BasicConstraints              string `json:"basicConstraints"`
				CertificatePolicies           string `json:"certificatePolicies"`
				CtlSignedCertificateTimestamp string `json:"ctlSignedCertificateTimestamp"`
				ExtendedKeyUsage              string `json:"extendedKeyUsage"`
				KeyUsage                      string `json:"keyUsage"`
				SubjectAltName                string `json:"subjectAltName"`
				SubjectKeyIdentifier          string `json:"subjectKeyIdentifier"`
			} `json:"extensions"`
			Fingerprint        string     `json:"fingerprint"`
			Issuer             EntityData `json:"issuer"`
			NotAfter           int        `json:"not_after"`
			NotBefore          int        `json:"not_before"`
			SerialNumber       string     `json:"serial_number"`
			SignatureAlgorithm string     `json:"signature_algorithm"`
			Subject            EntityData `json:"subject"`
		} `json:"leaf_cert"`
		Seen   float64 `json:"seen"`
		Source struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"source"`
		UpdateType string `json:"update_type"`
	} `json:"data"`
}

type EntityData struct {
	C            *string `json:"C"`
	CN           string  `json:"CN"`
	L            *string `json:"L"`
	O            *string `json:"O"`
	OU           *string `json:"OU"`
	ST           *string `json:"ST"`
	Aggregated   string  `json:"aggregated"`
	EmailAddress *string `json:"emailAddress"`
}

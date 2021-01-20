package templates

type PKCS12Template struct {
	FilePath    string
	IsEncrypted bool
}

func (t PKCS12Template) String() string {
	if t.IsEncrypted {
		return "Encrypted"
	}
	return "Unencrypted"
}

func (t PKCS12Template) Location() string {
	return t.FilePath
}

func (t PKCS12Template) MarshalBinary() (data []byte, err error) {
	panic("implement me")
}

func (t PKCS12Template) UnmarshalBinary(data []byte) error {
	panic("implement me")
}

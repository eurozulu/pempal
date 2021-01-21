package templates

import "strings"

type PKCS12Template struct {
	FilePath    string
	IsEncrypted bool
}

func (t PKCS12Template) String() string {
	e := "Unencrypted"
	if t.IsEncrypted {
		e = "Encrypted"
	}
	return strings.Join([]string{TemplateType(&t), e, t.Location()}, "\t")
}

func (t PKCS12Template) Location() string {
	return t.FilePath
}

func (t PKCS12Template) MarshalBinary() (_ []byte, err error) {
	panic("implement me")
}

func (t PKCS12Template) UnmarshalBinary(_ []byte) error {
	panic("implement me")
}

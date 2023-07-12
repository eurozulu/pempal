package identity

type Issuer interface {
	Identity() Identity
	Key() Key
	Certificate() Certificate
}

type issuer struct {
	key  Key
	cert Certificate
}

func (is issuer) Identity() Identity {
	if is.cert == nil {
		return ""
	}
	return is.cert.Identity()
}

func (is issuer) Key() Key {
	return is.key
}

func (is issuer) Certificate() Certificate {
	return is.cert
}

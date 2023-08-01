package identity

type User interface {
	Identity() Identity
	Key() Key
	Certificate() Certificate
}

type user struct {
	key  Key
	cert Certificate
}

func (is user) Identity() Identity {
	if is.cert == nil {
		return ""
	}
	return is.cert.Identity()
}

func (is user) Key() Key {
	return is.key
}

func (is user) Certificate() Certificate {
	return is.cert
}

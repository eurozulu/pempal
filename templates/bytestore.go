package templates

type ByteStore interface {
	Names() []string
	Contains(name string) bool
	Read(name string) ([]byte, error)
	Write(name string, blob []byte) error
	Delete(name string) error
}

func removeDuplicates(ss []string) []string {
	m := map[string]bool{}
	var found []string
	for _, s := range ss {
		if m[s] {
			continue
		}
		found = append(found, s)
	}
	return found
}

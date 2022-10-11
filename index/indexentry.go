package index

type IndexEntry interface {
	Key() string
	Value() interface{}
}

type indexEntry struct {
	k string
	v interface{}
}

func (i indexEntry) Key() string {
	return i.k
}

func (i indexEntry) Value() interface{} {
	return i.v
}

func NewIndexEntry(key string, value interface{}) IndexEntry {
	return &indexEntry{
		k: key,
		v: value,
	}
}

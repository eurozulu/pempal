package templates

import (
	"os"
)

type memoryStore struct {
	store map[string][]byte
}

func (m memoryStore) Names() []string {
	var names []string
	for n := range m.store {
		names = append(names, n)
	}
	return names
}

func (m memoryStore) Contains(name string) bool {
	_, ok := m.store[name]
	return ok
}

func (m memoryStore) Read(name string) ([]byte, error) {
	by, ok := m.store[name]
	if !ok {
		return nil, os.ErrNotExist
	}
	return by, nil
}

func (m *memoryStore) Write(name string, blob []byte) error {
	if m.Contains(name) {
		return os.ErrExist
	}
	m.store[name] = blob
	return nil
}

func (m *memoryStore) Delete(name string) error {
	if !m.Contains(name) {
		return os.ErrNotExist
	}
	delete(m.store, name)
	return nil
}

func newMemoryStore() ByteStore {
	return &memoryStore{store: map[string][]byte{}}
}

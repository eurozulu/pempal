package identity

import (
	"bytes"
	"context"
	"fmt"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/resourceio"
	"github.com/eurozulu/pempal/resources"
	"strings"
)

type Keys interface {
	AllKeys(ctx context.Context) <-chan Key
	KeyByIdentity(id string) (Key, error)
	KeysByName(name string) ([]Key, error)
}

type keys struct {
	NonRecursive bool
	Keypath      []string
}

func (ks keys) AllKeys(ctx context.Context) <-chan Key {
	ch := make(chan Key)
	go func() {
		for loc := range resourceio.NewResourceScanner(!ks.NonRecursive).Scan(ctx, ks.Keypath...) {
			k, err := NewKey(loc.Location(), locationKeysAsPem(loc))
			if err != nil {
				logger.Error("failed to read key at %s  %v", loc.Location(), err)
				continue
			}
			select {
			case <-ctx.Done():
				return
			case ch <- k:
			}
		}
	}()
	return ch
}

func (ks keys) KeyByIdentity(id string) (Key, error) {
	if IsIdentity(id) {
		id = Identity(id).String()
	}
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	for k := range ks.AllKeys(ctx) {
		if Identity(k.String()).String() == id {
			return k, nil
		}
	}
	return nil, fmt.Errorf("no key found with id %s", id)
}

func (ks keys) KeysByName(name string) ([]Key, error) {
	var found []Key
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	for k := range ks.AllKeys(ctx) {
		if strings.Contains(k.Location(), name) {
			found = append(found, k)
		}
	}
	if len(found) == 0 {
		return nil, fmt.Errorf("no identity found with name %s", name)
	}
	return found, nil
}

func locationKeysAsPem(loc resourceio.ResourceLocation) []byte {
	keyRes := append(loc.Resources(resources.PrivateKey), loc.Resources(resources.PublicKey)...)
	buf := bytes.NewBuffer(nil)
	for _, kr := range keyRes {
		buf.WriteString(kr.String())
	}
	return buf.Bytes()
}

func NewKeys(keypath []string) Keys {
	return &keys{
		Keypath: keypath,
	}
}

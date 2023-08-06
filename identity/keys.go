package identity

import (
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
	FindKeys(ctx context.Context, s string) []Key
}

type keys struct {
	NonRecursive bool
	Keypath      []string
}

func (ks keys) AllKeys(ctx context.Context) <-chan Key {
	ch := make(chan Key)
	go func() {
		defer close(ch)
		for loc := range resourceio.NewResourceScanner(!ks.NonRecursive).Scan(ctx, ks.Keypath...) {
			keyspem := loc.ResourcesAsPem(resources.PrivateKey, resources.PublicKey)
			if len(keyspem) == 0 {
				continue
			}
			k, err := NewKey(loc.Location(), keyspem)
			if err != nil {
				logger.Warning("failed to read key at %s  %v", loc.Location(), err)
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

func (ks keys) FindKeys(ctx context.Context, s string) []Key {
	var id Identity
	if IsIdentity(s) {
		id = Identity(s)
	}
	var kez []Key
	ctx, cnl := context.WithCancel(ctx)
	defer cnl()
	for k := range ks.AllKeys(ctx) {
		kid := k.Identity()
		if kid.String() != id.String() && kid.String() != s && !strings.Contains(k.Location(), s) {
			continue
		}
		kez = append(kez, k)
	}
	return kez
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

func NewKeys(keypath []string) Keys {
	return &keys{
		Keypath: keypath,
	}
}

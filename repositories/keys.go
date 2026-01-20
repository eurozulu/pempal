package repositories

import (
	"context"
	"fmt"
	"github.com/eurozulu/pempal/logging"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/resourcefiles"
)

type Keys string

type keyFilter func(key *model.PrivateKey) bool

func (kz Keys) ByPublicKey(puk *model.PublicKey) (*model.PrivateKey, error) {
	ps := puk.Fingerprint().String()
	return kz.FindFirst(func(k *model.PrivateKey) bool {
		return k.Public().Fingerprint().Equals(ps)
	})
}

func (kz Keys) ByFingerPrint(fingerprint string) (*model.PrivateKey, error) {
	return kz.FindFirst(func(key *model.PrivateKey) bool {
		return key.Fingerprint().Equals(fingerprint)
	})
}

func (kz Keys) MatchByFingerPrint(fingerprint string) []*model.PrivateKey {
	return kz.FindAll(func(key *model.PrivateKey) bool {
		return key.Fingerprint().Match(fingerprint)
	})
}
func (kz Keys) MatchByPublicKeyFingerPrint(fingerprint string) []*model.PrivateKey {
	return kz.FindAll(func(key *model.PrivateKey) bool {
		return key.Public().Fingerprint().Match(fingerprint)
	})
}

func (kz Keys) MatchByAnyFingerPrint(fingerprint string) []*model.PrivateKey {
	keyz := kz.MatchByPublicKeyFingerPrint(fingerprint)
	if len(keyz) == 0 {
		keyz = kz.MatchByFingerPrint(fingerprint)
	}
	return keyz
}

func (kz Keys) FindFirst(filter keyFilter) (*model.PrivateKey, error) {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	found, ok := <-kz.Find(ctx, filter)
	if !ok {
		return nil, fmt.Errorf("no keys found")
	}
	return found, nil
}

func (kz Keys) FindAll(filter keyFilter) []*model.PrivateKey {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	var found []*model.PrivateKey
	for key := range kz.Find(ctx, filter) {
		found = append(found, key)
	}
	return found
}

func (kz Keys) Find(ctx context.Context, filter keyFilter) <-chan *model.PrivateKey {
	found := make(chan *model.PrivateKey)
	go func() {
		defer close(found)
		keyFiles := resourcefiles.PemFiles(kz).FindByType(ctx, model.ResourceTypePrivateKey)
		for pf := range keyFiles {
			for _, blk := range pf.Blocks {
				k, err := model.NewPrivateKeyFromPem(blk)
				if err != nil {
					logging.Warning("could not parse pem in %s %v", pf.Path, err)
					continue
				}
				if filter != nil && !filter(k) {
					continue
				}
				select {
				case <-ctx.Done():
				case found <- k:
				}
			}
		}
	}()
	return found
}

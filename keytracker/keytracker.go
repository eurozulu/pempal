package keytracker

import (
	"context"
	"encoding/pem"
	"pempal/keytools"
	"pempal/pemreader"
	"sync"
)

type KeyTracker struct {
	ShowLogs  bool
	// IgnoreAnonymous ignores encrypted prvate keys that are not paired with a public key
	IgnoreAnonymous bool
	Recursive bool
}

func (kt KeyTracker) FindKeys(ctx context.Context, rootPaths ...string) <-chan Key {
	ch := make(chan Key)
	go func(ch chan<- Key, paths []string) {
		defer close(ch)
		var wg sync.WaitGroup
		for _, p := range rootPaths {
			wg.Add(1)
			// scan each path independently, all feeding into shared output channel
			go kt.findKeys(ctx, p, ch, &wg)
		}
		wg.Wait()
	}(ch, rootPaths)
	return ch
}

func (kt KeyTracker) findKeys(ctx context.Context, rootpath string, keyCh chan<- Key, wg *sync.WaitGroup) {
	defer wg.Done()
	pr := &pemreader.PemScanner{
		AddLocationHeader: true,
		PemTypes:          keytools.PrivateKeyTypes,
		Recursive:         kt.Recursive,
	}
	chIn := pr.Find(ctx, rootpath)
	collect := &collector{
		anons: map[string]*key{},
		puks:  map[string]*pem.Block{},
	}
pemFeedLoop:
	for {
		select {
		case <-ctx.Done():
			return
		case blk, ok := <-chIn:
			if !ok {
				break pemFeedLoop
			}

			k := collect.AddBlock(blk)
			if k == nil {
				continue
			}
			// here key has a public key, so send it out.
			select {
			case <-ctx.Done():
				return
			case keyCh <- k:
			}
		}
	}
	// send out the anonymous keys
	if !kt.IgnoreAnonymous {
		for _, v := range collect.UnknownKeys() {
			select {
			case <-ctx.Done():
				return
			case keyCh <- v:
			}
		}
	}
}

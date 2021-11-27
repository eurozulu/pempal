package keytracker

import (
	"context"
	"encoding/pem"
	"fmt"
	"pempal/keytools"
	"pempal/pemreader"
	"reflect"
	"sync"
)

var keyTypes = keytools.CombineMaps(keytools.PrivateKeyTypes, keytools.PublicKeyTypes)
var keysAndCertTypes = keytools.CombineMaps(keyTypes, keytools.CertificateTypes)

// KeyTracker tracks all the private keys and attempts to match them to any certificates with the same public key.
// Matched Certificates and keys are paired into 'Identities', one cert/key per ID.
type KeyTracker struct {
	ShowLogs  bool
	Recursive bool
}

// FindIdentities tracks all the private keys with known public key and matches them to any certificate with the same public key.
// The resulting channel contains all the matched certificates, linked to the key which signed them.
func (it KeyTracker) FindIdentities(ctx context.Context, rootPaths ...string) <-chan Identity {
	ch := make(chan Identity)
	go func() {
		defer close(ch)
		chPEms := it.findPems(ctx, keysAndCertTypes, rootPaths...)
		matcher := newIDMatcher()
		for {
			select {
			case <-ctx.Done():
				return
			case blk, ok := <-chPEms:
				if !ok {
					return
				}
				ids, err := matcher.AddPem(blk)
				if !it.handleError(err) || len(ids) == 0 {
					continue
				}
				for _, id := range ids {
					select {
					case <-ctx.Done():
						return
					case ch <- id:
					}
				}
			}
		}
	}()
	return ch
}

func (it KeyTracker) FindKeys(ctx context.Context, rootPaths ...string) <-chan Key {
	ch := make(chan Key)
	go func() {
		defer close(ch)
		chPEms := it.findPems(ctx, keyTypes, rootPaths...)
		matcher := newKeyMatcher()
		for {
			select {
			case <-ctx.Done():
				return
			case blk, ok := <-chPEms:
				if !ok {
					return
				}
				k := matcher.AddPem(blk)
				if k == nil || reflect.ValueOf(k).IsNil() {
					continue
				}
				select {
				case <-ctx.Done():
					return
				case ch <- k:
				}
			}
		}
	}()
	return ch
}

// findPems combnines all the pem blocks found in all the given root paths, into a single channel of pems.
// each root path is scaned in its own routine.  The results are unordered
func (it KeyTracker) findPems(ctx context.Context, pemType map[string]bool, rootPaths ...string) <-chan *pem.Block {
	ch := make(chan *pem.Block)
	go func(chPems chan<- *pem.Block, paths []string) {
		defer close(chPems)
		var wg sync.WaitGroup
		for _, p := range rootPaths {
			wg.Add(1)
			// scan each path independently, all feeding into shared output channel
			go func(rootpath string, pems chan<- *pem.Block, wg *sync.WaitGroup) {
				defer wg.Done()
				ps := &pemreader.PemScanner{
					AddLocationHeader: true,
					PemTypes:          pemType,
					Recursive:         it.Recursive,
				}
				chp := ps.Find(ctx, rootpath)
				for {
					select {
					case <-ctx.Done():
						return
					case blk, ok := <-chp:
						if !ok {
							return
						}
						select {
						case <-ctx.Done():
							return
						case pems <- blk:
						}
					}
				}
			}(p, chPems, &wg)
		}
		wg.Wait()
	}(ch, rootPaths)
	return ch
}

func (it KeyTracker) handleError(err error) bool {
	if err == nil {
		return true
	}
	if it.ShowLogs {
		fmt.Println(err)
	}
	return false
}

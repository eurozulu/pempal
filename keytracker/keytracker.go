package keytracker

import (
	"context"
	"encoding/pem"
	"fmt"
	"pempal/keytools"
	"pempal/pemreader"
	"sync"
)

// KeyTracker tracks all the private keys and attempts to match them to any certificates with the same public key.
// Matched Certificates and keys are paired into 'Identities', one cert/key per ID.
type KeyTracker struct {
	ShowLogs bool
	// IgnoreAnonymous ignores encrypted prvate keys that are not paired with a public key
	IgnoreAnonymous bool
	Recursive       bool
}

func (it KeyTracker) FindIdentities(ctx context.Context, rootPaths ...string) <-chan Identity {
	ch := make(chan Identity)
	go func(ch chan<- Identity, paths []string) {
		defer close(ch)
		kt := &KeyScanner{
			ShowLogs:        it.ShowLogs,
			IgnoreAnonymous: false,
			Recursive:       it.Recursive,
		}
		chKeys := kt.FindKeys(ctx, rootPaths...)

		chCerts := make(chan *pem.Block)
		go func() {
			defer close(chCerts)
			var wg sync.WaitGroup
			for _, p := range rootPaths {
				wg.Add(1)
				// scan each path independently, all feeding into shared output channel
				go it.findCertificates(ctx, p, chCerts, &wg)
			}
			wg.Wait()
		}()

		collect := newIDCollector()
		// Complete when both channels have completed.
		openChannels := 2

		for {
			var ids []Identity

			select {
			case <-ctx.Done():
				return
			case blk, ok := <-chCerts:
				if !ok {
					openChannels--
					if openChannels == 0 {
						return
					}
				}
				id, err := collect.AddCertificate(blk)
				if !it.handleError(err) {
					continue
				}
				ids = append(ids, id)

			case k, ok := <-chKeys:
				if !ok {
					openChannels--
					if openChannels == 0 {
						return
					}
				}
				ids = append(ids, collect.AddKey(k)...)
			}

			// Matched ID(s) found, push to channel
			for _, id := range ids {
				select {
				case <-ctx.Done():
					return
				case ch <- id:
				}
			}
		}
	}(ch, rootPaths)
	return ch
}

func (it KeyTracker) findCertificates(ctx context.Context, rootpath string, certs chan<- *pem.Block, wg *sync.WaitGroup) {
	defer wg.Done()
	ps := &pemreader.PemScanner{
		AddLocationHeader: true,
		PemTypes:          keytools.CertificateTypes,
		Recursive:         it.Recursive,
	}
	chCerts := ps.Find(ctx, rootpath)
	for {
		select {
		case <-ctx.Done():
			return
		case blk, ok := <-chCerts:
			if !ok {
				return
			}
			select {
			case <-ctx.Done():
				return
			case certs <- blk:
			}
		}
	}
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

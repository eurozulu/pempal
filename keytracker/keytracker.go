package keytracker

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"pempal/keytools"
	"pempal/pemreader"
	"strings"
	"sync"
)

type KeyTracker struct {
	ShowLogs bool
}

func (kt KeyTracker) FindIdentities(ctx context.Context, rootPaths ...string) <-chan []Identity {
	ch := make(chan []Identity, 10)
	go func(ch chan<- []Identity, paths []string) {
		defer close(ch)
		pr := &pemreader.PemReader{
			Verbose:           true,
			AddLocationHeader: true,
			PemTypes:          certsAndKeys,
			Recursive:         true,
		}
		var wg sync.WaitGroup
		wg.Add(len(paths))
		for _, p := range paths {
			pemCh := pr.Find(ctx, p)
			go kt.collectPems(ctx, pemCh, ch, &wg)
		}
		fmt.Println("==================================================== starting wait")
		wg.Wait()
		fmt.Println("==================================================== wait complete")

	}(ch, rootPaths)
	return ch
}

func (kt KeyTracker) FindKeys(ctx context.Context, rootPaths ...string) <-chan *key {
	ch := make(chan *key, 10)
	go func(ch chan<- *key, paths []string) {
		defer close(ch)
		pr := &pemreader.PemReader{
			AddLocationHeader: true,
			PemTypes:          keytools.PrivateKeyTypes,
			Recursive:         true,
		}
		var wg sync.WaitGroup

		for _, p := range rootPaths {
			wg.Add(1)
			// scan each path independantly, all feeding into shared output channel
			go func(p string, chOut chan<- *key, wg *sync.WaitGroup) {
				defer wg.Done()
				chIn := pr.Find(ctx, p)
				for {
					select {
					case <-ctx.Done():
						return
					case blk, ok := <-chIn:
						if !ok {
							return
						}
						select {
						case <-ctx.Done():
							return
						case chOut <- &key{pemBlock: blk}:
						}
					}
				}
			}(p, ch, &wg)
		}
		wg.Wait()
	}(ch, rootPaths)
	return ch
}

// Issuers returns all the issuer ID's available.
// Issuers are Identitess with at least one certificate bound to it with the 'KeyUsageCertSign' key usage.
// Resulting Identity only contain the relevant certificates for issuing.
func (kt KeyTracker) Issuers(ctx context.Context, dn string, keypath []string) <-chan Identity {
	ch := make(chan Identity)

	go func(ch chan<- Identity) {
		defer close(ch)
		idCh := kt.FindIdentities(ctx, keypath...)
		for {
			select {
			case <-ctx.Done():
				return
			case ids, ok := <-idCh:
				if !ok {
					return
				}
				for _, id := range ids {
					certs := certsContainDN(dn, id.Certificates(x509.KeyUsageCertSign))
					if len(certs) == 0 {
						continue
					}

					// found ID with matching cert(s).  Reform ID to just the relevant certificates
					id = &identity{
						prk:  id.(*identity).prk,
						blks: certsToPemBlocks(certs),
					}
					select {
					case <-ctx.Done():
						return
					case ch <- id:
					}
				}
			}
		}
	}(ch)
	return ch
}

func (kt KeyTracker) collectPems(ctx context.Context, chIn <-chan *pem.Block, chOut chan<- []Identity, wg *sync.WaitGroup) {
	defer wg.Done()
	collect := newCollector()
outerLoop:
	for {
		select {
		case <-ctx.Done():
			return
		case blk, ok := <-chIn:
			if !ok {
				break outerLoop
			}
			fmt.Printf("adding %s block to %v\n", blk.Type, collect)
			if err := collect.AddBlock(blk); err != nil && kt.ShowLogs {
				log.Println(err)
			}
		}
	}
	select {
	case <-ctx.Done():
		return
	case chOut <- collect.Identities():
	}
}

func certsToPemBlocks(certs []*x509.Certificate) []*pem.Block {
	pbs := make([]*pem.Block, len(certs))
	for i, c := range certs {
		pbs[i] = &pem.Block{
			Type:  keytools.PEM_CERTIFICATE,
			Bytes: c.Raw,
		}
	}
	return pbs
}

func certsContainDN(dn string, certs []*x509.Certificate) []*x509.Certificate {
	if dn == "" {
		return certs
	}
	var found []*x509.Certificate
	dn = strings.ToLower(dn)
	for _, c := range certs {
		if !strings.Contains(strings.ToLower(c.Subject.String()), dn) {
			continue
		}
		found = append(found, c)
	}
	return found
}

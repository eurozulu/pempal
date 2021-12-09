package pemresources

import (
	"context"
	"log"
	"pempal/fileformats"
)

// Keys locates the private keys and matches them to their public key counterpart
type Keys struct {
	NonRecursive  bool
	ShowLogs      bool
	HideAnonymous bool
}

// Find locates the key with the given keyid.
func (ks Keys) Find(ctx context.Context, keyids []string, keypath ...string) <-chan *PrivateKey {
	ch := make(chan *PrivateKey)
	go func(ch chan<- *PrivateKey) {
		defer close(ch)
		// make subcontext so we can bail out when all are found
		sctx, cnl := context.WithCancel(ctx)
		defer cnl()

		keyCh := ks.ScanKeys(sctx, keypath...)
		for len(keyids) > 0 {
			select {
			case <-ctx.Done():
				return

			case k, ok := <-keyCh:
				if !ok {
					return
				}

				i := indexOf(k.PublicKeyHash, keyids)
				if i < 0 {
					continue
				}
				// found id, remove it from the ids slice
				if i == len(keyids)-1 {
					// if last one, just trim it off
					keyids = keyids[:i]
				} else {
					// in the middle, shift later ones down on top of it.
					keyids = append(keyids[:i], keyids[i+1:]...)
				}
				select {
				case <-ctx.Done():
					return
				case ch <- k:
				}
			}
		}
	}(ch)
	return ch
}

// ScanKeys returns all the keys found in the given keypath
// All matched keys are returned first, followed by any unmatched key, depending on the HideAnonymous flag
func (ks Keys) ScanKeys(ctx context.Context, keypath ...string) <-chan *PrivateKey {
	ch := make(chan *PrivateKey)
	go func(ch chan<- *PrivateKey) {
		defer close(ch)

		pemCh := PemScanner{
			Recursive: !ks.NonRecursive,
			Verbose:   ks.ShowLogs,
			Reader:    fileformats.NewFormatReader(),
			TypeFilter: fileformats.CombineMaps(fileformats.PemTypesPrivateKey,
				fileformats.PemTypesPublicKey,
				fileformats.PemTypesCertificate,
				fileformats.PemTypesCertificateRequest),
		}.Scan(ctx, keypath...)

		kt := newKeyTracker()
	outerLoop:
		for {
			select {
			case <-ctx.Done():
				return

			case blk, ok := <-pemCh:
				if !ok {
					break outerLoop
				}
				k, err := kt.AddBlock(blk)
				if !ks.handleError(err) || k == nil {
					continue
				}
				select {
				case <-ctx.Done():
					return
				case ch <- k:
				}
			}
		}
		// finally, send out the unmatched private keys
		if ks.HideAnonymous {
			return
		}
		for _, k := range kt.AnonymousKeys() {
			select {
			case <-ctx.Done():
				return
			case ch <- k:
			}
		}

	}(ch)
	return ch
}

func (ks Keys) handleError(err error) bool {
	if err == nil {
		return true
	}
	if ks.ShowLogs {
		log.Println(err)
	}
	return false
}
func indexOf(s string, ss []string) int {
	for i, sz := range ss {
		if s == sz {
			return i
		}
	}
	return -1
}

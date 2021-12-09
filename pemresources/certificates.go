package pemresources

import (
	"context"
	"pempal/fileformats"
	"pempal/pemscanner"
)

// Certificates locates certificates and matches them to the privste key which signed them
type Certificates struct {
	KeyPath   []string
	ShowLogs  bool
	Recursive bool
}

// Identity is a single certificate matched to its key
type Identity struct {
	Key  *PrivateKey
	Cert *Certificate
}

func (cs Certificates) Find(ctx context.Context, keyids []string, certpath ...string) <-chan *Identity {
	ch := make(chan *Identity)
	go func(ch chan<- *Identity) {
		defer close(ch)
		// make subcontext so we can bail out when all are found
		sctx, cnl := context.WithCancel(ctx)
		defer cnl()

		idCh := cs.AllCertificates(sctx, certpath...)
		for len(keyids) > 0 {
			select {
			case <-ctx.Done():
				return

			case id, ok := <-idCh:
				if !ok {
					return
				}

				i := stringIndex(id.Cert.PublicKeyHash, keyids)
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
				case ch <- id:
				}
			}
		}
	}(ch)
	return ch
}

func (cs Certificates) AllCertificates(ctx context.Context, certpath ...string) <-chan *Identity {
	ch := make(chan *Identity)
	go func(ch chan<- *Identity) {
		defer close(ch)
		pp := pemscanner.NewPemParser(fileformats.CombineMaps(fileformats.PemTypesPrivateKey, fileformats.PemTypesPublicKey, fileformats.PemTypesCertificate, fileformats.PemTypesCertificateRequest))
		pp.AddLocationHeader = true
		pp.ShowLog = cs.ShowLogs
		pp.Recursive = cs.Recursive
		blkCh := pp.FindAll(ctx, certpath...)

		cc := newCertificateCollector()
	outerLoop:
		for {
			select {
			case <-ctx.Done():
				return

			case blks, ok := <-blkCh:
				if !ok {
					break outerLoop
				}

				ids, err := cc.AddBlocks(blks)
				if err != nil {

				}
				for _, k := range ids {
					select {
					case <-ctx.Done():
						return
					case ch <- k:
					}
				}
			}
		}

	}(ch)
	return ch
}

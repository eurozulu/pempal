package pemresources

import (
	"context"
	"encoding/pem"
	"io/ioutil"
	"log"
	"pempal/fileformats"
	"pempal/filepathscanner"

	"sync"
)

// PemScanner scans one or more locations for resources which can be formatted into pem blocks.
type PemScanner struct {
	// Recursive, when true, scans sub directories of locations (Assuming they're directories, otherwise its ignored)
	Recursive bool
	// Verbose outputs additional logging information
	Verbose bool
	// Reader should be set to the FormatReader(s) to use for parsing resources. See NewFormatReader
	Reader fileformats.FormatReader

	// TypeFilter is a pem type filter, to limit output to only the PEM types in the filter map.
	// when set, only PEM types in the map are returned, when not set all (parsable) pems are returned.
	TypeFilter map[string]bool
}

// Scan scans all of the given root paths simultaneously returning any found pems in no particular order.
func (rs PemScanner) Scan(ctx context.Context, rootpaths ...string) <-chan *pem.Block {
	ch := make(chan *pem.Block)
	go func(ch chan<- *pem.Block) {
		defer close(ch)
		fps := filepathscanner.FilePathScanner{
			Recursive: rs.Recursive,
			Verbose:   rs.Verbose,
			ExtFilter: filepathscanner.X509FileTypes,
		}
		ps := fps.Scan(ctx, rootpaths...)
		var wg sync.WaitGroup
		for {
			select {
			case <-ctx.Done():
				return

			case p, ok := <-ps:
				if !ok {
					wg.Wait()
					return
				}
				wg.Add(1)
				go rs.parseFile(ctx, ch, &wg, p)
			}
		}
	}(ch)
	return ch
}
func (rs PemScanner) parseFile(ctx context.Context, ch chan<- *pem.Block, wg *sync.WaitGroup, p string) {
	defer wg.Done()
	by, err := ioutil.ReadFile(p)
	if !rs.handleError(err) {
		return
	}
	blks, err := rs.Reader.Unmarshal(by)
	if !rs.handleError(err) {
		return
	}
	hasFilter := len(rs.TypeFilter) > 0
	for _, blk := range blks {
		if hasFilter && !rs.TypeFilter[blk.Type] {
			continue
		}
		select {
		case <-ctx.Done():
			return
		case ch <- blk:
		}
	}
}

func (rs PemScanner) handleError(err error) bool {
	if err == nil {
		return true
	}
	if rs.Verbose {
		log.Println(err)
	}
	return false
}

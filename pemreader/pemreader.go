// pemscan scans a list of paths looking for known file fileformats which can be read into PEM blocks.
package pemreader

import (
	"context"
	"fmt"
	"log"
	"path"
	"pempal/pemreader/fileformats"
	"strconv"
	"strings"
	"sync"

	"encoding/pem"
	"io/ioutil"
)

const LocationHeaderKey = "location"
const PemBlockBufferSize = 64
const maxOpenFiles = 64

var pemExtensionMap = extensionMapFromFormats()

var CustomExtensions = map[string]bool{}

var PemTypes = map[string]bool{}

// PemReader scans a single file location and, optionally, if its a directory, all its sub directories, for Pem Blocks.
type PemReader struct {
	// When true, displays errors incurred during scan. By default errors are ignored
	Verbose bool

	// AddLocationHeader adds the file location into the PEM header in the form of "location": <filepath>[:index]
	// where the suffix ':' is present when a file contains more than one PEM block.
	// The number following the ':' indicates the index of the block within that file.
	// Single block files omit the index altogether.
	AddLocationHeader bool

	// Recursive, when true, will scan sub directories of directories.  By default subdirectories are ignored.
	Recursive bool

	// Map of pemtypes to filter.  key should be std PEM type string (see PEMType))
	// Value should be true to include the resource in the result, false to exclude it.
	// If a type is not listed it will be excluded.
	// If the map is empty/nil (the default) ALL types are included
	PemTypes map[string]bool
}

func (p PemReader) Find(ctx context.Context, rootPath string) <-chan *pem.Block {
	ch := make(chan *pem.Block)

	go func(rootPath string, chOut chan<- *pem.Block) {
		defer close(chOut)

		fs := &FileScanner{
			Filter:       completeFileFormatsMap(),
			NonRecursive: !p.Recursive,
		}
		chFilePaths := fs.Find(ctx, rootPath)
		wg := &sync.WaitGroup{}
	outerLoop:
		for {
			select {
			case <-ctx.Done():
				// context concelled, return without waiting
				return
			case fp, ok := <-chFilePaths:
				if !ok {
					break outerLoop
				}
				wg.Add(1)
				go p.parseResource(ctx, fp, chOut, wg)
			}
		}
		// wait for parsers to complete before returning
		wg.Wait()

	}(rootPath, ch)
	return ch
}

// parseResource attempts to read the given filepath as a PEM, using its file extension as the Format.
// If successfully parsed into Pem(s), it is passed to the given channel as pem blocks
// If unsuccessful, the file is ignored
func (p PemReader) parseResource(ctx context.Context, fp string, outCh chan<- *pem.Block, wg *sync.WaitGroup) {
	defer wg.Done()
	by, err := ioutil.ReadFile(fp)
	if !p.handleError(err) {
		return
	}
	// Find a FileFormat based on file extension
	ext := strings.ToLower(strings.TrimLeft(path.Ext(fp), "."))
	f := fileformats.FileFormats[ext]
	// If not a known format, use the default format (Which attemps to try all the others, slowly!)
	if f == nil {
		f = fileformats.FileFormats[""]
	}
	blks, err := f.Format(by)
	if err != nil {
		err = fmt.Errorf("%s %w", fp, err)
	}
	if !p.handleError(err) || len(blks) == 0 {
		return
	}
	// Filter found blocks to with PEMType filter map
	blks = p.filterBlockTypes(blks)
	if len(blks) == 0 {
		// nothing found, ignore it
		return
	}

	if p.AddLocationHeader {
		blks = p.addLocationHeader(fp, blks)
	}
	for _, blk := range blks {
		select {
		case <-ctx.Done():
			return
		case outCh <- blk:
		}
	}
}

// filterBlockTypes returns only those given blocks with a type valid in the PemTypes maps.
func (p PemReader) filterBlockTypes(blks []*pem.Block) []*pem.Block {
	// If no types stated then do not filter any, return ALL blocks
	if len(p.PemTypes) == 0 {
		return blks
	}
	var fbs []*pem.Block
	for _, b := range blks {
		if !p.PemTypes[b.Type] {
			continue
		}
		fbs = append(fbs, b)
	}
	return fbs
}

// addLocationHeader adds the given filepath to the headers of all the given pem blocks under the LocationHeaderKey.
// If blocks length > 1, then an index as appended to the filepath for each index in the slice
func (p PemReader) addLocationHeader(location string, blks []*pem.Block) []*pem.Block {
	showIndex := len(blks) > 1
	for i, blk := range blks {
		lc := location
		if showIndex {
			lc = strings.Join([]string{lc, strconv.Itoa(i + 1)}, ":")
		}
		if blk.Headers == nil {
			blk.Headers = map[string]string{}
		}
		blk.Headers[LocationHeaderKey] = lc
	}
	return blks
}

func (p PemReader) handleError(err error) bool {
	if err == nil {
		return true
	}
	if p.Verbose {
		log.Println(err)
	}
	return false
}

// AddCustomFileTypes adds the given list of file extensions to the Custom scan map.
// File extensions may be excluded by preceeding the extension with a "!".
// If not preceeded with the !, the extension will be included in the search
func AddCustomFileTypes(types []string) {
	for _, ft := range types {
		v := !strings.HasPrefix(ft, "!")
		if !v {
			ft = ft[1:]
		}
		CustomExtensions[ft] = v
	}
}

// AddPemTypes adds new file extensions to view as a PEM type.
// each type element shoudk be the file extension EXcluding any preceeding "."
// The type may have a preceeding "!" to negate the extension and exclude it as a pem type.
func AddPemTypes(types []string) {
	for _, pt := range types {
		v := !strings.HasPrefix(pt, "!")
		if !v {
			// trim off !
			pt = pt[1:]
		}
		pt = strings.TrimLeft(pt, ".")
		PemTypes[strings.ToUpper(pt)] = v
	}
}

func extensionMapFromFormats() ExtensionFilter {
	em := ExtensionFilter{}
	for k := range fileformats.FileFormats {
		em[k] = true
	}
	return em
}

func completeFileFormatsMap() ExtensionFilter {
	em := ExtensionFilter{}
	for k, v := range pemExtensionMap {
		em[k] = v
	}
	for k, v := range CustomExtensions {
		em[k] = v
	}
	return em
}

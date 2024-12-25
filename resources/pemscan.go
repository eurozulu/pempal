package resources

import (
	"context"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/logging"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/utils"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

const workerCount = 6
const defaultFileExtensions = "pem,der,crt,cer,key,pub,prv"

type PemScan struct {
	FileExtensions []string
	ResourceTypes  []model.ResourceType
	Parsers        []PemParser
	NonRecursive   bool
}

func (sc PemScan) ScanPath(ctx context.Context, path ...string) <-chan PemResource {
	ch := make(chan PemResource)
	go func(result chan<- PemResource, paths []string) {
		defer close(result)
		pathFeed := make(chan string, workerCount)
		var allDone sync.WaitGroup
		allDone.Add(workerCount)

		// kick off the workers
		for i := 0; i < workerCount; i++ {
			go sc.parseFiles(pathFeed, result, &allDone)
		}

		usedPaths := make(map[string]bool)
		for _, p := range paths {
			ap, err := filepath.Abs(p)
			if err != nil {
				logging.Error("PemScan", "Failed to find path %s. %v", p, err)
				continue
			}
			if utils.FileExists(ap) {
				pathFeed <- ap
				continue
			}
			if !utils.DirExists(ap) {
				logging.Warning("PemScan", "ignoring path %q as not found", ap)
				continue
			}
			if err := sc.walkPath(ctx, pathFeed, usedPaths, ap); err != nil {
				logging.Error("PemScan", "scan path %s: %v", path, err)
			}
		}
		close(pathFeed)
		allDone.Wait() // wait for workers to complete
	}(ch, path)

	return ch
}

func (sc PemScan) walkPath(ctx context.Context, feed chan<- string, usedPaths map[string]bool, root string) error {
	return fs.WalkDir(os.DirFS(root), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == "." {
			return nil
		}
		// ignore hidden stuff
		if strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if sc.NonRecursive && d.IsDir() {
			return fs.SkipDir
		}

		// filter by extension if extentions set
		if len(sc.FileExtensions) > 0 && !sc.hasExtension(d.Name()) {
			return nil
		}
		p := filepath.Join(root, path)
		if usedPaths[p] {
			logging.Debug("PemScan", "ignoring path %q as duplicated", p)
			return nil
		}
		usedPaths[p] = true
		select {
		case <-ctx.Done():
			return ctx.Err()
		case feed <- p:
		}
		return nil
	})
}

// parseFiles is the worker thread, accepting path strings and parsing them into pemresources
// runs until paths is closed.
func (sc PemScan) parseFiles(paths <-chan string, result chan<- PemResource, done *sync.WaitGroup) {
	defer done.Done()
	for path := range paths {
		pr, err := sc.parseFile(path)
		if err != nil {
			log.Printf("Error reading %s: %v", path, err)
			continue
		}
		if pr == nil || len(pr.Content) == 0 {
			continue
		}
		result <- *pr
	}
}

// parseFile parsed a single file into PEm blocks from the given path.
func (sc PemScan) parseFile(path string) (*PemResource, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var pemz []*pem.Block
	for _, parser := range sc.Parsers {
		if !parser.CanParse(data) {
			continue
		}

		pemz, err = parser.FormatAsPem(data)
		if err != nil {
			return nil, err
		}
		pemz = sc.filterByResourceType(pemz)
		if len(pemz) == 0 {
			return nil, nil
		}
		return &PemResource{
			Path:    path,
			Content: pemz,
		}, nil
	}
	return nil, fmt.Errorf("unknown format of %s", path)
}

func (sc PemScan) filterByResourceType(pemz []*pem.Block) []*pem.Block {
	if len(sc.ResourceTypes) == 0 {
		return pemz
	}
	var found []*pem.Block
	for _, blk := range pemz {
		brt := model.ParseResourceTypeFromPEMType(blk.Type)
		if brt == model.UnknownResourceType {
			logging.Error("PemScan", "unknown pem type %q", blk.Type)
			continue
		}
		if !model.ContainsResourceType(sc.ResourceTypes, brt) {
			continue
		}
		found = append(found, blk)
	}
	return found
}

func (sc PemScan) hasExtension(name string) bool {
	return slices.Contains(sc.FileExtensions, strings.TrimLeft(filepath.Ext(name), "."))
}

func NewPemScan(types ...model.ResourceType) *PemScan {
	return &PemScan{
		FileExtensions: strings.Split(defaultFileExtensions, ","),
		ResourceTypes:  types,
		Parsers:        []PemParser{&PemResourceParser{}},
	}
}

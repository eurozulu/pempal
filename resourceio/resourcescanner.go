package resourceio

import (
	"context"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/model"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ResourceScanner scans one or more locations for x509 resources.
type ResourceScanner interface {
	Recursive() bool
	Scan(ctx context.Context, paths ...string) <-chan ResourceLocation
}

type resourceScanner struct {
	recursive bool
}

func (f resourceScanner) Recursive() bool {
	return f.recursive
}

func (f resourceScanner) Scan(ctx context.Context, paths ...string) <-chan ResourceLocation {
	ch := make(chan ResourceLocation)
	go func(paths []string) {
		defer close(ch)
		var wg sync.WaitGroup
		wg.Add(len(paths))
		for _, p := range paths {
			go f.scan(ctx, p, ch, &wg)
		}
		wg.Wait()
	}(paths)
	return ch
}

func (f resourceScanner) scan(ctx context.Context, root string, out chan<- ResourceLocation, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != root && (!f.recursive || strings.HasPrefix(d.Name(), ".")) {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasPrefix(d.Name(), ".") {
			return nil
		}
		// todo: add file extension filter
		loc, err := f.parseLocation(path)
		if err != nil {
			logger.Error("Failed to parse location %s  %v", path, err)
			return nil
		}
		if loc != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- loc:
				logger.Debug("pushed location %s with %d resources", loc.Location(), len(loc.Resources()))
			}
		}
		return nil

	}); err != nil {
		logger.Error("Failed to scan directory %s  %v", root, err)
	}
}

func (f resourceScanner) parseLocation(path string) (ResourceLocation, error) {
	logger.Debug("reading location %s", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var found []model.Resource
	for _, rp := range ResourceParsers {
		if !rp.CanParse(data) {
			continue
		}
		res, err := rp.ParseResources(data)
		if err != nil {
			logger.Warning("resource parsing error %s as %v", path, rp)
			continue
		}
		found = append(found, res...)
		break
	}
	if len(found) == 0 {
		logger.Debug("no resources found in %s", path)
		return nil, nil
	}
	logger.Debug("found %d resources in location %s", len(found), path)
	return NewResourceLocation(path, found), nil
}

func NewResourceScanner(recursive bool) ResourceScanner {
	return &resourceScanner{
		recursive: recursive,
	}
}

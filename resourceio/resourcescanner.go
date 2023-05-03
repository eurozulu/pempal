package resourceio

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"pempal/logger"
	"pempal/model"
	"strings"
	"sync"
)

var ResourceFileExtensions = map[string]bool{
	"":     true,
	"pem":  true,
	"cer":  true,
	"cert": true,
	"csr":  true,
	"key":  true,
	"pub":  true,
	"prk":  true,
	"req":  true,
}

type ResourceScanner interface {
	Scan(ctx context.Context, paths ...string) <-chan *ResourceLocation
}

type resourceScanner struct {
	recursive bool
	filters   []LocationFilter
}

func (r resourceScanner) Scan(ctx context.Context, paths ...string) <-chan *ResourceLocation {
	ch := make(chan *ResourceLocation)
	go func(paths []string) {
		defer close(ch)
		var wg sync.WaitGroup
		for _, p := range paths {
			wg.Add(1)
			go r.scan(ctx, p, ch, &wg)
		}
		wg.Wait()
	}(paths)
	return ch
}

func (r resourceScanner) scan(ctx context.Context, root string, out chan<- *ResourceLocation, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if d.IsDir() {
			if path == root {
				return nil
			}
			if !r.recursive || strings.HasPrefix(filepath.Base(path), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}

		ext := strings.ToLower(strings.TrimLeft(filepath.Ext(d.Name()), "."))
		if !ResourceFileExtensions[ext] {
			logger.Log(logger.Debug, "Skipping file %s as no file extension match", path)
			return nil
		}
		loc, err := r.parseResourceLocation(path)
		if err != nil {
			logger.Log(logger.Debug, "Skipping file %s as failed to parse as resources  %v", path, err)
			return nil
		}
		if len(loc.Resources) == 0 {
			logger.Log(logger.Debug, "Skipping file %s as no matching resources found", path)
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- loc:
		}
		return nil

	}); err != nil {
		logger.Log(logger.Error, "reading directory %s failed:  %v", root, err)
	}
}

func (r resourceScanner) parseResourceLocation(path string) (*ResourceLocation, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	res, err := model.ParseResources(data)
	if err != nil {
		return nil, err
	}
	loc := ResourceLocation{
		Path:      path,
		Resources: res,
	}
	loc.Resources = r.filterLocation(loc)
	return &loc, nil
}

func (r resourceScanner) filterLocation(loc ResourceLocation) []model.PEMResource {
	for _, f := range r.filters {
		loc.Resources = f.Filter(loc)
		if len(loc.Resources) == 0 {
			break
		}
	}
	return loc.Resources
}

func NewResourceScanner(recursive bool, filters ...LocationFilter) ResourceScanner {
	return &resourceScanner{
		recursive: recursive,
		filters:   filters,
	}
}

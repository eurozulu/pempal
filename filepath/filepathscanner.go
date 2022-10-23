package filepath

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type FilepathScanner interface {
	Find(ctx context.Context, locations ...string) <-chan string
}

type FileFilter interface {
	Accept(path string) bool
}

type DirFilter interface {
	AcceptDir(path string) bool
}

type filepathScanner struct {
	filter FileFilter
}

func (r filepathScanner) Find(ctx context.Context, locations ...string) <-chan string {
	ch := make(chan string)
	go func(ch chan<- string) {
		defer close(ch)
		for _, p := range locations {
			if ctx.Err() != nil {
				return
			}
			err := r.scanDirectory(ctx, p, ch)
			if err != nil {
				if err != filepath.SkipDir {
					log.Printf("%v\n", err)
				}
			}
		}
	}(ch)
	return ch
}

func (r filepathScanner) scanDirectory(ctx context.Context, p string, out chan<- string) error {
	return filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			if r.filter == nil {
				return fs.SkipDir
			}
			df, ok := r.filter.(DirFilter)
			if !ok || !df.AcceptDir(path) {
				return fs.SkipDir
			}
			return nil
		}

		if r.filter != nil && !r.filter.Accept(path) {
			return nil
		}
		select {
		case <-ctx.Done():
			return nil
		case out <- path:
		}
		return nil
	})
}

func NewFilepathScanner(path []string, filter FileFilter) FilepathScanner {
	var index int
	for _, p := range path {
		cp := filepath.Clean(os.ExpandEnv(p))
		if _, err := os.Stat(cp); err != nil {
			log.Printf("ignoring location %s, as %v", p, err)
			continue
		}
		path[index] = cp
		index++
	}
	return &filepathScanner{
		filter: filter,
	}
}

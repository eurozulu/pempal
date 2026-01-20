package resourcefiles

import (
	"context"
	"github.com/eurozulu/pempal/logging"
	"github.com/eurozulu/pempal/tools"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type FilePathScanner struct {
	Filter FilePathFilter
}

func (s *FilePathScanner) ScanPath(ctx context.Context, fsys fs.FS, path ...string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for _, p := range path {
			root := cleanPathExists(fsys, p)
			if root == "" {
				logging.Warning("path does not exist: %s", p)
				continue
			}
			if err := s.walkPath(ctx, fsys, root, out); err != nil {
				logging.Error("Failed to walk path %s: %v", root, err)
			}
		}
	}()
	return out
}
func (s *FilePathScanner) walkPath(ctx context.Context, fsys fs.FS, root string, out chan string) error {
	return fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// ignore hidden files & dirs
		if path != root && strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if !s.filterPath(path) {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- path:
		}
		return nil
	})
}

func (s FilePathScanner) filterPath(path string) bool {
	if s.Filter == nil {
		return true
	}
	return s.Filter.Filter(path)
}

// ensure path exist
func cleanPathExists(fsys fs.FS, path string) string {
	if path == "" {
		return path
	}
	path = filepath.Clean(os.ExpandEnv(path))
	if !tools.IsPathExistsFs(fsys, path) {
		path = ""
	}
	return path
}

type FilePathFilter interface {
	Filter(path string) bool
}

type FileExtensionFilter struct {
	Extensions []string
}

func (f FileExtensionFilter) Filter(path string) bool {
	if len(f.Extensions) == 0 {
		return true
	}
	for _, e := range f.Extensions {
		if strings.EqualFold(e, filepath.Ext(path)) {
			return true
		}
	}
	return false
}
func NewFileExtensionFilter(extensions ...string) *FileExtensionFilter {
	return &FileExtensionFilter{Extensions: extensions}
}

package finder

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"pempal/pemtypes"
	"strings"
	"sync"
)

const ENV_PATH = "PP_PATH"

var PPPath = os.ExpandEnv(os.Getenv(ENV_PATH))

type FileFilter interface {
	Accept(path string, d fs.DirEntry) bool
}

type Finder interface {
	Find(ctx context.Context, path ...string) (<-chan PemLocation, error)
}

type finder struct {
	parser    PemParser
	recursive bool
	filter    FileFilter

	scanHidden   bool
	reportErrors bool
}

func (f finder) Find(ctx context.Context, path ...string) (<-chan PemLocation, error) {
	cp, err := cleanPath(path)
	if err != nil {
		return nil, err
	}
	ch := make(chan PemLocation)
	go func() {
		defer close(ch)
		var wg sync.WaitGroup
		wg.Add(len(cp))
		for _, p := range cp {
			go f.searchPath(ctx, p, &wg, ch)
		}
		wg.Wait()
	}()
	return ch, nil
}

func (fd finder) searchPath(ctx context.Context, path string, wg *sync.WaitGroup, ch chan<- PemLocation) {
	defer wg.Done()
	hasFilter := fd.filter != nil

	err := filepath.WalkDir(path, func(entryPath string, d fs.DirEntry, err error) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if d == nil {
			// invalid file
			return fmt.Errorf("invalid entry in path! %s", entryPath)
		}
		// avoid 'hidden' (filenames preceeded with a dot)
		if !fd.scanHidden && entryPath != path && strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			// if not 'root path' skip it
			if !fd.recursive && entryPath != path {
				return fs.SkipDir
			}
			return nil
		}

		// skip loading files which don't comply
		if hasFilter && !fd.filter.Accept(entryPath, d) {
			return nil
		}

		data, err := os.ReadFile(entryPath)
		if err != nil {
			if fd.reportErrors {
				fmt.Fprintln(os.Stderr, "failed to read %s  %v", entryPath, err)
			}
			return nil
		}

		res, err := fd.parser.Parse(data)
		if err != nil {
			if fd.reportErrors {
				fmt.Fprintln(os.Stderr, "failed to read %s  %v", entryPath, err)
			}
			return nil
		}
		if len(res) == 0 {
			return nil
		}
		// Finally has something to shout about!
		// post the entry path with the resources found
		select {
		case <-ctx.Done():
			return nil
		case ch <- &pemLocation{
			location:  entryPath,
			resources: res,
		}:
		}
		return nil
	})
	if fd.reportErrors && err != nil {
		fmt.Fprintln(os.Stderr, "%v", err)
	}
}

func CleanPPPath() ([]string, error) {
	return cleanPath(strings.Split(PPPath, ":"))
}

func cleanPath(path []string) ([]string, error) {
	var paths []string
	for _, s := range path {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		// Ensure it exists
		fi, err := os.Stat(s)
		if err != nil {
			return nil, fmt.Errorf("invalid path %s  %v", s, err)
		}
		// Check if trailing slash missing from dir
		if fi.IsDir() && !strings.HasSuffix(s, "/") {
			s = strings.Join([]string{s, "/"}, "")
		}
		paths = append(paths, s)
	}
	if len(path) == 0 {
		return nil, fmt.Errorf("no search path found")
	}
	return paths, nil
}

func NewFinder(filter FileFilter, recursive, reportErrors bool, pemtype ...pemtypes.PEMType) Finder {
	return &finder{
		parser:       &pemParser{types: pemtype},
		recursive:    recursive,
		reportErrors: reportErrors,
		filter:       filter,
	}
}

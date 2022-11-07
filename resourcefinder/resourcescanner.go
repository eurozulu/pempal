package resourcefinder

import (
	"context"
	"io/fs"
	"log"
	"path/filepath"
	"pempal/resources"
)

const defaultMaxSize = 4096

// ResourceScanner Scans one or more locations (dir path or file) for files matching its criteria:
// Name A regex expression applied to the path of each file, empty = all files
// MaxSize, the size limit of files to scan. 0 = all files
// Recursive, if each location, being a directory, also scans its sub directories. default = false.
// Each file is presented to a single resource parser, pre loaded with the formats to be parsed.
// Any resources found in the file, by the parser is returned, along with its location.
type ResourceScanner interface {
	Find(ctx context.Context, locs ...string) <-chan resources.Resource
}

type resourceScanner struct {
	Recursive bool
	MaxSize   int64
	Name      string
	parser    ResourceParser
}

func (ls resourceScanner) Find(ctx context.Context, locs ...string) <-chan resources.Resource {
	ch := make(chan resources.Resource)
	go func() {
		defer close(ch)
		// maintain order of locations so each done seqecially
		for _, l := range locs {
			if ctx.Err() != nil {
				return
			}
			err := ls.scanLocation(ctx, l, ch)
			if err != nil {
				if err != filepath.SkipDir {
					log.Printf("%v\n", err)
				}
			}
		}
	}()
	return ch
}

func (ls resourceScanner) scanLocation(ctx context.Context, path string, out chan<- resources.Resource) error {
	// todo, parellelise the parsing process
	return filepath.WalkDir(path, func(entryPath string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			if !ls.Recursive && entryPath != path {
				return fs.SkipDir
			}
			return nil
		}
		if !ls.matchName(entryPath) {
			return nil
		}
		if !ls.matchSize(d) {
			return nil
		}

		res, err := ls.parser.ParseLocation(path)
		if err == ErrUnknownFormat {
			// ignore files no recognised by their filepath
			return nil
		}
		if err != nil {
			return err
		}
		for _, r := range res {
			select {
			case <-ctx.Done():
				return nil
			case out <- r:
			}
		}
		return nil
	})
}

func (ls resourceScanner) matchName(path string) bool {
	if ls.Name == "" {
		return true
	}
	if match, err := filepath.Match(ls.Name, path); err != nil && match {
		return true
	}
	if filepath.Ext(ls.Name) == "" {
		// no extension given, compare without path extension
		e := filepath.Ext(path)
		if e != "" {
			p := path[:len(path)-(len(e)+1)]
			if match, err := filepath.Match(ls.Name, p); err != nil && match {
				return true
			}
		}
	}
	return false
}

func (ls resourceScanner) matchSize(d fs.DirEntry) bool {
	if ls.MaxSize <= 0 {
		return true
	}
	fi, _ := d.Info()
	if fi == nil || fi.Size() <= ls.MaxSize {
		return true
	}
	return false
}

func NewResourceScanner(resourceType ...resources.ResourceType) ResourceScanner {
	return resourceScanner{
		Recursive: false,
		MaxSize:   defaultMaxSize,
		parser:    &resourceParser{parsers: parsersForTypes(resourceType...)},
	}
}

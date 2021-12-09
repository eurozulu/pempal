package filepathscanner

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

// FilePathScanner scans one or more file locations for filenames with a specific set of file extensions
// A Location may be a file or a directory path.  For a file, just the file path is returned (if it has the correct extension)
// For a directory, all the files in the directory with the correct extensions are returned, excluding any sub directories.
// Set the Recursive flag true to scan the subdirectories for files.
// locations may be absolute, relative to the current working directory or include "~" as a Home directory indicator, or any Environment variables such as $HOME etc
// Only the full path of files (matching the extension filter) are returned, directory paths are never returned.
// Scanner ignores non regular files (pipes, symlinks) and any location beginning with a dot 'hidden' files.
// However, hidden files or directories explicitly given as a location will be scanned, but any hidden subdirectories or files within that location will still be ignored, regardless of Recursive.
// i.e. To include any hidden location, it must be explicitly given as a location.  This usually applies to directories which do not contain hidden resources, such as often found in the home directory or config directories.
type FilePathScanner struct {
	// Recursive flag controls if sub directories are searched
	Recursive bool

	// Verbose displays error logs found during the scan.  By defqult, these are ignored.
	Verbose bool

	// ExtFilter is an optional way to filter each filename by its filename extension.
	// use the map keys for the extensions to find, WITHOUT any leading dot.  and a 'true' value.
	// e.g. { "xml": true, "json": true, "yaml": true }
	// If empty or nil, ALL files (except 'hidden' files) are returned.
	ExtFilter map[string]bool
}

// Scan finds all the file paths found in the given locations.
// locations may be a single file or a directory.  When its a directory, all the files found in that directory are returned in the channel.
// Each location path is scanned independantly, returning filepaths as they're found so no order can be assumed from the items returned in the channel
func (fps FilePathScanner) Scan(ctx context.Context, locations ...string) <-chan string {
	ch := make(chan string)
	go func(ch chan<- string) {
		defer close(ch)
		var wg sync.WaitGroup
		for _, l := range locations {
			wg.Add(1)
			go func(l string, wg *sync.WaitGroup) {
				chIn := fps.scan(ctx, l)
				for {
					select {
					case <-ctx.Done():
						return
					case fp, ok := <-chIn:
						if !ok {
							return
						}
						select {
						case <-ctx.Done():
							return
						case ch <- fp:
						}
					}
				}
			}(l, &wg)
		}
		wg.Wait()
	}(ch)
	return ch
}

// scan scans a single location for filepaths.
func (fps FilePathScanner) scan(ctx context.Context, location string) <-chan string {
	ch := make(chan string)
	go func(location string, ch chan<- string) {
		defer close(ch)

		err := filepath.WalkDir(location,
			func(fp string, d fs.DirEntry, err error) error {
				if !fps.handleError(err) {
					return nil
				}
				// Ignore 'hidden' files/dirs (Unless explicitly specified as the location)
				if fp != location && strings.HasPrefix(d.Name(), ".") {
					if d.IsDir() {
						return fs.SkipDir
					}
					return nil
				}
				// dir, skip if nonrecursive
				if d.IsDir() {
					if !fps.Recursive && fp != location {
						return fs.SkipDir
					}
					return nil
				}
				// avoid the weirdos like symlinks and pipes
				if !d.Type().IsRegular() {
					return nil
				}

				// filter by filename extension
				if len(fps.ExtFilter) > 0 &&
					!fps.ExtFilter[strings.ToLower(strings.TrimLeft(path.Ext(fp), "."))] {
					return nil
				}
				// Its passed the test, post it into the channel
				select {
				case <-ctx.Done():
					return ctx.Err()
				case ch <- fp:
				}
				return nil
			})
		// show any err (except a context cancellation)
		if err != ctx.Err() {
			fps.handleError(err)
		}
	}(resolvePath(location), ch)
	return ch
}

func (fps FilePathScanner) handleError(err error) bool {
	if err == nil {
		return true
	}
	if fps.Verbose {
		log.Println(err)
	}
	return false
}

func resolvePath(p string) string {
	if p == "." || strings.HasPrefix(p, "./") || strings.HasPrefix(p, "../") {
		p = strings.Join([]string{"$PWD", p}, "/")
	}
	if strings.Contains(p, "~") {
		p = strings.Replace(p, "~", "$HOME", -1)
	}
	return path.Clean(os.ExpandEnv(p))
}

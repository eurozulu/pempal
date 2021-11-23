package pemreader

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// FileScanner scans file locations for file pathnames.
// optionally filtered by the filename or file extension
type FileScanner struct {
	// Default is aa recursive search.  Set true to scan just one file/directory
	NonRecursive bool

	// Verbose displays error logs found during the scan.  By defqult, these are ignored.
	Verbose bool

	// Filter is an optional way to filter each filename.
	Filter FilePathFilter

	// FollowSymLinks will attempt to resolve symlinks to the actual location.  Off by default.
	FollowSymLinks bool
}

// FilePathFilter filters a given filepath
type FilePathFilter interface {
	Filter(p string) bool
}

type NameFilter struct {
	Name string
}

func (n NameFilter) Filter(p string) bool {
	return strings.EqualFold(path.Base(p), n.Name)
}

type ExtensionFilter map[string]bool

func (e ExtensionFilter) Filter(p string) bool {
	return e[strings.TrimLeft(path.Ext(p), ".")]
}

// Find finds all the file paths found in the given root.
// root may be a single file or a directory.  When its a directory, all the files found in that directory are returned in the channel.
func (p FileScanner) Find(ctx context.Context, root string) <-chan string {
	ch := make(chan string)
	go func(root string, ch chan<- string) {
		defer close(ch)

		err := filepath.WalkDir(resolvePath(root),
			func(fp string, d fs.DirEntry, err error) error {
				if !p.handleError(err) {
					return nil
				}

				// Ignore 'hidden' files/dirs (Unless specified as the root)
				if strings.HasPrefix(d.Name(), ".") && fp != root {
					if d.IsDir() {
						return fs.SkipDir
					}
					return nil
				}
				if p.FollowSymLinks {
					fp = p.resolveSymLink(d.Type(), fp)
					if fp == "" {
						// dead link, ignore it
						return nil
					}
				}

				// dir, skip if nonrecursive
				if d.IsDir() {
					if p.NonRecursive && fp != root {
						return fs.SkipDir
					}
					return nil
				}

				if !d.Type().IsRegular() {
					return nil
				}

				if p.Filter != nil && !p.Filter.Filter(fp) {
					return nil
				}

				select {
				case <-ctx.Done():
					return ctx.Err()
				case ch <- fp:
				}
				return nil
			})
		if err != nil && err != ctx.Err() {
			log.Println(err)
		}
	}(root, ch)
	return ch
}

func (p FileScanner) handleError(err error) bool {
	if err == nil {
		return true
	}
	if p.Verbose {
		log.Println(err)
	}
	return false
}

func (p FileScanner) resolveSymLink(mode os.FileMode, fp string) string {
	for mode&os.ModeSymlink == os.ModeSymlink {
		rp, err := filepath.EvalSymlinks(fp)
		if !p.handleError(err) {
			return ""
		}
		fi, err := os.Stat(rp)
		if !p.handleError(err) {
			return ""
		}
		fp = rp
		mode = fi.Mode()
	}
	return fp
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

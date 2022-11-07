package finder

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// FileScanner scans a given filepath for finder.
// It has an optional filter, which filters files by filepath, to decide if the file should be read.
type FileScanner interface {
	Scan(ctx context.Context, filter FilepathFilter, path string) <-chan Location
}

type FilepathFilter interface {
	MatchPath(path string) bool
}

type fileLocation struct {
	path string
	data []byte
	err  error
}

type fileScanner struct {
	recursive  bool
	scanHidden bool
}

func (fd fileScanner) Scan(ctx context.Context, filter FilepathFilter, path string) <-chan Location {
	ch := make(chan Location)
	go func(ch chan<- Location) {
		if path == "-" || path == "" {
			fd.scanStandardInput(ctx, ch)
		} else {
			fd.scanFilePath(ctx, filter, path, ch)
		}
	}(ch)
	return ch
}

func (fd fileScanner) scanFilePath(ctx context.Context, filter FilepathFilter, path string, out chan<- Location) {
	defer close(out)
	var err error
	// Check if trailing slash missing from dir
	if !strings.HasSuffix(path, "/") {
		fi, e := os.Stat(path)
		if e == nil && fi.IsDir() {
			path = strings.Join([]string{path, "/"}, "")
		}
	}

	hasFilter := filter != nil
	err = filepath.WalkDir(path, func(entryPath string, d fs.DirEntry, err error) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if d == nil {
			// invalid file
			return os.ErrNotExist
		}
		if !fd.scanHidden && entryPath != path && strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			if !fd.recursive && entryPath != path {
				return fs.SkipDir
			}
			return nil
		}

		if hasFilter && !filter.MatchPath(entryPath) {
			return nil
		}

		var fl *fileLocation
		data, err := os.ReadFile(entryPath)
		if err != nil {
			fl = &fileLocation{
				path: entryPath,
				data: nil,
				err:  fmt.Errorf("failed to read %s  %v", entryPath, err),
			}
		} else {
			fl = &fileLocation{
				path: entryPath,
				data: data,
				err:  nil,
			}
		}
		select {
		case <-ctx.Done():
			return nil
		case out <- fl:
		}
		return nil
	})
	if err != nil {
		select {
		case <-ctx.Done():
			return
		case out <- &fileLocation{
			path: path,
			err:  err,
		}:
		}
	}
}

func (fd fileScanner) scanStandardInput(ctx context.Context, out chan<- Location) {
	defer close(out)
	by, err := readStandardInput()
	if err != nil {
		if _, err := fmt.Fprintf(os.Stderr, "error reading standard input %v", err); err != nil {
			log.Println(err)
		}
		return
	}
	if len(by) == 0 {
		// nothing to see here, move along
		return
	}

	select {
	case <-ctx.Done():
		return
	case out <- &fileLocation{
		path: "-",
		data: by,
	}:
	}
}

func (fl fileLocation) Path() string {
	return fl.path
}

func (fl fileLocation) MarshalText() (text []byte, err error) {
	return []byte(fl.path), nil
}

func readStandardInput() ([]byte, error) {
	st, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}
	if st.Size() == 0 {
		return nil, nil
	}
	buf := bytes.NewBuffer(nil)
	_, err = buf.ReadFrom(os.Stdin)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func PathExists(l string) bool {
	_, err := os.Stat(l)
	if err != nil {
		return false
	}
	return true
}

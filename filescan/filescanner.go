package filescan

import (
	"io/ioutil"
	"path"
	"strings"
)

// NewFileScanner creates a new file scanner with the given root path of a directory to scan.
func NewFileScanner(dirname string) (*FileScanner, error) {
	files, err := directoryFiles(dirname)
	if err != nil {
		if !strings.HasSuffix(err.Error(), "not a directory") {
			return nil, err
		}
		// If filename (not dir) given, use just this in the files.
		files = []string{dirname}
	}
	return &FileScanner{
		files: files,
		index: -1,
	}, nil
}

type FileFilter func(p string) bool

// FileScanner scans a directory structure reporting each file
type FileScanner struct {
	Recursive bool
	Filter    FileFilter

	index int
	files []string
}

// FileName returns the full path of the current file
func (fs FileScanner) FileName() string {
	if fs.HasComplete() || fs.index < 0 {
		return ""
	}
	return fs.files[fs.index]
}

// Bytes returns the bytes read from the current file
func (fs FileScanner) Bytes() ([]byte, error) {
	p := fs.FileName()
	if p == "" {
		return nil, nil
	}
	return ioutil.ReadFile(p)
}

// HasComplete indicates of the scan has reached the end of the directory
func (fs FileScanner) HasComplete() bool {
	return fs.index >= len(fs.files)
}

// Scan will check for the next file. Return true if a file was found or false if no further files found.
func (fs *FileScanner) Scan() bool {
	if fs.HasComplete() {
		return false
	}

	fs.index++
	for i := fs.index; i < len(fs.files); i++ {
		if strings.HasSuffix(fs.files[i], "/") {
			if !fs.Recursive {
				continue
			}
			// insert sub directory files in place of directory entry
			df, err := directoryFiles(fs.files[i])
			if err != nil {
				continue
			}
			if len(df) == 0 {
				continue
			}
			if i+1 < len(fs.files) {
				df = append(df, fs.files[i+1:]...)
			}
			fs.files = df
			i = -1
			continue
		}

		if fs.Filter != nil && !fs.Filter(fs.files[i]) {
			continue
		}
		fs.index = i
		break
	}
	return !fs.HasComplete()
}

func directoryFiles(dirname string) ([]string, error) {
	fis, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	files := make([]string, len(fis))
	for i, fi := range fis {
		p := path.Join(dirname, fi.Name())
		if fi.IsDir() {
			p = string(append([]rune(p), '/'))

		} else if !fi.Mode().IsRegular() {
			continue
		}
		files[i] = p
	}
	return files, nil
}

package templates

import (
	"io/fs"
	"os"
	"time"
)

// LayerFs combines multiple FS locations into a single FileSystem.
// Layers are combined such that the entries for any given directory are listed
// for all the layers which contain that directory name.
// For the root, the combinatation of all of the layers root directories is presented.
// When File entries share the same path on two or more layers, the newest file is returned.
type LayerFs []fs.FS

func (l LayerFs) ReadDir(name string) ([]fs.DirEntry, error) {
	var entries []fs.DirEntry
	unique := map[string]fs.DirEntry{}

	for _, layer := range l {
		des, err := fs.ReadDir(layer, name)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		for _, de := range des {
			if ee, ok := unique[de.Name()]; ok {
				// existing entry already known
				ne, err := getNewestEntry(de, ee)
				if err != nil {
					return nil, err
				}
				if ne != ee {
					// this is newer than known so replace
					entries = replaceEntry(entries, ee, de)
				}
				continue
			}
			unique[de.Name()] = de
			entries = append(entries, de)
		}
	}
	if len(entries) == 0 {
		return nil, fs.ErrNotExist
	}
	return entries, nil
}

func (l LayerFs) Stat(name string) (fs.FileInfo, error) {
	if name == "." {
		return &rootFileInfo{}, nil
	}
	layer, err := l.findFSForName(name)
	if err != nil {
		return nil, err
	}
	return fs.Stat(layer, name)
}

func (l LayerFs) Open(name string) (fs.File, error) {
	layer, err := l.findFSForName(name)
	if err != nil {
		return nil, err
	}
	return layer.Open(name)
}

func (l LayerFs) findFSForName(name string) (fs.FS, error) {
	index := -1
	var found fs.FileInfo
	for i, layer := range l {
		fi, err := fs.Stat(layer, name)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		if found == nil || fi.ModTime().After(found.ModTime()) {
			index = i
			found = fi
		}
	}
	if index == -1 {
		return nil, fs.ErrNotExist
	}
	return l[index], nil
}

func replaceEntry(entries []fs.DirEntry, old, new fs.DirEntry) []fs.DirEntry {
	for i, e := range entries {
		if e != old {
			continue
		}
		end := []fs.DirEntry{new}
		if i+1 < len(entries) {
			end = append(end, entries[i+1:]...)
		}
		entries = append(entries[:i], end...)
	}
	return entries
}

func getNewestEntry(e, ee fs.DirEntry) (fs.DirEntry, error) {
	fie, err := e.Info()
	if err != nil {
		return nil, err
	}
	fiee, err := e.Info()
	if err != nil {
		return nil, err
	}
	if fie.ModTime().After(fiee.ModTime()) {
		return e, nil
	}
	return ee, nil

}

type rootFileInfo struct {
}

func (r rootFileInfo) Name() string {
	return ""
}
func (r rootFileInfo) Size() int64 {
	return 0
}
func (r rootFileInfo) Mode() fs.FileMode {
	return 0777
}
func (r rootFileInfo) ModTime() time.Time {
	return time.Time{}
}
func (r rootFileInfo) IsDir() bool {
	return true
}
func (r rootFileInfo) Sys() any {
	return nil
}

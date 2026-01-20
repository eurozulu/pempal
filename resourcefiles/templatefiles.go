package resourcefiles

import (
	"context"
	"embed"
	"fmt"
	"github.com/eurozulu/pempal/logging"
	"github.com/eurozulu/pempal/model"
	"io/fs"
	"os"
	"strings"
)

//go:embed embedded/*
var embeddedTemplates embed.FS

var TemplateFileExtensions = []string{
	".yml", ".yaml", ".yam", ".template",
}

// TemplateFiles represents all the available template files.
// It uses its value as a single path to the root of the template directory
// In addition, it adds the embedded templates, templates provided by the allplication.
// When file system templates have the same relative name as an embedded template,
// the file system template will take precedence, masking/replacing the embedded template.
type TemplateFiles string

type TemplateFileFilter func(file *model.TemplateFile) bool

func (fz TemplateFiles) Find(ctx context.Context, filter TemplateFileFilter) <-chan *model.TemplateFile {
	found := make(chan *model.TemplateFile)
	go func() {
		defer close(found)

		// Read embedded templates first and push only after all the files are done.
		// This eliminates the embedded templates which are being overridden by file templates.

		scanner := &FilePathScanner{Filter: &FileExtensionFilter{Extensions: TemplateFileExtensions}}
		var embedded []*model.TemplateFile
		for path := range scanner.ScanPath(ctx, embeddedTemplates, "embedded/") {
			file, err := readTemplateFile(embeddedTemplates, path)
			if err != nil {
				logging.Warning("failed to open template %v", err)
				continue
			}
			if filter != nil && !filter(file) {
				continue
			}
			embedded = append(embedded, file)
		}

		tempFS := os.DirFS(string(fz))
		for path := range scanner.ScanPath(ctx, tempFS, ".") {
			embedded = removePath(embedded, path)
			file, err := readTemplateFile(tempFS, path)
			if err != nil {
				logging.Warning("failed to open %s  %v", path, err)
				continue
			}
			if filter != nil && !filter(file) {
				continue
			}
			select {
			case <-ctx.Done():
				return
			case found <- file:
			}
		}
		// post any remaining embedded
		for _, et := range embedded {
			select {
			case <-ctx.Done():
				return
			case found <- et:
			}
		}
	}()
	return found
}

func removePath(files []*model.TemplateFile, path string) []*model.TemplateFile {
	var found []*model.TemplateFile
	for _, f := range files {
		if strings.EqualFold(f.Path, path) {
			continue
		}
		found = append(found, f)
	}
	return found
}

func readTemplateFile(fsys fs.FS, path string) (*model.TemplateFile, error) {
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s  %v", path, err)
	}
	file := &model.TemplateFile{
		Path: path,
		Data: data,
	}
	if err := file.IsValid(); err != nil {
		return nil, fmt.Errorf("file %s not a valid template  %v", path, err)
	}
	return file, nil
}

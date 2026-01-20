package resourcefiles

import (
	"context"
	"encoding/pem"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/logging"
	"github.com/eurozulu/pempal/model"
	"os"
	"slices"
	"strings"
)

type PemFiles string

type PemFileFilter func(file *model.PemFile) *model.PemFile

func (p PemFiles) FindByType(ctx context.Context, pemtype ...model.ResourceType) <-chan *model.PemFile {
	return p.Find(ctx, func(file *model.PemFile) *model.PemFile {
		if len(pemtype) > 0 {
			file.Blocks = filterPemsByType(file.Blocks, pemtype...)
		}
		return file
	})
}

func (p PemFiles) Find(ctx context.Context, filter ...PemFileFilter) <-chan *model.PemFile {
	files := make(chan *model.PemFile)
	go func() {
		defer close(files)
		hasFilter := len(filter) > 0 && filter[0] != nil
		scanner := &FilePathScanner{Filter: FileExtensionFilter{Extensions: config.FileExtensions()}}
		path := strings.Split(string(p), string(os.PathListSeparator))
		pemFilez := scanner.ScanPath(ctx, os.DirFS("."), path...)
		for pPAth := range pemFilez {
			data, err := os.ReadFile(pPAth)
			if err != nil {
				logging.Warning("failed to open %s  %v", pPAth, err)
				continue
			}
			blks, err := FormatAsPems(data)
			if err != nil {
				logging.Error("failed to parse %s  %v", pPAth, err)
				continue
			}
			if len(blks) == 0 {
				continue
			}
			file := &model.PemFile{
				Path:   pPAth,
				Blocks: blks,
			}
			if hasFilter {
				for _, flt := range filter {
					file = flt(file)
					if file == nil || len(file.Blocks) == 0 {
						file = nil
						break
					}
				}
				if file == nil {
					continue
				}
			}

			select {
			case <-ctx.Done():
				return
			case files <- file:
			}
		}
	}()
	return files
}

func filterPemsByType(pems []*pem.Block, types ...model.ResourceType) []*pem.Block {
	var filtered []*pem.Block
	for _, b := range pems {
		if !slices.Contains(types, model.ParseResourceType(b.Type)) {
			continue
		}
		filtered = append(filtered, b)
	}
	return filtered
}

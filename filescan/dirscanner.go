package filescan

import (
	"context"
	"fmt"
	"github.com/eurozulu/pempal/encoding"
	"github.com/eurozulu/pempal/templates"
	"log"
	"sync"
)

type DirectoryScanner struct {
	Recursive bool
	PrintErrors bool
}

func (ds DirectoryScanner) ScanDirectories(ctx context.Context, args []string) <-chan []templates.Template {
	ch := make(chan []templates.Template)
	var wg sync.WaitGroup
	wg.Add(len(args))
	go func(wg *sync.WaitGroup) {
		defer close(ch)
		wg.Wait()
	}(&wg)
	// Scan each argument for files to be parsed as templates
	for _, p := range args {
		go ds.ScanDirectory(ctx, p, ch, &wg)
	}
	return ch
}

func (ds DirectoryScanner) ScanDirectory(ctx context.Context, p string, ch chan<- []templates.Template, wg *sync.WaitGroup) {
	defer wg.Done()
	fs, err := NewFileScanner(p)
	if err != nil {
		ds.showError(err)
		return
	}
	fs.Recursive = ds.Recursive
	for fs.Scan() {
		if ctx.Err() != nil {
			return
		}
		by, err := fs.Bytes()
		if err != nil {
			ds.showError(err)
			continue
		}

		tpls, err := encoding.ParseTemplates(fs.FileName(), by)
		if err != nil {
			ds.showError(fmt.Errorf("%s  %v", fs.FileName(), err))
			continue
		}
		select {
		case <-ctx.Done():
			return
		case ch <- tpls:
		}
	}
}

func (ds DirectoryScanner) showError(err error) {
	if ds.PrintErrors {
		log.Println(err)
	}
}
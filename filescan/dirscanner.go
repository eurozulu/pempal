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
	Recursive   bool
	PrintErrors bool
	Password    string
}

func (ds DirectoryScanner) ScanDirectories(ctx context.Context, args []string) <-chan []templates.Template {
	ch := make(chan []templates.Template)
	go func(ch chan<- []templates.Template) {
		defer close(ch)
		var wg sync.WaitGroup
		for _, p := range args {
			wg.Add(len(args))
			go ds.ScanDirectory(ctx, p, ch, &wg)
		}
		wg.Wait()
	}(ch)
	return ch
}

func (ds DirectoryScanner) ScanDirectory(ctx context.Context, p string, ch chan<- []templates.Template, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
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

		tpls, err := encoding.ParseTemplates(fs.FileName(), by, ds.Password)
		if err != nil {
			ds.showError(fmt.Errorf("%s  %v", fs.FileName(), err))
			continue
		}
		if len(tpls) == 0 {
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

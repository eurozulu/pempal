package finder

import (
	"context"
	"encoding"
	"log"
	"pempal/pemtypes"
	"sync"
)

type Location interface {
	encoding.TextMarshaler
	Path() string
}

type Finder interface {
	Find(ctx context.Context, path ...string) <-chan Location
}

type finder struct {
	verboseOutput bool
	scanner       FileScanner
	parser        LocationParser
	filter        LocationFilter
}

func (rs finder) Find(ctx context.Context, path ...string) <-chan Location {
	ch := make(chan Location)
	go func(ch chan<- Location) {
		var wg sync.WaitGroup
		wg.Add(len(path))
		go func(wg *sync.WaitGroup) {
			defer close(ch)
			wg.Wait()
		}(&wg)

		for _, p := range path {
			go rs.find(ctx, p, ch, &wg)
		}

	}(ch)
	return ch
}

func (rs finder) find(ctx context.Context, path string, out chan<- Location, wg *sync.WaitGroup) {
	defer wg.Done()
	hasFilter := rs.filter != nil
	for l := range rs.scanner.Scan(ctx, rs.filter, path) {
		fl := l.(*fileLocation)
		if fl.err != nil {
			if rs.verboseOutput {
				log.Println(fl.err)
			}
			continue
		}
		rl, err := rs.parser.Parse(fl.Path(), fl.data)
		if err != nil {
			if rs.verboseOutput {
				log.Println(err)
			}
			continue
		}
		if hasFilter {
			rl = rs.filter.FilterLocation(rl)
			if rl == nil {
				continue
			}
		}
		select {
		case <-ctx.Done():
			return
		case out <- rl:
		}
	}
}

func newResources(p LocationParser, recursive, verbose bool) Finder {
	return &finder{
		scanner:       &fileScanner{recursive: recursive},
		parser:        p,
		filter:        p.(LocationFilter),
		verboseOutput: verbose,
	}
}

func NewPemFinderResources(recursive, verbose bool, types ...pemtypes.PEMType) Finder {
	return newResources(newPemParser(types...), recursive, verbose)
}

func NewTransformerFinder(query PemQuery, recursive, verbose bool, types ...pemtypes.PEMType) Finder {
	return newResources(newPemTransformerParser(query, types...), recursive, verbose)
}

func NewTemplateFinder(recursive, verbose bool) Finder {
	return newResources(newTemplateParser(), recursive, verbose)
}

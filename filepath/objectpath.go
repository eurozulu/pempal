package filepath

import (
	"context"
	"log"
)

// ObjectPath wraps the filepath to parse the relevant files found into specific objects.
type ObjectPath[T any] interface {
	Path() []string
	Find(name string) T
	FindAll(ctx context.Context) <-chan T
}

type ObjectParser[T any] interface {
	Parse(p string) (T, error)
}

type ObjectMatcher[T any] interface {
	Match(name string, t T) bool
}

type objectPath[T any] struct {
	parser  ObjectParser[T]
	matcher ObjectMatcher[T]
	fPath   FilePath
	filter  FileFilter
}

func (op objectPath[T]) Path() []string {
	return op.fPath.Path()
}

func (op objectPath[T]) Find(name string) T {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	for t := range op.FindAll(ctx) {
		if op.matcher.Match(name, t) {
			return t
		}
	}
	return nil
}

func (op objectPath[T]) FindAll(ctx context.Context) <-chan T {
	ch := make(chan T)
	go func(ch chan<- T) {
		defer close(ch)
		// use subcontext to control cancel over filepath
		ctxx, cnl := context.WithCancel(ctx)
		defer cnl()

		chLocations := op.fPath.Find(ctxx, op.filter)
		for {
			select {
			case <-ctx.Done(): // parent ctx cancelling
				return
			case l, ok := <-chLocations:
				if !ok {
					return
				}
				o, err := op.parser.Parse(l)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				select {
				case <-ctx.Done(): // parent ctx cancelling
					return
				case ch <- o:
				}
			}
		}
	}(ch)
	return ch
}

func NewObjectPath[T any](filePath string, parser ObjectParser[T], matcher ObjectMatcher[T], fileExtensions ...string) ObjectPath[T] {
	var filter FileFilter
	if len(fileExtensions) > 0 {
		filter = NewFileExtFilter(fileExtensions)
	}
	return &objectPath[T]{
		parser:  parser,
		matcher: matcher,
		fPath:   NewFilePath(filePath),
		filter:  filter,
	}
}

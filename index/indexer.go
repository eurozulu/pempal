package index

import "context"

type Indexer interface {
	Index(ctx context.Context, srcPath ...string) <-chan IndexEntry
}

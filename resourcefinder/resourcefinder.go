package resourcefinder

import (
	"context"
	"pempal/resources"
)

type ResourceFinder interface {
	Find(ctx context.Context, q Query, resourceType ...resources.ResourceType) <-chan resources.Resource
}

type resourceFinder struct {
	locations []string
}

func (rf resourceFinder) Find(ctx context.Context, q Query, resourceTypes ...resources.ResourceType) <-chan resources.Resource {
	ch := make(chan resources.Resource)
	go func() {
		defer close(ch)
		emptyQuery := len(q) == 0
		rs := NewResourceScanner(resourceTypes...)

		for r := range rs.Find(ctx, rf.locations...) {
			if !emptyQuery && !rf.filterResource(q, r) {
				continue
			}
			select {
			case <-ctx.Done():
				return
			case ch <- r:
			}
		}
	}()
	return ch
}

func (rf resourceFinder) filterResource(q Query, r resources.Resource) bool {

}

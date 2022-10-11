package index

import (
	"context"
	"sync"
)

type Index interface {
	AddPath(srcPath ...string)
	Get(key string) IndexEntry
	GetAll() []IndexEntry
}

type index struct {
	indexer Indexer

	chRequests chan *indexReq
	chRefresh  chan []string
}

func (ki index) AddPath(srcPath ...string) {
	ki.chRefresh <- srcPath
}

func (ki index) Get(key string) IndexEntry {
	entries := ki.makeRequest(key)
	if len(entries) > 0 {
		return entries[0]
	}
	return nil
}

func (ki index) GetAll() []IndexEntry {
	return ki.makeRequest("")
}

func (ki index) serveRequests(ctx context.Context) {
	defer close(ki.chRequests)
	defer close(ki.chRefresh)
	chNewEntries := make(chan IndexEntry, 10)
	defer close(chNewEntries)

	cache := map[string][]IndexEntry{}
	monitor := &StateMonitor{}
	var active bool
	waiting := map[string][]*indexReq{}

	for {
		select {
		case <-ctx.Done():
			return

		case src := <-ki.chRefresh:
			ki.startIndexing(ctx, chNewEntries, monitor.SetActive(), src...)

		case e := <-chNewEntries:
			cache[e.Key()] = append(cache[e.Key()], e)
			// Check if anyone waiting for this key
			if wrs, ok := waiting[e.Key()]; ok {
				// fulfil and close the waiting requests with newly updated cache
				for _, wr := range wrs {
					_ = ki.fulfilRequestFromCache(ctx, cache, wr)
					close(wr.Response)
				}
				delete(waiting, e.Key())
			}

		case r := <-ki.chRequests:
			// request satisfied or not satisfied and index not active (no further entries expected), close request and continue
			if ki.fulfilRequestFromCache(ctx, cache, r) || !active {
				close(r.Response)
				continue
			}
			// unfulfilled but still active, place request on hold until we're no longer active (or fulfilled)
			waiting[r.Key] = append(waiting[r.Key], r)

		case state := <-monitor.State():
			active = state
			if !active && len(waiting) > 0 {
				// nolonger active, release any waiting requests
				closeWaitingRequests(waiting)
				waiting = map[string][]*indexReq{}
			}
		}
	}

}

// startIndexing kicks off the indexer with the given srcPath(s) and feeds the result into the given out channel.
// When indexer completes, this completes. (or context cancels)
func (ki index) startIndexing(ctx context.Context, out chan<- IndexEntry, wg *sync.WaitGroup, srcPath ...string) {
	entries := ki.indexer.Index(ctx, srcPath...)
	for {
		select {
		case <-ctx.Done():
			return
		case e, ok := <-entries:
			if !ok {
				return
			}
			select {
			case <-ctx.Done():
				return
			case out <- e:
			}
		}
	}
}

func (ki index) fulfilRequestFromCache(ctx context.Context, cache map[string][]IndexEntry, r *indexReq) bool {
	var entries []IndexEntry
	// check for specific key
	if v, ok := cache[r.Key]; ok {
		entries = append(entries, v...)
	}
	// empty key collects all entries
	if r.Key == "" {
		for _, v := range cache {
			entries = append(entries, v...)
		}
		return true
	}
	if len(entries) == 0 {
		return false
	}
	for _, v := range entries {
		select {
		case <-ctx.Done():
			return false
		case r.Response <- v:
		}
	}
	return true
}

func (ki index) makeRequest(key string) []IndexEntry {
	r := &indexReq{
		Key:      key,
		Response: make(chan IndexEntry, 1),
	}
	ki.chRequests <- r
	var keys []IndexEntry
	for k := range r.Response {
		keys = append(keys, k)
	}
	return keys
}

type indexReq struct {
	Key      string
	Response chan IndexEntry
}

// sendWaitingRequests sends the given key to the slice of requests before closing them.
// If k is nil, then nothing is sent and the Response is closed anyway.
func closeWaitingRequests(waiting map[string][]*indexReq) {
	for _, rs := range waiting {
		for _, r := range rs {
			close(r.Response)
		}
	}
}

func NewKeyCache(ctx context.Context) *index {
	kc := &index{
		chRequests: make(chan *indexReq),
		chRefresh:  make(chan []string, 1),
	}
	go kc.serveRequests(ctx)
	return kc
}

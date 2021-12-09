package keycache

import (
	"context"
	"pempal/pemresources"
)

type cacheRequest struct {
	Request          string
	IncludeAnonymous bool
	Response         chan *pemresources.PrivateKey
}

type KeyCache struct {
	keyPath    []string
	chRequests chan *cacheRequest
	chRefresh  chan bool
}

func (ki KeyCache) Refresh() {
	ki.chRefresh <- true
}

func (ki KeyCache) KeyByID(publicKeyHash string) *pemresources.PrivateKey {
	keys := ki.makeRequest(publicKeyHash, false)
	if len(keys) == 0 {
		return nil
	}
	return keys[0]
}

func (ki KeyCache) Keys(includeAnonymous bool) []*pemresources.PrivateKey {
	return ki.makeRequest("*", includeAnonymous)
}

func (ki KeyCache) makeRequest(request string, anonymous bool) []*pemresources.PrivateKey {
	r := &cacheRequest{
		Request:          request,
		IncludeAnonymous: anonymous,
		Response:         make(chan *pemresources.PrivateKey, 1),
	}
	ki.chRequests <- r
	var keys []*pemresources.PrivateKey
	for k := range r.Response {
		keys = append(keys, k)
	}
	return keys
}

func (ki KeyCache) serveRequestsActive(ctx context.Context) {
	defer close(ki.chRequests)
	defer close(ki.chRefresh)

	keys := pemresources.Keys{}
	cache := map[string]*pemresources.PrivateKey{}

	// Main loop run whilst indexing (Active) (chIndexer is open)
	for {
		chIndexer := keys.ScanKeys(ctx, ki.keyPath...)
		waiting := map[string][]*cacheRequest{}
		ki.flushRefreshChannel()

		select {
		case <-ctx.Done():
			return

		case k, ok := <-chIndexer:
			// indexer has read a new key (or closed)
			if !ok {
				// finished indexing, close all waiting and enter passive mode until refresh or ctx done
				clearWaitingRequests(waiting)
				ki.serveRequestsPassive(ctx, cache)
				// on return, refresh cache again
				continue
			}

			cache[k.PublicKeyHash] = k

			// Check if any outstanding requests are waiting for this key
			wrs, ok := waiting[k.PublicKeyHash]
			if ok {
				delete(waiting, k.PublicKeyHash)
				sendWaitingRequests(k, wrs)
			}

		case r := <-ki.chRequests:
			k, ok := cache[r.Request]
			if !ok {
				// not found, place request in waiting until finished indexing
				waiting[r.Request] = append(waiting[r.Request], r)
				continue
			}
			// hit the cache, send response and close
			r.Response <- k
			close(r.Response)
		}
	}
}

// serveRequestsPassive serves requests with the static cache.
func (ki KeyCache) serveRequestsPassive(ctx context.Context, cache map[string]*pemresources.PrivateKey) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-ki.chRefresh:
			// refresh signal ends passive serving and returns to active serving
			return
		case r := <-ki.chRequests:
			if k, ok := cache[r.Request]; ok {
				r.Response <- k
			}
			close(r.Response)
		}
	}
}

// clear any outstanding refresh requests
func (ki KeyCache) flushRefreshChannel() {
	for len(ki.chRefresh) > 0 {
		<-ki.chRefresh
	}
}

// clearWaitingRequests closes all the outstanding requests in the given map without sending anything.
func clearWaitingRequests(waiting map[string][]*cacheRequest) {
	for _, v := range waiting {
		sendWaitingRequests(nil, v)
	}
}

// sendWaitingRequests sends the given key to the slice of requests before closing them.
// If k is nil, then nothing is sent and the Response is closed anyway.
func sendWaitingRequests(k *pemresources.PrivateKey, reqs []*cacheRequest) {
	for _, r := range reqs {
		if k != nil {
			r.Response <- k
		}
		close(r.Response)
	}
}

func NewKeyCache(ctx context.Context, keypath ...string) *KeyCache {
	kc := &KeyCache{
		keyPath:    keypath,
		chRequests: make(chan *cacheRequest, 5),
		chRefresh:  make(chan bool, 1),
	}
	go kc.serveRequestsActive(ctx)
	return kc
}

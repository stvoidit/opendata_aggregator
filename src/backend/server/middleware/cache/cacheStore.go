package cache

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

type cacheResponse struct {
	header     http.Header
	statusCode int
	value      []byte
	created    time.Time
	lastAccess time.Time
	TTL        time.Duration
}

func (cr cacheResponse) WriteResponse(w http.ResponseWriter) (err error) {
	for key, values := range cr.header {
		w.Header().Set(key, strings.Join(values, ","))
	}
	w.WriteHeader(cr.statusCode)
	_, err = w.Write(cr.value)
	return
}

type cacheStore struct {
	mu    sync.RWMutex
	store map[string]cacheResponse
}

func newCacheStore() *cacheStore { return &cacheStore{store: make(map[string]cacheResponse)} }

func (cs *cacheStore) Set(key string, cr cacheResponse) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cr.created = time.Now()
	cs.store[key] = cr
}

func (cs *cacheStore) Get(key string) *cacheResponse {
	cs.mu.RLock()
	cr, ok := cs.store[key]
	cs.mu.RUnlock()
	if ok {
		if cr.TTL > 0 && time.Since(cr.created) > cr.TTL {
			cs.Delete(key)
			return nil
		}
		cr.lastAccess = time.Now()
		return &cr
	}
	return nil
}

func (cs *cacheStore) Delete(key string) {
	cs.mu.RLock()
	_, ok := cs.store[key]
	cs.mu.RUnlock()
	if ok {
		cs.mu.Lock()
		delete(cs.store, key)
		cs.mu.Unlock()
	}
}

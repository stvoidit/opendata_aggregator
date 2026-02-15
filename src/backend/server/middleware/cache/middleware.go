package cache

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type CacheRule struct {
	Path string
	TTL  time.Duration
}

func (cr CacheRule) Validate(path string) bool {
	return strings.EqualFold(cr.Path, path)
}

func NewCacheMiddleware(rules ...CacheRule) func(http.Handler) http.Handler {
	var cs = newCacheStore()
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for i := range rules {
				if rules[i].Validate(r.URL.Path) {
					if cr := cs.Get(r.URL.Path); cr != nil {
						if err := cr.WriteResponse(w); err != nil {
							log.Err(err).Str("url", r.URL.Path).Msg("cache.NewCacheMiddleware.WriteResponse")
						} else {
							log.Debug().Str("url", r.URL.Path).Msg("from cache")
						}
						return
					}
					rec := httptest.NewRecorder()
					h.ServeHTTP(rec, r)
					result := rec.Result()
					res := cacheResponse{
						header:     result.Header.Clone(),
						statusCode: result.StatusCode,
						value:      rec.Body.Bytes(),
						TTL:        rules[i].TTL,
					}
					if result.StatusCode == http.StatusOK {
						cs.Set(r.URL.Path, res)
					}
					res.WriteResponse(w)
					return
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}

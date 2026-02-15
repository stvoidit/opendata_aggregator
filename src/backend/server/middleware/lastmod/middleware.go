package lastmod

import (
	"net/http"
	"time"
)

func CheckLastModified(actualTime time.Time, w http.ResponseWriter, r *http.Request) bool {
	if actualTime.IsZero() {
		return false
	}
	w.Header().Add("Last-Modified", actualTime.Format(time.RFC1123))
	if ims := r.Header.Get("If-Modified-Since"); ims != "" {
		if cacheLM, err := time.Parse(time.RFC1123, ims); err == nil {
			if actualTime.Equal(cacheLM) {
				w.WriteHeader(http.StatusNotModified)
				return true
			}
		}
	}
	return false
}

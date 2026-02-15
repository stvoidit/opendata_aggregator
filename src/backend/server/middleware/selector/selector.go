package selector

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"
)

// func SelectorAuth() func(h http.Handler) http.Handler {

// 	return func(h http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		})

// 	}
// }

type User struct {
	UserID   uint64 `json:"employee_id"`
	FullName string `json:"name"`
	Position string `json:"position"`
	EMail    string `json:"email"`
}

type ProviderUser = func(uint64) (*User, error)

func unzipResponse(response *http.Response) (err error) {
	if response.Uncompressed {
		return nil
	}
	if response.Header.Get("Content-Encoding") == "gzip" {
		var body []byte
		body, err = io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		if err := response.Body.Close(); err != nil {
			return err
		}
		response.Body, err = gzip.NewReader(bytes.NewReader(body))
	}
	return
}
func ProviderApi(providerApi string) ProviderUser {
	jar, _ := cookiejar.New(nil)
	c := http.Client{Timeout: 10 * time.Second, Jar: jar}
	etag := ""
	usersCache := make([]User, 0)
	var mu sync.RWMutex
	var ttl *time.Timer
	var matchUser = func(uid uint64, arr []User) *User {
		mu.RLock()
		defer mu.RUnlock()
		for i := range arr {
			if arr[i].UserID == uid {
				return &arr[i]
			}
		}
		return nil
	}
	return func(uid uint64) (*User, error) {
		defer func(t time.Time) {
			slog.Info("ProviderApi", slog.String("duration", time.Since(t).String()))
		}(time.Now())
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/mp_api/employee", providerApi), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Accept-Encoding", "gzip")
		req.Header.Add("Accept", "application/json")
		if len(usersCache) > 0 && etag != "" {
			req.Header.Add("If-None-Match", etag)
		}
		res, err := c.Do(req)
		if err != nil {
			return nil, err
		}
		if err := unzipResponse(res); err != nil {
			return nil, err
		}
		defer res.Body.Close()
		if res.StatusCode == http.StatusNotModified {
			slog.Debug("ProviderApi.useCache", slog.Int("statusCode", res.StatusCode))
			user := matchUser(uid, usersCache)
			if user == nil {
				return nil, errors.New("no match user id")
			}
			return user, nil
		}
		var users []User
		if err := json.NewDecoder(res.Body).Decode(&users); err != nil {
			return nil, err
		}
		mu.Lock()
		etag = res.Header.Get("Etag")
		usersCache = users
		if ttl != nil {
			ttl.Stop()
		}
		ttl = time.AfterFunc(time.Minute*20, func() {
			mu.Lock()
			defer mu.Unlock()
			etag = ""
			usersCache = make([]User, 0)
		})
		mu.Unlock()
		user := matchUser(uid, usersCache)
		if user == nil {
			return nil, errors.New("no match user id")
		}
		return user, nil
	}
}

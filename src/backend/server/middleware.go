package server

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"opendataaggregator/server/middleware/lastmod"
	"opendataaggregator/server/middleware/selector"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/klauspost/compress/gzhttp"
	"github.com/stvoidit/megaplan"
)

var (
	// defaultProfile - заглушка для дебага
	defaultProfile = UserProfile{UserID: 1000005, FullName: "Dev", CanLogin: true}
	// ErrCookiesNotFound - не найдены куки в заголовке
	ErrCookiesNotFound = errors.New("coockiesNotFound")
)

type contextKey string // ключ контекста - должен быть производным от базового типа

// UserProfile - cookies данные
type UserProfile struct {
	UserID    uint64 `json:"id,string"`
	FullName  string `json:"name"`
	Position  string `json:"position"`
	SessionID string `json:"sessionID"`
	CanLogin  bool   `json:"canLogin"`
	// Department Department `json:"department"`
}

// // Department - отдел пользователя
// type Department struct {
// 	ID   int64  `json:"id,string"`
// 	Name string `json:"name"`
// }

func (app *Application) setCookies(w http.ResponseWriter, r *http.Request, u *UserProfile) error {
	b, err := json.Marshal(u)
	if err != nil {
		return err
	}
	// зашифрованный данные для cookies
	encriptedValue, err := app.sg.Encrypt(b)
	if err != nil {
		return err
	}
	const oneDay = time.Hour * 24
	var _, domain, secure = extractHostForCookie(r)

	// cookies для бэка - данные о профиле
	cookies := http.Cookie{
		Path:     "/",
		Domain:   domain,
		Name:     string(app.contextKey),
		Value:    base64.RawStdEncoding.EncodeToString(encriptedValue),
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
		Secure:   secure,
		Expires:  time.Now().UTC().Truncate(oneDay).Add(oneDay).Add(-1 * time.Second),
	}
	if !cookies.Secure {
		cookies.SameSite = http.SameSiteLaxMode
	}
	http.SetCookie(w, &cookies)
	return nil
}

func (app *Application) readCookies(r *http.Request, u *UserProfile) error {
	cookiesValue, err := r.Cookie(string(app.contextKey))
	if err != nil {
		return err
	}
	b, err := base64.RawStdEncoding.DecodeString(cookiesValue.Value)
	if err != nil {
		return err
	}
	dec, err := app.sg.Decrypt(b)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dec, u)
	return err
}

func (app Application) readContextValue(ctx context.Context) (user *UserProfile, err error) {
	var ok bool
	user, ok = ctx.Value(app.contextKey).(*UserProfile)
	if !ok {
		err = errors.New("invalide context value")
	}
	return
}

// https://dev.megaplan.ru/apps/common.html#application-page
func (app *Application) megaplanverify(allowPaths ...string) func(http.Handler) http.Handler {
	var allowRgx = make([]*regexp.Regexp, 0)
	for _, ap := range allowPaths {
		rgx, err := regexp.Compile(ap)
		if err != nil {
			panic(err)
		}
		allowRgx = append(allowRgx, rgx)
	}
	var checkAllowedPath = func(p string) bool {
		for i := range allowRgx {
			if allowRgx[i].MatchString(p) {
				return true
			}
		}
		return false
	}
	var responseAuthError = func(w http.ResponseWriter, err error) {
		http.SetCookie(w, &http.Cookie{
			Path:     "/",
			Name:     string(app.contextKey),
			MaxAge:   -1,
			SameSite: http.SameSiteNoneMode,
			Secure:   true})
		Jsonify(w, map[string]string{
			"megaplanDomain": app.config.Megaplan.Domain,
			"appUUID":        app.config.Megaplan.UUID,
			"error":          err.Error(),
		}, http.StatusUnauthorized)
	}
	var selectorProvider = selector.ProviderApi(app.config.Megaplan.ProviderHost)
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем наличие наших куков по ключу, если есть, то передаем в контексте
			// если нет (будет ошибка), то идем дальше и проверяем наличие ключей для установки куков
			// fmt.Println(r.Header.Get("x-user-id"))
			if checkAllowedPath(r.URL.Path) {
				h.ServeHTTP(w, r)
				return
			}
			ctx := r.Context()
			if app.config.Debug {
				h.ServeHTTP(w, r.WithContext(context.WithValue(ctx, app.contextKey, &defaultProfile)))
				return
			}
			if username, password, ok := r.BasicAuth(); ok && strings.EqualFold(password, app.config.Megaplan.UUID) {
				serviceProfile := UserProfile{UserID: 1, FullName: username, CanLogin: true}
				h.ServeHTTP(w, r.WithContext(context.WithValue(ctx, app.contextKey, &serviceProfile)))
				return
			}
			ap, err := extractAuthParams(r)
			if err != nil {
				slog.Error("extractAuthParams", slog.String("error", err.Error()))
				responseAuthError(w, err)
				return
			}
			if ap.IsTest() {
				// фэйковый юзер от фронта селектора
				puser, err := selectorProvider(ap.signAsUint64())
				if err != nil {
					slog.Error("selectorProvider", slog.String("error", err.Error()))
					responseAuthError(w, err)
					return
				}
				profile := UserProfile{
					UserID:   puser.UserID,
					FullName: puser.FullName,
					Position: puser.Position,
					CanLogin: true,
				}
				if err := app.setCookies(w, r, &profile); err != nil {
					slog.Error("setCookies", slog.String("error", err.Error()))
					w.WriteHeader(500)
					return
				}
				h.ServeHTTP(w, r.WithContext(context.WithValue(ctx, app.contextKey, &profile)))
				return
			}
			// обычная аутентификация через v1 мегаплана
			var profile UserProfile
			if err := app.readCookies(r, &profile); err == nil {
				if profile.UserID != 0 {
					h.ServeHTTP(w, r.WithContext(context.WithValue(ctx, app.contextKey, &profile)))
					return
				}
			}
			response, err := app.mpapi.CheckUser(ap.UserSign)
			if err != nil {
				responseAuthError(w, fmt.Errorf("check userSign ERROR: %s", err))
				return
			}
			defer response.Body.Close()
			if err := json.NewDecoder(response.Body).Decode(megaplan.ExpectedResponse(&profile)); err != nil {
				slog.Error("Decode.profile", slog.String("error", err.Error()))
				w.WriteHeader(500)
				return
			}
			if !profile.CanLogin {
				slog.Error("CanLogin", slog.String("error", "user can't login"))
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if err := app.setCookies(w, r, &profile); err != nil {
				slog.Error("setCookies", slog.String("error", err.Error()))
				w.WriteHeader(500)
				return
			}
			h.ServeHTTP(w, r.WithContext(context.WithValue(ctx, app.contextKey, &profile)))
		})
	}
}

type authParams struct {
	ApplicationUuid string
	UserSign        string
	SessionId       string
}

func (ap authParams) signAsUint64() uint64 {
	n, _ := strconv.ParseUint(ap.UserSign, 10, 64)
	return n
}
func (ap authParams) IsTest() bool {
	return strings.EqualFold(ap.ApplicationUuid, "cndTest")
}

func containsAuthParams(args url.Values) bool {
	var keys = [...]string{
		"applicationUuid",
		"userSign",
		"sessionId",
	}
	for _, k := range keys {
		if _, ok := args[k]; !ok {
			return false
		}
	}
	return true
}
func extractAuthParams(r *http.Request) (*authParams, error) {
	args := r.URL.Query()
	if !containsAuthParams(args) {
		ref, err := url.Parse(r.Referer())
		if err != nil {
			return nil, errors.New("not enough authorization parameters")
		}
		args = ref.Query()
		if !containsAuthParams(args) {
			return nil, errors.New("not enough authorization parameters")
		}
	}
	return &authParams{
		ApplicationUuid: args.Get("applicationUuid"),
		UserSign:        args.Get("userSign"),
		SessionId:       args.Get("sessionId"),
	}, nil
}

func extractHostForCookie(r *http.Request) (addr, domain string, secure bool) {
	var (
		scheme = "http"
		host   = ""
		origin = r.Header.Get("Origin")
	)
	if strings.EqualFold(origin, "") {
		fmt.Println(r.Host)
		return "", "", false
	} else {
		fmt.Println("ORIGIN:", origin)
		URL, err := url.Parse(origin)
		if err != nil {
			return "", "", false
		}
		scheme = URL.Scheme
		host = URL.Host
	}

	// origin, r.Host, strings.HasPrefix(origin, "https")
	return origin, host, strings.EqualFold(scheme, "https")
}

func setCORS(w http.ResponseWriter, origin string) {
	if strings.EqualFold(origin, "") {
		return
	}
	headers := w.Header()
	headers.Add("Access-Control-Allow-Methods", "GET,POST,UPDATE,OPTIONS")
	headers.Add("Access-Control-Allow-Origin", origin)
	headers.Add("Access-Control-Allow-Credentials", "true")
	headers.Add("Access-Control-Allow-Headers", "Authorization,Cookie,Accept-Encodin,Accept,Host,Origin,Referer,User-Agent,Content-Type,Content-Length,Content-Disposition,Cache-Control")
}

func (app *Application) CrosDomainHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addr, domain, _ := extractHostForCookie(r)
		slog.Debug("extractHostForCookie", slog.String("addr", addr), slog.String("domain", domain))
		setCORS(w, addr)
		h.ServeHTTP(w, r)
	})
}

func (app *Application) LastModHeader(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/db_stats":
			if t, err := app.db.LastModUpdates(r.Context()); err == nil {
				if lastmod.CheckLastModified(t, w, r) {
					return
				}
			}
		}
		h.ServeHTTP(w, r)
	})
}

func CompressHandler(h http.Handler) http.Handler {
	wrapper, err := gzhttp.NewWrapper(
		gzhttp.AllowCompressedRequests(true),
		gzhttp.CompressionLevel(8),
	)
	if err != nil {
		panic(err)
	}
	return wrapper(h)
}

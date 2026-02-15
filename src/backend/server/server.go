package server

import (
	"context"
	"fmt"
	stdlog "log"
	"net/http"
	"opendataaggregator/config"
	"opendataaggregator/server/middleware/cache"
	"opendataaggregator/store"
	"opendataaggregator/utils"
	"os"
	"runtime"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/stvoidit/megaplan"
)

// Application - сборник всего необходимого в одну структуру
type Application struct {
	srv        *http.Server         // http.Server
	db         *store.DB            // подключение к БД
	config     *config.Config       // конфиг
	contextKey contextKey           // ключ контекста для context.Values
	mpapi      *megaplan.API        // мегаплан API v1
	sg         *utils.Cryptographer // шифрование
}

// NewApplication - фабрика
func NewApplication(cnf *config.Config) *Application {
	log.Info().Str("branch", config.Branch).Str("versionTag", config.VersionTag).Str("commitHash", config.CommitHash).Send()
	db, err := store.NewDB(cnf, runtime.NumCPU()*4)
	if err != nil {
		log.Fatal().Err(err).Msg("NewApplication")
	}
	mpapi := megaplan.NewAPI(cnf.Megaplan.UUID, cnf.Megaplan.Secret, "https://"+cnf.Megaplan.Domain)
	mpapi.EnableCompression(true)
	return &Application{
		contextKey: "profile",
		sg:         utils.NewCryptographer(cnf.Megaplan.Secret),
		db:         db,
		config:     cnf,
		mpapi:      mpapi}
}

// ListenAndServe - http.ListenAndServe + graceful shutdown
func (app *Application) ListenAndServe(ctx context.Context) (err error) {
	r := mux.NewRouter()
	r.Use(CompressHandler)
	app.setHandlers(r)
	app.setStatic(r)
	app.srv = &http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      r,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: 0,
		ErrorLog:     stdlog.Default(),
	}
	app.srv.RegisterOnShutdown(app.db.Close)
	go func() {
		log.Info().Str("listen on", fmt.Sprintf("http://%s", app.srv.Addr)).Send()
		err = app.srv.ListenAndServe()
	}()
	app.shutdown(ctx)
	<-ctx.Done()
	return
}

func (app *Application) setHandlers(r *mux.Router) {
	logginghandler := func(h http.Handler) http.Handler { return handlers.LoggingHandler(os.Stdout, h) }
	const defaultTLL = time.Hour * 2
	sr := r.PathPrefix("/api/").Subrouter()
	sr.Use(
		app.CrosDomainHeaders,
		app.LastModHeader,
		logginghandler,
		// app.megaplanverify(
		// 	`^/api/version`,
		// 	`^/api/db_stats`,
		// 	`^/api/json_schema`,
		// 	`^/[^(api)]`,
		// ),
	)
	sr.Use(cache.NewCacheMiddleware(
		cache.CacheRule{Path: "/api/handbook_okved", TTL: defaultTLL},
		cache.CacheRule{Path: "/api/handbook_tax_authority", TTL: defaultTLL},
		cache.CacheRule{Path: "/api/hotels", TTL: defaultTLL},
		cache.CacheRule{Path: "/api/legal_statuses", TTL: defaultTLL},
		cache.CacheRule{Path: "/api/db_stats", TTL: defaultTLL},
		cache.CacheRule{Path: "/api/db_stats/sources", TTL: defaultTLL},
		cache.CacheRule{Path: "/api/db_stats/last_updates_sources", TTL: defaultTLL},
		cache.CacheRule{Path: "/api/db_stats/stats_egr", TTL: defaultTLL},
		cache.CacheRule{Path: "/api/db_stats/stats_balance", TTL: defaultTLL},
		cache.CacheRule{Path: "/api/db_stats/stats_hotels", TTL: defaultTLL},
		cache.CacheRule{Path: "/api/db_stats/stats_ross_accreditation", TTL: defaultTLL},
		cache.CacheRule{Path: "/api/db_stats/stats_tax_offenses", TTL: defaultTLL},
		cache.CacheRule{Path: "/api/db_stats/stats_register_of_trademarks", TTL: defaultTLL},
		cache.CacheRule{Path: "/api/db_stats/stats_fssp", TTL: defaultTLL},
		cache.CacheRule{Path: "/api/db_stats/stats_smp", TTL: defaultTLL},
		cache.CacheRule{Path: "/api/db_stats/stats_debtam", TTL: defaultTLL},
		cache.CacheRule{Path: "/api/db_stats/stats_fgis", TTL: defaultTLL},
		cache.CacheRule{Path: "/api/db_stats/stats_tax_regime", TTL: defaultTLL},
		cache.CacheRule{Path: "/api/db_stats/stat_avg_employes_number", TTL: defaultTLL},
	))
	sr.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		type versionInfo struct {
			CommitHash string `json:"commitHash"`
			VersionTag string `json:"versionTag"`
			Branch     string `json:"branch"`
		}
		Jsonify(w, versionInfo{
			CommitHash: config.CommitHash,
			VersionTag: config.VersionTag,
			Branch:     config.Branch},
			http.StatusOK)
	})
	sr.HandleFunc("/init", app.InitHandler)
	sr.HandleFunc("/total_info", app.TotalInfoHandler)
	sr.HandleFunc("/info/{inn:[0-9]+}", app.TotalInfoHandler)
	sr.HandleFunc("/info/{inn:[0-9]+}/{ogrn:[0-9]+}", app.InfoCardHandler)
	sr.HandleFunc("/json_schema", app.JSONSchemaHandler).Methods(http.MethodGet)
	sr.HandleFunc("/search", app.Search)
	sr.HandleFunc("/count", app.CountSearch)
	sr.HandleFunc("/handbook_okved", app.HandbookOKVED).Methods(http.MethodGet)
	sr.HandleFunc("/search_okved", app.SearchOKVED).Methods(http.MethodGet)
	sr.HandleFunc("/handbook_tax_authority", app.HandbookTaxAuthority).Methods(http.MethodGet)
	sr.HandleFunc("/search_tax_authority", app.SearchHandbookTaxAuthority).Methods(http.MethodGet)
	sr.HandleFunc("/download", app.Download)
	sr.HandleFunc("/service_log_sources", app.ServiceLogSources)
	sr.HandleFunc("/hotels", app.GetHotelsView)
	sr.HandleFunc("/download/hotels", app.DownloadHotels)
	sr.HandleFunc("/handbook_categories_ip", app.HandbookCategoriesIP)
	sr.HandleFunc("/legal_statuses", app.HandbookLegalStatuses)
	sr.HandleFunc("/upload", app.UploadSourceFile).Methods(http.MethodPost)
	sr.HandleFunc("/db_stats", app.DatabaseStat)

	sr.HandleFunc("/db_stats/sources", app.DatabaseStatParts)
	sr.HandleFunc("/db_stats/last_updates_sources", app.DatabaseStatParts)
	sr.HandleFunc("/db_stats/stats_egr", app.DatabaseStatParts)
	sr.HandleFunc("/db_stats/stats_balance", app.DatabaseStatParts)
	sr.HandleFunc("/db_stats/stats_hotels", app.DatabaseStatParts)
	sr.HandleFunc("/db_stats/stats_ross_accreditation", app.DatabaseStatParts)
	sr.HandleFunc("/db_stats/stats_tax_offenses", app.DatabaseStatParts)
	sr.HandleFunc("/db_stats/stats_register_of_trademarks", app.DatabaseStatParts)
	sr.HandleFunc("/db_stats/stats_fssp", app.DatabaseStatParts)
	sr.HandleFunc("/db_stats/stats_smp", app.DatabaseStatParts)
	sr.HandleFunc("/db_stats/stats_debtam", app.DatabaseStatParts)
	sr.HandleFunc("/db_stats/stats_fgis", app.DatabaseStatParts)
	sr.HandleFunc("/db_stats/stats_tax_regime", app.DatabaseStatParts)
	sr.HandleFunc("/db_stats/stat_avg_employes_number", app.DatabaseStatParts)
}

func (app Application) setStatic(router *mux.Router) {
	router.PathPrefix("/").Handler(http.FileServer(http.FS(FileSystemSPA("static"))))
}

func (app *Application) shutdown(ctx context.Context) {
	go func() {
		<-ctx.Done()
		log.Info().Msg("app.Shutdown")
		if err := app.srv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
			log.Info().Err(err).Msg("app.srv.Shutdown")
			return
		}
	}()
}

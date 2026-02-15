// Package server - ...
package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"opendataaggregator/config"
	"opendataaggregator/store"
	"time"

	"github.com/gorilla/mux"
	"github.com/invopop/jsonschema"
	"github.com/rs/zerolog/log"
)

var dataSchema []byte

func init() {
	var r jsonschema.Reflector
	var buf = bytes.NewBuffer(nil)
	defer buf.Reset()
	e := json.NewEncoder(buf)
	schema := r.Reflect(new(store.AllTypes))
	schema.Version = config.VersionTag
	if err := e.Encode(schema); err != nil {
		log.Err(err).Msg("init.jsonschema")
		return
	}
	dataSchema = buf.Bytes()
}

// InitHandler - точка входа для SPA - авторизация, получение конфигурации с бэка
func (app *Application) InitHandler(w http.ResponseWriter, r *http.Request) {
	type applicationInfo struct {
		Domain string       `json:"megaplanDomain"`
		UUID   string       `json:"appUUID"`
		User   *UserProfile `json:"user"`
	}
	var ctx = r.Context()
	profile, err := app.readContextValue(ctx)
	if err != nil || profile == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	Jsonify(w, applicationInfo{
		Domain: app.config.Megaplan.Domain,
		UUID:   app.config.Megaplan.UUID,
		User:   profile},
		http.StatusOK)
}

// TotalInfoHandler - получение полной информации
func (app *Application) TotalInfoHandler(w http.ResponseWriter, r *http.Request) {
	var inn string
	if r.URL.Query().Has("inn") {
		inn = r.URL.Query().Get("inn")
	} else {
		inn = mux.Vars(r)["inn"]
	}
	if len(inn) == 0 {
		Jsonify(w, errors.New("ИНН не указан"), http.StatusBadRequest)
		return
	}
	result, err := app.db.SelectTotalInfo(r.Context(), inn)
	if result == nil {
		Jsonify(w, fmt.Errorf("данные по ИНН %s не найдены", inn), http.StatusNotFound)
		return
	}
	if err != nil {
		Jsonify(w, err, http.StatusInternalServerError)
		return
	}
	Jsonify(w, result, http.StatusOK)
}

// InfoCardHandler - карточка ЮЛ или ИП
func (app *Application) InfoCardHandler(w http.ResponseWriter, r *http.Request) {
	var inn string
	if r.URL.Query().Has("inn") {
		inn = r.URL.Query().Get("inn")
	} else {
		inn = mux.Vars(r)["inn"]
	}
	if len(inn) == 0 {
		Jsonify(w, errors.New("ИНН не указан"), http.StatusBadRequest)
		return
	}
	var ogrn string
	if r.URL.Query().Has("ogrn") {
		ogrn = r.URL.Query().Get("ogrn")
	} else {
		ogrn = mux.Vars(r)["ogrn"]
	}
	if len(ogrn) == 0 {
		Jsonify(w, errors.New("ОГРН не указан"), http.StatusBadRequest)
		return
	}
	result, err := app.db.SelectInfo(r.Context(), inn, ogrn)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			Jsonify(w, err, http.StatusNotFound)
		} else {
			Jsonify(w, err, http.StatusInternalServerError)
		}
		return
	}
	Jsonify(w, result, http.StatusOK)
}

// JSONSchemaHandler - JsonSchema данных с TotalInfoHandler
func (app *Application) JSONSchemaHandler(w http.ResponseWriter, r *http.Request) {
	Jsonify(w, dataSchema, http.StatusOK)
}

// Search - поиск
func (app *Application) Search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := r.URL.Query()
	result, err := app.db.Search(ctx, params)
	if err != nil {
		Jsonify(w, err, http.StatusInternalServerError)
	} else {
		Jsonify(w, result, http.StatusOK)
	}
}

// HandbookOKVED - весь справочник ОКВЭД
func (app *Application) HandbookOKVED(w http.ResponseWriter, r *http.Request) {
	result, err := app.db.SelectHandbookOKVED(r.Context())
	if err != nil {
		Jsonify(w, err, http.StatusInternalServerError)
	} else {
		Jsonify(w, result, http.StatusOK)
	}
}

// SearchOKVED - поиск ОКВЭД по названию
func (app *Application) SearchOKVED(w http.ResponseWriter, r *http.Request) {
	var query string
	var params = r.URL.Query()
	if !params.Has("q") {
		Jsonify(w, errors.New(`параметр "q" отсутствует`), http.StatusBadRequest)
		return
	}
	query = params.Get("q")
	if len(query) < 2 {
		Jsonify(w, errors.New(`параметр "q" должен содержать минимум 2 символа`), http.StatusBadRequest)
		return
	}
	result, err := app.db.SearchOKVED(r.Context(), query)
	if err != nil {
		Jsonify(w, err, http.StatusInternalServerError)
	} else {
		Jsonify(w, result, http.StatusOK)
	}
}

// HandbookTaxAuthority - весь справочник Налоговых органов
func (app *Application) HandbookTaxAuthority(w http.ResponseWriter, r *http.Request) {
	result, err := app.db.SelectHandbookTaxAuthority(r.Context())
	if err != nil {
		Jsonify(w, err, http.StatusInternalServerError)
	} else {
		Jsonify(w, result, http.StatusOK)
	}
}

// SearchHandbookTaxAuthority - поиск налогового органа
func (app *Application) SearchHandbookTaxAuthority(w http.ResponseWriter, r *http.Request) {
	var query string
	var params = r.URL.Query()
	if !params.Has("q") {
		Jsonify(w, errors.New(`параметр "q" отсутствует`), http.StatusBadRequest)
		return
	}
	query = params.Get("q")
	if len(query) < 2 {
		Jsonify(w, errors.New(`параметр "q" должен содержать минимум 2 символа`), http.StatusBadRequest)
		return
	}
	result, err := app.db.SearchHandbookTaxAuthority(r.Context(), query)
	if err != nil {
		Jsonify(w, err, http.StatusInternalServerError)
	} else {
		Jsonify(w, result, http.StatusOK)
	}
}

// CountSearch - кол-во строк в результате поиска
func (app *Application) CountSearch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := r.URL.Query()
	count, err := app.db.CountSearch(ctx, params)
	if err != nil {
		Jsonify(w, err, http.StatusInternalServerError)
	} else {
		Jsonify(w, count, http.StatusOK)
	}
}

func (app *Application) Download(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Add(headerContentType, "text/csv")
	w.WriteHeader(http.StatusOK)
	if err := app.db.SearchCSV(ctx, w); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		Jsonify(w, err, http.StatusInternalServerError)
	}
}

func (app *Application) ServiceLogSources(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if data, err := app.db.GetServiceLogSources(ctx); err != nil {
		Jsonify(w, err, http.StatusInternalServerError)
	} else {
		Jsonify(w, &data, http.StatusOK)
	}
}

func (app *Application) DatabaseStat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if data, err := app.db.SelectDatabaseStat(ctx); err != nil {
		Jsonify(w, err, http.StatusInternalServerError)
	} else {
		data.Sources = app.config.Sources
		Jsonify(w, &data, http.StatusOK)
	}
}
func (app *Application) DatabaseStatParts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var maxAge = time.Hour * 2
	var err error
	var v any
	switch r.URL.Path {
	case "/api/db_stats/sources":
		v = app.config.Sources
	case "/api/db_stats/last_updates_sources":
		v, err = app.db.SelectLastUpdatesSources(ctx)
	case "/api/db_stats/stats_egr":
		v, err = app.db.SelectStatEGR(ctx)
	case "/api/db_stats/stats_balance":
		v, err = app.db.SelectStatBalance(ctx)
	case "/api/db_stats/stats_hotels":
		v, err = app.db.SelectStatsHotels(ctx)
	case "/api/db_stats/stats_ross_accreditation":
		v, err = app.db.SelectStatsRossAccreditation(ctx)
	case "/api/db_stats/stats_tax_offenses":
		v, err = app.db.SelectStatsTaxOffenses(ctx)
	case "/api/db_stats/stats_register_of_trademarks":
		v, err = app.db.SelectStatsRegisterOfTrademarks(ctx)
	case "/api/db_stats/stats_fssp":
		v, err = app.db.SelectStatsFSSP(ctx)
	case "/api/db_stats/stats_smp":
		v, err = app.db.SelectStatsSMP(ctx)
	case "/api/db_stats/stats_debtam":
		v, err = app.db.SelectStatsDEBTAM(ctx)
	case "/api/db_stats/stats_fgis":
		v, err = app.db.SelectStatsFGIS(ctx)
	case "/api/db_stats/stats_tax_regime":
		v, err = app.db.SelectStatsTaxRegime(ctx)
	case "/api/db_stats/stat_avg_employes_number":
		v, err = app.db.SelectStatAvgEmployesNumber(ctx)
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		Jsonify(w, err, http.StatusInternalServerError)
	} else {
		w.Header().Add(headerCacheControl, "public")
		w.Header().Add(headerCacheControl, fmt.Sprintf("max-age=%d", int64(maxAge.Seconds())))
		Jsonify(w, &v, http.StatusOK)
	}
}

func (app *Application) GetHotelsView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if data, err := app.db.GetHotelsView(ctx); err != nil {
		Jsonify(w, err, http.StatusInternalServerError)
	} else {
		Jsonify(w, &data, http.StatusOK)
	}
}
func (app *Application) DownloadHotels(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Add(headerContentType, octetStream)
	w.Header().Add(headerContentDisposition, `attachment; filename="hotels.xlsx"`)
	w.WriteHeader(http.StatusOK)
	if err := app.db.DownloadHotels(ctx, w); err != nil {
		log.Err(err).Msg("handlers.DownloadHotels")
	}
}

func (app *Application) HandbookCategoriesIP(w http.ResponseWriter, r *http.Request) {
	if data, err := app.db.SelectHandbookCategoriesIP(r.Context()); err != nil {
		Jsonify(w, err, http.StatusInternalServerError)
	} else {
		Jsonify(w, data, http.StatusOK)
	}
}

func (app *Application) HandbookLegalStatuses(w http.ResponseWriter, r *http.Request) {
	if data, err := app.db.SelectLegalStatuses(r.Context()); err != nil {
		Jsonify(w, err, http.StatusInternalServerError)
	} else {
		Jsonify(w, data, http.StatusOK)
	}
}

// UploadSourceFile - загрузка файла исходников данных
// TODO: т.к. исходники лежат на другом сервере, то возможно стоит сделать отправку по http тут
func (app *Application) UploadSourceFile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var rc = http.NewResponseController(w)
	rc.SetReadDeadline(time.Time{})
	fmt.Println("sourceType:", r.FormValue("sourceType"))
	f, fh, err := r.FormFile("sourceFile")
	if err != nil {
		Jsonify(w, err, http.StatusBadRequest)
		return
	}
	defer f.Close()
	fmt.Println(fh.Filename, fh.Size)
	// io.Copy(io.Discard, f)
	Jsonify(w, "uploaded", http.StatusCreated)
}

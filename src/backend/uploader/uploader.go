// Package uploader - загрузчик в базу и парсинг данных
package uploader

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"opendataaggregator/models"
	"opendataaggregator/models/egr"
	"opendataaggregator/parsers/csvparser"
	"opendataaggregator/parsers/parsexlsx"
	"opendataaggregator/parsers/xmlparser"
	"opendataaggregator/store"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/text/encoding/charmap"
)

// NewParserUploader - инициализация
func NewParserUploader(db *store.DB, workersCount int) *ParserUploader {
	return &ParserUploader{db: db, workersCount: workersCount}
}

// ParserUploader - парсинг и добавление данных в БД
type ParserUploader struct {
	db           *store.DB
	workersCount int
}

// tickPrinter - для экономии ресурсов stdout вывод в консоль по таймеру
func tickPrinter(ctx context.Context, delay int, counter *atomic.Int64) {
	go func() {
		var (
			count    int
			duration time.Duration
			start    = time.Now()
			ticker   = time.NewTicker(time.Second * time.Duration(delay))
			avg      int
		)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				count = int(counter.Load())
				duration = time.Since(start)
				if sec := int(duration.Seconds()); sec > 0 {
					avg = count / int(duration.Seconds())
				}
				log.Info().Int("count", count).Int("avg/sec", avg).Send()
			case <-ctx.Done():
				count = int(counter.Load())
				duration = time.Since(start)
				if sec := int(duration.Seconds()); sec > 0 {
					avg = count / int(duration.Seconds())
				}
				log.Info().Str("Duration", duration.String()).Int("Total", count).Int("Avg/sec", avg).Msg("DONE")
				return
			}
		}
	}()
}

// splitFilanemAndCheck - извлечение названия файла и проверка в БД
func (pu ParserUploader) splitFilanemAndCheck(ctx context.Context, sourcetype, filename string) (fname string, err error) {
	_, fname = filepath.Split(filename)
	yes, err := pu.db.ChecFileIsParsed(ctx, sourcetype, fname)
	if err != nil {
		return fname, err
	}
	if yes {
		return fname, fmt.Errorf("file alredy parsed: %s", filename)
	}
	return
}

// ParseБухОтчетность - ...
func (pu ParserUploader) ParseБухОтчетность(ctx context.Context, filename string) {
	fname, err := pu.splitFilanemAndCheck(ctx, "balance", filename)
	if err != nil {
		log.Err(err).Send()
		return
	}
	var counter atomic.Int64
	zr, err := zip.OpenReader(filename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer zr.Close()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	tickPrinter(ctx, 10, &counter)
	log.Info().Str("file", filename).Msg("read started")
	var docsChan = xmlparser.ParseБухОтчетность(ctx, zr)
	var wg sync.WaitGroup
	wg.Add(pu.workersCount)
	for i := 0; i < pu.workersCount; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			var arr = make([]models.БухОтчетность, 0, 1000)
			for doc := range docsChan {
				if doc.INN() == "" {
					continue
				}
				arr = append(arr, doc)
				if len(arr) == 1000 {
					if err := pu.db.BatchБухОтчетность(ctx, arr); err != nil {
						if errors.Is(err, context.Canceled) {
							return
						}
						log.Err(err).Send()
					}
					counter.Add(int64(len(arr)))
					arr = make([]models.БухОтчетность, 0, 1000)
				}
			}
			if err := pu.db.BatchБухОтчетность(ctx, arr); err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				log.Err(err).Send()
			}
			counter.Add(int64(len(arr)))
		}(&wg)
	}
	wg.Wait()
	if err := pu.db.MarkFileIsParsed(ctx, "balance", fname); err != nil {
		log.Fatal().Err(err).Send()
	}
}

// ParseРеестрДисквалифицированныхЛиц - ...
func (pu ParserUploader) ParseРеестрДисквалифицированныхЛиц(ctx context.Context, filename string) {
	fname, err := pu.splitFilanemAndCheck(ctx, "registerdisqualified", filename)
	if err != nil {
		log.Err(err).Send()
		return
	}
	var count int
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer f.Close()
	var docs = make([]models.РеестрДисквалифицированныхЛиц, 0)
	docsChan, err := csvparser.ParseРеестрДисквалифицированныхЛиц(csv.NewReader(f))
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	for doc := range docsChan {
		docs = append(docs, doc)
		count++
	}
	if err := pu.db.BatchРеестрДисквалифицированныхЛиц(ctx, docs); err != nil {
		log.Fatal().Err(err).Send()
	}
	if err := pu.db.MarkFileIsParsed(ctx, "registerdisqualified", fname); err != nil {
		log.Fatal().Err(err).Send()
	}
}

// ParseНалоговыхПравонарушенияхИМерахОтветственности - ...
func (pu ParserUploader) ParseНалоговыхПравонарушенияхИМерахОтветственности(ctx context.Context, filename string) {
	fname, err := pu.splitFilanemAndCheck(ctx, "taxoffence", filename)
	if err != nil {
		log.Err(err).Send()
		return
	}
	var count int
	zr, err := zip.OpenReader(filename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer zr.Close()
	var batchSize = 100_000
	var docs = make([]models.НалоговыхПравонарушенияхИМерахОтветственности, 0, batchSize)
	for doc := range xmlparser.ParseParallelIterFunc[models.НалоговыхПравонарушенияхИМерахОтветственности](ctx, zr) {
		count++
		docs = append(docs, doc)
		if len(docs) == batchSize {
			if err := pu.db.BatchНалоговыхПравонарушенияхИМерахОтветственности(ctx, docs); err != nil {
				log.Fatal().Err(err).Send()
			}
			docs = make([]models.НалоговыхПравонарушенияхИМерахОтветственности, 0, batchSize)
			log.Info().Int("batch", count).Send()
		}
	}
	if err := pu.db.BatchНалоговыхПравонарушенияхИМерахОтветственности(ctx, docs); err != nil {
		log.Fatal().Err(err).Send()
	}
	if err := pu.db.MarkFileIsParsed(ctx, "taxoffence", fname); err != nil {
		log.Fatal().Err(err).Send()
	}
}

// ParseРеестрСубъектовМалогоИСреднегоПредпринимательства - ...
func (pu ParserUploader) ParseРеестрСубъектовМалогоИСреднегоПредпринимательства(ctx context.Context, filename string) {
	fname, err := pu.splitFilanemAndCheck(ctx, "rsmp", filename)
	if err != nil {
		log.Err(err).Send()
		return
	}
	var count int
	zr, err := zip.OpenReader(filename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer zr.Close()
	var batchSize = 100_000
	var docs = make([]models.РеестрСубъектовМалогоИСреднегоПредпринимательства, 0, batchSize)
	for doc := range xmlparser.ParseParallelIterFunc[models.РеестрСубъектовМалогоИСреднегоПредпринимательства](ctx, zr) {
		count++
		docs = append(docs, doc)
		if len(docs) == batchSize {
			if err := pu.db.BatchСубъектМалогоИСреднегоПредпринимательства(ctx, docs); err != nil {
				log.Fatal().Err(err).Send()
			}
			docs = make([]models.РеестрСубъектовМалогоИСреднегоПредпринимательства, 0, batchSize)
			log.Info().Int("batch", count).Send()
		}
	}
	if err := pu.db.BatchСубъектМалогоИСреднегоПредпринимательства(ctx, docs); err != nil {
		log.Fatal().Err(err).Send()
	}
	if err := pu.db.MarkFileIsParsed(ctx, "rsmp", fname); err != nil {
		log.Fatal().Err(err).Send()
	}
}

// ParseОКВЭД - ...
func (pu ParserUploader) ParseОКВЭД(ctx context.Context, filename string) {
	fname, err := pu.splitFilanemAndCheck(ctx, "okved2", filename)
	if err != nil {
		log.Err(err).Send()
		return
	}
	var count int
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer f.Close()
	r := csv.NewReader(charmap.Windows1251.NewDecoder().Reader(f))
	docsChan, err := csvparser.ParseОКВЭД(r)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	var docs = make([]models.ОКВЭД, 0)
	for doc := range docsChan {
		count++
		docs = append(docs, doc)
	}
	if err := pu.db.BatchОКВЭД(ctx, docs); err != nil {
		log.Fatal().Err(err).Send()
	}
	if err := pu.db.MarkFileIsParsed(ctx, "okved2", fname); err != nil {
		log.Fatal().Err(err).Send()
	}
}

// ParseОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков - ...
func (pu ParserUploader) ParseОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков(ctx context.Context, filename string) {
	fname, err := pu.splitFilanemAndCheck(ctx, "otz", filename)
	if err != nil {
		log.Err(err).Send()
		return
	}
	var count int
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer f.Close()
	f.Seek(3, 3) // BOM
	r := csv.NewReader(f)
	docsChan, err := csvparser.ParseОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков(r)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	var batchSize = 10_000
	var docs = make([]models.ОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков, 0, batchSize)
	for doc := range docsChan {
		count++
		docs = append(docs, doc)
		if len(docs) == batchSize {
			if err := pu.db.BatchОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков(ctx, docs); err != nil {
				log.Fatal().Err(err).Send()
			}
			docs = make([]models.ОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков, 0, batchSize)
			log.Info().Int("batch", count).Send()
		}
	}
	if err := pu.db.BatchОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков(ctx, docs); err != nil {
		log.Fatal().Err(err).Send()
	}
	if err := pu.db.MarkFileIsParsed(ctx, "otz", fname); err != nil {
		log.Fatal().Err(err).Send()
	}
}

// ParseОткрытыйРеестрТоварныхЗнаков - ...
func (pu ParserUploader) ParseОткрытыйРеестрТоварныхЗнаков(ctx context.Context, filename string) {
	fname, err := pu.splitFilanemAndCheck(ctx, "tz", filename)
	if err != nil {
		log.Err(err).Send()
		return
	}
	var count int
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer f.Close()
	f.Seek(3, 3) // BOM
	r := csv.NewReader(f)
	docsChan, err := csvparser.ParseОткрытыйРеестрТоварныхЗнаков(r)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	var batchSize = 10_000
	var docs = make([]models.ОткрытыйРеестрТоварныхЗнаков, 0, batchSize)
	for doc := range docsChan {
		count++
		docs = append(docs, doc)
		if len(docs) == batchSize {
			if err := pu.db.BatchОткрытыйРеестрТоварныхЗнаков(ctx, docs); err != nil {
				log.Fatal().Err(err).Send()
			}
			docs = make([]models.ОткрытыйРеестрТоварныхЗнаков, 0, batchSize)
			log.Info().Int("batch", count).Send()
		}
	}
	if err := pu.db.BatchОткрытыйРеестрТоварныхЗнаков(ctx, docs); err != nil {
		log.Fatal().Err(err).Send()
	}
	if err := pu.db.MarkFileIsParsed(ctx, "tz", fname); err != nil {
		log.Fatal().Err(err).Send()
	}
}

// ParseРосаккредитация - ...
// TODO: файл csv лежит в архиве z7
func (pu ParserUploader) ParseРосаккредитация(ctx context.Context, filename string) {
	if !strings.HasSuffix(filename, ".7z") {
		log.Warn().Str("filename", filename).Msg("file is not 7z")
		return
	}
	fname, err := pu.splitFilanemAndCheck(ctx, "rss", filename)
	if err != nil {
		log.Err(err).Send()
		return
	}
	binPath, err := exec.LookPath("7z")
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			log.Fatal().Msg("please install 7z")
		}
		log.Fatal().Err(err).Send()
	}
	var stdErrBud bytes.Buffer
	dir, _ := filepath.Split(filename)
	cmd := exec.Command(binPath, "x", filename, "-y", fmt.Sprintf("-o%s", dir))
	cmd.Stderr = &stdErrBud
	log.Info().Str("command", cmd.String()).Send()
	if err := cmd.Run(); err != nil {
		log.Fatal().Msg(stdErrBud.String())
	}

	matches, err := fs.Glob(os.DirFS(dir), "*.csv")
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	if len(matches) == 0 {
		log.Fatal().Msg("csv file not found")
	}
	for i := range matches {
		matches[i] = filepath.Join(dir, matches[i])
	}
	defer func() {
		// даже в случае ошибки где-либо - удаляем csv файл из папки
		if err := recover(); err != nil {
			for _, _filename := range matches {
				os.Remove(_filename)
			}
			log.Fatal().Any("err", err).Send()
		} else {
			for _, _filename := range matches {
				os.Remove(_filename)
			}
		}
	}()

	for _, _filename := range matches {
		log.Info().Str("filename", _filename).Msg("start parsing csv")
		f, err := os.Open(_filename)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		defer f.Close()
		docs, err := csvparser.ParseРосаккредитация(csv.NewReader(f))
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		log.Info().Str("filename", fname).Int("count", len(docs)).Msg("start insert batch")
		if err := pu.db.BatchInsertРосаккредитация(ctx, docs); err != nil {
			log.Fatal().Err(err).Send()
		}
	}
	if err := pu.db.MarkFileIsParsed(ctx, "rss", fname); err != nil {
		log.Fatal().Err(err).Send()
	}
}

// ParseRDS - Сведения из Реестра деклараций о соответствии
func (pu ParserUploader) ParseRDS(ctx context.Context, filename string) {
	if !strings.HasSuffix(filename, ".7z") {
		log.Warn().Str("filename", filename).Msg("file is not 7z")
		return
	}
	dir, _ := filepath.Split(filename)
	fname, err := pu.splitFilanemAndCheck(ctx, "rds", filename)
	if err != nil {
		log.Err(err).Send()
		return
	}
	binPath, err := exec.LookPath("7z")
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			log.Fatal().Msg("please install 7z")
		}
		log.Fatal().Err(err).Send()
	}
	var stdErrBud bytes.Buffer
	cmd := exec.Command(binPath, "x", filename, "-y", fmt.Sprintf("-o%s", dir))
	cmd.Stderr = &stdErrBud
	log.Info().Str("command", cmd.String()).Send()
	if err := cmd.Run(); err != nil {
		log.Fatal().Msg(stdErrBud.String())
	}

	matches, err := fs.Glob(os.DirFS(dir), "*.csv")
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	if len(matches) == 0 {
		log.Fatal().Msg("csv file not found")
	}
	for i := range matches {
		matches[i] = filepath.Join(dir, matches[i])
	}
	defer func() {
		// даже в случае ошибки где-либо - удаляем csv файл из папки
		if err := recover(); err != nil {
			for _, _filename := range matches {
				if err := os.Remove(_filename); err != nil {
					log.Error().Err(err).Send()
				}
			}
			log.Fatal().Any("err", err).Send()
		} else {
			for _, _filename := range matches {
				if err := os.Remove(_filename); err != nil {
					log.Error().Err(err).Send()
				}
			}
		}
	}()

	for _, _filename := range matches {
		log.Info().Str("filename", _filename).Msg("start parsing csv")
		f, err := os.Open(_filename)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		var count uint64
		var batchSize = 10_000
		var docs = make([]models.RDS, 0, batchSize)
		for doc := range csvparser.ParseRSD(ctx, csv.NewReader(f)) {
			count++
			docs = append(docs, doc)
			if len(docs) == batchSize {
				if err := pu.db.BatchInsertRDS(ctx, docs); err != nil {
					log.Fatal().Err(err).Send()
				}
				log.Info().Uint64("count", count).Msg("insert butch")
				docs = make([]models.RDS, 0, batchSize)
			}
		}
		if err := pu.db.BatchInsertRDS(ctx, docs); err != nil {
			log.Fatal().Err(err).Send()
		}
		log.Info().Uint64("count", count).Msg("insert butch")
		f.Close()
	}
	if err := pu.db.MarkFileIsParsed(ctx, "rds", fname); err != nil {
		log.Fatal().Err(err).Send()
	}
}

// ParseИсполнительныеПроизводстваВОтношенииЮридическихЛиц - ...
func (pu ParserUploader) ParseIpLegalList(ctx context.Context, filename string) {
	fname, err := pu.splitFilanemAndCheck(ctx, "iplegallist", filename)
	if err != nil {
		os.Remove(filename)
		log.Err(err).Send()
		return
	}
	var count int
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer f.Close()
	r := csv.NewReader(f)
	docsChan, err := csvparser.ParseИсполнительныеПроизводстваВОтношенииЮридическихЛиц(r)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Info().Msg("start update repaid flag for all to TRUE")
	if err := pu.db.RepaidFlagIpLegalList(ctx); err != nil {
		log.Fatal().Err(err).Msg("RepaidFlagIpLegalList FAIL")
	}
	log.Info().Msg("end update repaid flag for all are TRUE")
	var batchSize = 25_000
	var docs = make([]models.ИсполнительныеПроизводстваВОтношенииЮридическихЛиц, 0, batchSize)
	for doc := range docsChan {
		count++
		docs = append(docs, doc)
		if len(docs) == batchSize {
			if err := pu.db.BatchIpLegalList(ctx, docs); err != nil {
				log.Fatal().Err(err).Send()
			}
			docs = make([]models.ИсполнительныеПроизводстваВОтношенииЮридическихЛиц, 0, batchSize)
			log.Info().Int("batch", count).Send()
		}
	}
	if err := pu.db.BatchIpLegalList(ctx, docs); err != nil {
		log.Fatal().Err(err).Send()
	}
	if err := pu.db.UpdateIpLegalListDebTremainingBalanceZero(ctx); err != nil {
		log.Fatal().Err(err).Send()
	}
	if err := pu.db.MarkFileIsParsed(ctx, "iplegallist", fname); err != nil {
		log.Fatal().Err(err).Send()
	}
	if err := os.Remove(filename); err != nil {
		log.Fatal().Err(err).Msg("delete source file: iplegallist")
	}
}

// Parseiplegallistcomplete - https://opendata.fssp.gov.ru/7709576929-iplegallistcomplete
// Исполнительные производства в отношении юридических лиц,
// оконченные в соответствии с пунктами 3 и 4 части 1 статьи 46 и пунктами 6 и 7 части 1 статьи 47
// Федерального закона от 2 октября 2007 г. № 229-ФЗ «Об исполнительном производстве»
func (pu ParserUploader) ParseIpLegalListComplete(ctx context.Context, filename string) {
	fname, err := pu.splitFilanemAndCheck(ctx, "iplegallistcomplete", filename)
	if err != nil {
		log.Err(err).Send()
		return
	}
	var count int
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer f.Close()
	r := csv.NewReader(f)
	docsChan, err := csvparser.ParseIpLegalListComplete(r)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	var batchSize = 25_000
	var docs = make([]models.IpLegalListComplete, 0, batchSize)
	for doc := range docsChan {
		count++
		docs = append(docs, doc)
		if len(docs) == batchSize {
			if err := pu.db.BatchIpLegalListComplete(ctx, docs); err != nil {
				log.Fatal().Err(err).Send()
			}
			docs = make([]models.IpLegalListComplete, 0, batchSize)
			log.Info().Int("batch", count).Send()
		}
	}
	if err := pu.db.BatchIpLegalListComplete(ctx, docs); err != nil {
		log.Fatal().Err(err).Send()
	}
	if err := pu.db.MarkFileIsParsed(ctx, "iplegallistcomplete", fname); err != nil {
		log.Fatal().Err(err).Send()
	}
	if err := os.Remove(filename); err != nil {
		log.Fatal().Err(err).Msg("delete source file: iplegallistcomplete")
	}
}

// ParseНалоговыйРежимНалогоплательщика - ...
func (pu ParserUploader) ParseНалоговыйРежимНалогоплательщика(ctx context.Context, filename string) {
	fname, err := pu.splitFilanemAndCheck(ctx, "snr", filename)
	if err != nil {
		log.Err(err).Send()
		return
	}
	var count int
	zr, err := zip.OpenReader(filename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer zr.Close()
	var batchSize = 100_000
	var docs = make([]models.НалоговыйРежимНалогоплательщика, 0, batchSize)
	for doc := range xmlparser.ParseParallelIterFunc[models.НалоговыйРежимНалогоплательщика](ctx, zr) {
		count++
		docs = append(docs, doc)
		if len(docs) == batchSize {
			if err := pu.db.BatchInsertНалоговыйРежимНалогоплательщика(ctx, docs); err != nil {
				log.Fatal().Err(err).Send()
			}
			docs = make([]models.НалоговыйРежимНалогоплательщика, 0, batchSize)
			log.Info().Int("batch", count).Send()
		}
	}
	if err := pu.db.BatchInsertНалоговыйРежимНалогоплательщика(ctx, docs); err != nil {
		log.Fatal().Err(err).Send()
	}
	if err := pu.db.MarkFileIsParsed(ctx, "snr", fname); err != nil {
		log.Fatal().Err(err).Send()
	}
}

// ParseСведенияОСуммахНедоимки - ...
func (pu ParserUploader) ParseСведенияОСуммахНедоимки(ctx context.Context, filename string) {
	fname, err := pu.splitFilanemAndCheck(ctx, "debtam", filename)
	if err != nil {
		log.Err(err).Send()
		return
	}
	var count int
	zr, err := zip.OpenReader(filename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer zr.Close()
	var batchSize = 10_000
	var docs = make([]models.СведенияОСуммахНедоимки, 0, batchSize)
	for doc := range xmlparser.ParseParallelIterFunc[models.СведенияОСуммахНедоимки](ctx, zr) {
		count++
		docs = append(docs, doc)
		if len(docs) == batchSize {
			log.Debug().Msg("start BatchСведенияОСуммахНедоимки")
			if err := pu.db.BatchСведенияОСуммахНедоимки(ctx, docs); err != nil {
				log.Fatal().Err(err).Send()
			}
			docs = make([]models.СведенияОСуммахНедоимки, 0, batchSize)
			log.Info().Int("batch", count).Send()
		}
	}
	if err := pu.db.BatchСведенияОСуммахНедоимки(ctx, docs); err != nil {
		log.Fatal().Err(err).Send()
	}
	if err := pu.db.MarkFileIsParsed(ctx, "debtam", fname); err != nil {
		log.Fatal().Err(err).Send()
	}
}

// ParseСведенияОбУчастииВКонсГруппе - ...
func (pu ParserUploader) ParseСведенияОбУчастииВКонсГруппе(ctx context.Context, filename string) {
	fname, err := pu.splitFilanemAndCheck(ctx, "kgn", filename)
	if err != nil {
		log.Err(err).Send()
		return
	}
	var count int
	zr, err := zip.OpenReader(filename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer zr.Close()
	var batchSize = 100_000
	var docs = make([]models.СведенияОбУчастииВКонсГруппе, 0, batchSize)
	for doc := range xmlparser.ParseParallelIterFunc[models.СведенияОбУчастииВКонсГруппе](ctx, zr) {
		count++
		docs = append(docs, doc)
		if len(docs) == batchSize {
			if err := pu.db.BatchСведенияОбУчастииВКонсГруппе(ctx, docs); err != nil {
				log.Fatal().Err(err).Send()
			}
			docs = make([]models.СведенияОбУчастииВКонсГруппе, 0, batchSize)
			log.Info().Int("batch", count).Send()
		}
	}
	if err := pu.db.BatchСведенияОбУчастииВКонсГруппе(ctx, docs); err != nil {
		log.Fatal().Err(err).Send()
	}
	if err := pu.db.MarkFileIsParsed(ctx, "kgn", fname); err != nil {
		log.Fatal().Err(err).Send()
	}
}

// ParseСведенияОСреднесписочнойЧисленностиРаботников - ...
func (pu ParserUploader) ParseСведенияОСреднесписочнойЧисленностиРаботников(ctx context.Context, filename string) {
	fname, err := pu.splitFilanemAndCheck(ctx, "sshr", filename)
	if err != nil {
		log.Err(err).Send()
		return
	}
	var count int
	zr, err := zip.OpenReader(filename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer zr.Close()
	var batchSize = 100_000
	var docs = make([]models.СведенияОСреднесписочнойЧисленностиРаботников, 0, batchSize)
	for doc := range xmlparser.ParseParallelIterFunc[models.СведенияОСреднесписочнойЧисленностиРаботников](ctx, zr) {
		count++
		docs = append(docs, doc)
		if len(docs) == batchSize {
			if err := pu.db.BatchСведенияОСреднесписочнойЧисленностиРаботников(ctx, docs); err != nil {
				log.Fatal().Err(err).Send()
			}
			docs = make([]models.СведенияОСреднесписочнойЧисленностиРаботников, 0, batchSize)
			log.Info().Int("batch", count).Send()
		}
	}
	if err := pu.db.BatchСведенияОСреднесписочнойЧисленностиРаботников(ctx, docs); err != nil {
		log.Fatal().Err(err).Send()
	}
	if err := pu.db.MarkFileIsParsed(ctx, "sshr", fname); err != nil {
		log.Fatal().Err(err).Send()
	}
}

// ParseСведенияОбУплаченныхОрганизациейНалогов - ...
func (pu ParserUploader) ParseСведенияОбУплаченныхОрганизациейНалогов(ctx context.Context, filename string) {
	fname, err := pu.splitFilanemAndCheck(ctx, "paytax", filename)
	if err != nil {
		log.Err(err).Send()
		return
	}
	var count int
	zr, err := zip.OpenReader(filename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer zr.Close()
	var batchSize = 20_000
	var docs = make([]models.СведенияОбУплаченныхОрганизациейНалогов, 0, batchSize)
	for doc := range xmlparser.ParseParallelIterFunc[models.СведенияОбУплаченныхОрганизациейНалогов](ctx, zr) {
		count++
		docs = append(docs, doc)
		if len(docs) == batchSize {
			if err := pu.db.BatchСведенияОбУплаченныхОрганизациейНалогов(ctx, docs); err != nil {
				log.Fatal().Err(err).Send()
			}
			docs = make([]models.СведенияОбУплаченныхОрганизациейНалогов, 0, batchSize)
			log.Info().Int("batch", count).Send()
		}
	}
	if err := pu.db.BatchСведенияОбУплаченныхОрганизациейНалогов(ctx, docs); err != nil {
		log.Fatal().Err(err).Send()
	}
	if err := pu.db.MarkFileIsParsed(ctx, "paytax", fname); err != nil {
		log.Fatal().Err(err).Send()
	}
}

// ParseОКТМО - ...
func (pu ParserUploader) ParseОКТМО(ctx context.Context, filename string) {
	fname, err := pu.splitFilanemAndCheck(ctx, "oktmo", filename)
	if err != nil {
		log.Err(err).Send()
		return
	}
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer f.Close()
	r := csv.NewReader(charmap.Windows1251.NewDecoder().Reader(f))
	docsChan, err := csvparser.ParseОКТМО(r)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	var counter int
	var docs = make([]models.OKT, 0)
	for doc := range docsChan {
		docs = append(docs, doc)
		counter++
		log.Info().Int("counter", counter).Send()
	}
	if err := pu.db.BatchОКТМО(ctx, docs); err != nil {
		log.Fatal().Err(err).Send()
	}
	if err := pu.db.MarkFileIsParsed(ctx, "oktmo", fname); err != nil {
		log.Fatal().Err(err).Send()
	}
}

// ParseОКТАО - ...
func (pu ParserUploader) ParseОКТАО(ctx context.Context, filename string) {
	fname, err := pu.splitFilanemAndCheck(ctx, "okato", filename)
	if err != nil {
		log.Err(err).Send()
		return
	}
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer f.Close()
	r := csv.NewReader(charmap.Windows1251.NewDecoder().Reader(f))
	docsChan, err := csvparser.ParseОКТАО(r)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	var counter int
	var docs = make([]models.OKT, 0)
	for doc := range docsChan {
		docs = append(docs, doc)
		counter++
		log.Info().Int("counter", counter).Send()
	}
	if err := pu.db.BatchОКТАО(ctx, docs); err != nil {
		log.Fatal().Err(err).Send()
	}
	if err := pu.db.MarkFileIsParsed(ctx, "okato", fname); err != nil {
		log.Fatal().Err(err).Send()
	}
}

// ParseEGRUL - парсинг xml ЕГРЮЛ
func (pu ParserUploader) ParseEGRUL(ctx context.Context, filename string) {
	fname, err := pu.splitFilanemAndCheck(ctx, "egrul", filename)
	if err != nil {
		log.Err(err).Send()
		return
	}
	zr, err := zip.OpenReader(filename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer zr.Close()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var counter atomic.Int64
	tickPrinter(ctx, 10, &counter)
	var docsChan = xmlparser.ParseEGRUL(ctx, zr)
	var wg sync.WaitGroup
	wg.Add(pu.workersCount)
	for i := 0; i < pu.workersCount; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			var list = make([]egr.EGRUL, 0, 100)
			for doc := range docsChan {
				list = append(list, doc)
				if len(list) == 100 {
					if err := pu.db.InsertЕГРЮЛ(ctx, list); err != nil {
						log.Fatal().Err(err).Send()
					}
					counter.Add(int64(len(list)))
					list = make([]egr.EGRUL, 0, 100)
				}
			}
			if err := pu.db.InsertЕГРЮЛ(ctx, list); err != nil {
				log.Fatal().Err(err).Send()
			}
			counter.Add(int64(len(list)))
		}(&wg)
	}
	wg.Wait()
	if err := pu.db.MarkFileIsParsed(ctx, "egrul", fname); err != nil {
		log.Fatal().Err(err).Send()
	}
}

// ParseEGRIP - парсинг xml ЕГРИП
func (pu ParserUploader) ParseEGRIP(ctx context.Context, filename string) {
	fname, err := pu.splitFilanemAndCheck(ctx, "egrip", filename)
	if err != nil {
		log.Err(err).Send()
		return
	}
	zr, err := zip.OpenReader(filename)
	if err != nil {
		log.Err(err).Str("filename", filename).Send()
		return
	}
	defer zr.Close()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var counter atomic.Int64
	tickPrinter(ctx, 10, &counter)
	var docsChan = xmlparser.ParseEGRIP(ctx, zr)
	var wg sync.WaitGroup
	wg.Add(pu.workersCount)
	for i := 0; i < pu.workersCount; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			var list = make([]egr.EGRIP, 0, 100)
			for doc := range docsChan {
				list = append(list, doc)
				if len(list) == 100 {
					if err := pu.db.InsertЕГРИП(ctx, list); err != nil {
						log.Fatal().Err(err).Send()
					}
					counter.Add(int64(len(list)))
					list = make([]egr.EGRIP, 0, 100)
				}
			}
			if err := pu.db.InsertЕГРИП(ctx, list); err != nil {
				log.Fatal().Err(err).Send()
			}
			counter.Add(int64(len(list)))
		}(&wg)
	}
	wg.Wait()
	if err := pu.db.MarkFileIsParsed(ctx, "egrip", fname); err != nil {
		log.Fatal().Err(err).Send()
	}
}

// ParseПривлеченииУчастникаЗакупкиКАдминистративнойОтветственности - ...
func (pu ParserUploader) ParseПривлеченииУчастникаЗакупкиКАдминистративнойОтветственности(ctx context.Context, filename string) {
	fname, err := pu.splitFilanemAndCheck(ctx, "zakupki", filename)
	if err != nil {
		log.Err(err).Send()
		return
	}
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	var data []models.ПривлеченииУчастникаЗакупкиКАдминистративнойОтветственности
	switch {
	case strings.HasSuffix(filename, ".xlsx"):
		data, err = parsexlsx.ParseXLSX[models.ПривлеченииУчастникаЗакупкиКАдминистративнойОтветственности](f)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	case strings.HasSuffix(filename, ".xls"):
		data, err = parsexlsx.ParseXLS[models.ПривлеченииУчастникаЗакупкиКАдминистративнойОтветственности](f)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	default:
		log.Fatal().Msg("unknown file format")
	}
	for _, v := range data {
		if err := pu.db.InsertПривлеченииУчастникаЗакупкиКАдминистративнойОтветственности(ctx, v); err != nil {
			log.Fatal().Err(err).Send()
		}
	}
	if err := pu.db.MarkFileIsParsed(ctx, "zakupki", fname); err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Info().Int("count", len(data)).Msg("DONE")
}

// ParseHotels - добавление данных о гостиницах в БД
func (pu ParserUploader) ParseHotels(ctx context.Context, filename string) {
	fname, err := pu.splitFilanemAndCheck(ctx, "hotels", filename)
	if err != nil {
		log.Err(err).Send()
		return
	}
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer f.Close()
	var hs models.Hotels
	if err := json.NewDecoder(f).Decode(&hs); err != nil {
		log.Fatal().Err(err).Send()
	}
	if err := pu.db.TruncateHotels(ctx); err != nil {
		log.Fatal().Err(err).Send()
	}
	var ch = make(chan models.HotelData)
	var wg sync.WaitGroup
	for i := 0; i < pu.workersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for h := range ch {
				if err := pu.db.InsertHotelData(ctx, h); err != nil {
					log.Fatal().Err(err).Send()
				}
			}
		}()
	}
	go func() {
		defer close(ch)
		for i := range hs {
			ch <- hs[i]
		}
	}()
	wg.Wait()
	if err := pu.db.MarkFileIsParsed(ctx, "hotels", fname); err != nil {
		log.Fatal().Err(err).Send()
	}
}

func (pu ParserUploader) ParseUnscheduledInspectionsFGIS(ctx context.Context, filename string) {
	fname, err := pu.splitFilanemAndCheck(ctx, "fgis_unscheduled", filename)
	if err != nil {
		log.Err(err).Send()
		return
	}
	zr, err := zip.OpenReader(filename)
	if err != nil {
		log.Err(err).Str("filename", filename).Send()
		return
	}
	defer zr.Close()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var counter atomic.Int64
	tickPrinter(ctx, 10, &counter)
	var docsChan = xmlparser.ParseInspectionsFGIS(ctx, zr)
	var wg sync.WaitGroup
	wg.Add(pu.workersCount)
	for i := 0; i < pu.workersCount; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			var list = make([]models.InspectionFGIS, 0, 1000)
			for doc := range docsChan {
				list = append(list, doc)
				if len(list) == 1000 {
					if err := pu.db.InsertUnscheduledInspections(ctx, list); err != nil {
						log.Fatal().Err(err).Send()
					}
					counter.Add(int64(len(list)))
					list = make([]models.InspectionFGIS, 0, 1000)
				}
			}
			if err := pu.db.InsertUnscheduledInspections(ctx, list); err != nil {
				log.Fatal().Err(err).Send()
			}
			counter.Add(int64(len(list)))
		}(&wg)
	}
	wg.Wait()
	if err := pu.db.MarkFileIsParsed(ctx, "fgis_unscheduled", fname); err != nil {
		log.Fatal().Err(err).Send()
	}
	os.Remove(filename)
}

func (pu ParserUploader) ParseScheduledInspectionsFGIS(ctx context.Context, filename string) {
	fname, err := pu.splitFilanemAndCheck(ctx, "fgis_plan", filename)
	if err != nil {
		log.Err(err).Send()
		return
	}
	zr, err := zip.OpenReader(filename)
	if err != nil {
		log.Err(err).Str("filename", filename).Send()
		return
	}
	defer zr.Close()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var counter atomic.Int64
	tickPrinter(ctx, 10, &counter)
	var docsChan = xmlparser.ParseInspectionsFGIS(ctx, zr)
	var wg sync.WaitGroup
	wg.Add(pu.workersCount)
	for i := 0; i < pu.workersCount; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			var list = make([]models.InspectionFGIS, 0, 1000)
			for doc := range docsChan {
				list = append(list, doc)
				if len(list) == 1000 {
					if err := pu.db.InsertScheduledInspections(ctx, list); err != nil {
						log.Fatal().Err(err).Send()
					}
					counter.Add(int64(len(list)))
					list = make([]models.InspectionFGIS, 0, 1000)
				}
			}
			if err := pu.db.InsertScheduledInspections(ctx, list); err != nil {
				log.Fatal().Err(err).Send()
			}
			counter.Add(int64(len(list)))
		}(&wg)
	}
	wg.Wait()
	if err := pu.db.MarkFileIsParsed(ctx, "fgis_plan", fname); err != nil {
		log.Fatal().Err(err).Send()
	}
	os.Remove(filename)
}

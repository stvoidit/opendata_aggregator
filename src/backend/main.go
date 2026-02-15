// Package main - entrypoint
package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/fs"
	"opendataaggregator/config"
	"opendataaggregator/downloader"
	"opendataaggregator/server"
	"opendataaggregator/store"
	"opendataaggregator/uploader"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/invopop/jsonschema"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	configFilename string
	sourceType     string
	sourceLink     string
	serve          bool
	genJsonschema  bool
	parseType      string
	workersCount   int
	// download       bool
)

const (
	longHelp = `Обозначения источников для флагов "--parse" и "--source":
  balance - Отчет об бухгалтерском балансе
  iplegallist - Исполнительные производства в отношении юридических лиц
  iplegallistcomplete - Оконченные производства в отношении юридических лиц
  snr - Сведения о специальных налоговых режимах, применяемых налогоплательщиками
  taxoffence - Сведения о налоговых правонарушениях и мерах ответственности за их совершение
  okved2 - Общероссийский классификатор видов экономической деятельности (ОКВЭД2)
  okato - Общероссийский классификатор объектов административно-территориального деления (ОКАТО)
  oktmo - Общероссийский классификатор территорий муниципальных образований (ОКТМО)
  otz - Открытый реестр общеизвестных в Российской Федерации товарных знаков
  tz - Открытый реестр товарных знаков и знаков обслуживания Российской Федерации
  zakupki - Информация о привлечении участника закупки к административной ответственности по ст. 19.28 КоАП
  registerdisqualified - Реестр дисквалифицированных лиц
  rsmp - Единый реестр субъектов малого и среднего предпринимательства
  rss - Сведения из Реестра сертификатов соответствия
  sshr - Сведения о среднесписочной численности работников организации
  debtam - Сведения о суммах недоимки и задолженности по пеням и штрафам
  paytax - Сведения об уплаченных организацией суммах налогов и сборов
  kgn - Сведения об участии в консолидированной группе налогоплательщиков
  hotels - Федералльный реестр туристских объектов
  egrul - Единый государственный реестр юридических лиц
  egrip - Единый государственный реестр индивидуальных предпринимателей
  fgis_unscheduled - ФГИС внеплановые проверки
  fgis_plan - ФГИС плановые проверки
		`
	exampleUse = `  ./opendataaggregator downloader --source=egrul --config=config.toml
  ./opendataaggregator parser --parse=egrul --config=config.toml
  ./opendataaggregator server --serve --config=config.toml`
)

var (
	rootCmd     = &cobra.Command{Use: "", Example: exampleUse, Long: longHelp, Run: userCmd}
	serverCmd   = &cobra.Command{Use: "server", Short: "Запуск сервера", Run: userCmd}
	downloadCmd = &cobra.Command{Use: "downloader", Short: "Скачивание исходников данных", Run: userCmd}
	parserCmd   = &cobra.Command{Use: "parser", Short: "Парсер исходников данных", Run: userCmd}
	scandirCmd  = &cobra.Command{Use: "scan", Short: "Проверка файлов в папке, если скачивались отдельно", Run: userCmd}
)

var enumParserTypeName = [...]string{
	"balance",
	"iplegallist",
	"iplegallistcomplete",
	"snr",
	"taxoffence",
	"okved2",
	"okato",
	"oktmo",
	"otz",
	"tz",
	"zakupki",
	"registerdisqualified",
	"rsmp",
	"rss",
	"sshr",
	"debtam",
	"paytax",
	"kgn",
	"hotels",
	"egrul",
	"egrip",
	"fgis_unscheduled",
	"fgis_plan",
}

func init() {
	initLogger()
	rootCmd.Version = config.VersionTag + " => " + config.CommitHash
	rootCmd.PersistentFlags().StringVarP(&configFilename, "config", "c", "config.toml", "файл конфига в формате .toml")
	serverCmd.Flags().BoolVar(&serve, "serve", false, "Запуск сервера")
	serverCmd.Flags().BoolVar(&genJsonschema, "jsonschema", false, "Сгенерировать jsonschema файл ответа")
	downloadCmd.Flags().StringVarP(&sourceType, "source", "s", "all", "Номер источника")
	scandirCmd.Flags().StringVarP(&sourceType, "source", "s", "all", "Номер источника")
	parserCmd.Flags().StringVarP(&parseType, "source", "s", "all", "Номер источника")
	downloadCmd.Flags().StringVarP(&sourceLink, "link", "l", "", "Ссылка на источник")
	parserCmd.Flags().StringVarP(&parseType, "parse", "p", "all", "Выбор функции парсера")
	parserCmd.Flags().IntVarP(&workersCount, "workers", "w", runtime.NumCPU(), "кол-во воркеров для парсера")
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(parserCmd)
	rootCmd.AddCommand(scandirCmd)
}

func initLogger() {
	log.Logger = zerolog.New(zerolog.MultiLevelWriter(
		zerolog.ConsoleWriter{
			Out:         os.Stdout,
			NoColor:     true,
			TimeFormat:  time.RFC3339,
			FormatLevel: func(i any) string { return strings.ToUpper(i.(string)) },
		})).With().Caller().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func main() {
	rootCmd.Execute()
}

func interruptContext(ctx context.Context) (interruptCtx context.Context) {
	var once sync.Once
	var cancel context.CancelFunc
	once.Do(func() {
		interruptCtx, cancel = context.WithCancel(ctx)
		var sigChan = make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			defer cancel()
			sig := <-sigChan
			close(sigChan)
			log.Info().Str("signal", sig.String()).Send()
		}()
	})
	return
}

func userCmd(cmd *cobra.Command, _ []string) {
	var ctx = interruptContext(context.Background())
	cnf, err := config.LoadConfigFromFile(configFilename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	zerolog.SetGlobalLevel(zerolog.Level(cnf.LogLevel))
	log.Info().Stringer("loglevel", zerolog.Level(cnf.LogLevel)).Send()
	switch cmd.Use {
	case "server":
		switch {
		case serve:
			if err := server.NewApplication(cnf).ListenAndServe(ctx); err != nil {
				log.Fatal().Err(err).Send()
			}
		case genJsonschema:
			createJsonschema()
		}
	case "downloader":
		downloadSource(ctx, cnf)
	case "parser":
		parseSource(ctx, cnf)
	case "scan":
		scanDirSource(ctx, cnf)
	default:
		log.Fatal().Msg("unknown CLI command")
	}
}

func createJsonschema() {
	schema := new(jsonschema.Reflector).Reflect(new(store.AllTypes))
	f, err := os.Create("TotalResult.json")
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer f.Close()
	e := json.NewEncoder(f)
	e.SetIndent("", "  ")
	if err := e.Encode(*schema); err != nil {
		log.Fatal().Err(err).Send()
	}
}

func filterSourceType() (list []string) {
	if parseType == "all" {
		return enumParserTypeName[2:]
	}
	return append(list, parseType)
}

func downloadSource(ctx context.Context, cnf *config.Config) {
	db, err := store.NewDB(cnf, workersCount)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer db.Close()
	c := downloader.NewDownloadClient(db, cnf)
	defer c.Close()
	if sourceLink != "" {
		if err := c.Download(ctx, sourceType, sourceLink); err != nil {
			log.Error().Err(err).Send()
		}
		return
	}
	doDownload := func(st string) {
		sourceURL, ok := cnf.Sources[st]
		if !ok {
			log.Warn().Str("sourceType", st).Msg("sourceType not found")
			return
		}
		log.Info().Str("sourceType", st).Str("SourceURL", sourceURL).Msg("DOWNLOAD SOURCE")
		if err := c.Download(ctx, st, sourceURL); err != nil {
			log.Error().Err(err).Send()
		}
	}
	switch sourceType {
	case "all":
		for _, st := range enumParserTypeName {
			doDownload(st)
		}
	default:
		doDownload(sourceType)
	}
}

func parseSource(ctx context.Context, cnf *config.Config) {
	db, err := store.NewDB(cnf, workersCount)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer db.Close()
	var up = uploader.NewParserUploader(db, workersCount)
	var start = time.Now()
	var counterF int
	for _, st := range filterSourceType() {
		if err := fs.WalkDir(os.DirFS(cnf.FS.DownloadFolder), st, func(_path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}
			filepath := path.Join(cnf.FS.DownloadFolder, _path)
			switch st {
			default:
				log.Fatal().Msg("unknown parser type")
			case "balance":
				// ! FTP
				up.ParseБухОтчетность(ctx, filepath)
			case "egrul":
				// https://ftp.egrul.nalog.ru/?dir=EGRUL_406
				up.ParseEGRUL(ctx, filepath)
			case "egrip":
				// https://ftp.egrul.nalog.ru/?dir=EGRIP_405
				up.ParseEGRIP(ctx, filepath)
			case "iplegallist":
				// https://opendata.fssp.gov.ru/7709576929-iplegallist
				up.ParseIpLegalList(ctx, filepath)
			case "iplegallistcomplete":
				// https://opendata.fssp.gov.ru/7709576929-iplegallistcomplete
				up.ParseIpLegalListComplete(ctx, filepath)
			case "snr":
				// https://www.nalog.gov.ru/opendata/7707329152-snr/
				up.ParseНалоговыйРежимНалогоплательщика(ctx, filepath)
			case "taxoffence":
				// https://www.nalog.gov.ru/opendata/7707329152-taxoffence/
				up.ParseНалоговыхПравонарушенияхИМерахОтветственности(ctx, filepath)
			case "okved2":
				// https://rosstat.gov.ru/opendata/7708234640-okved2
				up.ParseОКВЭД(ctx, filepath)
			case "okato":
				// https://rosstat.gov.ru/opendata/7708234640-okato
				up.ParseОКТАО(ctx, filepath)
			case "oktmo":
				// https://rosstat.gov.ru/opendata/7708234640-oktmo
				up.ParseОКТМО(ctx, filepath)
			case "otz":
				// https://rospatent.gov.ru/opendata/7730176088-otz
				up.ParseОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков(ctx, filepath)
			case "tz":
				// https://rospatent.gov.ru/opendata/7730176088-tz
				up.ParseОткрытыйРеестрТоварныхЗнаков(ctx, filepath)
			case "zakupki":
				// https://zakupki.gov.ru/epz/main/public/document/view.html?searchString=&sectionId=2369&strictEqual=false
				up.ParseПривлеченииУчастникаЗакупкиКАдминистративнойОтветственности(ctx, filepath)
			case "registerdisqualified":
				// https://www.nalog.gov.ru/opendata/7707329152-registerdisqualified/
				up.ParseРеестрДисквалифицированныхЛиц(ctx, filepath)
			case "rsmp":
				// https://www.nalog.gov.ru/opendata/7707329152-rsmp/
				up.ParseРеестрСубъектовМалогоИСреднегоПредпринимательства(ctx, filepath)
			case "rss":
				// https://fsa.gov.ru/opendata/7736638268-rss/
				up.ParseРосаккредитация(ctx, filepath)
			case "sshr":
				// https://www.nalog.gov.ru/opendata/7707329152-sshr2019/
				up.ParseСведенияОСреднесписочнойЧисленностиРаботников(ctx, filepath)
			case "debtam":
				// https://www.nalog.gov.ru/opendata/7707329152-debtam/
				up.ParseСведенияОСуммахНедоимки(ctx, filepath)
			case "paytax":
				// https://www.nalog.gov.ru/opendata/7707329152-paytax/
				up.ParseСведенияОбУплаченныхОрганизациейНалогов(ctx, filepath)
			case "kgn":
				// https://www.nalog.gov.ru/opendata/7707329152-kgn/
				up.ParseСведенияОбУчастииВКонсГруппе(ctx, filepath)
			case "hotels":
				up.ParseHotels(ctx, filepath)
			case "rds":
				// https://fsa.gov.ru/opendata/7736638268-rds/
				up.ParseRDS(ctx, filepath)
			case "fgis_unscheduled":
				up.ParseUnscheduledInspectionsFGIS(ctx, filepath)
			case "fgis_plan":
				up.ParseScheduledInspectionsFGIS(ctx, filepath)
			}
			counterF++
			log.Info().Int("file counter", counterF).Str("path", _path).Send()
			return err
		}); err != nil {
			log.Fatal().Err(err).Send()
		}
	}
	log.Info().Str("completion", time.Since(start).String()).Send()
}

func scanDirSource(ctx context.Context, cnf *config.Config) {
	db, err := store.NewDB(cnf, workersCount)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer db.Close()
	if err := fs.WalkDir(os.DirFS(cnf.FS.DownloadFolder), sourceType, func(_path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		filepath := path.Join(cnf.FS.DownloadFolder, _path)
		exists, err := db.ExistsDownloadFile(ctx, sourceType, d.Name())
		if err != nil {
			return err
		}
		if exists {
			return nil
		}
		b, err := os.ReadFile(filepath)
		if err != nil {
			return err
		}
		h := sha256.New()
		h.Write(b)
		hexstr := hex.EncodeToString(h.Sum(nil))
		df := config.DownloadedFile{
			Filename:   d.Name(),
			Filepath:   filepath,
			SHA265SUM:  hexstr,
			SourceType: sourceType,
			SourceLink: "manual",
		}
		if err := db.InsertDownloadedFileInfo(ctx, &df); err != nil {
			return err
		}
		log.Info().Str("filename", d.Name()).Str("hash", hexstr).Msg("new file found")
		return err
	}); err != nil {
		log.Fatal().Err(err).Send()
	}
}

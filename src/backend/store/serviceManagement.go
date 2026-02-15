package store

import (
	"context"
	"errors"
	"fmt"
	"io"
	"opendataaggregator/config"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/xuri/excelize/v2"
	"golang.org/x/sync/errgroup"
)

// InsertDownloadedFileInfo - фиксирование информации о скаченном файле
func (db *DB) InsertDownloadedFileInfo(ctx context.Context, df *config.DownloadedFile) error {
	const q = `INSERT
	INTO
	service_management.source_files
	(
		source_type
		, source_link
		, filename
		, sha256sum
		, downloaded
	)
	VALUES(
		$1
		, $2
		, $3
		, $4
		, TRUE
	)`
	if _, err := db.pool.Exec(ctx, q, df.SourceType, df.SourceLink, df.Filename, df.SHA265SUM); err != nil {
		if strings.Contains(err.Error(), "(SQLSTATE 23505)") {
			return nil
		}
		return err
	}
	return nil
}

// ExistsDownloadFile - проверка наличия файла, чтобы не скачивать лишний раз
func (db *DB) ExistsDownloadFile(ctx context.Context, sourcetype string, filename string) (exists bool, err error) {
	const q = `SELECT EXISTS (SELECT 1 FROM service_management.source_files sf WHERE sf.source_type = $1 AND sf.filename = $2)`
	err = db.pool.QueryRow(ctx, q, sourcetype, filename).Scan(&exists)
	return
}

// MarkFileIsParsed - пометить файл как распаршенный
func (db *DB) MarkFileIsParsed(ctx context.Context, sourcetype string, filename string) error {
	const q = `UPDATE service_management.source_files SET uploaded=TRUE, task_datetime=now() WHERE source_type=$1 AND filename=$2`
	_, err := db.pool.Exec(ctx, q, sourcetype, filename)
	return err
}

// ChecFileIsParsed - проверка был ли файл уже распаршен
func (db *DB) ChecFileIsParsed(ctx context.Context, sourcetype string, filename string) (exists bool, err error) {
	const q = `SELECT sf.uploaded FROM service_management.source_files sf WHERE sf.source_type = $1 AND sf.filename = $2`
	err = db.pool.QueryRow(ctx, q, sourcetype, filename).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return
}

// SelectDatabaseStat - статистика по данным таблиц в БД открытых источников
func (db *DB) SelectDatabaseStat(ctx context.Context) (dbs DatabaseStat, err error) {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return db.pool.
			QueryRow(ctx,
				`SELECT
					count(*) FILTER (WHERE es.is_legal IS TRUE) AS is_legals
					, count(*) FILTER (WHERE es.is_legal IS FALSE) AS is_ip
				FROM
					egr.egr_search AS es`).Scan(&dbs.StatEGR.EGRUL, &dbs.StatEGR.EGRIP)
	})
	g.Go(func() error {
		return db.pool.
			QueryRow(ctx, `SELECT max(es.date_discharge)::text FROM egr.egr_search AS es`).
			Scan(&dbs.StatEGR.LastDateDischarge)
	})
	g.Go(func() error {
		return db.pool.
			QueryRow(ctx, `SELECT jsonb_agg(b.st) FROM (SELECT jsonb_build_object('year', bmv."year" , 'count', count(*)) AS st FROM accounting_statements.balance_mat_view AS bmv GROUP BY bmv."year" ORDER BY bmv."year" DESC) b`).
			Scan(&dbs.StatBalance.CountYears)
	})
	g.Go(func() error {
		return db.pool.
			QueryRow(ctx, `SELECT max(sf.task_datetime) FROM service_management.source_files AS sf WHERE sf.source_type = 'balance'`).
			Scan(&dbs.StatBalance.LastDocDate)
	})
	g.Go(func() error {
		return db.pool.
			QueryRow(ctx, `SELECT jsonb_agg(row_to_json(su.*)) FROM (SELECT sf.source_type , max(sf.task_datetime) AS datetime FROM service_management.source_files AS sf GROUP BY sf.source_type ORDER BY sf.source_type) su`).
			Scan(&dbs.LastUpdatesSources)
	})
	g.Go(func() error {
		return db.pool.
			QueryRow(ctx, `SELECT count(*) FROM public.hotels AS h`).
			Scan(&dbs.StatsHotels.TotalCount)
	})
	g.Go(func() error {
		return db.pool.
			QueryRow(ctx, `SELECT jsonb_agg(hj.*) FROM (SELECT h."type", count(*)FROM public.hotels AS h GROUP BY h."type" ORDER BY h."type") hj`).
			Scan(&dbs.StatsHotels.CountHotelsTypes)
	})
	g.Go(func() error {
		return db.pool.
			QueryRow(ctx, `SELECT count(*) FROM public.росаккредитация AS r`).
			Scan(&dbs.StatsRossAccreditation.TotalCount)
	})
	g.Go(func() error {
		return db.pool.
			QueryRow(ctx, `SELECT json_agg(row_to_json(cs.*)) FROM (SELECT r.cert_status, count(*) FROM public.росаккредитация AS r GROUP BY r.cert_status ORDER BY r.cert_status) cs`).
			Scan(&dbs.StatsRossAccreditation.RossAccreditationStatuses)
	})
	g.Go(func() error {
		return db.pool.
			QueryRow(ctx, `SELECT count(*), sum(npr.сумштраф) FROM public.налоговыеправонарушенияиштрафы AS npr`).
			Scan(&dbs.StatsTaxOffenses.TotalCount, &dbs.StatsTaxOffenses.TotalSum)
	})
	g.Go(func() error {
		return db.pool.
			QueryRow(ctx, `SELECT json_agg(row_to_json(nr.*)) FROM (SELECT EXTRACT(YEAR FROM npr.датасост)::int4 AS "year",  sum(npr.сумштраф), count(*) FROM public.налоговыеправонарушенияиштрафы AS npr GROUP BY 1 ORDER BY 1 DESC) AS nr`).
			Scan(&dbs.StatsTaxOffenses.SumsByYears)
	})
	g.Go(func() error {
		return db.pool.
			QueryRow(ctx, `SELECT count(*) FROM public.открытыйреестртоварныхзнаков`).
			Scan(&dbs.StatsRegisterOfTrademarks.OpenRegistryCount)
	})
	g.Go(func() error {
		return db.pool.
			QueryRow(ctx, `SELECT count(*) FROM public.реестробщеизвестныхтоварныхзнак`).
			Scan(&dbs.StatsRegisterOfTrademarks.WellKnownRegistryCount)
	})
	g.Go(func() error {
		return db.pool.
			QueryRow(ctx, `SELECT count(*) FROM public.исппроизввотнюрлиц`).
			Scan(&dbs.StatsFSSP.IpLegalList)
	})
	g.Go(func() error {
		return db.pool.
			QueryRow(ctx, `SELECT count(*) FROM public.iplegallistcomplete`).
			Scan(&dbs.StatsFSSP.IpLegalListComplite)
	})
	g.Go(func() error {
		return db.pool.
			QueryRow(ctx, `SELECT count(*) FROM public.субъектымалогоисреднегопредприн`).
			Scan(&dbs.StatsSMP.TotalCount)
	})
	g.Go(func() error {
		return db.pool.
			QueryRow(ctx, `SELECT count(*), sum(общсумнедоим) FROM public.сведенияосуммахнедоимки`).
			Scan(&dbs.StatsDEBTAM.TotalCount, &dbs.StatsDEBTAM.TotalSum)
	})
	g.Go(func() error {
		return db.pool.
			QueryRow(ctx, `SELECT count(*) FROM fgis.unscheduled_inspections`).
			Scan(&dbs.StatsFGIS.UnscheduledCount)
	})
	g.Go(func() error {
		return db.pool.
			QueryRow(ctx, `SELECT count(*) FROM fgis.scheduled_inspections`).
			Scan(&dbs.StatsFGIS.ScheduledCount)
	})

	g.Go(func() error {
		return db.pool.
			QueryRow(ctx, `SELECT count(*) FROM public.сведенияосреднчислработников`).
			Scan(&dbs.StatAvgEmployesNumber.Count)
	})

	g.Go(func() error {
		return db.pool.
			QueryRow(ctx, `SELECT
			count(*) AS total
			, count(*) FILTER (WHERE rn."есхн" )
			, count(*) FILTER (WHERE rn."усн" )
			, count(*) FILTER (WHERE rn."енвд" )
			, count(*) FILTER (WHERE rn."срп" )
		FROM
			public."режимналогоплательщика" AS rn`).
			Scan(&dbs.StatsTaxRegime.Count, &dbs.StatsTaxRegime.ЕСХН, &dbs.StatsTaxRegime.УСН, &dbs.StatsTaxRegime.ЕНВД, &dbs.StatsTaxRegime.СРП)
	})
	if err := g.Wait(); err != nil {
		return dbs, err
	}
	dbs.StatEGR.TotalCount = dbs.StatEGR.EGRIP + dbs.StatEGR.EGRUL
	for _, y := range dbs.StatBalance.CountYears {
		dbs.StatBalance.TotalCount += y.Count
	}
	return
}

type Hotels []Hotel

func (hs Hotels) GetCategories() []string {
	var unc = make(map[string]struct{})
	for i := range hs {
		for _, room := range hs[i].HotelRooms {
			category := room.Category
			if len(category) == 0 {
				category = "Без категории"
			}
			unc[fmt.Sprintf("%s (комнаты)", category)] = struct{}{}
			unc[fmt.Sprintf("%s (места)", category)] = struct{}{}
		}
	}
	var list = make([]string, 0, len(unc))
	for k := range unc {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}

func ToExcel(w io.Writer, hs Hotels) error {
	var cord = func(col, row int) string {
		coordinate, _ := excelize.CoordinatesToCellName(col, row)
		return coordinate
	}
	f := excelize.NewFile()
	const sheetname = "hotels"
	f.SetSheetName("Sheet1", sheetname)
	var headers = []string{
		"Порядковый номер в Федеральном перечне",
		"Вид",
		"Полное наименование классифицированного объекта",
		"Cокращенное наименование классифицированного объекта",
		"Наименование юридического лица/индивидуального предпринимателя",
		"Регион",
		"ИНН",
		"ОГРН/ОГРНИП",
		"Адрес места нахождения",
		"Телефон",
		"Факс",
		"E-mail",
		"Адрес сайта",
		"Присвоенная категория",
		"Регистрационный номер",
		"Регистрационный номер свидетельства",
		"Дата выдачи свидетельства",
		"Срок действия свидетельства",
	}
	categories := hs.GetCategories()
	headers = append(headers, categories...)
	categoriesMap := make(map[string]int)
	for _, cat := range categories {
		for i, h := range headers {
			if h == cat {
				categoriesMap[cat] = i
				break
			}
		}
	}
	for i, h := range headers {
		if err := f.SetCellStr(sheetname, cord(i+1, 1), h); err != nil {
			return err
		}
	}
	for n, h := range hs {
		f.SetCellStr(sheetname, cord(1, n+2), h.FederalNumber)
		f.SetCellStr(sheetname, cord(2, n+2), h.Type)
		f.SetCellStr(sheetname, cord(3, n+2), h.Fullname)
		f.SetCellStr(sheetname, cord(4, n+2), h.ShortName)
		f.SetCellStr(sheetname, cord(5, n+2), h.Owner)
		f.SetCellStr(sheetname, cord(6, n+2), h.Region)
		f.SetCellStr(sheetname, cord(7, n+2), h.INN)
		f.SetCellStr(sheetname, cord(8, n+2), h.OGRN)
		f.SetCellStr(sheetname, cord(9, n+2), h.Address)
		f.SetCellStr(sheetname, cord(10, n+2), h.Phone)
		f.SetCellStr(sheetname, cord(11, n+2), h.Fax)
		f.SetCellStr(sheetname, cord(12, n+2), h.Email)
		f.SetCellStr(sheetname, cord(13, n+2), h.Site)
		f.SetCellStr(sheetname, cord(14, n+2), h.HotelClassification[0].Category)
		f.SetCellStr(sheetname, cord(15, n+2), h.HotelClassification[0].RegistrationNumber)
		f.SetCellStr(sheetname, cord(16, n+2), h.HotelClassification[0].LicenseNumber)
		f.SetCellStr(sheetname, cord(17, n+2), h.HotelClassification[0].DateIssued)
		f.SetCellStr(sheetname, cord(18, n+2), h.HotelClassification[0].DateEnd)
		for _, room := range h.HotelRooms {
			category := room.Category
			if len(category) == 0 {
				category = "Без категории"
			}
			if colN, ok := categoriesMap[fmt.Sprintf("%s (комнаты)", category)]; ok {
				f.SetCellInt(sheetname, cord(colN+1, n+2), room.Rooms)
			}
			if colN, ok := categoriesMap[fmt.Sprintf("%s (места)", category)]; ok {
				f.SetCellInt(sheetname, cord(colN+1, n+2), room.Seats)
			}
		}
	}
	return f.Write(w)
}

func (db *DB) LastModUpdates(ctx context.Context) (t time.Time, err error) {
	err = db.pool.QueryRow(ctx, `SELECT max(sf.task_datetime) FROM service_management.source_files AS sf WHERE downloaded IS TRUE AND uploaded IS TRUE`).Scan(&t)
	t = t.Truncate(time.Second)
	return
}

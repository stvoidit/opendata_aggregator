package store

import "context"

func (db *DB) SelectStatEGR(ctx context.Context) (v StatEGR, err error) {
	const q = `
	SELECT
		count(*) FILTER (WHERE es.is_legal IS TRUE) AS is_legals
		, count(*) FILTER (WHERE es.is_legal IS FALSE) AS is_ip
	FROM
		egr.egr_search AS es`
	err = db.pool.QueryRow(ctx, q).Scan(&v.EGRUL, &v.EGRIP)
	if err != nil {
		return
	}
	err = db.pool.
		QueryRow(ctx, `SELECT max(es.date_discharge)::text FROM egr.egr_search AS es`).
		Scan(&v.LastDateDischarge)
	if err != nil {
		return
	}
	v.TotalCount = v.EGRIP + v.EGRUL
	return
}
func (db *DB) SelectStatBalance(ctx context.Context) (v StatBalance, err error) {
	err = db.pool.
		QueryRow(ctx, `SELECT jsonb_agg(b.st) FROM (SELECT jsonb_build_object('year', bmv."year" , 'count', count(*)) AS st FROM accounting_statements.balance_mat_view AS bmv GROUP BY bmv."year" ORDER BY bmv."year" DESC) b`).
		Scan(&v.CountYears)
	if err != nil {
		return
	}
	err = db.pool.
		QueryRow(ctx, `SELECT max(sf.task_datetime) FROM service_management.source_files AS sf WHERE sf.source_type = 'balance'`).
		Scan(&v.LastDocDate)
	if err != nil {
		return
	}
	for _, y := range v.CountYears {
		v.TotalCount += y.Count
	}
	return
}
func (db *DB) SelectLastUpdatesSources(ctx context.Context) (v []LastUpdateSource, err error) {
	err = db.pool.
		QueryRow(ctx, `SELECT jsonb_agg(row_to_json(su.*)) FROM (SELECT sf.source_type , max(sf.task_datetime) AS datetime FROM service_management.source_files AS sf GROUP BY sf.source_type ORDER BY sf.source_type) su`).
		Scan(&v)
	return
}
func (db *DB) SelectStatsHotels(ctx context.Context) (v StatsHotels, err error) {
	err = db.pool.
		QueryRow(ctx, `SELECT jsonb_agg(hj.*) FROM (SELECT h."type", count(*)FROM public.hotels AS h GROUP BY h."type" ORDER BY h."type") hj`).
		Scan(&v.CountHotelsTypes)
	if err != nil {
		return
	}
	err = db.pool.
		QueryRow(ctx, `SELECT count(*) FROM public.росаккредитация AS r`).
		Scan(&v.TotalCount)
	return
}
func (db *DB) SelectStatsRossAccreditation(ctx context.Context) (v StatsRossAccreditation, err error) {
	err = db.pool.
		QueryRow(ctx, `SELECT count(*) FROM public.росаккредитация AS r`).
		Scan(&v.TotalCount)
	if err != nil {
		return
	}
	err = db.pool.
		QueryRow(ctx, `SELECT json_agg(row_to_json(cs.*)) FROM (SELECT r.cert_status, count(*) FROM public.росаккредитация AS r GROUP BY r.cert_status ORDER BY r.cert_status) cs`).
		Scan(&v.RossAccreditationStatuses)
	return
}
func (db *DB) SelectStatsTaxOffenses(ctx context.Context) (v TaxOffenses, err error) {
	err = db.pool.
		QueryRow(ctx, `SELECT count(*), sum(npr.сумштраф) FROM public.налоговыеправонарушенияиштрафы AS npr`).
		Scan(&v.TotalCount, &v.TotalSum)
	if err != nil {
		return
	}
	err = db.pool.
		QueryRow(ctx, `SELECT json_agg(row_to_json(nr.*)) FROM (SELECT EXTRACT(YEAR FROM npr.датасост)::int4 AS "year",  sum(npr.сумштраф), count(*) FROM public.налоговыеправонарушенияиштрафы AS npr GROUP BY 1 ORDER BY 1 DESC) AS nr`).
		Scan(&v.SumsByYears)
	return
}
func (db *DB) SelectStatsRegisterOfTrademarks(ctx context.Context) (v StatsRegisterOfTrademarks, err error) {
	err = db.pool.
		QueryRow(ctx, `SELECT count(*) FROM public.открытыйреестртоварныхзнаков`).
		Scan(&v.OpenRegistryCount)
	if err != nil {
		return
	}
	err = db.pool.
		QueryRow(ctx, `SELECT count(*) FROM public.реестробщеизвестныхтоварныхзнак`).
		Scan(&v.WellKnownRegistryCount)
	return
}
func (db *DB) SelectStatsFSSP(ctx context.Context) (v StatsFSSP, err error) {
	err = db.pool.
		QueryRow(ctx, `SELECT count(*) FROM public.исппроизввотнюрлиц`).
		Scan(&v.IpLegalList)
	if err != nil {
		return
	}
	err = db.pool.
		QueryRow(ctx, `SELECT count(*) FROM public.iplegallistcomplete`).
		Scan(&v.IpLegalListComplite)
	return
}
func (db *DB) SelectStatsSMP(ctx context.Context) (v StatsSMP, err error) {
	err = db.pool.
		QueryRow(ctx, `SELECT count(*) FROM public.субъектымалогоисреднегопредприн`).
		Scan(&v.TotalCount)
	return
}
func (db *DB) SelectStatsDEBTAM(ctx context.Context) (v StatsDEBTAM, err error) {
	err = db.pool.
		QueryRow(ctx, `SELECT count(*), sum(общсумнедоим) FROM public.сведенияосуммахнедоимки`).
		Scan(&v.TotalCount, &v.TotalSum)
	return
}
func (db *DB) SelectStatsFGIS(ctx context.Context) (v StatsFGIS, err error) {
	err = db.pool.
		QueryRow(ctx, `SELECT count(*) FROM fgis.scheduled_inspections`).
		Scan(&v.ScheduledCount)
	if err != nil {
		return
	}
	err = db.pool.
		QueryRow(ctx, `SELECT count(*) FROM fgis.unscheduled_inspections`).
		Scan(&v.UnscheduledCount)
	return
}
func (db *DB) SelectStatAvgEmployesNumber(ctx context.Context) (v StatAvgEmployesNumber, err error) {
	err = db.pool.
		QueryRow(ctx, `SELECT count(*) FROM public.сведенияосреднчислработников`).
		Scan(&v.Count)
	return
}
func (db *DB) SelectStatsTaxRegime(ctx context.Context) (v StatsTaxRegime, err error) {
	err = db.pool.
		QueryRow(ctx, `
		SELECT
			count(*) AS total
			, count(*) FILTER (WHERE rn."есхн" )
			, count(*) FILTER (WHERE rn."усн" )
			, count(*) FILTER (WHERE rn."енвд" )
			, count(*) FILTER (WHERE rn."срп" )
		FROM
			public."режимналогоплательщика" AS rn`).
		Scan(&v.Count, &v.ЕСХН, &v.УСН, &v.ЕНВД, &v.СРП)
	return
}

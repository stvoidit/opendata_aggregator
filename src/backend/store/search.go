package store

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"opendataaggregator/models/egr"
	"strconv"
	"strings"
)

// | Поле                         | Комментарий                              |
// |------------------------------|------------------------------------------|
// | income                       | Доходы                                   |
// | expenses                     | Расходы                                  |
// | income_tax                   | Налог на прибыль                         |
// | tax_usn                      | Налог по УСН                             |
// | intangible_assets            | Нематериальные активы                    |
// | basic_assets                 | Основные средства                        |
// | other_non_current_assets     | Прочие внеоборотные активы               |
// | non_current_assets           | Внеоборотные активы                      |
// | stocks                       | Запасы                                   |
// | net_aassets                  | Чистые активы                            |
// | accounts_receivable          | Дебиторская задолженность                |
// | cash_and_equivalents         | Денежные средства и денежные эквиваленты |
// | current_assets               | Оборотные активы                         |
// | total_assets                 | Активы всего                             |
// | capital_and_reserves         | Капитал и резервы                        |
// | borrowed_funds_long_term     | Заёмные средства (долгосрочные)          |
// | borrowed_funds_short_term    | Заёмные средства (краткосрочные)         |
// | accounts_payable             | Кредиторская задолженность               |
// | other_short_term_liabilities | Прочие краткосрочные обязательства       |
// | total_liabilities            | Пассивы всего                            |
// | revenue                      | Выручка                                  |
// | cost_of_sales                | Себестоимость продаж                     |
// | gross_profit                 | Валовая прибыль                          |
// | commercial_expenses          | Коммерческие расходы                     |
// | management_expenses          | Управленческие расходы                   |
// | profit_from_sale             | Прибыль (убыток) от продажи              |
// | interest_payable             | Проценты к уплате                        |
// | other_income                 | Прочие доходы                            |
// | other_expenses               | Прочие расходы                           |
// | profit_before_taxation       | Прибыль (убыток) до налогообложения      |
// | сurrent_income_tax           | Текущий налог на прибыль                 |
// | net_profit                   | Чистая прибыль (убыток)                  |

var (
	financeFields = [...]string{
		"income",
		"expenses",
		"income_tax",
		"tax_usn",
		"intangible_assets",
		"basic_assets",
		"other_non_current_assets",
		"non_current_assets",
		"stocks",
		"net_aassets",
		"accounts_receivable",
		"cash_and_equivalents",
		"current_assets",
		"total_assets",
		"capital_and_reserves",
		"borrowed_funds_long_term",
		"borrowed_funds_short_term",
		"accounts_payable",
		"other_short_term_liabilities",
		"total_liabilities",
		"revenue",
		"cost_of_sales",
		"gross_profit",
		"commercial_expenses",
		"management_expenses",
		"profit_from_sale",
		"interest_payable",
		"other_income",
		"other_expenses",
		"profit_before_taxation",
		"сurrent_income_tax",
		"net_profit",
	}
	taxRegimeFields = map[string]string{
		"usn":  "усн",
		"eshn": "есхн",
		"envd": "енвд",
		"srp":  "срп",
	}
	paramsSuffix = [...]string{"from", "to"}
)

func parseParamsURL(params url.Values, oncount bool) (cte, join, where string, args []any, err error) {
	const (
		// withBalance - CTE фильтра по бух.балансу
		withBalance = `cte_balance AS (
			SELECT
				DISTINCT inn
			FROM
				accounting_statements.balance_mat_view
			WHERE
				%s
		)`
		withTaxAuthority = `cte_ta AS (
			SELECT
				inn, ogrn
			FROM
				egr.tax_authority ta
			WHERE
				ta.code = $%d::text
		)`
		cteBalanceJoin      = `INNER JOIN cte_balance ON cte_balance.inn = es.inn`                  // join бух.баланса фильтра
		cteTaxAuthorityJoin = `INNER JOIN cte_ta ON cte_ta.inn = es.inn AND cte_ta.ogrn = es.ogrn ` // join Налогового органа для фильтра
		// taxRegimeJoin - join фильтра по налоговым режимам
		taxRegimeJoin = `INNER JOIN LATERAL (
			SELECT
				rn.иннюл AS inn
			FROM
				public.режимналогоплательщика rn
			WHERE
				%s
		) AS rn ON
			rn.inn = es.inn`
	)
	var (
		cteBalance      bool                             // флаг для определения были ли финансовые фильтры
		limit           = uint64(20)                     // лимит на страницу
		offset          = uint64(0)                      // оффсет пагинации
		page            = uint64(1)                      // номер страницы
		conditions      = make([]string, 0, len(params)) // строковые условия для CTE баланса
		ctes            = make([]string, 0, 2)           // список из CTE
		joins           = make([]string, 0, 4)           // список из JOIN
		whereConditions = make([]string, 0, len(params)) // строковые условия для WHERE egr_search
		trConditions    = make([]string, 0, len(params)) // строковые условия для налоговых режимов
	)

	if !oncount {
		// определение номера страницы и расчет offset
		if params.Has("page") {
			if n, err := strconv.ParseUint(params.Get("page"), 10, 64); err == nil {
				page = n
			} else {
				return cte, join, where, args, nil
			}
		}
		if params.Has("limit") {
			if n, err := strconv.ParseUint(params.Get("limit"), 10, 64); err == nil {
				limit = n
			} else {
				return cte, join, where, args, nil
			}
		}
		offset = (page - 1) * limit
		// список аргументов для подстановки в запрос, сразу добавляются обязательные параметры
		args = append(args, limit, offset)
	}

	// обход фильтров бух.баланса
	for _, field := range financeFields {
		for _, suffix := range paramsSuffix {
			fieldName := fmt.Sprintf("%s_%s", field, suffix)
			if !params.Has(fieldName) {
				continue
			}
			n, err := strconv.ParseUint(params.Get(fieldName), 10, 64)
			if err != nil {
				return cte, join, where, args, nil
			}
			var term string
			if strings.EqualFold(suffix, "from") {
				term = ">="
			} else {
				term = "<="
			}
			args = append(args, n)
			condition := fmt.Sprintf("%s %s $%d::bigint\n", field, term, len(args))
			conditions = append(conditions, condition)
			cteBalance = true
		}
	}
	// фильтр по годам бухгалтерского баланса
	if params.Has("year") {
		var years = make([]uint64, 0, len(params["year"]))
		for _, val := range params["year"] {
			year, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				return cte, join, where, args, nil
			}
			years = append(years, year)
		}
		args = append(args, years)
		condition := fmt.Sprintf("year = ANY($%d::int4[])\n", len(args))
		conditions = append(conditions, condition)
		cteBalance = true
	}
	// проверка условий для налоговых режимов
	for k, v := range taxRegimeFields {
		if params.Has(k) {
			trConditions = append(trConditions, fmt.Sprintf("rn.%s IS TRUE", v))
		}
	}
	if len(trConditions) > 0 {
		trJoin := fmt.Sprintf(taxRegimeJoin, strings.Join(trConditions, "\n\t\t\t\tOR\n\t\t\t\t"))
		joins = append(joins, trJoin)
	}
	// если есть финансовые фильтры, то добавляем CTE в список
	if cteBalance {
		ctes = append(ctes, fmt.Sprintf(withBalance, strings.Join(conditions, "\t\t\tAND\n\t\t\t\t")))
		joins = append(joins, cteBalanceJoin)
	}
	// фильтр по основному коду ОКВЭД
	if params.Has("okved") && len(params["okved"]) > 0 {
		args = append(args, params["okved"])
		whereConditions = append(whereConditions, fmt.Sprintf("(okved ->> 'code') = ANY($%d::TEXT[])", len(args)))
	}
	// если есть фильтр по Налоговому органу, то добавляем CTE в список
	if params.Has("ta") {
		args = append(args, params.Get("ta"))
		ctes = append(ctes, fmt.Sprintf(withTaxAuthority, len(args)))
		joins = append(joins, cteTaxAuthorityJoin)
	}
	if len(ctes) > 0 {
		cte = "WITH " + strings.Join(ctes, ",\n") // формируем строку с нужными CTE
	}
	if len(joins) > 0 {
		join = strings.Join(joins, "\n") // формируем строку с нужными JOIN
	}
	// фильтры для главной таблицы
	if params.Has("date_registration_from") {
		value := params.Get("date_registration_from")
		args = append(args, value)
		whereConditions = append(whereConditions, fmt.Sprintf("date_registration >= $%d::date", len(args)))
	}
	if params.Has("date_registration_to") {
		value := params.Get("date_registration_to")
		args = append(args, value)
		whereConditions = append(whereConditions, fmt.Sprintf("date_registration <= $%d::date", len(args)))
	}
	if params.Has("date_liquidation_from") {
		value := params.Get("date_liquidation_from")
		args = append(args, value)
		whereConditions = append(whereConditions, fmt.Sprintf("date_liquidation >= $%d::date", len(args)))
	}
	if params.Has("date_liquidation_to") {
		value := params.Get("date_liquidation_to")
		args = append(args, value)
		whereConditions = append(whereConditions, fmt.Sprintf("date_liquidation <= $%d::date", len(args)))
	}
	if params.Has("is_legal") {
		switch params.Get("is_legal") {
		case "1":
			whereConditions = append(whereConditions, `es.is_legal IS TRUE`)
		case "0":
			whereConditions = append(whereConditions, `es.is_legal IS FALSE`)
		}
	}
	if params.Has("q") {
		const conditionTextSearch = `(
		es.idx_name_full @@ $%d::tsquery
		OR
		es.inn = $%d::text
		OR
		es.ogrn = $%d::text
	)`
		value := strings.Join(strings.Split(strings.ReplaceAll(strings.ToLower(params.Get("q")), `"`, ""), " "), " & ") // разбивка по словам
		args = append(args, value, value, value)
		argsCount := len(args)
		// драйвер не может положить одно значения в несколько переменных, нужно дублировать
		whereConditions = append(whereConditions, fmt.Sprintf(conditionTextSearch, argsCount-2, argsCount-1, argsCount))
	}
	if params.Has("kpp") {
		value := params.Get("kpp")
		args = append(args, value)
		whereConditions = append(whereConditions, fmt.Sprintf("kpp = $%d::text", len(args)))
	}
	if params.Has("status") {
		var statusParamsCondition = make([]string, 0, len(params["status"]))
		//? a - действующие (active)
		//? d - недействующие (disabled)
		//? l - ликвидация (liquidation)
		//? b - банкротство (bankruptcy)
		//? b - реорганизация (restructuring)
		const (
			active = `(
				es.date_liquidation IS NULL
					AND
				NOT EXISTS (
					SELECT
						1
					FROM
						egr.statuses s
					WHERE
						s.status_name != 'Действующая' AND s.status_full_name = es.status
				)
			)`
			disabled = `(
				es.date_liquidation IS NOT NULL
					OR
				(
						es.date_liquidation IS NULL
							AND
						es.status IS NOT NULL
							AND
						NOT EXISTS (
							SELECT
								1
							FROM
								egr.statuses s
							WHERE
								s.status_name != 'Недействующее' AND s.status_full_name = es.status
						)
					)
			)`
			liquidation = `(
				es.status IS NOT NULL
					AND
				EXISTS (
					SELECT
						1
					FROM
						egr.statuses s
					WHERE
						s.status_name = 'Ликвидация' AND s.status_full_name = es.status
				)
			)`
			bankruptcy = `(
				es.status IS NOT NULL
					AND
				EXISTS (
					SELECT
						1
					FROM
						egr.statuses s
					WHERE
						s.status_name = 'Банкротство' AND s.status_full_name = es.status
				)
			)`
			restructuring = `(
				es.status IS NOT NULL
					AND
				EXISTS (
					SELECT
						1
					FROM
						egr.statuses s
					WHERE
						s.status_name = 'Реорганизация' AND s.status_full_name = es.status
				)
			)`
		)
		for _, status := range params["status"] {
			switch status {
			case "a":
				statusParamsCondition = append(statusParamsCondition, active)
			case "d":
				statusParamsCondition = append(statusParamsCondition, disabled)
			case "l":
				statusParamsCondition = append(statusParamsCondition, liquidation)
			case "b":
				statusParamsCondition = append(statusParamsCondition, bankruptcy)
			case "r":
				statusParamsCondition = append(statusParamsCondition, restructuring)
			}
		}
		if len(statusParamsCondition) > 0 {
			whereConditions = append(whereConditions, "("+strings.Join(statusParamsCondition, " OR ")+")")
		}
	}
	// сбор выражения where для главной таблицы
	if len(whereConditions) > 0 {
		where = "WHERE\n\t" + strings.Join(whereConditions, "\n\tAND\n\t")
	}
	return cte, join, where, args, nil
}

// Search - поиск
func (db *DB) Search(ctx context.Context, params url.Values) ([]SearchResult, error) {
	cte, join, where, args, err := parseParamsURL(params, false)
	if err != nil {
		return nil, err
	}
	const q = `
%s

SELECT
	esp.date_discharge::text
	, esp.ogrn
	, esp.inn
	, esp.is_legal
	, EXISTS (SELECT 1 FROM public.реестрдисквалифицированныхлиц AS rd WHERE rd.inn = esp.inn) AS is_disqualified
    , EXISTS (SELECT 1 FROM public.информация1928коап AS koap WHERE koap.инн = esp.inn AND koap.огрн = esp.ogrn) AS is_koap
	, esp.avg_number_employees
	, esp.name
	, esp.name_full
	, pgn.status
	, esp.address
	, esp.date_registration::text
	, esp.date_liquidation::text
	, esp.chief
	, esp.kpp
	, esp.tax_authority
	, esp.tax_regime
	, esp.okved
	, esp.okpo
FROM
	egr.egr_search esp
INNER JOIN (
	SELECT
		es.inn
		, es.ogrn
		, CASE
			WHEN es.date_liquidation IS NOT NULL OR (ss.status_name IS NULL AND es.status IS NOT NULL)
				THEN 'Недействующее'
			WHEN es.date_liquidation IS NULL AND ss.status_name IS NULL
				THEN 'Действующая'
			ELSE
				COALESCE(ss.status_name, 'Действующая')
		END AS status
	FROM egr.egr_search AS es
	LEFT JOIN (
		SELECT
			s.status_name
			, s.status_full_name
		FROM
			egr.statuses AS s
	) AS ss ON ss.status_full_name = es.status
	%s
	%s
	ORDER BY es.date_discharge DESC
	LIMIT $1
	OFFSET $2
	) AS pgn ON pgn.inn = esp.inn AND pgn.ogrn = esp.ogrn`
	var query = fmt.Sprintf(q, cte, join, where)
	var result = make([]SearchResult, 0, 20)
	rows, err := db.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var v SearchResult
		if err := rows.Scan(
			&v.DateDischarge,
			&v.OGRN,
			&v.INN,
			&v.IsLegal,
			&v.IsDisqualified,
			&v.IsKOAP,
			&v.AvgNumberEmployees,
			&v.Name,
			&v.NameFull,
			&v.Status,
			&v.Address,
			&v.DateRegistration,
			&v.DateLiquidation,
			&v.Chief,
			&v.KPP,
			&v.TaxAuthority,
			&v.TaxRegime,
			&v.OKVED,
			&v.OKPO,
		); err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}

// CountSearch - кол-во строк в результате поиска
func (db *DB) CountSearch(ctx context.Context, params url.Values) (count SearchCount, err error) {
	cte, join, where, args, err := parseParamsURL(params, true)
	if err != nil {
		return count, err
	}
	const q = `
	%s

	SELECT
		count(*)
	FROM
		egr.egr_search es
	%s
	%s`
	var query = fmt.Sprintf(q, cte, join, where)
	err = db.pool.QueryRow(ctx, query, args...).Scan(&count.Count)
	return
}

type SearchCount struct {
	Count uint64 `json:"count"`
}

// ShortFIO - данные о физ.лице
type ShortFIO struct {
	Имя      string  `json:"Имя" jsonschema_description:"Имя"`
	Фамилия  string  `json:"Фамилия" jsonschema_description:"Фамилия"`
	Отчество string  `json:"Отчество" jsonschema_description:"Отчество"`
	ИННФЛ    *string `json:"ИННФЛ,omitempty" jsonschema_description:"ИНН (отсутствует для ИП)"`
}

// ShortPosition - должность физ.лица (только для ЮД)
type ShortPosition struct {
	ВидДолжн     string
	НаимДолжн    string
	НаимВидДолжн string
}

// Cheif - руководитель (ИП или ЮЛ)
type Cheif struct {
	СвФЛ    *ShortFIO      `json:"СвФЛ,omitempty" jsonschema_description:"ФИО (только ЮЛ)"`
	ФИОРус  *ShortFIO      `json:"ФИОРус,omitempty" jsonschema_description:"ФИО (только ИП)"`
	СвДолжн *ShortPosition `json:"СвДолжн,omitempty" jsonschema_description:"Должность (только ЮЛ)"`
}

// TaxRegimeShort - краткое инфо о налоговых режимах
type TaxRegimeShort struct {
	СРП  bool `json:"срп"`
	УСН  bool `json:"усн"`
	ЕНВД bool `json:"енвд"`
	ЕСХН bool `json:"есхн"`
}

type ShortOKVED struct {
	Code  string `json:"code" jsonschema_description:"Код ОКВЭД"`
	Title string `json:"title" jsonschema_description:"Название ОКВЭД"`
}

// SearchResult - строка результата поиска
type SearchResult struct {
	DateDischarge      string          `json:"date_discharge" jsonschema_description:"дата выписки"`
	OGRN               string          `json:"ogrn" jsonschema_description:"ОГРН"`
	INN                string          `json:"inn" jsonschema_description:"ИНН"`
	IsLegal            bool            `json:"is_legal" jsonschema_description:"является юридическим лицом"`
	AvgNumberEmployees *uint32         `json:"avg_number_employees,omitempty" jsonschema_description:"среднесписочная численность сотрудников"`
	IsDisqualified     bool            `json:"is_disqualified" jsonschema_description:"Есть в реестре дисквалифицированных лиц"`
	IsKOAP             bool            `json:"is_koap" jsonschema_description:"Есть в реестре КОАП"`
	Name               *string         `json:"name,omitempty" jsonschema_description:"название короткое (или ФИО)"`
	NameFull           string          `json:"name_full" jsonschema_description:"название полное (или ФИО)"`
	Status             *string         `json:"status,omitempty" jsonschema_description:"статус"`
	DateRegistration   *string         `json:"date_registration,omitempty" jsonschema_description:"дата регистрации"`
	DateLiquidation    *string         `json:"date_liquidation,omitempty" jsonschema_description:"дата ликвидации"`
	KPP                *string         `json:"kpp,omitempty" jsonschema_description:"КПП"`
	Chief              *Cheif          `json:"chief" jsonschema_description:"руководитель"`
	TaxAuthority       *egr.СвНОТип    `json:"tax_authority,omitempty" jsonschema_description:"налоговый орган"`
	Address            *egr.СвАдресЮЛ  `json:"address,omitempty" jsonschema_description:"адрес"`
	TaxRegime          *TaxRegimeShort `json:"tax_regime,omitempty" jsonschema_description:"налоговые режимы"`
	OKVED              *ShortOKVED     `json:"okved,omitempty" jsonschema_description:"основной ОКВЭД"`
	OKPO               *string         `json:"okpo,omitempty" jsonschema_description:"код ОКПО"`
}

// CodeOKVED - код ОКВЭД из справочника
type CodeOKVED struct {
	Code    string  `json:"code" jsonschema_description:"Код"`              // код ОКВЭД
	Title   string  `json:"title" jsonschema_description:"Название"`        // название ОКВЭД
	Version *string `json:"vers,omitempty" jsonschema_description:"Версия"` // версия ОКВЭД (2014 или пусто)
}

// SelectHandbookOKVED - весь справочник ОКВЭД
func (db *DB) SelectHandbookOKVED(ctx context.Context) ([]CodeOKVED, error) {
	const q = `SELECT code, title, vers FROM egr.handbook_okved ORDER BY code ASC`
	rows, err := db.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result = make([]CodeOKVED, 0)
	for rows.Next() {
		var v CodeOKVED
		if err := rows.Scan(
			&v.Code,
			&v.Title,
			&v.Version,
		); err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, rows.Err()
}

// SearchOKVED - поиск ОКВЭД по названию
func (db *DB) SearchOKVED(ctx context.Context, find string) ([]CodeOKVED, error) {
	const q = `
	SELECT
		code
		, title
		, vers
	FROM
		egr.handbook_okved
	WHERE
		code::TEXT OPERATOR(egr.%>) $1::TEXT
		OR
		title::TEXT OPERATOR(egr.%>) $1::TEXT
	ORDER BY
		egr.similarity(
			title
			, $1::TEXT
		) DESC`
	rows, err := db.pool.Query(ctx, q, find)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result = make([]CodeOKVED, 0)
	for rows.Next() {
		var v CodeOKVED
		if err := rows.Scan(
			&v.Code,
			&v.Title,
			&v.Version,
		); err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, rows.Err()
}

// CodeTaxAuthority - код налогового органа
type CodeTaxAuthority struct {
	Code string `json:"code" jsonschema_description:"Код"`
	Name string `json:"name" jsonschema_description:"Наименование"`
}

// SelectHandbookTaxAuthority - весь справочник налоговых органов
func (db *DB) SelectHandbookTaxAuthority(ctx context.Context) ([]CodeTaxAuthority, error) {
	const q = `SELECT code, name FROM egr.handbook_tax_authority`
	rows, err := db.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result = make([]CodeTaxAuthority, 0)
	for rows.Next() {
		var v CodeTaxAuthority
		if err := rows.Scan(
			&v.Code,
			&v.Name,
		); err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, rows.Err()
}

// SearchHandbookTaxAuthority - поиск Налогового органа
func (db *DB) SearchHandbookTaxAuthority(ctx context.Context, find string) ([]CodeTaxAuthority, error) {
	const q = `SELECT code, name FROM egr.handbook_tax_authority WHERE concat(code, name) ILIKE concat('%', $1::text, '%')`
	rows, err := db.pool.Query(ctx, q, find)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result = make([]CodeTaxAuthority, 0)
	for rows.Next() {
		var v CodeTaxAuthority
		if err := rows.Scan(
			&v.Code,
			&v.Name,
		); err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, rows.Err()
}

func (db *DB) SearchCSV(ctx context.Context, w io.Writer) error {
	const q = `COPY (SELECT
		es.*
		FROM egr.egr_search es
		WHERE replace(es.name_full, '"', '') ILIKE '%магнит%'
		) TO STDOUT WITH CSV HEADER DELIMITER ';'`
	conn, err := db.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	if _, err := conn.Conn().PgConn().CopyTo(ctx, w, q); err != nil {
		return err
	}
	return nil
}

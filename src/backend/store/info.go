package store

import (
	"context"
	"errors"
	"fmt"
	"opendataaggregator/models"
	"opendataaggregator/models/egr"

	"github.com/jackc/pgx/v5"
)

var (
	ErrNotFound = errors.New("не найдено")
)

type FGIS struct {
	UnscheduledInspections []models.InspectionFGIS `json:"UnscheduledInspections" jsonschema_description:"Внеплановые проверки"`
	ScheduledInspections   []models.InspectionFGIS `json:"ScheduledInspections" jsonschema_description:"Плановые проверки"`
}

type SSHR struct {
	Date  string `json:"датасост" jsonschema_description:"Дата составления отчета"`
	Count int64  `json:"count" jsonschema_description:"Кол-во работников"`
}

// InfoResult - карточка ИП или ЮЛ
type InfoResult struct {
	// B503 *models.БухОтчетностьV503                                            `json:"typeБухОтчетностьV503,omitempty" jsonschema_description:"Тип бухотчета 5.03"`
	// B508 *models.БухОтчетностьV508                                            `json:"typeБухОтчетностьV508,omitempty" jsonschema_description:"Тип бухотчета 5.08"`
	F1  *egr.EGRUL                                                           `json:"ЕГРЮЛ,omitempty" jsonschema_description:"ЕГРЮЛ"`
	F2  *egr.EGRIP                                                           `json:"ЕГРИП,omitempty" jsonschema_description:"ЕГРИП"`
	F3  []БухОтчет                                                           `json:"Бухгалтерская отчетность,omitempty" jsonschema_description:"Бухгалтерская отчетность [ TypeBudget503[] или TypeBudget508[] ]"`
	F4  *models.НалоговыйРежимНалогоплательщика                              `json:"Режим налогоплательщика,omitempty" jsonschema_description:"Режим налогоплательщика"`
	F5  []models.СведенияОСуммахНедоимки                                     `json:"Сведения о суммах недоимки и задолженности по пеням и штрафам,omitempty" jsonschema_description:"Сведения о суммах недоимки и задолженности по пеням и штрафам"`
	F6  []models.Росаккредитация                                             `json:"Росаккредитация,omitempty" jsonschema_description:"Росаккредитация"`
	F7  *models.СведенияОбУплаченныхОрганизациейНалогов                      `json:"Сведения об уплаченных организацией в календарном году налогов и сборов,omitempty" jsonschema_description:"Сведения об уплаченных организацией в календарном году налогов и сборов"`
	F8  *models.РеестрСубъектовМалогоИСреднегоПредпринимательства            `json:"Субъект малого или среднего предпринимательства,omitempty" jsonschema_description:"Субъект малого или среднего предпринимательства"`
	F9  []models.РеестрДисквалифицированныхЛиц                               `json:"Реестр дисквалифицированных лиц,omitempty" jsonschema_description:"Реестр дисквалифицированных лиц"`
	F10 []models.ИсполнительныеПроизводстваВОтношенииЮридическихЛиц          `json:"Исполнительные производства в отношении юридических лиц,omitempty" jsonschema_description:"Исполнительные производства в отношении юридических лиц"`
	F11 []models.ОткрытыйРеестрТоварныхЗнаков                                `json:"Открытый реестр товарных знаков и знаков обслуживания Российской Федерации,omitempty" jsonschema_description:"Открытый реестр товарных знаков и знаков обслуживания Российской Федерации"`
	F12 []models.ОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков                 `json:"Открытый реестр общеизвестных в Российской Федерации товарных знаков,omitempty" jsonschema_description:"Открытый реестр общеизвестных в Российской Федерации товарных знаков"`
	F13 []models.НалоговыхПравонарушенияхИМерахОтветственности               `json:"Сведения о налоговых правонарушениях и мерах ответственности за их совершение,omitempty" jsonschema_description:"Сведения о налоговых правонарушениях и мерах ответственности за их совершение"`
	F14 *models.СведенияОбУчастииВКонсГруппе                                 `json:"Сведения об участии в консолидированной группе налогоплательщиков,omitempty" jsonschema_description:"Сведения об участии в консолидированной группе налогоплательщиков"`
	F15 []Hotel                                                              `json:"Данные о гостиницах,omitempty" jsonschema_description:"Данные о гостиницах"`
	F16 []models.ПривлеченииУчастникаЗакупкиКАдминистративнойОтветственности `json:"Информация о привлечении участника закупки к административной ответственности по ст. 19.28 КоАП,omitempty" jsonschema_description:"Информация о привлечении участника закупки к административной ответственности по ст. 19.28 КоАП"`
	F17 []models.RDS                                                         `json:"Сведения из Реестра деклараций о соответствии,omitempty" jsonschema_description:"Сведения из Реестра деклараций о соответстви"`
	F18 FGIS                                                                 `json:"ФГИС" jsonschema_description:"Реестр контрольных (надзорных) и профилактических мероприятий"`
	F19 *SSHR                                                                `json:"ССЧС,omitempty" jsonschema_description:"Сведения о среднесписочной численности работников организации"`
}

// SelectInfo - получить карточку ИП или ЮЛ
func (db *DB) SelectInfo(ctx context.Context, inn, ogrn string) (*InfoResult, error) {
	var m = new(InfoResult)
	var err error
	if m.F1, err = db.SelectInfoEGRUL(ctx, inn, ogrn); err != nil {
		return nil, err
	}
	if m.F2, err = db.SelectInfoEGRIP(ctx, inn, ogrn); err != nil {
		return nil, err
	}
	if m.F1 == nil && m.F2 == nil {
		return nil, ErrNotFound
	}
	if m.F3, err = db.SelectAccountingStatements(ctx, inn); err != nil {
		return nil, fmt.Errorf("F3 %w", err)
	}
	if m.F4, err = db.SelectTaxpayersRegime(ctx, inn); err != nil {
		return nil, fmt.Errorf("F4 %w", err)
	}
	if m.F5, err = db.SelectInformationAboutAmountsArrears(ctx, inn); err != nil {
		return nil, fmt.Errorf("F5 %w", err)
	}
	if m.F6, err = db.SelectRossacreditaciya(ctx, inn, ogrn); err != nil {
		return nil, fmt.Errorf("F6 %w", err)
	}
	if m.F7, err = db.SelectInformationAboutTaxesAndFeesPaid(ctx, inn); err != nil {
		return nil, fmt.Errorf("F7 %w", err)
	}
	if m.F8, err = db.SelectSMP(ctx, inn, ogrn); err != nil {
		return nil, fmt.Errorf("F8 %w", err)
	}
	if m.F9, err = db.SelectRegisterDisqualifiedPersons(ctx, inn); err != nil {
		return nil, fmt.Errorf("F9 %w", err)
	}
	if m.F10, err = db.SelectИсполнительныеПроизводстваВОтношенииЮридическихЛиц(ctx, inn); err != nil {
		return nil, fmt.Errorf("F10 %w", err)
	}
	if m.F11, err = db.SelectОткрытыйРеестрТоварныхЗнаков(ctx, inn); err != nil {
		return nil, fmt.Errorf("F11 %w", err)
	}
	if m.F12, err = db.SelectОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков(ctx, inn); err != nil {
		return nil, fmt.Errorf("F12 %w", err)
	}
	if m.F13, err = db.SelectНалоговыхПравонарушенияхИМерахОтветственности(ctx, inn); err != nil {
		return nil, fmt.Errorf("F13 %w", err)
	}
	if m.F14, err = db.SelectСведенияОбУчастииВКонсГруппе(ctx, inn); err != nil {
		return nil, fmt.Errorf("F14 %w", err)
	}
	if m.F15, err = db.SelectГостиница(ctx, inn, ogrn); err != nil {
		return nil, fmt.Errorf("F15 %w", err)
	}
	if m.F16, err = db.SelectПривлеченииУчастникаЗакупкиКАдминистративнойОтветственности(ctx, inn); err != nil {
		return nil, fmt.Errorf("F16 %w", err)
	}
	if m.F17, err = db.SelectRDS(ctx, inn, ogrn); err != nil {
		return nil, fmt.Errorf("F17 %w", err)
	}
	if m.F18.UnscheduledInspections, err = db.SelectUnscheduledInspections(ctx, inn, ogrn); err != nil {
		return nil, fmt.Errorf("F18.1 %w", err)
	}
	if m.F18.ScheduledInspections, err = db.SelectScheduledInspections(ctx, inn, ogrn); err != nil {
		return nil, fmt.Errorf("F18.2 %w", err)
	}
	if m.F19, err = db.SelectSSHR(ctx, inn); err != nil {
		return nil, fmt.Errorf("F19 %w", err)
	}
	return m, nil
}

// SelectInfoEGRUL - получение данных о ЕГРЮЛ
func (db *DB) SelectInfoEGRUL(ctx context.Context, inn, ogrn string) (data *egr.EGRUL, err error) {
	const q = `
	SELECT
		egrul_data
	FROM
		egr.егрюл
	WHERE
		инн = $1::text
		AND огрн = $2::text
	LIMIT 1`
	data = new(egr.EGRUL)
	err = db.pool.QueryRow(ctx, q, inn, ogrn).Scan(data)
	if errors.Is(err, pgx.ErrNoRows) {
		err = nil
		data = nil
	}
	return
}

// SelectInfoEGRIP - получение данных о ЕГРИП
func (db *DB) SelectInfoEGRIP(ctx context.Context, inn, ogrn string) (data *egr.EGRIP, err error) {
	const q = `
	SELECT
		egrip_data
	FROM
		egr.егрип
	WHERE
		инн = $1::text
		AND огрн = $2::text
	LIMIT 1`
	data = new(egr.EGRIP)
	err = db.pool.QueryRow(ctx, q, inn, ogrn).Scan(data)
	if errors.Is(err, pgx.ErrNoRows) {
		err = nil
		data = nil
	}
	return
}

// GetServiceLogSources - service_management.source_files
func (db *DB) GetServiceLogSources(ctx context.Context) ([]ServiceLogSource, error) {
	const q = `SELECT
    source_type
    , source_link
    , filename
    , sha256sum
    , downloaded
    , uploaded
    , task_datetime
    , id
FROM
    service_management.source_files sf
    ORDER BY sf.task_datetime DESC`
	rows, err := db.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var data = make([]ServiceLogSource, 0, 10000)
	for rows.Next() {
		var sls ServiceLogSource
		if err := rows.Scan(
			&sls.SourceType,
			&sls.SourceLink,
			&sls.Filename,
			&sls.SHA256Sum,
			&sls.Downloaded,
			&sls.Uploaded,
			&sls.TaskDatetime,
			&sls.ID,
		); err != nil {
			return nil, err
		}
		data = append(data, sls)
	}
	return data, nil
}

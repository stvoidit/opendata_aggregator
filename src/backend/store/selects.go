package store

import (
	"context"
	"errors"
	"io"
	"opendataaggregator/models"
	"opendataaggregator/models/egr"

	"github.com/jackc/pgx/v5"
)

// AllTypes - список всех используемых типов
type AllTypes struct {
	TypeDatabaseStat     DatabaseStat              `json:"TypeDatabaseStat" jsonschema_description:"Статистика источников ( /api/db_stats )"`
	TypeSearchResult     []SearchResult            `json:"TypeSearchResult" jsonschema_description:"Результат поиска ( /api/search )"`
	TypeInfoResult       InfoResult                `json:"TypeInfoResult" jsonschema_description:"Карточка ЮЛ/ИП ( /api/info/{inn:[0-9]+}/{ogrn:[0-9]+} )"`
	TypeResultCard       ResultCard                `json:"TypeResultCard" jsonschema_description:"Карточка компании ( /api/info/{inn:[0-9]+} ) [deprecated]"`
	TypeHandbookOKVED    []CodeOKVED               `json:"TypeHandbookOKVED" jsonschema_description:"Справочник ОКВЭД ( /api/handbook_okved )"`
	TypeCodeTaxAuthority []CodeTaxAuthority        `json:"TypeCodeTaxAuthority" jsonschema_description:"Справочник Налоговых органов ( /api/handbook_tax_authority )"`
	TypeCategoriesIP     []CategoriesIP            `json:"TypeCategoriesIP" jsonschema_description:"Справочник категорий ИП ( /api/handbook_categories_ip )"`
	TypeLegalStatus      LegalStatus               `json:"TypeLegalStatus" jsonschema_description:"Справочник расшифровки статусов ЮЛ ( /api/legal_statuses ) [ Map<string, string> ]"`
	TypeBudget503        *models.БухОтчетностьV503 `json:"TypeBudget503,omitempty" jsonschema_description:"Тип бухотчета 5.03"`
	TypeBudget508        *models.БухОтчетностьV508 `json:"TypeBudget508,omitempty" jsonschema_description:"Тип бухотчета 5.08"`
}

// ResultCard - карточка ИП или ЮЛ
type ResultCard struct {
	// B503 *models.БухОтчетностьV503                                            `json:"typeБухОтчетностьV503,omitempty" jsonschema_description:"Тип бухотчета 5.03"`
	// B508 *models.БухОтчетностьV508                                            `json:"typeБухОтчетностьV508,omitempty" jsonschema_description:"Тип бухотчета 5.08"`
	F1  []egr.EGRUL                                                          `json:"ЕГРЮЛ,omitempty" jsonschema_description:"ЕГРЮЛ"`
	F2  []egr.EGRIP                                                          `json:"ЕГРИП,omitempty" jsonschema_description:"ЕГРИП"`
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
}

// SelectTotalInfo - получить всю инфу
func (db *DB) SelectTotalInfo(ctx context.Context, inn string) (*ResultCard, error) {
	var m = new(ResultCard)
	var err error
	if m.F1, err = db.SelectEGRUL(ctx, inn); err != nil {
		return nil, err
	}
	if m.F2, err = db.SelectEGRIP(ctx, inn); err != nil {
		return nil, err
	}
	if len(m.F1) == 0 && len(m.F2) == 0 {
		return nil, nil
	}
	if m.F3, err = db.SelectAccountingStatements(ctx, inn); err != nil {
		return nil, err
	}
	if m.F4, err = db.SelectTaxpayersRegime(ctx, inn); err != nil {
		return nil, err
	}
	if m.F5, err = db.SelectInformationAboutAmountsArrears(ctx, inn); err != nil {
		return nil, err
	}
	if m.F6, err = db.SelectRossacreditaciya(ctx, inn, ""); err != nil {
		return nil, err
	}
	if m.F7, err = db.SelectInformationAboutTaxesAndFeesPaid(ctx, inn); err != nil {
		return nil, err
	}
	if m.F8, err = db.SelectSMP(ctx, inn, ""); err != nil {
		return nil, err
	}
	if m.F9, err = db.SelectRegisterDisqualifiedPersons(ctx, inn); err != nil {
		return nil, err
	}
	if m.F10, err = db.SelectИсполнительныеПроизводстваВОтношенииЮридическихЛиц(ctx, inn); err != nil {
		return nil, err
	}
	if m.F11, err = db.SelectОткрытыйРеестрТоварныхЗнаков(ctx, inn); err != nil {
		return nil, err
	}
	if m.F12, err = db.SelectОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков(ctx, inn); err != nil {
		return nil, err
	}
	if m.F13, err = db.SelectНалоговыхПравонарушенияхИМерахОтветственности(ctx, inn); err != nil {
		return nil, err
	}
	if m.F14, err = db.SelectСведенияОбУчастииВКонсГруппе(ctx, inn); err != nil {
		return nil, err
	}
	if m.F15, err = db.SelectГостиница(ctx, inn, ""); err != nil {
		return nil, err
	}
	if m.F16, err = db.SelectПривлеченииУчастникаЗакупкиКАдминистративнойОтветственности(ctx, inn); err != nil {
		return nil, err
	}
	return m, nil
}

// SelectEGRUL - получение данных о ЕГРЮЛ
func (db *DB) SelectEGRUL(ctx context.Context, inn string) ([]egr.EGRUL, error) {
	const q = `SELECT egrul_data FROM egr.егрюл WHERE инн = $1::text ORDER BY огрн`
	rows, err := db.pool.Query(ctx, q, inn)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var data = make([]egr.EGRUL, 0)
	for rows.Next() {
		var v egr.EGRUL
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		data = append(data, v)
	}
	return data, nil
}

// SelectAccountingStatements - получение данных о Бухгалтерской отчетности
func (db *DB) SelectAccountingStatements(ctx context.Context, inn string) ([]БухОтчет, error) {
	const q = `SELECT balance_data FROM accounting_statements.баланс WHERE иннюл = $1::text ORDER BY год ASC`
	rows, err := db.pool.Query(ctx, q, inn)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var data = make([]БухОтчет, 0)
	for rows.Next() {
		var v = make(БухОтчет)
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		data = append(data, v)
	}
	return data, nil
}

// SelectTaxpayersRegime - получение данных о режиме налогоплательщика
func (db *DB) SelectTaxpayersRegime(ctx context.Context, inn string) (*models.НалоговыйРежимНалогоплательщика, error) {
	const q = `
	SELECT
		датадок::text
		, наиморг
		, иннюл
		, есхн
		, усн
		, енвд
		, срп
	FROM
		public.режимналогоплательщика
	WHERE иннюл = $1::text`
	var v = new(models.НалоговыйРежимНалогоплательщика)
	if err := db.pool.QueryRow(ctx, q, inn).Scan(
		&v.ДатаДок,
		&v.СведНП.НаимОрг,
		&v.СведНП.ИННЮЛ,
		&v.СведСНР.ПризнЕСХН,
		&v.СведСНР.ПризнУСН,
		&v.СведСНР.ПризнЕНВД,
		&v.СведСНР.ПризнСРП); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

// SelectEGRIP - получение данных о ЕГРИП
func (db *DB) SelectEGRIP(ctx context.Context, inn string) ([]egr.EGRIP, error) {
	const q = `SELECT egrip_data FROM egr.егрип WHERE инн = $1::text ORDER BY огрн`
	rows, err := db.pool.Query(ctx, q, inn)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var data = make([]egr.EGRIP, 0)
	for rows.Next() {
		var v egr.EGRIP
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		data = append(data, v)
	}
	return data, nil
}

// SelectInformationAboutAmountsArrears - Сведения о суммах недоимки и задолженности по пеням и штрафам
func (db *DB) SelectInformationAboutAmountsArrears(ctx context.Context, inn string) (values []models.СведенияОСуммахНедоимки, err error) {
	const q = `
	SELECT
		иддок
		, датадок::TEXT
		, датасост::TEXT
		, наиморг
		, иннюл
		, json_agg(json_build_object(
		'НаимНалог', наимналог,
		'СумНедНалог', сумнедналог,
		'СумПени', сумпени,
		'СумШтраф', сумштраф,
		'ОбщСумНедоим', общсумнедоим) ORDER BY наимналог)
	FROM
		public.сведенияосуммахнедоимки
	WHERE
		иннюл = $1::text
	GROUP BY 1,2,3,4,5
	ORDER BY датасост ASC, датадок ASC`
	rows, err := db.pool.Query(ctx, q, inn)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	values = make([]models.СведенияОСуммахНедоимки, 0, 8)
	for rows.Next() {
		var v models.СведенияОСуммахНедоимки
		if err := rows.Scan(
			&v.ИдДок,
			&v.ДатаДок,
			&v.ДатаСост,
			&v.СведНП.НаимОрг,
			&v.СведНП.ИННЮЛ,
			&v.СведНедоим,
		); err != nil {
			return nil, err
		}
		values = append(values, v)
	}
	return values, nil
}

// SelectRossacreditaciya - Росаккредитация
func (db *DB) SelectRossacreditaciya(ctx context.Context, inn, ogrn string) ([]models.Росаккредитация, error) {
	const q = `
	SELECT
		id_cert
		, cert_status
		, cert_type
		, reg_number
		, date_begining::text
		, coalesce(date_finish::TEXT, '') AS date_finish
		, product_scheme
		, product_object_type_cert
		, product_type
		, product_okpd2
		, product_tn_ved
		, product_tech_reg
		, product_group
		, product_name
		, product_info
		, applicant_type
		, person_applicant_type
		, applicant_ogrn
		, applicant_inn
		, applicant_phone
		, applicant_fax
		, applicant_email
		, applicant_website
		, applicant_name
		, applicant_director_name
		, applicant_address
		, applicant_address_actual
		, manufacturer_type
		, manufacturer_ogrn
		, manufacturer_inn
		, manufacturer_phone
		, manufacturer_fax
		, manufacturer_email
		, manufacturer_website
		, manufacturer_name
		, manufacturer_director_name
		, manufacturer_country
		, manufacturer_address
		, manufacturer_address_actual
		, manufacturer_address_filial
		, organ_to_certification_name
		, organ_to_certification_reg_number
		, organ_to_certification_head_name
		, basis_for_certificate
		, old_basis_for_certificate
		, fio_expert
		, fio_signatory
		, product_national_standart
		, production_analysis_for_act
		, production_analysis_for_act_number
		, COALESCE(production_analysis_for_act_date::text, '')
	FROM
		public.росаккредитация
	WHERE
		applicant_inn = $1::text
		AND
		applicant_ogrn = $2::text
	ORDER BY id_cert`
	rows, err := db.pool.Query(ctx, q, inn, ogrn)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var data = make([]models.Росаккредитация, 0)
	for rows.Next() {
		var rk models.Росаккредитация
		if err := rows.Scan(
			&rk.IDcert,
			&rk.CertStatus,
			&rk.CertType,
			&rk.RegNumber,
			&rk.DateBegining,
			&rk.DateFinish,
			&rk.ProductScheme,
			&rk.ProductObjectTypeCert,
			&rk.ProductType,
			&rk.ProductOKPD2,
			&rk.ProductTnVed,
			&rk.ProductTechReg,
			&rk.ProductGroup,
			&rk.ProductName,
			&rk.ProductInfo,
			&rk.ApplicantType,
			&rk.PersonApplicantType,
			&rk.ApplicantOGRN,
			&rk.ApplicantINN,
			&rk.ApplicantPhone,
			&rk.ApplicantFax,
			&rk.ApplicantEmail,
			&rk.ApplicantWebsite,
			&rk.ApplicantName,
			&rk.ApplicantDirectorName,
			&rk.ApplicantAddress,
			&rk.ApplicantAddressActual,
			&rk.ManufacturerType,
			&rk.ManufacturerOGRN,
			&rk.ManufacturerINN,
			&rk.ManufacturerPhone,
			&rk.ManufacturerFax,
			&rk.ManufacturerEmail,
			&rk.ManufacturerWebsite,
			&rk.ManufacturerName,
			&rk.ManufacturerDirectorName,
			&rk.ManufacturerCountry,
			&rk.ManufacturerAddress,
			&rk.ManufacturerAddressActual,
			&rk.ManufacturerAddressFilial,
			&rk.OrganToCertificationName,
			&rk.OrganToCertificationRegNumber,
			&rk.OrganToCertificationHeadName,
			&rk.BasisForCertificate,
			&rk.OldBasisForCertificate,
			&rk.DioExpert,
			&rk.DioSignatory,
			&rk.ProductNationalStandart,
			&rk.ProductionAnalysisForAct,
			&rk.ProductionAnalysisForActNumber,
			&rk.ProductionAnalysisForActDate); err != nil {
			return nil, err
		}
		data = append(data, rk)
	}
	return data, nil
}

// SelectInformationAboutTaxesAndFeesPaid - Сведения об уплаченных организацией в календарном году налогов и сборов
func (db *DB) SelectInformationAboutTaxesAndFeesPaid(ctx context.Context, inn string) (*models.СведенияОбУплаченныхОрганизациейНалогов, error) {
	const q = `
	SELECT
		датасост::text
		, наиморг
		, иннюл
		, jsonb_agg(jsonb_build_object('НаимНалог', наимналог, 'СумУплНал', сумуплнал) ORDER BY наимналог) AS "СвУплСумНал"
	FROM
		public.сведенияобуплаченныхорганизацие
	WHERE иннюл = $1::text
	GROUP BY 1,2,3`
	var v = new(models.СведенияОбУплаченныхОрганизациейНалогов)
	if err := db.pool.QueryRow(ctx, q, inn).Scan(
		&v.ДатаДок,
		&v.СведНП.НаимОрг,
		&v.СведНП.ИННЮЛ,
		&v.СвУплСумНал); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

// SelectSMP - Единый реестр субъектов малого и среднего предпринимательства
func (db *DB) SelectSMP(ctx context.Context, inn, ogrn string) (*models.РеестрСубъектовМалогоИСреднегоПредпринимательства, error) {
	const q = `
	SELECT
		датасост::text
		, датавклмсп
		, видсубмсп
		, катсубмсп
		, призновмсп
		, сведсоцпред
		, ссчр
		, наиморг
		, наиморгсокр
		, иннюл
		, огрнюл
		, иннфл
		, огрнип
		, фамилия
		, имя
		, отчество
		, номлиценз
		, даталиценз
		, датаначлиценз
		, датаконлиценз
		, датаостлиценз
		, серлиценз
		, видлиценз
		, оргвыдлиценз
		, оргостлиценз
		, наимлицвд
		, кодрегион
		, регионтип
		, регионнаим
		, районтип
		, районнаим
		, городтип
		, городнаим
		, населпункттип
		, населпунктнаим
		, свпрод
		, свпрогпарт
		, свконтр
		, свдог
	FROM
		public.субъектымалогоисреднегопредприн
	WHERE инн = $1::text AND огрн = $2::text
	LIMIT 1`
	var v = new(models.РеестрСубъектовМалогоИСреднегоПредпринимательства)
	if err := db.pool.QueryRow(ctx, q, inn, ogrn).Scan(
		&v.ДатаСост,
		&v.ДатаВклМСП,
		&v.ВидСубМСП,
		&v.КатСубМСП,
		&v.ПризНовМСП,
		&v.СведСоцПред,
		&v.ССЧР,
		&v.ОргВклМСП.НаимОрг,
		&v.ОргВклМСП.НаимОргСокр,
		&v.ОргВклМСП.ИННЮЛ,
		&v.ОргВклМСП.ОГРН,
		&v.ИПВклМСП.ИННФЛ,
		&v.ИПВклМСП.ОГРНИП,
		&v.ИПВклМСП.ФИОИП.Фамилия,
		&v.ИПВклМСП.ФИОИП.Имя,
		&v.ИПВклМСП.ФИОИП.Отчество,
		&v.СвЛиценз.НомЛиценз,
		&v.СвЛиценз.ДатаЛиценз,
		&v.СвЛиценз.ДатаНачЛиценз,
		&v.СвЛиценз.ДатаКонЛиценз,
		&v.СвЛиценз.ДатаОстЛиценз,
		&v.СвЛиценз.СерЛиценз,
		&v.СвЛиценз.ВидЛиценз,
		&v.СвЛиценз.ОргВыдЛиценз,
		&v.СвЛиценз.ОргОстЛиценз,
		&v.СвЛиценз.НаимЛицВД,
		&v.СведМН.КодРегион,
		&v.СведМН.Регион.Тип,
		&v.СведМН.Регион.Наим,
		&v.СведМН.Район.Тип,
		&v.СведМН.Район.Наим,
		&v.СведМН.Город.Тип,
		&v.СведМН.Город.Наим,
		&v.СведМН.НаселПункт.Тип,
		&v.СведМН.НаселПункт.Наим,
		&v.СвПрод,
		&v.СвПрогПарт,
		&v.СвКонтр,
		&v.СвДог,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

// SelectRegisterDisqualifiedPersons - Реестр дисквалифицированных лиц
func (db *DB) SelectRegisterDisqualifiedPersons(ctx context.Context, inn string) ([]models.РеестрДисквалифицированныхЛиц, error) {
	const q = `
	SELECT
		id
		, fio
		, bdate::text
		, bplace
		, orgname
		, inn
		, positionfl
		, nkoap
		, gorgname
		, sudfio
		, sudposition
		, disqualificationduration
		, disstartdate::text
		, disenddate::text
	FROM
		public.реестрдисквалифицированныхлиц
	WHERE inn = $1::TEXT
	ORDER BY id`

	rows, err := db.pool.Query(ctx, q, inn)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var data = make([]models.РеестрДисквалифицированныхЛиц, 0)
	for rows.Next() {
		var v models.РеестрДисквалифицированныхЛиц
		if err := rows.Scan(
			&v.ID,
			&v.FIO,
			&v.BDate,
			&v.BPlace,
			&v.OrgName,
			&v.INN,
			&v.PositionFL,
			&v.NKOAP,
			&v.GOrgName,
			&v.SudFIO,
			&v.SudPosition,
			&v.DisqualificationDuration,
			&v.DisStartDate,
			&v.DisEndDate); err != nil {
			return nil, err
		}
		data = append(data, v)
	}
	return data, nil
}

// SelectИсполнительныеПроизводстваВОтношенииЮридическихЛиц - получение данных Набор открытых данных, содержащий общедоступные сведения,
// необходимые для осуществления задач по принудительному исполнению судебных актов,
// актов других органов и должностных лиц (в отношении юридических лиц).
func (db *DB) SelectИсполнительныеПроизводстваВОтношенииЮридическихЛиц(ctx context.Context, inn string) ([]models.ИсполнительныеПроизводстваВОтношенииЮридическихЛиц, error) {
	const q = `
	SELECT
		nameofdebtor
		, addressofdebtororganization
		, actualaddressofdebtororganization
		, numberofenforcementproceeding
		, dateofinstitutionproceeding
		, totalnumberofenforcementproceedings
		, executivedocumenttype
		, dateofexecutivedocument
		, numberofexecutivedocument
		, objectofexecutivedocuments
		, objectofexecution
		, amountdue
		, debtremainingbalance
		, departmentsofbailiffs
		, addressofdepartmentsofbailiff
		, debtortaxpayeridentificationnumber
		, taxpayeridentificationnumberoforganizationcollector
		, repaid
	FROM
		public.исппроизввотнюрлиц
	WHERE debtortaxpayeridentificationnumber = $1::TEXT
	ORDER BY dateofexecutivedocument`
	rows, err := db.pool.Query(ctx, q, inn)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var data = make([]models.ИсполнительныеПроизводстваВОтношенииЮридическихЛиц, 0)
	for rows.Next() {
		var v models.ИсполнительныеПроизводстваВОтношенииЮридическихЛиц
		if err := rows.Scan(
			&v.NameOfDebtor,
			&v.AddressOfDebtorOrganization,
			&v.ActualAddressOfDebtorOrganization,
			&v.NumberOfEnforcementProceeding,
			&v.DateOfInstitutionProceeding,
			&v.TotalNumberOfEnforcementProceedings,
			&v.ExecutiveDocumentType,
			&v.DateOfExecutiveDocument,
			&v.NumberOfExecutiveDocument,
			&v.ObjectOfExecutiveDocuments,
			&v.ObjectOfExecution,
			&v.AmountDue,
			&v.DebtRemainingBalance,
			&v.DepartmentsOfBailiffs,
			&v.AddressOfDepartmentsOfBailiff,
			&v.DebtorTaxpayerIdentificationNumber,
			&v.TaxpayerIdentificationNumberOfOrganizationCollector,
			&v.Repaid); err != nil {
			return nil, err
		}
		data = append(data, v)
	}
	return data, nil
}

// SelectОткрытыйРеестрТоварныхЗнаков - Открытый реестр товарных знаков и знаков обслуживания Российской Федерации
func (db *DB) SelectОткрытыйРеестрТоварныхЗнаков(ctx context.Context, inn string) ([]models.ОткрытыйРеестрТоварныхЗнаков, error) {
	const q = `
	SELECT
		registrationnumber
		, registrationdate
		, applicationnumber
		, applicationdate
		, prioritydate
		, exhibitionprioritydate
		, parisconventionprioritynumber
		, parisconventionprioritydate
		, parisconventionprioritycountrycode
		, initialapplicationnumber
		, initialapplicationorioritydate
		, initialregistrationnumber
		, initialregistrationdate
		, internationalregistrationnumber
		, internationalregistrationdate
		, internationalregistrationprioritydate
		, internationalregistrationentrydate
		, applicationnumberforrecognitionoftrademarkfromcrimea
		, applicationdateforrecognitionoftrademarkfromcrimea
		, crimeantrademarkapplicationnumberforstateregistrationinukraine
		, crimeantrademarkapplicationdateforstateregistrationinukraine
		, crimeantrademarkcertificatenumberinukraine
		, exclusiverightstransferagreementregistrationnumber
		, exclusiverightstransferagreementregistrationdate
		, legallyrelatedapplications
		, legallyrelatedregistrations
		, expirationdate
		, rightholdername
		, foreignrightholdername
		, rightholderaddress
		, rightholdercountrycode
		, rightholderogrn
		, rightholderinn
		, correspondenceaddress
		, collective
		, collectiveusers
		, extractionfromcharterofthecollectivetrademark
		, colorspecification
		, unprotectedelements
		, kindspecification
		, threedimensional
		, threedimensionalspecification
		, holographic
		, holographicspecification
		, sound
		, soundspecification
		, olfactory
		, olfactoryspecification
		, color
		, colortrademarkspecification
		, light
		, lightspecification
		, changing
		, changingspecification
		, positional
		, positionalspecification
		, actual
		, publicationurl
	FROM
		public.открытыйреестртоварныхзнаков
	WHERE rightholderinn = $1::TEXT
	ORDER BY registrationnumber`
	rows, err := db.pool.Query(ctx, q, inn)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var data = make([]models.ОткрытыйРеестрТоварныхЗнаков, 0)
	for rows.Next() {
		var v models.ОткрытыйРеестрТоварныхЗнаков
		if err := rows.Scan(
			&v.RegistrationNumber,
			&v.RegistrationDate,
			&v.ApplicationNumber,
			&v.ApplicationDate,
			&v.PriorityDate,
			&v.ExhibitionPriorityDate,
			&v.ParisConventionPriorityNumber,
			&v.ParisConventionPriorityDate,
			&v.ParisConventionPriorityCountryCode,
			&v.InitialApplicationNumber,
			&v.InitialApplicationOriorityDate,
			&v.InitialRegistrationNumber,
			&v.InitialRegistrationDate,
			&v.InternationalRegistrationNumber,
			&v.InternationalRegistrationDate,
			&v.InternationalRegistrationPriorityDate,
			&v.InternationalRegistrationEntryDate,
			&v.ApplicationNumberForRecognitionOfTrademarkFromCrimea,
			&v.ApplicationDateForRecognitionOfTrademarkFromCrimea,
			&v.CrimeanTrademarkApplicationNumberForStateRegistrationInUkraine,
			&v.CrimeanTrademarkApplicationDateForStateRegistrationInUkraine,
			&v.CrimeanTrademarkCertificateNumberInUkraine,
			&v.ExclusiveRightsTransferAgreementRegistrationNumber,
			&v.ExclusiveRightsTransferAgreementRegistrationDate,
			&v.LegallyRelatedApplications,
			&v.LegallyRelatedRegistrations,
			&v.ExpirationDate,
			&v.RightHolderName,
			&v.ForeignRightHolderName,
			&v.RightHolderAddress,
			&v.RightHolderCountryCode,
			&v.RightHolderOgrn,
			&v.RightHolderInn,
			&v.CorrespondenceAddress,
			&v.Collective,
			&v.CollectiveUsers,
			&v.ExtractionFromCharterOfTheCollectiveTrademark,
			&v.ColorSpecification,
			&v.UnprotectedElements,
			&v.KindSpecification,
			&v.Threedimensional,
			&v.ThreedimensionalSpecification,
			&v.Holographic,
			&v.HolographicSpecification,
			&v.Sound,
			&v.SoundSpecification,
			&v.Olfactory,
			&v.OlfactorySpecification,
			&v.Color,
			&v.ColorTrademarkSpecification,
			&v.Light,
			&v.LightSpecification,
			&v.Changing,
			&v.ChangingSpecification,
			&v.Positional,
			&v.PositionalSpecification,
			&v.Actual,
			&v.PublicationURL); err != nil {
			return nil, err
		}
		data = append(data, v)
	}
	return data, nil
}

// SelectОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков - Открытый реестр общеизвестных в Российской Федерации товарных знаков
func (db *DB) SelectОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков(ctx context.Context, inn string) ([]models.ОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков, error) {
	const q = `
	SELECT
		registrationnumber
		, registrationdate::text
		, wellknowntrademarkdate::text
		, legallyrelatedregistrations
		, rightholdername
		, foreignrightholdername
		, rightholderaddress
		, rightholdercountrycode
		, rightholderogrn
		, rightholderinn
		, correspondenceaddress
		, collective
		, collectiveusers
		, extractionfromcharterofcollectivetrademark
		, colorspecification
		, unprotectedelements
		, kindspecification
		, threedimensional
		, threedimensionalspecification
		, holographic
		, holographicspecification
		, sound
		, soundspecification
		, olfactory
		, olfactoryspecification
		, color
		, colortrademarkspecification
		, light
		, lightspecification
		, changing
		, changingspecification
		, positional
		, positionalspecification
		, actual
		, publicationurl
	FROM
		public.реестробщеизвестныхтоварныхзнак
	WHERE rightholderinn = $1::TEXT
	ORDER BY registrationnumber`
	rows, err := db.pool.Query(ctx, q, inn)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var data = make([]models.ОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков, 0)
	for rows.Next() {
		var v models.ОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков
		if err := rows.Scan(
			&v.RegistrationNumber,
			&v.RegistrationDate,
			&v.WellKnownTrademarkDate,
			&v.LegallyRelatedRegistrations,
			&v.RightHolderName,
			&v.ForeignRightHolderName,
			&v.RightHolderAddress,
			&v.RightHolderCountryCode,
			&v.RightHolderOgrn,
			&v.RightHolderInn,
			&v.CorrespondenceAddress,
			&v.Collective,
			&v.CollectiveUsers,
			&v.ExtractionFromCharterOfCollectiveTrademark,
			&v.ColorSpecification,
			&v.UnprotectedElements,
			&v.KindSpecification,
			&v.Threedimensional,
			&v.ThreedimensionalSpecification,
			&v.Holographic,
			&v.HolographicSpecification,
			&v.Sound,
			&v.SoundSpecification,
			&v.Olfactory,
			&v.OlfactorySpecification,
			&v.Color,
			&v.ColorTrademarkSpecification,
			&v.Light,
			&v.LightSpecification,
			&v.Changing,
			&v.ChangingSpecification,
			&v.Positional,
			&v.PositionalSpecification,
			&v.Actual,
			&v.PublicationURL); err != nil {
			return nil, err
		}
		data = append(data, v)
	}
	return data, nil
}

// SelectНалоговыхПравонарушенияхИМерахОтветственности - Сведения о налоговых правонарушениях и мерах ответственности за их совершение
func (db *DB) SelectНалоговыхПравонарушенияхИМерахОтветственности(ctx context.Context, inn string) ([]models.НалоговыхПравонарушенияхИМерахОтветственности, error) {
	const q = `
	SELECT
		датасост::text
		, иннюл
		, наиморг
		, сумштраф
	FROM
		public.налоговыеправонарушенияиштрафы
	WHERE иннюл = $1::TEXT
	ORDER BY датасост`
	rows, err := db.pool.Query(ctx, q, inn)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var data = make([]models.НалоговыхПравонарушенияхИМерахОтветственности, 0)
	for rows.Next() {
		var v models.НалоговыхПравонарушенияхИМерахОтветственности
		if err := rows.Scan(
			&v.ДатаДок,
			&v.СведНП.ИННЮЛ,
			&v.СведНП.НаимОрг,
			&v.СведНаруш.СумШтраф); err != nil {
			return nil, err
		}
		data = append(data, v)
	}
	return data, nil
}

// SelectСведенияОбУчастииВКонсГруппе - Сведения об участии в консолидированной группе налогоплательщиков
func (db *DB) SelectСведенияОбУчастииВКонсГруппе(ctx context.Context, inn string) (*models.СведенияОбУчастииВКонсГруппе, error) {
	const q = `
	SELECT
		датасост::text
		, наиморг
		, иннюл
		, признучкгн::text
	FROM
		public.сведенияобучастиивконсгруппе
	WHERE
		иннюл = $1::TEXT
	LIMIT 1`
	var v models.СведенияОбУчастииВКонсГруппе
	err := db.pool.QueryRow(ctx, q, inn).Scan(
		&v.ДатаДок,
		&v.СведНП.НаимОрг,
		&v.СведНП.ИННЮЛ,
		&v.СведКГН.ПризнУчКГН)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	switch v.СведКГН.ПризнУчКГН {
	case "1":
		v.СведКГН.ПризнУчКГН = "ответственный участник консолидированной группы налогоплательщиков"
	case "2":
		v.СведКГН.ПризнУчКГН = "участник консолидированной группы налогоплательщиков"
	}
	return &v, nil
}

// SelectГостиница - данные о гостиницах
// TODO: огрн может приходить пустой строкой - сделать CASE
func (db *DB) SelectГостиница(ctx context.Context, inn, ogrn string) ([]Hotel, error) {
	const q = `
	SELECT
	(row_to_json(h.*))::jsonb || jsonb_build_object('classification', kg_j.k) || jsonb_build_object('rooms', ng_j.n) AS hotel_data
	FROM
		public.hotels AS h
	LEFT JOIN LATERAL(
			SELECT
				coalesce(jsonb_agg(row_to_json(hc.*) ORDER BY hc.date_issued ASC), '[]'::jsonb) AS k
			FROM
				public.hotels_classification AS hc
			WHERE
				hc.federal_number = h.federal_number
		) AS kg_j ON
		TRUE
	LEFT JOIN LATERAL (
			SELECT
				coalesce(jsonb_agg(row_to_json(hr.*) ORDER BY hr.category), '[]'::jsonb) AS n
			FROM
				public.hotels_rooms AS hr
			WHERE
				hr.federal_number = h.federal_number
		) AS ng_j ON
		TRUE
	WHERE
		h.inn = $1::TEXT
		AND h.ogrn = $2::TEXT`
	rows, err := db.pool.Query(ctx, q, inn, ogrn)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var data = make([]Hotel, 0, 1000)
	for rows.Next() {
		var h Hotel
		if err := rows.Scan(&h); err != nil {
			return nil, err
		}
		data = append(data, h)
	}
	return data, nil
}

// SelectПривлеченииУчастникаЗакупкиКАдминистративнойОтветственности - получение данных Информация о привлечении участника закупки к административной ответственности по ст. 19.28 КоАП
func (db *DB) SelectПривлеченииУчастникаЗакупкиКАдминистративнойОтветственности(ctx context.Context, inn string) ([]models.ПривлеченииУчастникаЗакупкиКАдминистративнойОтветственности, error) {
	const q = `
	SELECT
		инн
		, огрн
		, кпп
		, наименованиеюл
		, типучастника
		, суд
		, номердела
		, датавынесенияпостановления::text
		, датавступлениявзаконнуюсилу::text
	FROM
		public.информация1928коап
	WHERE инн = $1::TEXT
	ORDER BY датавынесенияпостановления`
	rows, err := db.pool.Query(ctx, q, inn)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var data = make([]models.ПривлеченииУчастникаЗакупкиКАдминистративнойОтветственности, 0)
	for rows.Next() {
		var v models.ПривлеченииУчастникаЗакупкиКАдминистративнойОтветственности
		if err := rows.Scan(
			&v.ИНН,
			&v.ОГРН,
			&v.КПП,
			&v.НаименованиеЮЛ,
			&v.ТипУчастника,
			&v.Суд,
			&v.НомерДела,
			&v.ДатаВынесенияПостановления,
			&v.ДатаВступленияВЗаконнуюСилу); err != nil {
			return nil, err
		}
		data = append(data, v)
	}
	return data, nil
}

func (db *DB) GetHotelsView(ctx context.Context) ([]Hotel, error) {
	rows, err := db.pool.Query(ctx, `SELECT hotel_data FROM public.hotels_view`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var data = make([]Hotel, 0, 1000)
	for rows.Next() {
		var h Hotel
		if err := rows.Scan(&h); err != nil {
			return nil, err
		}
		data = append(data, h)
	}
	return data, nil
}

func (db *DB) DownloadHotels(ctx context.Context, w io.Writer) error {
	hotels, err := db.GetHotelsView(ctx)
	if err != nil {
		return err
	}
	return ToExcel(w, hotels)
}

// SelectHandbookCategoriesIP  -справочник категорий ИП
func (db *DB) SelectHandbookCategoriesIP(ctx context.Context) (data []CategoriesIP, err error) {
	const q = `
	SELECT
		jsonb_build_object(
			'category'
			, category
			, 'subcategories'
			, jsonb_agg(subcategory)
		)
	FROM
		public.handbook_categories_ip
	GROUP BY
		category`
	data = make([]CategoriesIP, 0, 20)
	rows, err := db.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var cip CategoriesIP
		if err := rows.Scan(&cip); err != nil {
			return nil, err
		}
		data = append(data, cip)
	}
	return
}

// SelectLegalStatuses - получить справочник расшифровки статусов ЮЛ
func (db *DB) SelectLegalStatuses(ctx context.Context) (data LegalStatus, err error) {
	const q = `SELECT status_full_name, status_name FROM egr.statuses`
	rows, err := db.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	data = make(LegalStatus)
	for rows.Next() {
		var (
			sfn string
			sn  string
		)
		if err := rows.Scan(&sfn, &sn); err != nil {
			return nil, err
		}
		data[sfn] = sn
	}
	return
}

// SelectRDS - Сведения из Реестра деклараций о соответствии
func (db *DB) SelectRDS(ctx context.Context, inn, ogrn string) ([]models.RDS, error) {
	const q = `
SELECT
	id_decl
	, reg_number
	, decl_status
	, decl_type
	, date_beginning::text
	, COALESCE(date_finish::TEXT, '') AS date_finish
	, declaration_scheme
	, product_object_type_decl
	, product_type
	, product_group
	, product_name
	, asproduct_info
	, product_tech_reg
	, organ_to_certification_name
	, organ_to_certification_reg_number
	, basis_for_decl
	, old_basis_for_decl
	, applicant_type
	, person_applicant_type
	, applicant_ogrn
	, applicant_inn
	, applicant_name
	, manufacturer_type
	, manufacturer_ogrn
	, manufacturer_inn
	, manufacturer_name
FROM
	public.rds AS r
WHERE
	r.applicant_inn = $1::text
	AND
	r.applicant_ogrn = $2::text
ORDER BY
	r.date_finish ASC
	, r.id_decl ASC`
	rows, err := db.pool.Query(ctx, q, inn, ogrn)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list = make([]models.RDS, 0)
	for rows.Next() {
		var v models.RDS
		if err := rows.Scan(
			&v.IDDecl,
			&v.RegNumber,
			&v.DeclStatus,
			&v.DeclType,
			&v.DateBeginning,
			&v.DateFinish,
			&v.DeclarationScheme,
			&v.ProductObjectTypeDecl,
			&v.ProductType,
			&v.ProductGroup,
			&v.ProductName,
			&v.AsproductInfo,
			&v.ProductTechReg,
			&v.OrganToCertificationName,
			&v.OrganToCertificationRegNumber,
			&v.BasisForDecl,
			&v.OldBasisForDecl,
			&v.ApplicantType,
			&v.PersonApplicantType,
			&v.ApplicantOGRN,
			&v.ApplicantINN,
			&v.ApplicantName,
			&v.ManufacturerType,
			&v.ManufacturerOGRN,
			&v.ManufacturerINN,
			&v.ManufacturerName,
		); err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	return list, nil
}

func (db *DB) SelectUnscheduledInspections(ctx context.Context, inn, ogrn string) ([]models.InspectionFGIS, error) {
	const q = `
	SELECT
		ui.value
	FROM
		fgis.unscheduled_inspections ui
	WHERE
		ui.inn = $1::text
		AND ui.ogrn = $2::text
	ORDER BY
		ui.start_date ASC
		, ui.erpid ASC`
	rows, err := db.pool.Query(ctx, q, inn, ogrn)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var data = make([]models.InspectionFGIS, 0)
	for rows.Next() {
		var insp models.InspectionFGIS
		if err := rows.Scan(&insp); err != nil {
			return nil, err
		}
		data = append(data, insp)
	}
	return data, nil
}

func (db *DB) SelectScheduledInspections(ctx context.Context, inn, ogrn string) ([]models.InspectionFGIS, error) {
	const q = `
	SELECT
		si.value
	FROM
		fgis.scheduled_inspections si
	WHERE
		si.inn = $1::text
		AND si.ogrn = $2::text
	ORDER BY
		si.start_date ASC
		, si.erpid ASC`
	rows, err := db.pool.Query(ctx, q, inn, ogrn)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var data = make([]models.InspectionFGIS, 0)
	for rows.Next() {
		var insp models.InspectionFGIS
		if err := rows.Scan(&insp); err != nil {
			return nil, err
		}
		data = append(data, insp)
	}
	return data, nil
}

func (db *DB) SelectSSHR(ctx context.Context, inn string) (*SSHR, error) {
	const q = `
	SELECT
		sshr.датасост::TEXT
		, sshr.колраб
	FROM
		public.сведенияосреднчислработников sshr
	WHERE
		sshr.иннюл = $1::TEXT`
	var sshr SSHR
	if err := db.pool.QueryRow(ctx, q, inn).Scan(&sshr.Date, &sshr.Count); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &sshr, nil
}

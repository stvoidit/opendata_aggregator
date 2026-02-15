package store

import (
	"context"
	"opendataaggregator/models"
	"opendataaggregator/models/egr"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

// BatchInsertНалоговыйРежимНалогоплательщика - пакетная вставка в БД объекта models.НалоговыйРежимНалогоплательщика.
// Запись добавляется, либо обновляются данные, если запись уже есть
func (db *DB) BatchInsertНалоговыйРежимНалогоплательщика(ctx context.Context, nrns []models.НалоговыйРежимНалогоплательщика) (err error) {
	const q = `INSERT INTO public.режимналогоплательщика
	(датадок, наиморг, иннюл, есхн, усн, енвд, срп)
	VALUES($1::date, $2::text, $3::text, $4::bool, $5::bool, $6::bool, $7::bool)
	ON CONFLICT (иннюл)
	DO UPDATE
	SET
		датадок = EXCLUDED.датадок,
		наиморг = EXCLUDED.наиморг,
		иннюл = EXCLUDED.иннюл,
		есхн = EXCLUDED.есхн,
		усн = EXCLUDED.усн,
		енвд = EXCLUDED.енвд,
		срп = EXCLUDED.срп`
	var b pgx.Batch
	for _, nrn := range nrns {
		b.Queue(q,
			nrn.DocDate(),
			nrn.СведНП.НаимОрг,
			nrn.СведНП.ИННЮЛ,
			nrn.СведСНР.ПризнЕСХН,
			nrn.СведСНР.ПризнУСН,
			nrn.СведСНР.ПризнЕНВД,
			nrn.СведСНР.ПризнСРП)
	}
	return db.pool.SendBatch(ctx, &b).Close()
}

// BatchСведенияОСуммахНедоимки - пакетная вставка в БД объекта models.СведенияОСуммахНедоимки
func (db *DB) BatchСведенияОСуммахНедоимки(ctx context.Context, ssns []models.СведенияОСуммахНедоимки) (err error) {
	const q = `INSERT INTO public.сведенияосуммахнедоимки
	(датадок, наиморг, иннюл, наимналог, сумнедналог, сумпени, сумштраф, общсумнедоим, иддок, датасост)
	VALUES($1::date, $2::text, $3::text, $4::text, $5, $6, $7, $8, $9, $10::date)
	ON CONFLICT (иддок, наимналог)
	DO UPDATE
	SET
		датадок = EXCLUDED.датадок
		, наиморг = EXCLUDED.наиморг
		, сумнедналог = EXCLUDED.сумнедналог
		, сумпени = EXCLUDED.сумпени
		, сумштраф = EXCLUDED.сумштраф
		, общсумнедоим = EXCLUDED.общсумнедоим
		, иддок = EXCLUDED.иддок
		, датасост = EXCLUDED.датасост`
	var b pgx.Batch
	for _, ssn := range ssns {
		docdate := ssn.DocDate()
		datebuild := ssn.DateBuild()
		for _, sn := range ssn.СведНедоим {
			b.Queue(q,
				docdate,
				ssn.СведНП.НаимОрг,
				ssn.СведНП.ИННЮЛ,
				sn.НаимНалог,
				sn.СумНедНалог,
				sn.СумПени,
				sn.СумШтраф,
				sn.ОбщСумНедоим,
				ssn.ИдДок,
				datebuild)
		}
	}
	start := time.Now()
	defer func() {
		log.Debug().Str("SendBatch", time.Since(start).String()).Msg("BatchСведенияОСуммахНедоимки")
	}()
	return db.pool.SendBatch(ctx, &b).Close()
}

// BatchInsertРосаккредитация - пакетная вставка в БД объекта models.НалоговыйРежимНалогоплательщика
func (db *DB) BatchInsertРосаккредитация(ctx context.Context, rks []models.Росаккредитация) (err error) {
	const q = `INSERT
	INTO
	public.росаккредитация
	(
		id_cert
		, cert_status
		, cert_type
		, reg_number
		, date_begining
		, date_finish
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
		, production_analysis_for_act_date
	)
	VALUES(
		$1,$2,$3,$4,$5::date,$6::date,$7,$8,$9,$10,
		$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,
		$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,
		$31,$32,$33,$34,$35,$36,$37,$38,$39,$40,
		$41,$42,$43,$44,$45,$46,$47,$48,$49,$50,$51::date
	)
	ON CONFLICT (id_cert)
	DO UPDATE
	SET
		cert_status =EXCLUDED.cert_status
		, cert_type =EXCLUDED.cert_type
		, reg_number =EXCLUDED.reg_number
		, date_begining =EXCLUDED.date_begining
		, date_finish =EXCLUDED.date_finish
		, product_scheme =EXCLUDED.product_scheme
		, product_object_type_cert =EXCLUDED.product_object_type_cert
		, product_type =EXCLUDED.product_type
		, product_okpd2 =EXCLUDED.product_okpd2
		, product_tn_ved =EXCLUDED.product_tn_ved
		, product_tech_reg =EXCLUDED.product_tech_reg
		, product_group =EXCLUDED.product_group
		, product_name =EXCLUDED.product_name
		, product_info =EXCLUDED.product_info
		, applicant_type =EXCLUDED.applicant_type
		, person_applicant_type =EXCLUDED.person_applicant_type
		, applicant_ogrn =EXCLUDED.applicant_ogrn
		, applicant_inn =EXCLUDED.applicant_inn
		, applicant_phone =EXCLUDED.applicant_phone
		, applicant_fax =EXCLUDED.applicant_fax
		, applicant_email =EXCLUDED.applicant_email
		, applicant_website =EXCLUDED.applicant_website
		, applicant_name =EXCLUDED.applicant_name
		, applicant_director_name =EXCLUDED.applicant_director_name
		, applicant_address =EXCLUDED.applicant_address
		, applicant_address_actual =EXCLUDED.applicant_address_actual
		, manufacturer_type =EXCLUDED.manufacturer_type
		, manufacturer_ogrn =EXCLUDED.manufacturer_ogrn
		, manufacturer_inn =EXCLUDED.manufacturer_inn
		, manufacturer_phone =EXCLUDED.manufacturer_phone
		, manufacturer_fax =EXCLUDED.manufacturer_fax
		, manufacturer_email =EXCLUDED.manufacturer_email
		, manufacturer_website =EXCLUDED.manufacturer_website
		, manufacturer_name =EXCLUDED.manufacturer_name
		, manufacturer_director_name =EXCLUDED.manufacturer_director_name
		, manufacturer_country =EXCLUDED.manufacturer_country
		, manufacturer_address =EXCLUDED.manufacturer_address
		, manufacturer_address_actual =EXCLUDED.manufacturer_address_actual
		, manufacturer_address_filial =EXCLUDED.manufacturer_address_filial
		, organ_to_certification_name =EXCLUDED.organ_to_certification_name
		, organ_to_certification_reg_number =EXCLUDED.organ_to_certification_reg_number
		, organ_to_certification_head_name =EXCLUDED.organ_to_certification_head_name
		, basis_for_certificate =EXCLUDED.basis_for_certificate
		, old_basis_for_certificate =EXCLUDED.old_basis_for_certificate
		, fio_expert =EXCLUDED.fio_expert
		, fio_signatory =EXCLUDED.fio_signatory
		, product_national_standart =EXCLUDED.product_national_standart
		, production_analysis_for_act =EXCLUDED.production_analysis_for_act
		, production_analysis_for_act_number =EXCLUDED.production_analysis_for_act_number
		, production_analysis_for_act_date =EXCLUDED.production_analysis_for_act_date`

	var b pgx.Batch
	for _, rk := range rks {
		b.Queue(q,
			rk.IDcert,
			rk.CertStatus,
			rk.CertType,
			rk.RegNumber,
			isNullStr(rk.DateBegining),
			isNullStr(rk.DateFinish),
			rk.ProductScheme,
			rk.ProductObjectTypeCert,
			rk.ProductType,
			rk.ProductOKPD2,
			rk.ProductTnVed,
			rk.ProductTechReg,
			rk.ProductGroup,
			rk.ProductName,
			rk.ProductInfo,
			rk.ApplicantType,
			rk.PersonApplicantType,
			rk.ApplicantOGRN,
			rk.ApplicantINN,
			rk.ApplicantPhone,
			rk.ApplicantFax,
			rk.ApplicantEmail,
			rk.ApplicantWebsite,
			rk.ApplicantName,
			rk.ApplicantDirectorName,
			rk.ApplicantAddress,
			rk.ApplicantAddressActual,
			rk.ManufacturerType,
			rk.ManufacturerOGRN,
			rk.ManufacturerINN,
			rk.ManufacturerPhone,
			rk.ManufacturerFax,
			rk.ManufacturerEmail,
			rk.ManufacturerWebsite,
			rk.ManufacturerName,
			rk.ManufacturerDirectorName,
			rk.ManufacturerCountry,
			rk.ManufacturerAddress,
			rk.ManufacturerAddressActual,
			rk.ManufacturerAddressFilial,
			rk.OrganToCertificationName,
			rk.OrganToCertificationRegNumber,
			rk.OrganToCertificationHeadName,
			rk.BasisForCertificate,
			rk.OldBasisForCertificate,
			rk.DioExpert,
			rk.DioSignatory,
			rk.ProductNationalStandart,
			rk.ProductionAnalysisForAct,
			rk.ProductionAnalysisForActNumber,
			isNullStr(rk.ProductionAnalysisForActDate))
	}
	return db.pool.SendBatch(ctx, &b).Close()
}

// BatchInsertRDS - вставка Сведения из Реестра деклараций о соответствии
func (db *DB) BatchInsertRDS(ctx context.Context, list []models.RDS) (err error) {
	const q = `INSERT
		INTO
		public.rds
	(
		id_decl
		, reg_number
		, decl_status
		, decl_type
		, date_beginning
		, date_finish
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
	)
VALUES(
	$1,$2,$3,$4,
	case when $5 != '' then $5::date else null end,
	case when $6 != '' then $6::date else null end,
	$7,$8,$9,$10,
	$11,$12,$13,$14,$15,$16,$17,$18,
	$19,$20,$21,$22,$23,$24,$25,$26
)
ON
	CONFLICT (id_decl)
	DO
UPDATE
SET
	reg_number = EXCLUDED.reg_number
	, decl_status = EXCLUDED.decl_status
	, decl_type = EXCLUDED.decl_type
	, date_beginning = EXCLUDED.date_beginning
	, date_finish = EXCLUDED.date_finish
	, declaration_scheme = EXCLUDED.declaration_scheme
	, product_object_type_decl = EXCLUDED.product_object_type_decl
	, product_type = EXCLUDED.product_type
	, product_group = EXCLUDED.product_group
	, product_name = EXCLUDED.product_name
	, asproduct_info = EXCLUDED.asproduct_info
	, product_tech_reg = EXCLUDED.product_tech_reg
	, organ_to_certification_name = EXCLUDED.organ_to_certification_name
	, organ_to_certification_reg_number = EXCLUDED.organ_to_certification_reg_number
	, basis_for_decl = EXCLUDED.basis_for_decl
	, old_basis_for_decl = EXCLUDED.old_basis_for_decl
	, applicant_type = EXCLUDED.applicant_type
	, person_applicant_type = EXCLUDED.person_applicant_type
	, applicant_ogrn = EXCLUDED.applicant_ogrn
	, applicant_inn = EXCLUDED.applicant_inn
	, applicant_name = EXCLUDED.applicant_name
	, manufacturer_type = EXCLUDED.manufacturer_type
	, manufacturer_ogrn = EXCLUDED.manufacturer_ogrn
	, manufacturer_inn = EXCLUDED.manufacturer_inn
	, manufacturer_name = EXCLUDED.manufacturer_name
`
	var b pgx.Batch
	for _, v := range list {
		b.Queue(q,
			v.IDDecl,
			v.RegNumber,
			v.DeclStatus,
			v.DeclType,
			v.DateBeginning,
			v.DateFinish,
			v.DeclarationScheme,
			v.ProductObjectTypeDecl,
			v.ProductType,
			v.ProductGroup,
			v.ProductName,
			v.AsproductInfo,
			v.ProductTechReg,
			v.OrganToCertificationName,
			v.OrganToCertificationRegNumber,
			v.BasisForDecl,
			v.OldBasisForDecl,
			v.ApplicantType,
			v.PersonApplicantType,
			v.ApplicantOGRN,
			v.ApplicantINN,
			v.ApplicantName,
			v.ManufacturerType,
			v.ManufacturerOGRN,
			v.ManufacturerINN,
			v.ManufacturerName,
		)
	}
	return db.pool.SendBatch(ctx, &b).Close()
}

// BatchСведенияОбУчастииВКонсГруппе - пакетная вставка в БД объекта models.СведенияОбУчастииВКонсГруппе
func (db *DB) BatchСведенияОбУчастииВКонсГруппе(ctx context.Context, list []models.СведенияОбУчастииВКонсГруппе) (err error) {
	const q = `INSERT INTO public.сведенияобучастиивконсгруппе
	(датасост, наиморг, иннюл, признучкгн)
	VALUES($1::date, $2::text, $3::text, $4::int4)
	ON CONFLICT (иннюл)
	DO UPDATE
	SET
		датасост = EXCLUDED.датасост
		, наиморг = EXCLUDED.наиморг
		, признучкгн = EXCLUDED.признучкгн`
	var b pgx.Batch
	for _, v := range list {
		b.Queue(q,
			v.DocDate(),
			v.СведНП.НаимОрг,
			v.СведНП.ИННЮЛ,
			v.СведКГН.ПризнУчКГН)
	}
	return db.pool.SendBatch(ctx, &b).Close()
}

// BatchСведенияОСреднесписочнойЧисленностиРаботников - пакетная вставка в БД объекта models.СведенияОСреднесписочнойЧисленностиРаботников
func (db *DB) BatchСведенияОСреднесписочнойЧисленностиРаботников(ctx context.Context, list []models.СведенияОСреднесписочнойЧисленностиРаботников) (err error) {
	const q = `INSERT INTO public.сведенияосреднчислработников
	(датасост, наиморг, иннюл, колраб)
	VALUES($1::date, $2::text, $3::text, $4)
	ON CONFLICT (иннюл)
	DO UPDATE
	SET
		датасост = EXCLUDED.датасост
		, наиморг = EXCLUDED.наиморг
		, колраб = EXCLUDED.колраб`
	var b pgx.Batch
	for _, v := range list {
		b.Queue(q,
			v.DocDate(),
			v.СведНП.НаимОрг,
			v.СведНП.ИННЮЛ,
			v.СведССЧР.КолРаб)
	}
	return db.pool.SendBatch(ctx, &b).Close()
}

// BatchСведенияОбУплаченныхОрганизациейНалогов - пакетная вставка в БД объекта models.СведенияОбУплаченныхОрганизациейНалогов
func (db *DB) BatchСведенияОбУплаченныхОрганизациейНалогов(ctx context.Context, list []models.СведенияОбУплаченныхОрганизациейНалогов) (err error) {
	const q = `INSERT INTO public.сведенияобуплаченныхорганизацие
	(датасост, наиморг, иннюл, наимналог, сумуплнал)
	VALUES($1::date, $2::text, $3::text, $4::text, $5)
	ON CONFLICT (иннюл, наимналог)
	DO UPDATE
	SET
		датасост = EXCLUDED.датасост
		, наиморг = EXCLUDED.наиморг
		, сумуплнал = EXCLUDED.сумуплнал`
	var b pgx.Batch
	for _, v := range list {
		for _, n := range v.СвУплСумНал {
			b.Queue(q,
				v.DocDate(),
				v.СведНП.НаимОрг,
				v.СведНП.ИННЮЛ,
				n.НаимНалог,
				n.СумУплНал)
		}
	}
	return db.pool.SendBatch(ctx, &b).Close()
}

// RepaidFlagIpLegalList - установка флага "погашено" в позицию TRUE перед обновлением
func (db *DB) RepaidFlagIpLegalList(ctx context.Context) error {
	_, err := db.pool.Exec(ctx, `UPDATE ONLY public.исппроизввотнюрлиц SET repaid = TRUE`)
	return err
}

// UpdateIpLegalListDebTremainingBalanceZero - обновить поле repaid на TRUE, если debtremainingbalance = 0
func (db *DB) UpdateIpLegalListDebTremainingBalanceZero(ctx context.Context) error {
	_, err := db.pool.Exec(ctx, `UPDATE
    public.исппроизввотнюрлиц
SET
    repaid = TRUE
WHERE
    NOT repaid
    AND debtremainingbalance = (0)::NUMERIC`)
	return err
}

// BatchIpLegalList - пакетная вставка в БД объекта models.ИсполнительныеПроизводстваВОтношенииЮридическихЛиц.
// Добавление или изменение записей по уникальному ключу из 2х полей: "Номер исполнительного производства" + "Номер исполнительного документа"
func (db *DB) BatchIpLegalList(ctx context.Context, list []models.ИсполнительныеПроизводстваВОтношенииЮридическихЛиц) (err error) {
	const q = `INSERT
	INTO
	public.исппроизввотнюрлиц
	(
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
	)
	VALUES(
		$1
		, $2
		, $3
		, $4
		, $5
		, $6
		, $7
		, $8
		, $9
		, $10
		, $11
		, $12
		, $13
		, $14
		, $15
		, $16
		, $17
	)
	ON CONFLICT (numberofenforcementproceeding, numberofexecutivedocument)
	DO UPDATE
	SET
		nameofdebtor = EXCLUDED.nameofdebtor
		, addressofdebtororganization = EXCLUDED.addressofdebtororganization
		, actualaddressofdebtororganization = EXCLUDED.actualaddressofdebtororganization
		, numberofenforcementproceeding = EXCLUDED.numberofenforcementproceeding
		, dateofinstitutionproceeding = EXCLUDED.dateofinstitutionproceeding
		, totalnumberofenforcementproceedings = EXCLUDED.totalnumberofenforcementproceedings
		, executivedocumenttype = EXCLUDED.executivedocumenttype
		, dateofexecutivedocument = EXCLUDED.dateofexecutivedocument
		, numberofexecutivedocument = EXCLUDED.numberofexecutivedocument
		, objectofexecutivedocuments = EXCLUDED.objectofexecutivedocuments
		, objectofexecution = EXCLUDED.objectofexecution
		, amountdue = EXCLUDED.amountdue
		, debtremainingbalance = EXCLUDED.debtremainingbalance
		, departmentsofbailiffs = EXCLUDED.departmentsofbailiffs
		, addressofdepartmentsofbailiff = EXCLUDED.addressofdepartmentsofbailiff
		, debtortaxpayeridentificationnumber = EXCLUDED.debtortaxpayeridentificationnumber
		, taxpayeridentificationnumberoforganizationcollector = EXCLUDED.taxpayeridentificationnumberoforganizationcollector
		, repaid = FALSE`
	var b pgx.Batch
	for _, v := range list {
		b.Queue(q,
			v.NameOfDebtor,
			v.AddressOfDebtorOrganization,
			v.ActualAddressOfDebtorOrganization,
			v.NumberOfEnforcementProceeding,
			v.DateOfInstitutionProceeding,
			v.TotalNumberOfEnforcementProceedings,
			v.ExecutiveDocumentType,
			v.DateOfExecutiveDocument,
			v.NumberOfExecutiveDocument,
			v.ObjectOfExecutiveDocuments,
			v.ObjectOfExecution,
			v.AmountDue,
			v.DebtRemainingBalance,
			v.DepartmentsOfBailiffs,
			v.AddressOfDepartmentsOfBailiff,
			v.DebtorTaxpayerIdentificationNumber,
			v.TaxpayerIdentificationNumberOfOrganizationCollector)
	}
	return db.pool.SendBatch(ctx, &b).Close()
}

// BatchIpLegalListComplete - iplegallistcomplete
func (db *DB) BatchIpLegalListComplete(ctx context.Context, list []models.IpLegalListComplete) (err error) {
	const q = `INSERT
	INTO
	public.iplegallistcomplete
	(
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
		, datecompleteipreason
		, departmentsofbailiffs
		, addressofdepartmentsofbailiff
		, debtortaxpayeridentificationnumber
		, taxpayeridentificationnumberoforganizationcollector
	)
	VALUES(
		$1
		, $2
		, $3
		, $4
		, $5
		, $6
		, $7
		, $8
		, $9
		, $10
		, $11
		, $12
		, $13
		, $14
		, $15
		, $16
	)
	ON CONFLICT (numberofenforcementproceeding, numberofexecutivedocument)
	DO UPDATE
	SET
		nameofdebtor = EXCLUDED.nameofdebtor
		, addressofdebtororganization = EXCLUDED.addressofdebtororganization
		, actualaddressofdebtororganization = EXCLUDED.actualaddressofdebtororganization
		, numberofenforcementproceeding = EXCLUDED.numberofenforcementproceeding
		, dateofinstitutionproceeding = EXCLUDED.dateofinstitutionproceeding
		, totalnumberofenforcementproceedings = EXCLUDED.totalnumberofenforcementproceedings
		, executivedocumenttype = EXCLUDED.executivedocumenttype
		, dateofexecutivedocument = EXCLUDED.dateofexecutivedocument
		, numberofexecutivedocument = EXCLUDED.numberofexecutivedocument
		, objectofexecutivedocuments = EXCLUDED.objectofexecutivedocuments
		, objectofexecution = EXCLUDED.objectofexecution
		, datecompleteipreason = EXCLUDED.datecompleteipreason
		, departmentsofbailiffs = EXCLUDED.departmentsofbailiffs
		, addressofdepartmentsofbailiff = EXCLUDED.addressofdepartmentsofbailiff
		, debtortaxpayeridentificationnumber = EXCLUDED.debtortaxpayeridentificationnumber
		, taxpayeridentificationnumberoforganizationcollector = EXCLUDED.taxpayeridentificationnumberoforganizationcollector`
	var b pgx.Batch
	for _, v := range list {
		b.Queue(q,
			v.NameOfDebtor,
			v.AddressOfDebtorOrganization,
			v.ActualAddressOfDebtorOrganization,
			v.NumberOfEnforcementProceeding,
			v.DateOfInstitutionProceeding,
			v.TotalNumberOfEnforcementProceedings,
			v.ExecutiveDocumentType,
			v.DateOfExecutiveDocument,
			v.NumberOfExecutiveDocument,
			v.ObjectOfExecutiveDocuments,
			v.ObjectOfExecution,
			v.DateCompleteIPreason,
			v.DepartmentsOfBailiffs,
			v.AddressOfDepartmentsOfBailiff,
			v.DebtorTaxpayerIdentificationNumber,
			v.TaxpayerIdentificationNumberOfOrganizationCollector)
	}
	return db.pool.SendBatch(ctx, &b).Close()
}

// BatchОткрытыйРеестрТоварныхЗнаков - пакетная вставка в БД объекта models.ОткрытыйРеестрТоварныхЗнаков
// TODO: очень много кривых дат в формате "20060102"
func (db *DB) BatchОткрытыйРеестрТоварныхЗнаков(ctx context.Context, list []models.ОткрытыйРеестрТоварныхЗнаков) (err error) {
	const q = `INSERT
	INTO
	public.открытыйреестртоварныхзнаков
	(
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
	)
	VALUES(
		$1
		, $2
		, $3
		, $4
		, $5
		, $6
		, $7
		, $8
		, $9
		, $10
		, $11
		, $12
		, $13
		, $14
		, $15
		, $16
		, $17
		, $18
		, $19
		, $20
		, $21
		, $22
		, $23
		, $24
		, $25
		, $26
		, $27
		, $28
		, $29
		, $30
		, $31
		, $32
		, $33
		, $34
		, $35
		, $36
		, $37
		, $38
		, $39
		, $40
		, $41
		, $42
		, $43
		, $44
		, $45
		, $46
		, $47
		, $48
		, $49
		, $50
		, $51
		, $52
		, $53
		, $54
		, $55
		, $56
		, $57
		, $58
	)
	ON CONFLICT (registrationnumber)
	DO UPDATE
	SET
		registrationdate = EXCLUDED.registrationdate
		, applicationnumber = EXCLUDED.applicationnumber
		, applicationdate = EXCLUDED.applicationdate
		, prioritydate = EXCLUDED.prioritydate
		, exhibitionprioritydate = EXCLUDED.exhibitionprioritydate
		, parisconventionprioritynumber = EXCLUDED.parisconventionprioritynumber
		, parisconventionprioritydate = EXCLUDED.parisconventionprioritydate
		, parisconventionprioritycountrycode = EXCLUDED.parisconventionprioritycountrycode
		, initialapplicationnumber = EXCLUDED.initialapplicationnumber
		, initialapplicationorioritydate = EXCLUDED.initialapplicationorioritydate
		, initialregistrationnumber = EXCLUDED.initialregistrationnumber
		, initialregistrationdate = EXCLUDED.initialregistrationdate
		, internationalregistrationnumber = EXCLUDED.internationalregistrationnumber
		, internationalregistrationdate = EXCLUDED.internationalregistrationdate
		, internationalregistrationprioritydate = EXCLUDED.internationalregistrationprioritydate
		, internationalregistrationentrydate = EXCLUDED.internationalregistrationentrydate
		, applicationnumberforrecognitionoftrademarkfromcrimea = EXCLUDED.applicationnumberforrecognitionoftrademarkfromcrimea
		, applicationdateforrecognitionoftrademarkfromcrimea = EXCLUDED.applicationdateforrecognitionoftrademarkfromcrimea
		, crimeantrademarkapplicationnumberforstateregistrationinukraine = EXCLUDED.crimeantrademarkapplicationnumberforstateregistrationinukraine
		, crimeantrademarkapplicationdateforstateregistrationinukraine = EXCLUDED.crimeantrademarkapplicationdateforstateregistrationinukraine
		, crimeantrademarkcertificatenumberinukraine = EXCLUDED.crimeantrademarkcertificatenumberinukraine
		, exclusiverightstransferagreementregistrationnumber = EXCLUDED.exclusiverightstransferagreementregistrationnumber
		, exclusiverightstransferagreementregistrationdate = EXCLUDED.exclusiverightstransferagreementregistrationdate
		, legallyrelatedapplications = EXCLUDED.legallyrelatedapplications
		, legallyrelatedregistrations = EXCLUDED.legallyrelatedregistrations
		, expirationdate = EXCLUDED.expirationdate
		, rightholdername = EXCLUDED.rightholdername
		, foreignrightholdername = EXCLUDED.foreignrightholdername
		, rightholderaddress = EXCLUDED.rightholderaddress
		, rightholdercountrycode = EXCLUDED.rightholdercountrycode
		, rightholderogrn = EXCLUDED.rightholderogrn
		, rightholderinn = EXCLUDED.rightholderinn
		, correspondenceaddress = EXCLUDED.correspondenceaddress
		, collective = EXCLUDED.collective
		, collectiveusers = EXCLUDED.collectiveusers
		, extractionfromcharterofthecollectivetrademark = EXCLUDED.extractionfromcharterofthecollectivetrademark
		, colorspecification = EXCLUDED.colorspecification
		, unprotectedelements = EXCLUDED.unprotectedelements
		, kindspecification = EXCLUDED.kindspecification
		, threedimensional = EXCLUDED.threedimensional
		, threedimensionalspecification = EXCLUDED.threedimensionalspecification
		, holographic = EXCLUDED.holographic
		, holographicspecification = EXCLUDED.holographicspecification
		, sound = EXCLUDED.sound
		, soundspecification = EXCLUDED.soundspecification
		, olfactory = EXCLUDED.olfactory
		, olfactoryspecification = EXCLUDED.olfactoryspecification
		, color = EXCLUDED.color
		, colortrademarkspecification = EXCLUDED.colortrademarkspecification
		, light = EXCLUDED.light
		, lightspecification = EXCLUDED.lightspecification
		, changing = EXCLUDED.changing
		, changingspecification = EXCLUDED.changingspecification
		, positional = EXCLUDED.positional
		, positionalspecification = EXCLUDED.positionalspecification
		, actual = EXCLUDED.actual
		, publicationurl = EXCLUDED.publicationurl`
	var b pgx.Batch
	for _, v := range list {
		b.Queue(q,
			v.RegistrationNumber,
			v.RegistrationDate,
			v.ApplicationNumber,
			v.ApplicationDate,
			v.PriorityDate,
			v.ExhibitionPriorityDate,
			v.ParisConventionPriorityNumber,
			v.ParisConventionPriorityDate,
			v.ParisConventionPriorityCountryCode,
			v.InitialApplicationNumber,
			v.InitialApplicationOriorityDate,
			v.InitialRegistrationNumber,
			v.InitialRegistrationDate,
			v.InternationalRegistrationNumber,
			v.InternationalRegistrationDate,
			v.InternationalRegistrationPriorityDate,
			v.InternationalRegistrationEntryDate,
			v.ApplicationNumberForRecognitionOfTrademarkFromCrimea,
			v.ApplicationDateForRecognitionOfTrademarkFromCrimea,
			v.CrimeanTrademarkApplicationNumberForStateRegistrationInUkraine,
			v.CrimeanTrademarkApplicationDateForStateRegistrationInUkraine,
			v.CrimeanTrademarkCertificateNumberInUkraine,
			v.ExclusiveRightsTransferAgreementRegistrationNumber,
			v.ExclusiveRightsTransferAgreementRegistrationDate,
			v.LegallyRelatedApplications,
			v.LegallyRelatedRegistrations,
			v.ExpirationDate,
			v.RightHolderName,
			v.ForeignRightHolderName,
			v.RightHolderAddress,
			v.RightHolderCountryCode,
			v.RightHolderOgrn,
			v.RightHolderInn,
			v.CorrespondenceAddress,
			v.Collective,
			v.CollectiveUsers,
			v.ExtractionFromCharterOfTheCollectiveTrademark,
			v.ColorSpecification,
			v.UnprotectedElements,
			v.KindSpecification,
			v.Threedimensional,
			v.ThreedimensionalSpecification,
			v.Holographic,
			v.HolographicSpecification,
			v.Sound,
			v.SoundSpecification,
			v.Olfactory,
			v.OlfactorySpecification,
			v.Color,
			v.ColorTrademarkSpecification,
			v.Light,
			v.LightSpecification,
			v.Changing,
			v.ChangingSpecification,
			v.Positional,
			v.PositionalSpecification,
			v.Actual,
			v.PublicationURL)
	}
	return db.pool.SendBatch(ctx, &b).Close()
}

// BatchОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков - пакетная вставка в БД объекта models.ОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков
func (db *DB) BatchОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков(ctx context.Context, list []models.ОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков) (err error) {
	const q = `INSERT
    INTO
    public.реестробщеизвестныхтоварныхзнак
	(
		registrationnumber
		, registrationdate
		, wellknowntrademarkdate
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
	)
	VALUES(
		$1
		, $2
		, $3
		, $4
		, $5
		, $6
		, $7
		, $8
		, $9
		, $10
		, $11
		, $12
		, $13
		, $14
		, $15
		, $16
		, $17
		, $18
		, $19
		, $20
		, $21
		, $22
		, $23
		, $24
		, $25
		, $26
		, $27
		, $28
		, $29
		, $30
		, $31
		, $32
		, $33
		, $34
		, $35
	)
	ON CONFLICT (registrationnumber)
	DO UPDATE
	SET
		registrationdate = EXCLUDED.registrationdate
		, wellknowntrademarkdate = EXCLUDED.wellknowntrademarkdate
		, legallyrelatedregistrations = EXCLUDED.legallyrelatedregistrations
		, rightholdername = EXCLUDED.rightholdername
		, foreignrightholdername = EXCLUDED.foreignrightholdername
		, rightholderaddress = EXCLUDED.rightholderaddress
		, rightholdercountrycode = EXCLUDED.rightholdercountrycode
		, rightholderogrn = EXCLUDED.rightholderogrn
		, rightholderinn = EXCLUDED.rightholderinn
		, correspondenceaddress = EXCLUDED.correspondenceaddress
		, collective = EXCLUDED.collective
		, collectiveusers = EXCLUDED.collectiveusers
		, extractionfromcharterofcollectivetrademark = EXCLUDED.extractionfromcharterofcollectivetrademark
		, colorspecification = EXCLUDED.colorspecification
		, unprotectedelements = EXCLUDED.unprotectedelements
		, kindspecification = EXCLUDED.kindspecification
		, threedimensional = EXCLUDED.threedimensional
		, threedimensionalspecification = EXCLUDED.threedimensionalspecification
		, holographic = EXCLUDED.holographic
		, holographicspecification = EXCLUDED.holographicspecification
		, sound = EXCLUDED.sound
		, soundspecification = EXCLUDED.soundspecification
		, olfactory = EXCLUDED.olfactory
		, olfactoryspecification = EXCLUDED.olfactoryspecification
		, color = EXCLUDED.color
		, colortrademarkspecification = EXCLUDED.colortrademarkspecification
		, light = EXCLUDED.light
		, lightspecification = EXCLUDED.lightspecification
		, changing = EXCLUDED.changing
		, changingspecification = EXCLUDED.changingspecification
		, positional = EXCLUDED.positional
		, positionalspecification = EXCLUDED.positionalspecification
		, actual = EXCLUDED.actual
		, publicationurl = EXCLUDED.publicationurl`
	var b pgx.Batch
	for _, v := range list {
		b.Queue(q,
			v.RegistrationNumber,
			v.RegistrationDateDateTime(),
			v.WellKnownTrademarkDateDateTime(),
			v.LegallyRelatedRegistrations,
			v.RightHolderName,
			v.ForeignRightHolderName,
			v.RightHolderAddress,
			v.RightHolderCountryCode,
			v.RightHolderOgrn,
			v.RightHolderInn,
			v.CorrespondenceAddress,
			v.Collective,
			v.CollectiveUsers,
			v.ExtractionFromCharterOfCollectiveTrademark,
			v.ColorSpecification,
			v.UnprotectedElements,
			v.KindSpecification,
			v.Threedimensional,
			v.ThreedimensionalSpecification,
			v.Holographic,
			v.HolographicSpecification,
			v.Sound,
			v.SoundSpecification,
			v.Olfactory,
			v.OlfactorySpecification,
			v.Color,
			v.ColorTrademarkSpecification,
			v.Light,
			v.LightSpecification,
			v.Changing,
			v.ChangingSpecification,
			v.Positional,
			v.PositionalSpecification,
			v.Actual,
			v.PublicationURL)
	}
	return db.pool.SendBatch(ctx, &b).Close()
}

// BatchОКВЭД - пакетная вставка в БД объекта models.ОКВЭД
func (db *DB) BatchОКВЭД(ctx context.Context, list []models.ОКВЭД) (err error) {
	const q = `INSERT INTO public.оквэд
	(кодоквэд, наимоквэд)
	VALUES($1, $2)
	ON CONFLICT (кодоквэд)
	DO UPDATE
	SET
		наимоквэд = EXCLUDED.наимоквэд`
	var b pgx.Batch
	for _, v := range list {
		b.Queue(q,
			v.КодОКВЭД,
			v.НаимОКВЭД)
	}
	return db.pool.SendBatch(ctx, &b).Close()
}

// BatchСубъектМалогоИСреднегоПредпринимательства - пакетная вставка в БД объекта models.СубъектовМалогоИСреднегоПредпринимательства
func (db *DB) BatchСубъектМалогоИСреднегоПредпринимательства(ctx context.Context, list []models.РеестрСубъектовМалогоИСреднегоПредпринимательства) (err error) {
	const q = `INSERT
		INTO
		public.субъектымалогоисреднегопредприн
	(
		датасост
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
		, инн
		, огрн
		, свпрод
		, свпрогпарт
		, свконтр
		, свдог
	)
	VALUES(
		$1
		, $2
		, $3
		, $4
		, $5
		, $6
		, $7
		, $8
		, $9
		, $10
		, $11
		, $12
		, $13
		, $14
		, $15
		, $16
		, $17
		, $18
		, $19
		, $20
		, $21
		, $22
		, $23
		, $24
		, $25
		, $26
		, $27
		, $28
		, $29
		, $30
		, $31
		, $32
		, $33
		, $34
		, $35
		, $36
		, $37
		, $38
		, $39
		, $40
		, $41
	)
	ON CONFLICT (инн, огрн)
	DO UPDATE
	SET
		датасост = EXCLUDED.датасост
		, датавклмсп = EXCLUDED.датавклмсп
		, видсубмсп = EXCLUDED.видсубмсп
		, катсубмсп = EXCLUDED.катсубмсп
		, призновмсп = EXCLUDED.призновмсп
		, сведсоцпред = EXCLUDED.сведсоцпред
		, ссчр = EXCLUDED.ссчр
		, наиморг = EXCLUDED.наиморг
		, наиморгсокр = EXCLUDED.наиморгсокр
		, иннюл = EXCLUDED.иннюл
		, огрнюл = EXCLUDED.огрнюл
		, иннфл = EXCLUDED.иннфл
		, огрнип = EXCLUDED.огрнип
		, фамилия = EXCLUDED.фамилия
		, имя = EXCLUDED.имя
		, отчество = EXCLUDED.отчество
		, номлиценз = EXCLUDED.номлиценз
		, даталиценз = EXCLUDED.даталиценз
		, датаначлиценз = EXCLUDED.датаначлиценз
		, датаконлиценз = EXCLUDED.датаконлиценз
		, датаостлиценз = EXCLUDED.датаостлиценз
		, серлиценз = EXCLUDED.серлиценз
		, видлиценз = EXCLUDED.видлиценз
		, оргвыдлиценз = EXCLUDED.оргвыдлиценз
		, оргостлиценз = EXCLUDED.оргостлиценз
		, наимлицвд = EXCLUDED.наимлицвд
		, кодрегион = EXCLUDED.кодрегион
		, регионтип = EXCLUDED.регионтип
		, регионнаим = EXCLUDED.регионнаим
		, районтип = EXCLUDED.районтип
		, районнаим = EXCLUDED.районнаим
		, городтип = EXCLUDED.городтип
		, городнаим = EXCLUDED.городнаим
		, населпункттип = EXCLUDED.населпункттип
		, населпунктнаим = EXCLUDED.населпунктнаим
		, инн = EXCLUDED.инн
		, огрн = EXCLUDED.огрн
		, свпрод = EXCLUDED.свпрод
		, свпрогпарт = EXCLUDED.свпрогпарт
		, свконтр = EXCLUDED.свконтр
		, свдог = EXCLUDED.свдог
	WHERE public.субъектымалогоисреднегопредприн.датасост < EXCLUDED.датасост`
	const qokveds = `
	INSERT INTO public.смпоквэды
	(инн, огрн, кодоквэд, основной)
	SELECT $1, $2, $3, $4
	WHERE NOT EXISTS (SELECT 1 FROM public.субъектымалогоисреднегопредприн smp WHERE smp.инн = $1 AND smp.огрн = $2)
	ON CONFLICT DO NOTHING`

	{
		// добавление СМП
		var b pgx.Batch
		for _, v := range list {
			docdate := v.DocDate()
			inn := v.INN()
			ogrn := v.OGRN()
			b.Queue(q,
				docdate,
				v.ДатаВклМСП,
				v.ВидСубМСП,
				v.КатСубМСП,
				v.ПризНовМСП,
				v.СведСоцПред,
				v.ССЧР,
				v.ОргВклМСП.НаимОрг,
				v.ОргВклМСП.НаимОргСокр,
				v.ОргВклМСП.ИННЮЛ,
				v.ОргВклМСП.ОГРН,
				v.ИПВклМСП.ИННФЛ,
				v.ИПВклМСП.ОГРНИП,
				v.ИПВклМСП.ФИОИП.Фамилия,
				v.ИПВклМСП.ФИОИП.Имя,
				v.ИПВклМСП.ФИОИП.Отчество,
				v.СвЛиценз.НомЛиценз,
				v.СвЛиценз.ДатаЛиценз,
				v.СвЛиценз.ДатаНачЛиценз,
				v.СвЛиценз.ДатаКонЛиценз,
				v.СвЛиценз.ДатаОстЛиценз,
				v.СвЛиценз.СерЛиценз,
				v.СвЛиценз.ВидЛиценз,
				v.СвЛиценз.ОргВыдЛиценз,
				v.СвЛиценз.ОргОстЛиценз,
				v.СвЛиценз.НаимЛицВД,
				v.СведМН.КодРегион,
				v.СведМН.Регион.Тип,
				v.СведМН.Регион.Наим,
				v.СведМН.Район.Тип,
				v.СведМН.Район.Наим,
				v.СведМН.Город.Тип,
				v.СведМН.Город.Наим,
				v.СведМН.НаселПункт.Тип,
				v.СведМН.НаселПункт.Наим,
				inn,
				ogrn,
				v.СвПрод,
				v.СвПрогПарт,
				v.СвКонтр,
				v.СвДог,
			)
		}
		if err := db.pool.SendBatch(ctx, &b).Close(); err != nil {
			return err
		}
	}
	{
		// добавление данных об ОКВЭД
		var b pgx.Batch
		for _, v := range list {
			inn := v.INN()
			ogrn := v.OGRN()
			if v.СвОКВЭД.СвОКВЭДОсн.КодОКВЭД != "" {
				b.Queue(qokveds,
					inn,
					ogrn,
					v.СвОКВЭД.СвОКВЭДОсн.КодОКВЭД,
					true)
			}
			for _, o := range v.СвОКВЭД.СвОКВЭДДоп {
				b.Queue(qokveds,
					inn,
					ogrn,
					o.КодОКВЭД,
					false)
			}
		}
		if err := db.pool.SendBatch(ctx, &b).Close(); err != nil {
			return err
		}
	}
	return
}

// BatchНалоговыхПравонарушенияхИМерахОтветственности - пакетная вставка в БД объекта models.НалоговыхПравонарушенияхИМерахОтветственности.
// Добавление, но на всякий случай обновление. По идее выгрузка не предполагает изменений в старых записях, но на всякий случай сделана обработка
func (db *DB) BatchНалоговыхПравонарушенияхИМерахОтветственности(ctx context.Context, list []models.НалоговыхПравонарушенияхИМерахОтветственности) (err error) {
	const q = `INSERT INTO public.налоговыеправонарушенияиштрафы
	(датасост, иннюл, наиморг, сумштраф)
	VALUES($1, $2, $3, $4)
	ON CONFLICT (датасост, иннюл)
	DO UPDATE
	SET
		датасост = EXCLUDED.датасост
		, иннюл = EXCLUDED.иннюл
		, наиморг = EXCLUDED.наиморг
		, сумштраф = EXCLUDED.сумштраф`
	var b pgx.Batch
	for _, v := range list {
		b.Queue(q,
			v.DocDate(),
			v.СведНП.ИННЮЛ,
			v.СведНП.НаимОрг,
			v.СведНаруш.СумШтраф)
	}
	return db.pool.SendBatch(ctx, &b).Close()
}

// BatchРеестрДисквалифицированныхЛиц - пакетная вставка в БД объекта models.РеестрДисквалифицированныхЛиц
func (db *DB) BatchРеестрДисквалифицированныхЛиц(ctx context.Context, list []models.РеестрДисквалифицированныхЛиц) (err error) {
	const q = `INSERT
	INTO
	public.реестрдисквалифицированныхлиц
	(
		id
		, fio
		, bdate
		, bplace
		, orgname
		, inn
		, positionfl
		, nkoap
		, gorgname
		, sudfio
		, sudposition
		, disqualificationduration
		, disstartdate
		, disenddate
	)
	VALUES(
		$1
		, $2
		, $3
		, $4
		, $5
		, $6
		, $7
		, $8
		, $9
		, $10
		, $11
		, $12
		, $13
		, $14
	)
	ON CONFLICT DO NOTHING`
	var b pgx.Batch
	for _, v := range list {
		b.Queue(q,
			v.ID,
			v.FIO,
			v.BDatetime(),
			v.BPlace,
			v.OrgName,
			v.INN,
			v.PositionFL,
			v.NKOAP,
			v.GOrgName,
			v.SudFIO,
			v.SudPosition,
			v.DisqualificationDuration,
			v.DisStartDatetime(),
			v.DisEndDatetime())
	}
	return db.pool.SendBatch(ctx, &b).Close()
}

// BatchБухОтчетность - пакетная вставка в БД объекта models.БухОтчетность
func (db *DB) BatchБухОтчетность(ctx context.Context, list []models.БухОтчетность) (err error) {
	const q = `INSERT INTO accounting_statements.баланс
	(год, иннюл, balance_data)
	VALUES($1::int4, $2::text, $3::jsonb)
	ON CONFLICT (год, иннюл)
	DO UPDATE
	SET balance_data=EXCLUDED.balance_data`
	var b pgx.Batch
	for _, v := range list {
		b.Queue(q,
			v.ReportYear(),
			v.INN(),
			v)
	}
	return db.pool.SendBatch(ctx, &b).Close()
}

// InsertБухОтчетность - добавление в БД объекта models.БухОтчетность
func (db *DB) InsertБухОтчетность(ctx context.Context, v models.БухОтчетность) (err error) {
	const q = `INSERT INTO accounting_statements.баланс
	(год, иннюл, balance_data)
	VALUES($1::int4, $2::text, $3::jsonb)
	ON CONFLICT (год, иннюл)
	DO UPDATE
	SET balance_data=EXCLUDED.balance_data`
	_, err = db.pool.Exec(ctx, q,
		v.ReportYear(),
		v.INN(),
		v)
	return
}

// BatchОКТМО - добавление в БД ОКТМО
// TODO: прописать ON CONFLICT действие
func (db *DB) BatchОКТМО(ctx context.Context, list []models.OKT) (err error) {
	const q = `INSERT
    INTO
    public.октмо
	(
		ter
		, kod1
		, kod2
		, kod3
		, razdel
		, "name"
		, centrum
		, nomdescr
		, nomakt
		, status
		, dateutv
		, datevved
	)
	VALUES(
		$1
		, $2
		, $3
		, $4
		, $5
		, $6
		, $7
		, $8
		, $9
		, $10
		, to_date($11, 'DD.MM.YYYY')
		, to_date($12, 'DD.MM.YYYY')
	)`
	var b pgx.Batch
	for _, v := range list {
		b.Queue(q,
			v.Ter,
			v.Kod1,
			v.Kod2,
			v.Kod3,
			v.Razdel,
			v.Name,
			v.Centrum,
			v.NomDescr,
			v.NomAkt,
			v.Status,
			v.Dateutv,
			v.Datevved,
		)
	}
	return db.pool.SendBatch(ctx, &b).Close()
}

// BatchОКТАО - добавление в БД ОКТАО
// TODO: прописать ON CONFLICT действие
func (db *DB) BatchОКТАО(ctx context.Context, list []models.OKT) (err error) {
	const q = `INSERT
    INTO
    public.октао
	(
		ter
		, kod1
		, kod2
		, kod3
		, razdel
		, "name"
		, centrum
		, nomdescr
		, nomakt
		, status
		, dateutv
		, datevved
	)
	VALUES(
		$1
		, $2
		, $3
		, $4
		, $5
		, $6
		, $7
		, $8
		, $9
		, $10
		, to_date($11, 'DD.MM.YYYY')
		, to_date($12, 'DD.MM.YYYY')
	)`
	var b pgx.Batch
	for _, v := range list {
		b.Queue(q,
			v.Ter,
			v.Kod1,
			v.Kod2,
			v.Kod3,
			v.Razdel,
			v.Name,
			v.Centrum,
			v.NomDescr,
			v.NomAkt,
			v.Status,
			v.Dateutv,
			v.Datevved,
		)
	}
	return db.pool.SendBatch(ctx, &b).Close()
}

// InsertЕГРЮЛ - добавление в БД объектов ЕГРЮЛ
func (db *DB) InsertЕГРЮЛ(ctx context.Context, vs []egr.EGRUL) (err error) {
	const q = `
	WITH cte AS (
		SELECT $1::date AS датавып, $2 AS инн, $3 AS огрн, $4::jsonb AS egrul_data
	)
	INSERT INTO egr.егрюл
	(датавып, инн, огрн, egrul_data)
	SELECT * FROM cte AS src
	ON CONFLICT (инн,огрн)
	DO UPDATE
	SET датавып=EXCLUDED.датавып, egrul_data=EXCLUDED.egrul_data
	WHERE egr.егрюл.датавып < EXCLUDED.датавып`
	var b pgx.Batch
	for i := range vs {
		if vs[i].ИНН == nil {
			continue
		}
		b.Queue(
			q,
			vs[i].ДатаВып,
			*vs[i].ИНН,
			vs[i].ОГРН,
			vs[i],
		)
	}
	if b.Len() == 0 {
		return nil
	}
	return db.pool.SendBatch(ctx, &b).Close()
}

// InsertЕГРИП - добавление в БД объектов ЕГРИП
func (db *DB) InsertЕГРИП(ctx context.Context, vs []egr.EGRIP) (err error) {
	const q = `
	WITH cte AS (
		SELECT $1::date AS датавып, $2 AS инн, $3 AS огрн, $4::jsonb AS egrip_data
	)
	INSERT INTO egr.егрип
	(датавып, инн, огрн, egrip_data)
	SELECT * FROM cte AS src
	ON CONFLICT (инн,огрн)
	DO UPDATE
	SET датавып=EXCLUDED.датавып, egrip_data=EXCLUDED.egrip_data
	WHERE egr.егрип.датавып < EXCLUDED.датавып`
	var b pgx.Batch
	for i := range vs {
		if vs[i].ИНН == "" {
			continue
		}
		b.Queue(
			q,
			vs[i].ДатаВып,
			vs[i].ИНН,
			vs[i].ОГРН,
			vs[i],
		)
	}
	if b.Len() == 0 {
		return nil
	}
	return db.pool.SendBatch(ctx, &b).Close()
}

// InsertПривлеченииУчастникаЗакупкиКАдминистративнойОтветственности - добавление в БД объектов ПривлеченииУчастникаЗакупкиКАдминистративнойОтветственности
func (db *DB) InsertПривлеченииУчастникаЗакупкиКАдминистративнойОтветственности(ctx context.Context, v models.ПривлеченииУчастникаЗакупкиКАдминистративнойОтветственности) (err error) {
	const q = `INSERT INTO public.информация1928коап
	(инн, огрн, кпп, наименованиеюл, типучастника, суд, номердела, датавынесенияпостановления, датавступлениявзаконнуюсилу)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT DO NOTHING`
	_, err = db.pool.Exec(ctx, q,
		v.ИНН,
		v.ОГРН,
		v.КПП,
		v.НаименованиеЮЛ,
		v.ТипУчастника,
		v.Суд,
		v.НомерДела,
		v.DateДатаВынесенияПостановления(),
		v.DateДатаВступленияВЗаконнуюСилу(),
	)
	return
}

// InsertHotelData - добавление в БД
func (db *DB) InsertHotelData(ctx context.Context, h models.HotelData) error {
	var hotelNumber = h.GetNumber()
	if hotelNumber == "" {
		return nil
	}
	const qh = `INSERT
	INTO
	public.hotels
	(
		federal_number
		, type
		, full_name
		, short_name
		, region
		, inn
		, ogrn
		, address
		, phone
		, fax
		, email
		, site
		, owner
	)
	VALUES(
		$1::text,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10,
		$11,
		$12,
		$13
	)
	ON CONFLICT (federal_number)
	DO UPDATE
	SET
		"type" = EXCLUDED."type"
		, full_name = EXCLUDED.full_name
		, short_name = EXCLUDED.short_name
		, region = EXCLUDED.region
		, inn = EXCLUDED.inn
		, ogrn = EXCLUDED.ogrn
		, address = EXCLUDED.address
		, phone = EXCLUDED.phone
		, fax = EXCLUDED.fax
		, email = EXCLUDED.email
		, site = EXCLUDED.site
		, "owner" = EXCLUDED."owner"`
	if _, err := db.pool.Exec(ctx, qh,
		hotelNumber,
		h.GetType(),
		h.GetFullName(),
		h.GetShortName(),
		h.GetRegion(),
		h.GetINN(),
		h.GetOGRN(),
		h.GetAddress(),
		h.GetPhone(),
		h.GetFax(),
		h.GetEmail(),
		h.GetSite(),
		h.GetOwner()); err != nil {
		if strings.Contains(err.Error(), `duplicate key value violates unique constraint`) {
			return nil
		}
		return err
	}

	const qr = `
	INSERT
	INTO
	public.hotels_rooms
	(
		federal_number
		, category
		, rooms
		, seats
	)
	VALUES(
		$1::text
		, $2::text
		, $3::int4
		, $4::int4
	)`
	for _, room := range h.GetRooms() {
		var category = room.RoomCategory
		if len(category) == 0 {
			category = "Без категории"
		}
		if _, err := db.pool.Exec(ctx, qr,
			hotelNumber,
			category,
			room.NumberRooms,
			room.NumberSeats); err != nil {
			return err
		}
	}

	const qc = `
	INSERT
		INTO
		public.hotels_classification
	(
		federal_number
		, date_issued
		, date_end
		, category
		, license_number
		, registration_number
	)
	VALUES(
		$1::text
		, to_timestamp($2::text, 'DD.MM.YYYY')::date
		, to_timestamp($3::text, 'DD.MM.YYYY')::date
		, $4::text
		, $5::text
		, $6::text
	)`

	if _, err := db.pool.Exec(ctx, qc,
		hotelNumber,
		h.GetLicenseDateIssued(),
		h.GetLicenseDateEnd(),
		h.GetCategoryStars(),
		h.GetLicenseNumber(),
		h.GetRegistrationNumber()); err != nil {
		return err
	}

	return nil
}

// InsertUnscheduledInspections - добавление данных по внеплановым проверкам ФГИС
func (db *DB) InsertUnscheduledInspections(ctx context.Context, data []models.InspectionFGIS) (err error) {
	const q = `
	INSERT
		INTO
		fgis.unscheduled_inspections
		(
			erpid
			, inn
			, ogrn
			, start_date
			, status
			, fz_name
			, value
		)
	VALUES($1,$2,$3,$4::date,$5,$6,$7)
	ON CONFLICT (erpid)
	DO UPDATE
	SET
		inn = EXCLUDED.inn
		, ogrn = EXCLUDED.ogrn
		, start_date = EXCLUDED.start_date
		, status = EXCLUDED.status
		, fz_name = EXCLUDED.fz_name
		, value = EXCLUDED.value`
	var b pgx.Batch
	for _, v := range data {
		b.Queue(q,
			v.ERPID,
			v.INN(),
			v.OGRN(),
			v.START_DATE,
			v.STATUS,
			v.FZ_NAME,
			v)
	}
	return db.pool.SendBatch(ctx, &b).Close()
}

// InsertScheduledInspections - добавление данных по плановым проверкам ФГИС
func (db *DB) InsertScheduledInspections(ctx context.Context, data []models.InspectionFGIS) (err error) {
	const q = `
	INSERT
		INTO
		fgis.scheduled_inspections
		(
			erpid
			, inn
			, ogrn
			, start_date
			, status
			, fz_name
			, value
		)
	VALUES($1,$2,$3,$4::date,$5,$6,$7)
	ON CONFLICT (erpid)
	DO UPDATE
	SET
		inn = EXCLUDED.inn
		, ogrn = EXCLUDED.ogrn
		, start_date = EXCLUDED.start_date
		, status = EXCLUDED.status
		, fz_name = EXCLUDED.fz_name
		, value = EXCLUDED.value`
	var b pgx.Batch
	for _, v := range data {
		b.Queue(q,
			v.ERPID,
			v.INN(),
			v.OGRN(),
			v.START_DATE,
			v.STATUS,
			v.FZ_NAME,
			v)
	}
	return db.pool.SendBatch(ctx, &b).Close()
}

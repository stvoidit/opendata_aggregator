// Package models - структуры объектов для парсинга из открытых данных
package models

import (
	"encoding/json"
	"strings"
	"time"
)

const dateLayout = `02.01.2006`

func transformDate(s string) time.Time {
	t, _ := time.Parse(dateLayout, s)
	return t
}

// СведСНР - ...
type СведСНР struct {
	ПризнЕСХН bool `xml:"ПризнЕСХН,attr" jsonschema_description:"Признак использования единого сельскохозяйственного налога"`                                 // Признак использования единого сельскохозяйственного налога
	ПризнУСН  bool `xml:"ПризнУСН,attr" jsonschema_description:"Признак использования упрощенной системы налогооблажения"`                                    // Признак использования упрощенной системы налогооблажения
	ПризнЕНВД bool `xml:"ПризнЕНВД,attr" jsonschema_description:"Признак использования единого налога на вмененый доход"`                                     // Признак использования единого налога на вмененый доход
	ПризнСРП  bool `xml:"ПризнСРП,attr" jsonschema_description:"Признак использования системы налогооблажения при выполеннии соглашений о разделе продукции"` // Признак использования системы налогооблажения при выполеннии соглашений о разделе продукции
}

// НалоговыйРежимНалогоплательщика - структура объекта из https://www.nalog.gov.ru/opendata/7707329152-snr/
type НалоговыйРежимНалогоплательщика struct {
	ИдДок   string  `xml:"ИдДок,attr" jsonschema_description:"ID документа"`               // ID документа
	ДатаДок string  `xml:"ДатаДок,attr" jsonschema_description:"Дата актуальности данных"` // Дата актуальности данных
	СведНП  СведНП  `xml:"СведНП" jsonschema_description:"Сведения о налогоплательщике"`
	СведСНР СведСНР `xml:"СведСНР" jsonschema_description:"Сведения о налоговых режимах"`
}

// DocDate - строка даты в формат time.Time
func (nrn НалоговыйРежимНалогоплательщика) DocDate() time.Time {
	return transformDate(nrn.ДатаДок)
}

// Росаккредитация - структура данных от https://fsa.gov.ru/opendata/7736638268-rss/
type Росаккредитация struct {
	IDcert                         string `json:"id_cert" jsonschema_description:"id"`                                                                        // id
	CertStatus                     string `json:"cert_status" jsonschema_description:"Статус сертификата"`                                                    // Статус сертификата
	CertType                       string `json:"cert_type" jsonschema_description:"Тип сертификата"`                                                         // Тип сертификата
	RegNumber                      string `json:"reg_number" jsonschema_description:"Регистрационный номер сертификата"`                                      // Регистрационный номер сертификата
	DateBegining                   string `json:"date_begining" jsonschema_description:"Дата регистрации сертификата"`                                        // Дата регистрации сертификата
	DateFinish                     string `json:"date_finish,omitempty" jsonschema_description:"Дата окончания действия сертификата"`                         // Дата окончания действия сертификата
	ProductScheme                  string `json:"product_scheme" jsonschema_description:"Схема сертификации"`                                                 // Схема сертификации
	ProductObjectTypeCert          string `json:"product_object_type_cert" jsonschema_description:"Тип объекта сертификации"`                                 // Тип объекта сертификации
	ProductType                    string `json:"product_type" jsonschema_description:"Вид продукции: Импортная/ Отечественная"`                              // Вид продукции: Импортная/ Отечественная
	ProductOKPD2                   string `json:"product_okpd2" jsonschema_description:"Коды ОКПД2"`                                                          // Коды ОКПД2
	ProductTnVed                   string `json:"product_tn_ved" jsonschema_description:"Коды ТН ВЭД ЕАЭС"`                                                   // Коды ТН ВЭД ЕАЭС
	ProductTechReg                 string `json:"product_tech_reg" jsonschema_description:"Технические регламенты"`                                           // Технические регламенты
	ProductGroup                   string `json:"product_group" jsonschema_description:"Группа продукции"`                                                    // Группа продукции
	ProductName                    string `json:"product_name" jsonschema_description:"Общее наименование продукции"`                                         // Общее наименование продукции
	ProductInfo                    string `json:"product_info" jsonschema_description:"Информация по продукции"`                                              // Информация по продукции
	ApplicantType                  string `json:"applicant_type" jsonschema_description:"Тип заявителя"`                                                      // Тип заявителя
	PersonApplicantType            string `json:"person_applicant_type" jsonschema_description:"Вид заявителя"`                                               // Вид заявителя
	ApplicantOGRN                  string `json:"applicant_ogrn" jsonschema_description:"ОГРН/ОГРНИП заявителя"`                                              // ОГРН/ОГРНИП заявителя
	ApplicantINN                   string `json:"applicant_inn" jsonschema_description:"ИНН заявителя"`                                                       // ИНН заявителя
	ApplicantPhone                 string `json:"applicant_phone" jsonschema_description:"Телефон заявителя"`                                                 // Телефон заявителя
	ApplicantFax                   string `json:"applicant_fax" jsonschema_description:"Факс заявителя"`                                                      // Факс заявителя
	ApplicantEmail                 string `json:"applicant_email" jsonschema_description:"EMAIL заявителя"`                                                   // EMAIL заявителя
	ApplicantWebsite               string `json:"applicant_website" jsonschema_description:"Сайт заявителя"`                                                  // Сайт заявителя
	ApplicantName                  string `json:"applicant_name" jsonschema_description:"Полное наименование заявителя"`                                      // Полное наименование заявителя
	ApplicantDirectorName          string `json:"applicant_director_name" jsonschema_description:"ФИО руководителя заявителя"`                                // ФИО руководителя заявителя
	ApplicantAddress               string `json:"applicant_address" jsonschema_description:"Адрес места нахождения заявителя"`                                // Адрес места нахождения заявителя
	ApplicantAddressActual         string `json:"applicant_address_actual" jsonschema_description:"Адрес места осуществления деятельности заявления"`         // Адрес места осуществления деятельности заявления
	ManufacturerType               string `json:"manufacturer_type" jsonschema_description:"Тип изготовителя"`                                                // Тип изготовителя
	ManufacturerOGRN               string `json:"manufacturer_ogrn" jsonschema_description:"ОГРН/ОГРНИП изготовителя"`                                        // ОГРН/ОГРНИП изготовителя
	ManufacturerINN                string `json:"manufacturer_inn" jsonschema_description:"ИНН изготовителя"`                                                 // ИНН изготовителя
	ManufacturerPhone              string `json:"manufacturer_phone" jsonschema_description:"Телефон изготовителя"`                                           // Телефон изготовителя
	ManufacturerFax                string `json:"manufacturer_fax" jsonschema_description:"Факс изготовителя"`                                                // Факс изготовителя
	ManufacturerEmail              string `json:"manufacturer_email" jsonschema_description:"EMAIL изготовителя"`                                             // EMAIL изготовителя
	ManufacturerWebsite            string `json:"manufacturer_website" jsonschema_description:"Сайт изготовителя"`                                            // Сайт изготовителя
	ManufacturerName               string `json:"manufacturer_name" jsonschema_description:"Полное наименование изготовителя"`                                // Полное наименование изготовителя
	ManufacturerDirectorName       string `json:"manufacturer_director_name" jsonschema_description:"ФИО руководителя изготовителя"`                          // ФИО руководителя изготовителя
	ManufacturerCountry            string `json:"manufacturer_country" jsonschema_description:"Страна места нахождения изготовителя"`                         // Страна места нахождения изготовителя
	ManufacturerAddress            string `json:"manufacturer_address" jsonschema_description:"Адрес места нахождения изготовителя"`                          // Адрес места нахождения изготовителя
	ManufacturerAddressActual      string `json:"manufacturer_address_actual" jsonschema_description:"Адрес места осуществления деятельности изготовителя"`   // Адрес места осуществления деятельности изготовителя
	ManufacturerAddressFilial      string `json:"manufacturer_address_filial" jsonschema_description:"Филиалы изготовителя"`                                  // Филиалы изготовителя
	OrganToCertificationName       string `json:"organ_to_certification_name" jsonschema_description:"Полное наименование ОС"`                                // Полное наименование ОС
	OrganToCertificationRegNumber  string `json:"organ_to_certification_reg_number" jsonschema_description:"Регистрационный номер аттестата аккредитации ОС"` // Регистрационный номер аттестата аккредитации ОС
	OrganToCertificationHeadName   string `json:"organ_to_certification_head_name" jsonschema_description:"ФИО руководителя ОС"`                              // ФИО руководителя ОС
	BasisForCertificate            string `json:"basis_for_certificate" jsonschema_description:"Основание выдачи СС"`                                         // Основание выдачи СС
	OldBasisForCertificate         string `json:"old_basis_for_certificate" jsonschema_description:"Основание выдачи СС (ФГИС 1.0)"`                          // Основание выдачи СС (ФГИС 1.0)
	DioExpert                      string `json:"fio_expert" jsonschema_description:"Эксперты"`                                                               // Эксперты
	DioSignatory                   string `json:"fio_signatory" jsonschema_description:"Лицо, подписавшее сертификат"`                                        // Лицо, подписавшее сертификат
	ProductNationalStandart        string `json:"product_national_standart" jsonschema_description:"Стандарты"`                                               // Стандарты
	ProductionAnalysisForAct       string `json:"production_analysis_for_act" jsonschema_description:"Акт анализа производства"`                              // Акт анализа производства
	ProductionAnalysisForActNumber string `json:"production_analysis_for_act_number" jsonschema_description:"Номер акта анализа производства"`                // Номер акта анализа производства
	ProductionAnalysisForActDate   string `json:"production_analysis_for_act_date,omitempty" jsonschema_description:"Дата акта анализа производства"`         // Дата акта анализа производства
}

// RDS - Сведения из Реестра деклараций о соответствии
// https://fsa.gov.ru/opendata/7736638268-rds/
type RDS struct {
	IDDecl                        string  `json:"id_decl" jsonschema_description:"id"`
	RegNumber                     string  `json:"reg_number" jsonschema_description:"Рег. Номер"`
	DeclStatus                    string  `json:"decl_status" jsonschema_description:"Статус"`
	DeclType                      string  `json:"decl_type" jsonschema_description:"Тип декларации"`
	DateBeginning                 *string `json:"date_beginning" jsonschema_description:"Дата начала действия"`
	DateFinish                    string  `json:"date_finish,omitempty" jsonschema_description:"Дата окончания действия"`
	DeclarationScheme             string  `json:"declaration_scheme" jsonschema_description:"Схема декларирования"`
	ProductObjectTypeDecl         string  `json:"product_object_type_decl" jsonschema_description:"Тип объекта декларирования"`
	ProductType                   string  `json:"product_type" jsonschema_description:"Вид продукции"`
	ProductGroup                  string  `json:"product_group" jsonschema_description:"Группа продукции"`
	ProductName                   string  `json:"product_name" jsonschema_description:"Общее наименование продукции"`
	AsproductInfo                 string  `json:"asproduct_info" jsonschema_description:"Информация по продукции"`
	ProductTechReg                string  `json:"product_tech_reg" jsonschema_description:"Технический регламент"`
	OrganToCertificationName      string  `json:"organ_to_certification_name" jsonschema_description:"Наименование ОС"`
	OrganToCertificationRegNumber string  `json:"organ_to_certification_reg_number" jsonschema_description:"Номер аттестата аккредитации ОС"`
	BasisForDecl                  string  `json:"basis_for_decl" jsonschema_description:"Основание выдачи ДС"`
	OldBasisForDecl               string  `json:"old_basis_for_decl" jsonschema_description:"Основание выдачи ДС (ФГИС 1.0)"`
	ApplicantType                 string  `json:"applicant_type" jsonschema_description:"Тип заявителя"`
	PersonApplicantType           string  `json:"person_applicant_type" jsonschema_description:"Вид заявителя"`
	ApplicantOGRN                 string  `json:"applicant_ogrn" jsonschema_description:"ОГРН/ОГРНИП заявителя"`
	ApplicantINN                  string  `json:"applicant_inn" jsonschema_description:"ИНН заявителя"`
	ApplicantName                 string  `json:"applicant_name" jsonschema_description:"Наименование Заявителя"`
	ManufacturerType              string  `json:"manufacturer_type" jsonschema_description:"Вид изготовителя"`
	ManufacturerOGRN              string  `json:"manufacturer_ogrn" jsonschema_description:"ОГРН/ОГРНИП изготовителя"`
	ManufacturerINN               string  `json:"manufacturer_inn" jsonschema_description:"ИНН изготовителя"`
	ManufacturerName              string  `json:"manufacturer_name" jsonschema_description:"Полное наименование изготовителя"`
}

// СведНедоим - ...
type СведНедоим struct {
	НаимНалог    string  `xml:"НаимНалог,attr" jsonschema_description:"Наименование налога (сбора, страховых взносов), денежного взыскания"` // Наименование налога (сбора, страховых взносов), денежного взыскания
	СумНедНалог  float64 `xml:"СумНедНалог,attr" jsonschema_description:"Сумма недоимки по налогу"`                                          // Сумма нелоимки по налогу
	СумПени      float64 `xml:"СумПени,attr" jsonschema_description:"Сумма пени"`                                                            // Сумма пени
	СумШтраф     float64 `xml:"СумШтраф,attr" jsonschema_description:"Сумма штрафа"`                                                         // Сумма штрафа
	ОбщСумНедоим float64 `xml:"ОбщСумНедоим,attr" jsonschema_description:"Общая сумма недоимки по налогу, пени и штрафу"`                    // Общая сумма недоимки по налогу, пени и штрафу
}

// СведенияОСуммахНедоимки - Сведения о суммах недоимки и задолженности по пеням и штрафам
// https://www.nalog.gov.ru/opendata/7707329152-debtam/
type СведенияОСуммахНедоимки struct {
	ИдДок      string       `xml:"ИдДок,attr" jsonschema_description:"ID документа"` // ID документа
	ДатаДок    string       `xml:"ДатаДок,attr" jsonschema_description:"Дата актуальности данных"`
	ДатаСост   string       `xml:"ДатаСост,attr" jsonschema_description:"Дата составления данных"`
	СведНП     СведНП       `xml:"СведНП" jsonschema_description:"Сведения о налогоплательщике"`
	СведНедоим []СведНедоим `xml:"СведНедоим" jsonschema_description:"Сведения о суммах недоимки"`
}

// DocDate - строка даты в формат time.Time
func (ssn СведенияОСуммахНедоимки) DocDate() time.Time {
	return transformDate(ssn.ДатаДок)
}

// DocDate - строка даты в формат time.Time
func (ssn СведенияОСуммахНедоимки) DateBuild() time.Time {
	return transformDate(ssn.ДатаСост)
}

// СведКГН - Признак участника (ответственного участника) консолидированной группы налогоплательщиков
type СведКГН struct {
	// Принимает значение:
	// 1 – ответственный участник консолидированной группы налогоплательщиков
	// 2 – участник консолидированной группы налогоплательщиков
	ПризнУчКГН string `xml:"ПризнУчКГН,attr"` // Признак участника (ответственного участника) консолидированной группы налогоплательщиков
}

// СведНП - ...
type СведНП struct {
	НаимОрг string `xml:"НаимОрг,attr" jsonschema_description:"Наименование организации"` // Наименование организации
	ИННЮЛ   string `xml:"ИННЮЛ,attr" jsonschema_description:"ИНН организации"`            // ИНН организации
}

// СведенияОбУчастииВКонсГруппе - Сведения об участии в консолидированной группе налогоплательщиков
// https://www.nalog.gov.ru/opendata/7707329152-kgn/
type СведенияОбУчастииВКонсГруппе struct {
	// ИдДок    string `xml:"ИдДок,attr"`
	ДатаДок string  `xml:"ДатаДок,attr" jsonschema_description:"Дата актуальности данных"` // Дата актуальности данных
	СведНП  СведНП  `xml:"СведНП" jsonschema_description:"Сведения о налогоплательщике"`
	СведКГН СведКГН `xml:"СведКГН" jsonschema_description:"Сведения об участии в консолидированной группе"`
}

// DocDate - строка даты в формат time.Time
func (skg СведенияОбУчастииВКонсГруппе) DocDate() time.Time {
	return transformDate(skg.ДатаДок)
}

// СведенияОСреднесписочнойЧисленностиРаботников - Сведения о среднесписочной численности работников организации
// https://www.nalog.gov.ru/opendata/7707329152-sshr2019/
type СведенияОСреднесписочнойЧисленностиРаботников struct {
	// ИдДок    string   `xml:"ИдДок,attr"`
	ДатаДок  string `xml:"ДатаДок,attr"` // Дата актуальности данных
	СведНП   СведНП `xml:"СведНП"`
	СведССЧР struct {
		КолРаб string `xml:"КолРаб,attr"` // кол-во работников
	} `xml:"СведССЧР"`
}

// DocDate - строка даты в формат time.Time
func (scr СведенияОСреднесписочнойЧисленностиРаботников) DocDate() time.Time {
	return transformDate(scr.ДатаДок)
}

// СвУплСумНал - ...
type СвУплСумНал struct {
	НаимНалог string  `xml:"НаимНалог,attr" jsonschema_description:"Наименование налога (сбора, страхового взноса)"`      // Наименование налога (сбора, страхового взноса)
	СумУплНал float64 `xml:"СумУплНал,attr" jsonschema_description:"Сумма уплаченного налога (сбора, страхового взноса)"` // Сумма уплаченного налога (сбора, страхового взноса)
}

// СведенияОбУплаченныхОрганизациейНалогов - Сведения об уплаченных организацией в календарном году налогов и сборов
// https://www.nalog.gov.ru/opendata/7707329152-paytax/
type СведенияОбУплаченныхОрганизациейНалогов struct {
	// ИдДок    string `xml:"ИдДок,attr"`
	ДатаДок     string        `xml:"ДатаДок,attr" jsonschema_description:"Дата актуальности данных"` // Дата актуальности данных
	СведНП      СведНП        `xml:"СведНП" jsonschema_description:"Сведения о налогоплательщике"`
	СвУплСумНал []СвУплСумНал `xml:"СвУплСумНал" jsonschema_description:"Сведения об уплаченных налогах"`
}

// DocDate - строка даты в формат time.Time
func (son СведенияОбУплаченныхОрганизациейНалогов) DocDate() time.Time {
	return transformDate(son.ДатаДок)
}

// ИсполнительныеПроизводстваВОтношенииЮридическихЛиц - Набор открытых данных, содержащий общедоступные сведения,
// необходимые для осуществления задач по принудительному исполнению судебных актов,
// актов других органов и должностных лиц (в отношении юридических лиц).
// https://opendata.fssp.gov.ru/7709576929-iplegallist
// TODO: имена полей json можно сделать нормальные - согласовать с фронтом
type ИсполнительныеПроизводстваВОтношенииЮридическихЛиц struct {
	NameOfDebtor                                        string  `json:"Наименование юридического лица" jsonschema_description:"Наименование юридического лица"`                                                                     // Наименование юридического лица
	AddressOfDebtorOrganization                         string  `json:"Адрес организации - должника" jsonschema_description:"Адрес организации - должника"`                                                                         // Адрес организации - должника
	ActualAddressOfDebtorOrganization                   string  `json:"Фактический адрес организации должника" jsonschema_description:"Фактический адрес организации должника"`                                                     // Фактический адрес организации должника
	NumberOfEnforcementProceeding                       string  `json:"Номер исполнительного производства" jsonschema_description:"Номер исполнительного производства"`                                                             // Номер исполнительного производства
	DateOfInstitutionProceeding                         string  `json:"Дата возбуждения исполнительного производства" jsonschema_description:"Дата возбуждения исполнительного производства"`                                       // Дата возбуждения исполнительного производства
	TotalNumberOfEnforcementProceedings                 string  `json:"Номер сводного производства по взыскателю или должнику" jsonschema_description:"Номер сводного производства по взыскателю или должнику"`                     // Номер сводного производства по взыскателю или должнику
	ExecutiveDocumentType                               string  `json:"Тип исполнительного документа" jsonschema_description:"Тип исполнительного документа"`                                                                       // Тип исполнительного документа
	DateOfExecutiveDocument                             string  `json:"Дата исполнительного документа" jsonschema_description:"Дата исполнительного документа"`                                                                     // Дата исполнительного документа
	NumberOfExecutiveDocument                           string  `json:"Номер исполнительного документа" jsonschema_description:"Номер исполнительного документа"`                                                                   // Номер исполнительного документа
	ObjectOfExecutiveDocuments                          string  `json:"Требования исполнительного документа" jsonschema_description:"Требования исполнительного документа"`                                                         // Требования исполнительного документа
	ObjectOfExecution                                   string  `json:"Предмет исполнения" jsonschema_description:"Предмет исполнения"`                                                                                             // Предмет исполнения
	AmountDue                                           float64 `json:"Сумма долга" jsonschema_description:"Сумма долга"`                                                                                                           // Сумма долга
	DebtRemainingBalance                                float64 `json:"Остаток непогашенной задолженности" jsonschema_description:"Остаток непогашенной задолженности"`                                                             // Остаток непогашенной задолженности
	DepartmentsOfBailiffs                               string  `json:"Наименование отдела" jsonschema_description:"Наименование отдела"`                                                                                           // Наименование отдела
	AddressOfDepartmentsOfBailiff                       string  `json:"Адрес отдела судебных приставов" jsonschema_description:"Адрес отдела судебных приставов"`                                                                   // Адрес отдела судебных приставов
	DebtorTaxpayerIdentificationNumber                  string  `json:"Идентификационный номер налогоплательщика должника" jsonschema_description:"Идентификационный номер налогоплательщика должника"`                             // Идентификационный номер налогоплательщика должника
	TaxpayerIdentificationNumberOfOrganizationCollector string  `json:"Идентификационный номер налогоплательщика взыскателя-организации" jsonschema_description:"Идентификационный номер налогоплательщика взыскателя-организации"` // Идентификационный номер налогоплательщика взыскателя-организации
	Repaid                                              bool    `json:"Погашено" jsonschema_description:"Погашено"`
}

// IpLegalListComplete - Исполнительные производства в отношении юридических лиц,
// оконченные в соответствии с пунктами 3 и 4 части 1 статьи 46 и пунктами 6 и 7 части 1 статьи 47
// Федерального закона от 2 октября 2007 г. № 229-ФЗ «Об исполнительном производстве»
// https://opendata.fssp.gov.ru/7709576929-iplegallistcomplete
// TODO: имена полей json можно сделать нормальные - согласовать с фронтом
type IpLegalListComplete struct {
	NameOfDebtor                                        string `json:"Наименование юридического лица" jsonschema_description:"Наименование юридического лица"`                                                                     // Наименование юридического лица
	AddressOfDebtorOrganization                         string `json:"Адрес организации - должника" jsonschema_description:"Адрес организации - должника"`                                                                         // Адрес организации - должника
	ActualAddressOfDebtorOrganization                   string `json:"Фактический адрес организации должника" jsonschema_description:"Фактический адрес организации должника"`                                                     // Фактический адрес организации должника
	NumberOfEnforcementProceeding                       string `json:"Номер исполнительного производства" jsonschema_description:"Номер исполнительного производства"`                                                             // Номер исполнительного производства
	DateOfInstitutionProceeding                         string `json:"Дата возбуждения исполнительного производства" jsonschema_description:"Дата возбуждения исполнительного производства"`                                       // Дата возбуждения исполнительного производства
	TotalNumberOfEnforcementProceedings                 string `json:"Номер сводного производства по взыскателю или должнику" jsonschema_description:"Номер сводного производства по взыскателю или должнику"`                     // Номер сводного производства по взыскателю или должнику
	ExecutiveDocumentType                               string `json:"Тип исполнительного документа" jsonschema_description:"Тип исполнительного документа"`                                                                       // Тип исполнительного документа
	DateOfExecutiveDocument                             string `json:"Дата исполнительного документа" jsonschema_description:"Дата исполнительного документа"`                                                                     // Дата исполнительного документа
	NumberOfExecutiveDocument                           string `json:"Номер исполнительного документа" jsonschema_description:"Номер исполнительного документа"`                                                                   // Номер исполнительного документа
	ObjectOfExecutiveDocuments                          string `json:"Требования исполнительного документа" jsonschema_description:"Требования исполнительного документа"`                                                         // Требования исполнительного документа
	ObjectOfExecution                                   string `json:"Предмет исполнения" jsonschema_description:"Предмет исполнения"`                                                                                             // Предмет исполнения
	DateCompleteIPreason                                string `json:"Дата причина окончания или прекращения ИП" jsonschema_description:"Дата, причина окончания или прекращения ИП (статья, часть, пункт основания)"`             // Дата, причина окончания или прекращения ИП (статья, часть, пункт основания)
	DepartmentsOfBailiffs                               string `json:"Наименование отдела" jsonschema_description:"Наименование отдела"`                                                                                           // Наименование отдела
	AddressOfDepartmentsOfBailiff                       string `json:"Адрес отдела судебных приставов" jsonschema_description:"Адрес отдела судебных приставов"`                                                                   // Адрес отдела судебных приставов
	DebtorTaxpayerIdentificationNumber                  string `json:"Идентификационный номер налогоплательщика должника" jsonschema_description:"Идентификационный номер налогоплательщика должника"`                             // Идентификационный номер налогоплательщика должника
	TaxpayerIdentificationNumberOfOrganizationCollector string `json:"Идентификационный номер налогоплательщика взыскателя-организации" jsonschema_description:"Идентификационный номер налогоплательщика взыскателя-организации"` // Идентификационный номер налогоплательщика взыскателя-организации
}

// ОткрытыйРеестрТоварныхЗнаков - Открытый реестр товарных знаков и знаков обслуживания Российской Федерации
// https://rospatent.gov.ru/opendata/7730176088-tz
type ОткрытыйРеестрТоварныхЗнаков struct {
	RegistrationNumber                                             string `json:"registration number" jsonschema_description:"Номер государственной регистрации"`                                                                                                                                                                                                                                                                                                                     // Номер государственной регистрации
	RegistrationDate                                               string `json:"registration date" jsonschema_description:"Дата государственной регистрации"`                                                                                                                                                                                                                                                                                                                        // Дата государственной регистрации
	ApplicationNumber                                              string `json:"application number" jsonschema_description:"Номер заявки на государственную регистрацию"`                                                                                                                                                                                                                                                                                                            // Номер заявки на государственную регистрацию
	ApplicationDate                                                string `json:"application date" jsonschema_description:"Дата подачи заявки на государственную регистрацию"`                                                                                                                                                                                                                                                                                                        // Дата подачи заявки на государственную регистрацию
	PriorityDate                                                   string `json:"priority date" jsonschema_description:"Дата приоритета"`                                                                                                                                                                                                                                                                                                                                             // Дата приоритета
	ExhibitionPriorityDate                                         string `json:"exhibition priority date" jsonschema_description:"Дата начала открытого показа экспоната на выставке"`                                                                                                                                                                                                                                                                                               // Дата начала открытого показа экспоната на выставке
	ParisConventionPriorityNumber                                  string `json:"paris convention priority number" jsonschema_description:"Дата подачи первой заявки в государстве - участнике Парижской конвенции по охране промышленной собственности"`                                                                                                                                                                                                                             // Дата подачи первой заявки в государстве - участнике Парижской конвенции по охране промышленной собственности
	ParisConventionPriorityDate                                    string `json:"paris convention priority date" jsonschema_description:"Номер первой заявки в государстве - участнике Парижской конвенции по охране промышленной собственности"`                                                                                                                                                                                                                                     // Номер первой заявки в государстве - участнике Парижской конвенции по охране промышленной собственности
	ParisConventionPriorityCountryCode                             string `json:"paris convention priority country code" jsonschema_description:"Код страны подачи первой заявки в государстве - участнике Парижской конвенции по охране промышленной собственности"`                                                                                                                                                                                                                 // Код страны подачи первой заявки в государстве - участнике Парижской конвенции по охране промышленной собственности
	InitialApplicationNumber                                       string `json:"initial application number" jsonschema_description:"Номер первоначальной заявки, из которой выделена заявка, по которой произведена государственная регистрация"`                                                                                                                                                                                                                                    // Номер первоначальной заявки, из которой выделена заявка, по которой произведена государственная регистрация
	InitialApplicationOriorityDate                                 string `json:"initial application priority date" jsonschema_description:"по которой произведена государственная регистрация"`                                                                                                                                                                                                                                                                                      // по которой произведена государственная регистрация
	InitialRegistrationNumber                                      string `json:"initial registration number" jsonschema_description:"Номер первоначальной регистрации, из которой выделена отдельная регистрация"`                                                                                                                                                                                                                                                                   // Номер первоначальной регистрации, из которой выделена отдельная регистрация
	InitialRegistrationDate                                        string `json:"initial registration date" jsonschema_description:"Дата первоначальной регистрации, из которой выделена отдельная регистрация"`                                                                                                                                                                                                                                                                      // Дата первоначальной регистрации, из которой выделена отдельная регистрация
	InternationalRegistrationNumber                                string `json:"international registration number" jsonschema_description:"Номер международной регистрации, преобразованной в национальную заявку на регистрацию"`                                                                                                                                                                                                                                                   // Номер международной регистрации, преобразованной в национальную заявку на регистрацию
	InternationalRegistrationDate                                  string `json:"international registration date" jsonschema_description:"Дата международной регистрации, преобразованной в национальную заявку на регистрацию"`                                                                                                                                                                                                                                                      // Дата международной регистрации, преобразованной в национальную заявку на регистрацию
	InternationalRegistrationPriorityDate                          string `json:"international registration priority date" jsonschema_description:"Дата приоритета международной регистрации, преобразованной в национальную заявку на регистрацию"`                                                                                                                                                                                                                                  // Дата приоритета международной регистрации, преобразованной в национальную заявку на регистрацию
	InternationalRegistrationEntryDate                             string `json:"international registration entry date" jsonschema_description:"Дата внесения в международный реестр записи о территориальном расширении по международной регистрации"`                                                                                                                                                                                                                               // Дата внесения в международный реестр записи о территориальном расширении по международной регистрации
	ApplicationNumberForRecognitionOfTrademarkFromCrimea           string `json:"application number for recognition of trademark from Crimea" jsonschema_description:"Номер заявления о признании действия исключительного права на товарный знак на территории Российской Федерации, удостоверенное официальными документами Украины, действовавшими на день принятия в Российскую Федерацию Республики Крым и образования в составе Российской Федерации новых субъектов"`          // Номер заявления о признании действия исключительного права на товарный знак на территории Российской Федерации, удостоверенное официальными документами Украины, действовавшими на день принятия в Российскую Федерацию Республики Крым и образования в составе Российской Федерации новых субъектов
	ApplicationDateForRecognitionOfTrademarkFromCrimea             string `json:"application date for recognition of trademark from Crimea" jsonschema_description:"Дата поступления заявления о признании действия исключительного права на товарный знак на территории Российской Федерации, удостоверенное официальными документами Украины, действовавшими на день принятия в Российскую Федерацию Республики Крым и образования в составе Российской Федерации новых субъектов"` // Дата поступления заявления о признании действия исключительного права на товарный знак на территории Российской Федерации, удостоверенное официальными документами Украины, действовавшими на день принятия в Российскую Федерацию Республики Крым и образования в составе Российской Федерации новых субъектов
	CrimeanTrademarkApplicationNumberForStateRegistrationInUkraine string `json:"Crimean trademark application number for state registration in Ukraine" jsonschema_description:"Номер заявки в Украинский институт интеллектуальной собственности товарного знака из Республики Крым, регистрация которого действовала на день принятия в Российскую Федерацию Республики Крым и образования в составе Российской Федерации новых субъектов"`                                        // Номер заявки в Украинский институт интеллектуальной собственности товарного знака из Республики Крым, регистрация которого действовала на день принятия в Российскую Федерацию Республики Крым и образования в составе Российской Федерации новых субъектов
	CrimeanTrademarkApplicationDateForStateRegistrationInUkraine   string `json:"Crimean trademark application date for state registration in Ukraine" jsonschema_description:"Дата подачи заявки в Украинский институт интеллектуальной собственности товарного знака из Республики Крым, регистрация которого действовала на день принятия в Российскую Федерацию Республики Крым и образования в составе Российской Федерации новых субъектов"`                                    // Дата подачи заявки в Украинский институт интеллектуальной собственности товарного знака из Республики Крым, регистрация которого действовала на день принятия в Российскую Федерацию Республики Крым и образования в составе Российской Федерации новых субъектов
	CrimeanTrademarkCertificateNumberInUkraine                     string `json:"Crimean trademark certificate number in Ukraine" jsonschema_description:"Номер свидетельства на товарный знак из Республики Крым, действующего на день принятия в Российскую Федерацию Республики Крым и образования в составе Российской Федерации новых субъектов"`                                                                                                                                // Номер свидетельства на товарный знак из Республики Крым, действующего на день принятия в Российскую Федерацию Республики Крым и образования в составе Российской Федерации новых субъектов
	ExclusiveRightsTransferAgreementRegistrationNumber             string `json:"exclusive rights transfer agreement registration number" jsonschema_description:"Номер государственной регистрации отчуждения исключительного права по договору, на основании которого выдано свидетельство"`                                                                                                                                                                                        // Номер государственной регистрации отчуждения исключительного права по договору, на основании которого выдано свидетельство
	ExclusiveRightsTransferAgreementRegistrationDate               string `json:"exclusive rights transfer agreement registration date" jsonschema_description:"Дата государственной регистрации отчуждения исключительного права по договору, на основании которого выдано свидетельство"`                                                                                                                                                                                           // Дата государственной регистрации отчуждения исключительного права по договору, на основании которого выдано свидетельство
	LegallyRelatedApplications                                     string `json:"legally related applications" jsonschema_description:"Номера и даты юридически связанных заявок"`                                                                                                                                                                                                                                                                                                    // Номера и даты юридически связанных заявок
	LegallyRelatedRegistrations                                    string `json:"legally related registrations" jsonschema_description:"Номера и даты юридически связанных регистраций"`                                                                                                                                                                                                                                                                                              // Номера и даты юридически связанных регистраций
	ExpirationDate                                                 string `json:"expiration date" jsonschema_description:"Дата истечения срока действия исключительного права"`                                                                                                                                                                                                                                                                                                       // Дата истечения срока действия исключительного права
	RightHolderName                                                string `json:"right holder name" jsonschema_description:"Наименование или ФИО правообладателя"`                                                                                                                                                                                                                                                                                                                    // Наименование или ФИО правообладателя
	ForeignRightHolderName                                         string `json:"foreign right holder name" jsonschema_description:"Наименование или ФИО правообладателя на иностранном языке"`                                                                                                                                                                                                                                                                                       // Наименование или ФИО правообладателя на иностранном языке
	RightHolderAddress                                             string `json:"right holder address" jsonschema_description:"Адрес правообладателя"`                                                                                                                                                                                                                                                                                                                                // Адрес правообладателя
	RightHolderCountryCode                                         string `json:"right holder country code" jsonschema_description:"Код страны правообладателя"`                                                                                                                                                                                                                                                                                                                      // Код страны правообладателя
	RightHolderOgrn                                                string `json:"right holder ogrn" jsonschema_description:"ОГРН или ОГРНИП правообладателя"`                                                                                                                                                                                                                                                                                                                         // ОГРН или ОГРНИП правообладателя
	RightHolderInn                                                 string `json:"right holder inn" jsonschema_description:"ИНН правообладателя"`                                                                                                                                                                                                                                                                                                                                      // ИНН правообладателя
	CorrespondenceAddress                                          string `json:"correspondence address" jsonschema_description:"Адрес для переписки"`                                                                                                                                                                                                                                                                                                                                // Адрес для переписки
	Collective                                                     bool   `json:"collective" jsonschema_description:"Указание на то, что товарный знак является коллективным"`                                                                                                                                                                                                                                                                                                        // Указание на то, что товарный знак является коллективным
	CollectiveUsers                                                string `json:"collective users" jsonschema_description:"Сведения о лицах, имеющих право использования коллективного знака"`                                                                                                                                                                                                                                                                                        // Сведения о лицах, имеющих право использования коллективного знака
	ExtractionFromCharterOfTheCollectiveTrademark                  string `json:"extraction from charter of the collective trademark" jsonschema_description:"Выписка из устава коллективного знака о единых характеристиках качества или иных общих характеристиках товаров"`                                                                                                                                                                                                        // Выписка из устава коллективного знака о единых характеристиках качества или иных общих характеристиках товаров
	ColorSpecification                                             string `json:"color specification" jsonschema_description:"Указание цвета или цветового сочетания"`                                                                                                                                                                                                                                                                                                                // Указание цвета или цветового сочетания
	UnprotectedElements                                            string `json:"unprotected elements" jsonschema_description:"Неохраняемые элементы"`                                                                                                                                                                                                                                                                                                                                // Неохраняемые элементы
	KindSpecification                                              string `json:"kind specification" jsonschema_description:"Указание, относящееся к виду знака"`                                                                                                                                                                                                                                                                                                                     // Указание, относящееся к виду знака
	Threedimensional                                               bool   `json:"threedimensional" jsonschema_description:"Указание на то, что товарный знак является объемным"`                                                                                                                                                                                                                                                                                                      // Указание на то, что товарный знак является объемным
	ThreedimensionalSpecification                                  string `json:"threedimensional specification" jsonschema_description:"Характеристики объемного товарного знака"`                                                                                                                                                                                                                                                                                                   // Характеристики объемного товарного знака
	Holographic                                                    bool   `json:"holographic" jsonschema_description:"Указание на то, что товарный знак является голографическим"`                                                                                                                                                                                                                                                                                                    // Указание на то, что товарный знак является голографическим
	HolographicSpecification                                       string `json:"holographic specification" jsonschema_description:"Характеристики голографического товарного знака"`                                                                                                                                                                                                                                                                                                 // Характеристики голографического товарного знака
	Sound                                                          bool   `json:"sound" jsonschema_description:"Указание на то, что товарный знак является звуковым"`                                                                                                                                                                                                                                                                                                                 // Указание на то, что товарный знак является звуковым
	SoundSpecification                                             string `json:"sound specification" jsonschema_description:"Характеристики звукового товарного знака"`                                                                                                                                                                                                                                                                                                              // Характеристики звукового товарного знака
	Olfactory                                                      bool   `json:"olfactory" jsonschema_description:"Указание на то, что товарный знак является обонятельным"`                                                                                                                                                                                                                                                                                                         // Указание на то, что товарный знак является обонятельным
	OlfactorySpecification                                         string `json:"olfactory specification" jsonschema_description:"Характеристики обонятельного товарного знака"`                                                                                                                                                                                                                                                                                                      // Характеристики обонятельного товарного знака
	Color                                                          bool   `json:"color" jsonschema_description:"Указание на то, что товарный знак состоит исключительно из одного или нескольких цветов"`                                                                                                                                                                                                                                                                             // Указание на то, что товарный знак состоит исключительно из одного или нескольких цветов
	ColorTrademarkSpecification                                    string `json:"color trademark specification" jsonschema_description:"Характеристики товарного знака, который состоит исключительно из одного или нескольких цветов"`                                                                                                                                                                                                                                               // Характеристики товарного знака, который состоит исключительно из одного или нескольких цветов
	Light                                                          bool   `json:"light" jsonschema_description:"Указание на то, что товарный знак является световым"`                                                                                                                                                                                                                                                                                                                 // Указание на то, что товарный знак является световым
	LightSpecification                                             string `json:"light specification" jsonschema_description:"Характеристики светового товарного знака"`                                                                                                                                                                                                                                                                                                              // Характеристики светового товарного знака
	Changing                                                       bool   `json:"changing" jsonschema_description:"Указание на то, что товарный знак является изменяющимся"`                                                                                                                                                                                                                                                                                                          // Указание на то, что товарный знак является изменяющимся
	ChangingSpecification                                          string `json:"changing specification" jsonschema_description:"Характеристики изменяющегося товарного знака"`                                                                                                                                                                                                                                                                                                       // Характеристики изменяющегося товарного знака
	Positional                                                     bool   `json:"positional" jsonschema_description:"Указание на то, что товарный знак является позиционным"`                                                                                                                                                                                                                                                                                                         // Указание на то, что товарный знак является позиционным
	PositionalSpecification                                        string `json:"positional specification" jsonschema_description:"Характеристики позиционного товарного знака"`                                                                                                                                                                                                                                                                                                      // Характеристики позиционного товарного знака
	Actual                                                         bool   `json:"actual" jsonschema_description:"Признак действия правовой охраны"`                                                                                                                                                                                                                                                                                                                                   // Признак действия правовой охраны
	PublicationURL                                                 string `json:"publication URL" jsonschema_description:"URL публикации в открытых реестрах сайта ФИПС"`                                                                                                                                                                                                                                                                                                             // URL публикации в открытых реестрах сайта ФИПС
}

// ОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков - Открытый реестр общеизвестных в Российской Федерации товарных знаков
// https://rospatent.gov.ru/opendata/7730176088-otz
type ОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков struct {
	RegistrationNumber                         string `json:"registration number" jsonschema_description:"Номер государственной регистрации"`                                                                                                                             // Номер государственной регистрации
	RegistrationDate                           string `json:"registration date" jsonschema_description:"Дата вступления в силу решения о признании товарного знака общеизвестным в РФ"`                                                                                   // Дата вступления в силу решения о признании товарного знака общеизвестным в РФ
	WellKnownTrademarkDate                     string `json:"well-known trademark date" jsonschema_description:"Дата, с которой товарный знак признан общеизвестным"`                                                                                                     // Дата, с которой товарный знак признан общеизвестным
	LegallyRelatedRegistrations                string `json:"legally related registrations" jsonschema_description:"Номера и даты юридически связанных регистраций"`                                                                                                      // Номера и даты юридически связанных регистраций
	RightHolderName                            string `json:"right holder name" jsonschema_description:"Наименование или ФИО правообладателя"`                                                                                                                            // Наименование или ФИО правообладателя
	ForeignRightHolderName                     string `json:"foreign right holder name" jsonschema_description:"Наименование или ФИО правообладателя на иностранном языке"`                                                                                               // Наименование или ФИО правообладателя на иностранном языке
	RightHolderAddress                         string `json:"right holder address" jsonschema_description:"Адрес правообладателя"`                                                                                                                                        // Адрес правообладателя
	RightHolderCountryCode                     string `json:"right holder country code" jsonschema_description:"Код страны правообладателя"`                                                                                                                              // Код страны правообладателя
	RightHolderOgrn                            string `json:"right holder ogrn" jsonschema_description:"ОГРН или ОГРНИП правообладателя"`                                                                                                                                 // ОГРН или ОГРНИП правообладателя
	RightHolderInn                             string `json:"right holder inn" jsonschema_description:"ИНН правообладателя"`                                                                                                                                              // ИНН правообладателя
	CorrespondenceAddress                      string `json:"correspondence address" jsonschema_description:"Адрес для переписки"`                                                                                                                                        // Адрес для переписки
	Collective                                 bool   `json:"collective" jsonschema_description:"Указание на то, что общеизвестный товарный знак является коллективным"`                                                                                                  // Указание на то, что общеизвестный товарный знак является коллективным
	CollectiveUsers                            string `json:"collective users" jsonschema_description:"Сведения о лицах, имеющих право использования коллективного общеизвестного знака"`                                                                                 // Сведения о лицах, имеющих право использования коллективного общеизвестного знака
	ExtractionFromCharterOfCollectiveTrademark string `json:"extraction from charter of the collective trademark" jsonschema_description:"Выписка из устава коллективного общеизвестного знака о единых характеристиках качества или иных общих характеристиках товаров"` // Выписка из устава коллективного общеизвестного знака о единых характеристиках качества или иных общих характеристиках товаров
	ColorSpecification                         string `json:"color specification" jsonschema_description:"Указание цвета или цветового сочетания"`                                                                                                                        // Указание цвета или цветового сочетания
	UnprotectedElements                        string `json:"unprotected elements" jsonschema_description:"Неохраняемые элементы"`                                                                                                                                        // Неохраняемые элементы
	KindSpecification                          string `json:"kind specification" jsonschema_description:"Указание, относящееся к виду знака"`                                                                                                                             // Указание, относящееся к виду знака
	Threedimensional                           bool   `json:"threedimensional" jsonschema_description:"Указание на то, что общеизвестный товарный знак является объемным"`                                                                                                // Указание на то, что общеизвестный товарный знак является объемным
	ThreedimensionalSpecification              string `json:"threedimensional specification" jsonschema_description:"Характеристики объемного общеизвестного товарного знака"`                                                                                            // Характеристики объемного общеизвестного товарного знака
	Holographic                                bool   `json:"holographic" jsonschema_description:"Указание на то, что общеизвестный товарный знак является голографическим"`                                                                                              // Указание на то, что общеизвестный товарный знак является голографическим
	HolographicSpecification                   string `json:"holographic specification" jsonschema_description:"Характеристики голографического общеизвестного товарного знака"`                                                                                          // Характеристики голографического общеизвестного товарного знака
	Sound                                      bool   `json:"sound" jsonschema_description:"Указание на то, что общеизвестный товарный знак является звуковым"`                                                                                                           // Указание на то, что общеизвестный товарный знак является звуковым
	SoundSpecification                         string `json:"sound specification" jsonschema_description:"Характеристики звукового общеизвестного товарного знака"`                                                                                                       // Характеристики звукового общеизвестного товарного знака
	Olfactory                                  bool   `json:"olfactory" jsonschema_description:"Указание на то, что общеизвестный товарный знак является обонятельным"`                                                                                                   // Указание на то, что общеизвестный товарный знак является обонятельным
	OlfactorySpecification                     string `json:"olfactory specification" jsonschema_description:"Характеристики обонятельного общеизвестного товарного знака"`                                                                                               // Характеристики обонятельного общеизвестного товарного знака
	Color                                      bool   `json:"color" jsonschema_description:"Указание на то, что общеизвестный товарный знак состоит исключительно из одного или нескольких цветов"`                                                                       // Указание на то, что общеизвестный товарный знак состоит исключительно из одного или нескольких цветов
	ColorTrademarkSpecification                string `json:"color trademark specification" jsonschema_description:"Характеристики общеизвестного товарного знака, который состоит исключительно из одного или нескольких цветов"`                                        // Характеристики общеизвестного товарного знака, который состоит исключительно из одного или нескольких цветов
	Light                                      bool   `json:"light" jsonschema_description:"Указание на то, что общеизвестный товарный знак является световым"`                                                                                                           // Указание на то, что общеизвестный товарный знак является световым
	LightSpecification                         string `json:"light specification" jsonschema_description:"Характеристики светового общеизвестного товарного знака"`                                                                                                       // Характеристики светового общеизвестного товарного знака
	Changing                                   bool   `json:"changing" jsonschema_description:"Указание на то, что общеизвестный товарный знак является изменяющимся"`                                                                                                    // Указание на то, что общеизвестный товарный знак является изменяющимся
	ChangingSpecification                      string `json:"changing specification" jsonschema_description:"Характеристики изменяющегося общеизвестного товарного знака"`                                                                                                // Характеристики изменяющегося общеизвестного товарного знака
	Positional                                 bool   `json:"positional" jsonschema_description:"Указание на то, что общеизвестный товарный знак является позиционным"`                                                                                                   // Указание на то, что общеизвестный товарный знак является позиционным
	PositionalSpecification                    string `json:"positional specification" jsonschema_description:"Характеристики позиционного общеизвестного товарного знака"`                                                                                               // Характеристики позиционного общеизвестного товарного знака
	Actual                                     bool   `json:"actual" jsonschema_description:"Признак действия правовой охраны"`                                                                                                                                           // Признак действия правовой охраны
	PublicationURL                             string `json:"publication URL" jsonschema_description:"URL публикации в открытых реестрах сайта ФИПС"`                                                                                                                     // URL публикации в открытых реестрах сайта ФИПС
}

// RegistrationDateDateTime - ...
func (otz ОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков) RegistrationDateDateTime() *time.Time {
	if len(otz.RegistrationDate) == 0 {
		return nil
	}
	t, _ := time.Parse("20060102", otz.RegistrationDate)
	return &t
}

// WellKnownTrademarkDateDateTime - ...
func (otz ОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков) WellKnownTrademarkDateDateTime() *time.Time {
	if len(otz.WellKnownTrademarkDate) == 0 {
		return nil
	}
	t, _ := time.Parse("20060102", otz.WellKnownTrademarkDate)
	return &t
}

// ОргВклМСП - Сведения о юридическом лице, включенном в реестр МСП
type ОргВклМСП struct {
	НаимОрг     string `xml:"НаимОрг,attr" jsonschema_description:"Полное наименование юридического лица на русском языке"`          // Полное наименование юридического лица на русском языке
	НаимОргСокр string `xml:"НаимОргСокр,attr" jsonschema_description:"Сокращенное наименование юридического лица на русском языке"` // Сокращенное наименование юридического лица на русском языке
	ИННЮЛ       string `xml:"ИННЮЛ,attr" jsonschema_description:"ИНН юридического лица"`                                             // ИНН юридического лица
	ОГРН        string `xml:"ОГРН,attr" jsonschema_description:"ОГРН юридического лица"`                                             // ОГРН юридического лица
}

// МестоТип - ...
type МестоТип struct {
	Тип  string `xml:"Тип,attr" jsonschema_description:"Тип места"`
	Наим string `xml:"Наим,attr" jsonschema_description:"Наименование места"`
}

// СведМН - Сведения о месте нахождения юридического лица / месте жительства индивидуального предпринимателя
type СведМН struct {
	КодРегион  int64    `xml:"КодРегион,attr" jsonschema_description:"Клд региона"`                // Клд региона
	Регион     МестоТип `xml:"Регион" jsonschema_description:"Субъект Российской Федерации"`       // Субъект Российской Федерации
	Район      МестоТип `xml:"Район" jsonschema_description:"Район (улус и т.п.)"`                 // Район (улус и т.п.)
	Город      МестоТип `xml:"Город" jsonschema_description:"Город (волость и т.п.)"`              // Город (волость и т.п.)
	НаселПункт МестоТип `xml:"НаселПункт" jsonschema_description:"Населенный пункт (село и т.п.)"` // Населенный пункт (село и т.п.)
}

// ИПВклМСП - Сведения об индивидуальном предпринимателе, включенном в реестр МСП
type ИПВклМСП struct {
	ИННФЛ  string `xml:"ИННФЛ,attr" jsonschema_description:"ИНН индивидуального предпринимателя"`               // ИНН индивидуального предпринимателя
	ОГРНИП string `xml:"ОГРНИП,attr" jsonschema_description:"ОГРНИП индивидуального предпринимателя"`           // ОГРНИП индивидуального предпринимателя
	ФИОИП  ФИО    `xml:"ФИОИП" jsonschema_description:"Фамилия, имя, отчество индивидуального предпринимателя"` // Фамилия, имя, отчество индивидуального предпринимателя
}

// СвЛиценз - Сведения о лицензиях, выданных субъекту МСП
type СвЛиценз struct {
	НомЛиценз     string `xml:"НомЛиценз,attr" jsonschema_description:"Номер лицензии"`                                                                // Номер лицензии
	ДатаЛиценз    string `xml:"ДатаЛиценз,attr" jsonschema_description:"Дата лицензии"`                                                                // Дата лицензии
	ДатаНачЛиценз string `xml:"ДатаНачЛиценз,attr" jsonschema_description:"Дата начала действия лицензии"`                                             // Дата начала действия лицензии
	ДатаКонЛиценз string `xml:"ДатаКонЛиценз,attr" jsonschema_description:"Дата окончания действия лицензии"`                                          // Дата окончания действия лицензии
	ДатаОстЛиценз string `xml:"ДатаОстЛиценз,attr" jsonschema_description:"Дата приостановления действия лицензии"`                                    // Дата приостановления действия лицензии
	СерЛиценз     string `xml:"СерЛиценз,attr" jsonschema_description:"Серия лицензии"`                                                                // Серия лицензии
	ВидЛиценз     string `xml:"ВидЛиценз,attr" jsonschema_description:"Вид лицензии"`                                                                  // Вид лицензии
	ОргВыдЛиценз  string `xml:"ОргВыдЛиценз,attr" jsonschema_description:"Наименование лицензирующего органа, выдавшего или переоформившего лицензию"` // Наименование лицензирующего органа, выдавшего или переоформившего лицензию
	ОргОстЛиценз  string `xml:"ОргОстЛиценз,attr" jsonschema_description:"Наименование лицензирующего органа, приостановившего действие лицензии"`     // Наименование лицензирующего органа, приостановившего действие лицензии
	НаимЛицВД     string `xml:"НаимЛицВД" jsonschema_description:"Наименование лицензируемого вида деятельности, на который выдана лицензия"`          // Наименование лицензируемого вида деятельности, на который выдана лицензия
}

// СвПрод - Сведения о производимой субъектом МСП продукции
type СвПрод struct {
	КодПрод   string `xml:"КодПрод,attr" jsonschema_description:"Код вида продукции"`                                                 // Код вида продукции
	НаимПрод  string `xml:"НаимПрод,attr" jsonschema_description:"Наименование вида продукции"`                                       // Наименование вида продукции
	ПрОтнПрод int64  `xml:"ПрОтнПрод,attr" jsonschema_description:"Признак отнесения продукции к инновационной, высокотехнологичной"` // Признак отнесения продукции к инновационной, высокотехнологичной
}

// СвПрогПарт - Сведения о включении субъекта МСП в реестры программ партнерства
type СвПрогПарт struct {
	НаимЮЛПП string `xml:"НаимЮЛ_ПП,attr" jsonschema_description:"Наименование заказчика, реализующего программу партнерства"`    // Наименование заказчика, реализующего программу партнерства
	ИННЮЛПП  string `xml:"ИННЮЛ_ПП,attr" jsonschema_description:"ИНН заказчика, реализующего программу партнерства"`              // ИНН заказчика, реализующего программу партнерства
	НомДог   string `xml:"НомДог,attr" jsonschema_description:"Номер договора о присоединении к выбранной программе партнерства"` // Номер договора о присоединении к выбранной программе партнерства
	ДатаДог  string `xml:"ДатаДог,attr" jsonschema_description:"Дата договора о присоединении к выбранной программе партнерства"` // Дата договора о присоединении к выбранной программе партнерства
}

// СвКонтр - Сведения о наличии у субъекта МСП в предшествующем календарном году контрактов, заключенных в соответствии с Федеральным законом от 5 апреля 2013 года №44-ФЗ
type СвКонтр struct {
	НаимЮЛЗК       string `xml:"НаимЮЛ_ЗК,attr" jsonschema_description:"Наименование заказчика по контракту"` // Наименование заказчика по контракту
	ИННЮЛЗК        string `xml:"ИННЮЛ_ЗК,attr" jsonschema_description:"ИНН заказчика по контракту"`           // ИНН заказчика по контракту
	ПредмКонтр     string `xml:"ПредмКонтр,attr" jsonschema_description:"Предмет контракта"`                  // Предмет контракта
	НомКонтрРеестр string `xml:"НомКонтрРеестр,attr" jsonschema_description:"Реестровый номер контракта"`     // Реестровый номер контракта
	ДатаКонтр      string `xml:"ДатаКонтр,attr" jsonschema_description:"Дата заключения контракта"`           // Дата заключения контракта
}

// СвДог - Сведения о наличии у субъекта МСП в предшествующем календарном году договоров, заключенных в соответствии с Федеральным законом от 18 июля 2011 года №223-ФЗ
type СвДог struct {
	НаимЮЛЗД     string `xml:"НаимЮЛ_ЗД,attr" jsonschema_description:"Наименование заказчика по договору"` // Наименование заказчика по договору
	ИННЮЛЗД      string `xml:"ИННЮЛ_ЗД,attr" jsonschema_description:"ИНН заказчика по договору"`           // ИНН заказчика по договору
	ПредмДог     string `xml:"ПредмДог,attr" jsonschema_description:"Предмет договора"`                    // Предмет договора
	НомДогРеестр string `xml:"НомДогРеестр,attr" jsonschema_description:"Реестровый номер договора"`       // Реестровый номер договора
	ДатаДог      string `xml:"ДатаДог,attr" jsonschema_description:"Дата заключения договора"`             // Дата заключения договора
}

// РеестрСубъектовМалогоИСреднегоПредпринимательства - Единый реестр субъектов малого и среднего предпринимательства
// https://www.nalog.gov.ru/opendata/7707329152-rsmp/
type РеестрСубъектовМалогоИСреднегоПредпринимательства struct {
	ДатаСост    string    `xml:"ДатаСост,attr" jsonschema_description:"По состоянию реестра на дату"`                                                                                    // По состоянию реестра на дату
	ДатаВклМСП  string    `xml:"ДатаВклМСП,attr" jsonschema_description:"Дата включения юридического лица / индивидуального предпринимателя в реестр МСП"`                               // Дата включения юридического лица / индивидуального предпринимателя в реестр МСП
	ВидСубМСП   int64     `xml:"ВидСубМСП,attr" jsonschema_description:"Вид субъекта МСП"`                                                                                               // Вид субъекта МСП
	КатСубМСП   int64     `xml:"КатСубМСП,attr" jsonschema_description:"Категория субъекта МСП"`                                                                                         // Категория субъекта МСП
	ПризНовМСП  int64     `xml:"ПризНовМСП,attr" jsonschema_description:"Признак сведений о вновь созданном юридическом лице / вновь зарегистрированном индивидуальном предпринимателе"` // Признак сведений о вновь созданном юридическом лице / вновь зарегистрированном индивидуальном предпринимателе
	СведСоцПред int64     `xml:"СведСоцПред,attr" jsonschema_description:"Сведения о том, что юридическое лицо / индивидуальный предприниматель является социальным предприятием"`       // Сведения о том, что юридическое лицо / индивидуальный предприниматель является социальным предприятием
	ССЧР        int64     `xml:"ССЧР,attr" jsonschema_description:"Сведения о среднесписочной численности работников"`                                                                   // Сведения о среднесписочной численности работников
	ОргВклМСП   ОргВклМСП `xml:"ОргВклМСП" jsonschema_description:"Сведения о юридическом лице, включенном в реестр МСП"`                                                                // Сведения о юридическом лице, включенном в реестр МСП
	СведМН      СведМН    `xml:"СведМН" jsonschema_description:"Сведения о месте нахождения юридического лица / месте жительства индивидуального предпринимателя"`                       // Сведения о месте нахождения юридического лица / месте жительства индивидуального предпринимателя
	СвОКВЭД     struct {
		СвОКВЭДОсн struct {
			КодОКВЭД  string `xml:"КодОКВЭД,attr" jsonschema_description:"Код вида деятельности по Общероссийскому классификатору видов экономической деятельности"`           // Код вида деятельности по Общероссийскому классификатору видов экономической деятельности
			НаимОКВЭД string `xml:"НаимОКВЭД,attr" jsonschema_description:"Наименование вида деятельности по Общероссийскому классификатору видов экономической деятельности"` // Наименование вида деятельности по Общероссийскому классификатору видов экономической деятельности
			ВерсОКВЭД int64  `xml:"ВерсОКВЭД,attr" jsonschema_description:"Признак версии Общероссийского классификатора видов экономической деятельности"`                    // Признак версии Общероссийского классификатора видов экономической деятельности
		} `xml:"СвОКВЭДОсн"` // Сведения о кодах по Общероссийскому классификатору видов экономической деятельности
		СвОКВЭДДоп []struct {
			КодОКВЭД  string `xml:"КодОКВЭД,attr" jsonschema_description:"Код вида деятельности по Общероссийскому классификатору видов экономической деятельности"`           // Код вида деятельности по Общероссийскому классификатору видов экономической деятельности
			НаимОКВЭД string `xml:"НаимОКВЭД,attr" jsonschema_description:"Наименование вида деятельности по Общероссийскому классификатору видов экономической деятельности"` // Наименование вида деятельности по Общероссийскому классификатору видов экономической деятельности
			ВерсОКВЭД int64  `xml:"ВерсОКВЭД,attr" jsonschema_description:"Признак версии Общероссийского классификатора видов экономической деятельности"`                    // Признак версии Общероссийского классификатора видов экономической деятельности
		} `xml:"СвОКВЭДДоп"`
	} `xml:"СвОКВЭД" json:"-"` // ! Сведения о видах экономической деятельности по Общероссийскому классификатору видов экономической деятельности
	ИПВклМСП   ИПВклМСП      `xml:"ИПВклМСП" jsonschema_description:"Сведения об индивидуальном предпринимателе, включенном в реестр МСП"`                                                                                          // Сведения об индивидуальном предпринимателе, включенном в реестр МСП
	СвЛиценз   СвЛиценз      `xml:"СвЛиценз" jsonschema_description:"Сведения о лицензиях, выданных субъекту МСП"`                                                                                                                  // Сведения о лицензиях, выданных субъекту МСП
	СвПрод     *[]СвПрод     `xml:"СвПрод" jsonschema_description:"Сведения о производимой субъектом МСП продукции"`                                                                                                                // Сведения о производимой субъектом МСП продукции
	СвПрогПарт *[]СвПрогПарт `xml:"СвПрогПарт" jsonschema_description:"? (list???) Сведения о включении субъекта МСП в реестры программ партнерства"`                                                                               // ? (list???) Сведения о включении субъекта МСП в реестры программ партнерства
	СвКонтр    *[]СвКонтр    `xml:"СвКонтр" jsonschema_description:"Сведения о наличии у субъекта МСП в предшествующем календарном году контрактов, заключенных в соответствии с Федеральным законом от 5 апреля 2013 года №44-ФЗ"` // Сведения о наличии у субъекта МСП в предшествующем календарном году контрактов, заключенных в соответствии с Федеральным законом от 5 апреля 2013 года №44-ФЗ
	СвДог      *[]СвДог      `xml:"СвДог" jsonschema_description:"Сведения о наличии у субъекта МСП в предшествующем календарном году договоров, заключенных в соответствии с Федеральным законом от 18 июля 2011 года №223-ФЗ"`    // Сведения о наличии у субъекта МСП в предшествующем календарном году договоров, заключенных в соответствии с Федеральным законом от 18 июля 2011 года №223-ФЗ
}

// DocDate - строка даты в формат time.Time
func (smp РеестрСубъектовМалогоИСреднегоПредпринимательства) DocDate() time.Time {
	return transformDate(smp.ДатаСост)
}

// INN - возвращает ИНН ИП или Брлица
func (smp РеестрСубъектовМалогоИСреднегоПредпринимательства) INN() string {
	if smp.ОргВклМСП.ИННЮЛ == "" {
		return smp.ИПВклМСП.ИННФЛ
	}
	return smp.ОргВклМСП.ИННЮЛ
}

// OGRN - возвращает ОГРН ИП или Брлица
func (smp РеестрСубъектовМалогоИСреднегоПредпринимательства) OGRN() string {
	if smp.ОргВклМСП.ОГРН == "" {
		return smp.ИПВклМСП.ОГРНИП
	}
	return smp.ОргВклМСП.ОГРН
}

// ОКВЭД - Общероссийский классификатор видов экономической деятельности (ОКВЭД2)
// https://rosstat.gov.ru/opendata/7708234640-okved2
type ОКВЭД struct {
	КодОКВЭД  string `jsonschema_description:"Код вида деятельности по Общероссийскому классификатору видов экономической деятельности"`          // Код вида деятельности по Общероссийскому классификатору видов экономической деятельности
	НаимОКВЭД string `jsonschema_description:"Наименование вида деятельности по Общероссийскому классификатору видов экономической деятельности"` // Наименование вида деятельности по Общероссийскому классификатору видов экономической деятельности
}

// СведНаруш - ...
type СведНаруш struct {
	СумШтраф float64 `xml:"СумШтраф,attr" jsonschema_description:"Сумма штрафа"` // Сумма штрафа
}

// НалоговыхПравонарушенияхИМерахОтветственности - Сведения о налоговых правонарушениях и мерах ответственности за их совершение
// https://www.nalog.gov.ru/opendata/7707329152-taxoffence/
type НалоговыхПравонарушенияхИМерахОтветственности struct {
	ДатаДок   string    `xml:"ДатаДок,attr" jsonschema_description:"Дата актуальности данных"` // Дата актуальности данных
	СведНП    СведНП    `xml:"СведНП" jsonschema_description:"Сведения о налогоплательщике"`
	СведНаруш СведНаруш `xml:"СведНаруш" jsonschema_description:"Сведения о правонарушениях"`
}

// DocDate - строка даты в формат time.Time
func (npo НалоговыхПравонарушенияхИМерахОтветственности) DocDate() time.Time {
	return transformDate(npo.ДатаДок)
}

// БухОтчетность - обобзенный тип для бух.отчета
type БухОтчетность interface {
	// БухОтчетностьV503 | БухОтчетностьV508
	DocDate() time.Time
	ReportYear() int64
	INN() string
}

// ВПокОПП - ...
type ВПокОПП struct {
	НаимПок  string  `xml:"НаимПок,attr" jsonschema_description:"Наименование показателя"`
	СумОтч   float64 `xml:"СумОтч,attr" jsonschema_description:"На отчетную дату отчетного периода"`
	СумПрдщ  float64 `xml:"СумПрдщ,attr" jsonschema_description:"На 31 декабря предыдущего года"`
	СумПрдшв float64 `xml:"СумПрдшв,attr" jsonschema_description:"На 31 декабря года, предшествующего предыдущему"`
}

// ДопПокОП - ...
type ДопПокОП struct {
	НаимПок string  `xml:"НаимПок,attr" jsonschema_description:"Наименование показателя"`
	СумОтч  float64 `xml:"СумОтч,attr" jsonschema_description:"На отчетную дату отчетного периода"`
	СумПред float64 `xml:"СумПред,attr" jsonschema_description:"На 31 декабря предыдущего года"`
}

// ФИО - ...
type ФИО struct {
	Фамилия  string `xml:"Фамилия,attr" jsonschema_description:"Фамилия"`
	Имя      string `xml:"Имя,attr" jsonschema_description:"Имя"`
	Отчество string `xml:"Отчество,attr" jsonschema_description:"Отчество"`
}

// СвПред - ...
type СвПред struct {
	НаимДок string `xml:"НаимДок,attr" jsonschema_description:"Наименование и реквизиты документа, подтверждающего полномочия уполномоченного представителя"`
}

// Подписант - подписан их бухотчета
type Подписант struct {
	ПрПодп string `xml:"ПрПодп,attr" jsonschema_description:"Признак лица, подписавшего документ"`
	ФИО    ФИО    `xml:"ФИО" jsonschema_description:"Фамилия, имя, отчество руководителя (уполномоченного представителя)"`
	СвПред СвПред `xml:"СвПред" jsonschema_description:"Сведения об уполномоченном представителе"`
}

// НПЮЛ - ...
type НПЮЛ struct {
	НаимОрг string `xml:"НаимОрг,attr" jsonschema_description:"Наименование организации"`
	ИННЮЛ   string `xml:"ИННЮЛ,attr" jsonschema_description:"ИНН"`
	КПП     string `xml:"КПП,attr" jsonschema_description:"КПП"`
	АдрМН   string `xml:"АдрМН,attr" jsonschema_description:"Адрес местанахождения"`
}

// СвНП - ...
type СвНП struct {
	ОКВЭД2 string `xml:"ОКВЭД2,attr" jsonschema_description:"Код вида экономической деятельности по ОКВЭД 2"`
	ОКПО   string `xml:"ОКПО,attr" jsonschema_description:"Код по ОКПО"`
	ОКФС   string `xml:"ОКФС,attr" jsonschema_description:"Форма собственности (по ОКФС)"`
	ОКОПФ  string `xml:"ОКОПФ,attr" jsonschema_description:"Организационно-правовая форма (по ОКОПФ)"`
	НПЮЛ   НПЮЛ   `xml:"НПЮЛ" jsonschema_description:"Организация"`
}

// Кап31ДекПред - ...
type Кап31ДекПред struct {
	УстКапитал  float64 `xml:"УстКапитал,attr" jsonschema_description:"Уставный капитал"`
	СобВыкупАкц float64 `xml:"СобВыкупАкц,attr" jsonschema_description:"Собственные акции, выкупленные у акционеров"`
	ДобКапитал  float64 `xml:"ДобКапитал,attr" jsonschema_description:"Добавочный капитал"`
	РезКапитал  float64 `xml:"РезКапитал,attr" jsonschema_description:"Резервный капитал"`
	НераспПриб  float64 `xml:"НераспПриб,attr" jsonschema_description:"Нераспределенная прибыль (непокрытый убыток)"`
	Итог        float64 `xml:"Итог,attr" jsonschema_description:"Итого"`
}

// АудитОрг - ...
type АудитОрг struct {
	НаимОрг string `xml:"НаимОрг,attr" jsonschema_description:"Наименование аудиторской организации"`
	ИННЮЛ   string `xml:"ИННЮЛ,attr" jsonschema_description:"ИНН аудиторской организации"`
	ОГРН    string `xml:"ОГРН,attr" jsonschema_description:"ОГРН аудиторской организации"`
}

// СвАудит - ...
type СвАудит struct {
	АудитОрг АудитОрг `xml:"АудитОрг" jsonschema_description:"Аудиторская организация"`
}

// РеестрДисквалифицированныхЛиц - Реестр дисквалифицированных лиц
// https://www.nalog.gov.ru/opendata/7707329152-registerdisqualified/
type РеестрДисквалифицированныхЛиц struct {
	ID                       uint64 `json:"НомерЗаписи" jsonschema_description:"Номер записи из реестра дисквалифицированных лиц"`                                         // Номер записи из реестра дисквалифицированных лиц
	FIO                      string `json:"ФИО" jsonschema_description:"ФИО"`                                                                                              // ФИО
	BDate                    string `json:"ДатаРожденияФЛ" jsonschema_description:"Дата рождения ФЛ"`                                                                      // Дата рождения ФЛ
	BPlace                   string `json:"МестоРожденияФЛ" jsonschema_description:"Место рождения ФЛ"`                                                                    // Место рождения ФЛ
	OrgName                  string `json:"НаименованиеОрганизации" jsonschema_description:"Наименование организации, где ФЛ работало во время совершения правонарушения"` // Наименование организации, где ФЛ работало во время совершения правонарушения
	INN                      string `json:"ИННОрганизации" jsonschema_description:"ИНН организации"`                                                                       // ИНН организации
	PositionFL               string `json:"ДолжностьФЛ" jsonschema_description:"Должность, в которой ФЛ работало во время совершения правонарушения"`                      // Должность, в которой ФЛ работало во время совершения правонарушения
	NKOAP                    string `json:"CтатьяКоАПРФ" jsonschema_description:"Cтатья КоАП РФ"`                                                                          // Cтатья КоАП РФ
	GOrgName                 string `json:"Наименование органа" jsonschema_description:"Наименование органа, составившего протокол об административном правонарушении"`    // Наименование органа, составившего протокол об административном правонарушении
	SudFIO                   string `json:"ФИОСудьи" jsonschema_description:"ФИО судьи, вынесшего постановление о дисквалификации"`                                        // ФИО судьи, вынесшего постановление о дисквалификации
	SudPosition              string `json:"ДолжностьСудьи" jsonschema_description:"Должность судьи"`                                                                       // Должность судьи
	DisqualificationDuration string `json:"CрокДисквалификации" jsonschema_description:"Cрок дисквалификации"`                                                             // Cрок дисквалификации
	DisStartDate             string `json:"ДатаНачала" jsonschema_description:"Дата начала"`                                                                               // Дата начала
	DisEndDate               string `json:"ДатаОкончания" jsonschema_description:"Дата окончания"`                                                                         // Дата окончания
}

func (dis РеестрДисквалифицированныхЛиц) String() string {
	var sb strings.Builder
	e := json.NewEncoder(&sb)
	e.SetIndent("", "\t")
	e.Encode(&dis)
	return sb.String()
}

// BDatetime - Дата рождения ФЛ
func (dis РеестрДисквалифицированныхЛиц) BDatetime() time.Time {
	return transformDate(dis.BDate)
}

// DisStartDatetime - Дата начала
func (dis РеестрДисквалифицированныхЛиц) DisStartDatetime() time.Time {
	return transformDate(dis.DisStartDate)
}

// DisEndDatetime - Дата окончания
func (dis РеестрДисквалифицированныхЛиц) DisEndDatetime() time.Time {
	return transformDate(dis.DisEndDate)
}

// OKT - ОКТМО / ОКТАО
// https://rosstat.gov.ru/opendata/7708234640-oktmo - ОКТМО
// https://rosstat.gov.ru/opendata/7708234640-okato - ОКТАО
type OKT struct {
	Ter      string `json:"ter"`      // Код региона
	Kod1     string `json:"kod1"`     // Код района/города
	Kod2     string `json:"kod2"`     // Код рабочего поселка/сельсовета
	Kod3     string `json:"kod3"`     // Код сельского населенного пункта
	Razdel   string `json:"razdel"`   // Код раздела
	Name     string `json:"name"`     // Наименование территории
	Centrum  string `json:"centrum"`  // Дополнительная информация
	NomDescr string `json:"nomdescr"` // Описание
	NomAkt   string `json:"nomakt"`   // Номер изменения
	Status   string `json:"status"`   // Тип изменения
	Dateutv  string `json:"dateutv"`  // Дата принятия
	Datevved string `json:"datevved"` // Дата введения
}

// DateutvDate - Дата
func (o OKT) DateutvDate() time.Time {
	return transformDate(o.Dateutv)
}

// DatevvedDate - Дата
func (o OKT) DatevvedDate() time.Time {
	return transformDate(o.Datevved)
}

// ПривлеченииУчастникаЗакупкиКАдминистративнойОтветственности - Информация о привлечении участника закупки к административной ответственности по ст. 19.28 КоАП
// https://zakupki.gov.ru/epz/main/public/document/view.html?searchString=&sectionId=2369&strictEqual=false
type ПривлеченииУчастникаЗакупкиКАдминистративнойОтветственности struct {
	ИНН                         string `json:"ИНН/Аналог ИНН" jsonschema_description:"ИНН/Аналог ИНН"`
	ОГРН                        string `json:"ОГРН" jsonschema_description:"ОГРН"`
	КПП                         string `json:"КПП" jsonschema_description:"КПП"`
	НаименованиеЮЛ              string `json:"Наименование ЮЛ" jsonschema_description:"Наименование ЮЛ"`
	ТипУчастника                string `json:"Тип участника" jsonschema_description:"Тип участника"`
	Суд                         string `json:"Суд" jsonschema_description:"Суд"`
	НомерДела                   string `json:"Номер дела" jsonschema_description:"Номер дела"`
	ДатаВынесенияПостановления  string `json:"Дата вынесения постановления" jsonschema_description:"Дата вынесения постановления"`
	ДатаВступленияВЗаконнуюСилу string `json:"Дата вступления в законную силу постановления о назначении административного наказания" jsonschema_description:"Дата вступления в законную силу постановления о назначении административного наказания"`
}

// DateДатаВынесенияПостановления - ...
func (puz ПривлеченииУчастникаЗакупкиКАдминистративнойОтветственности) DateДатаВынесенияПостановления() time.Time {
	t, _ := time.Parse("01.02.2006", puz.ДатаВынесенияПостановления)
	return t
}

// DateДатаВступленияВЗаконнуюСилу - ...
func (puz ПривлеченииУчастникаЗакупкиКАдминистративнойОтветственности) DateДатаВступленияВЗаконнуюСилу() time.Time {
	t, _ := time.Parse("01.02.2006", puz.ДатаВступленияВЗаконнуюСилу)
	return t
}

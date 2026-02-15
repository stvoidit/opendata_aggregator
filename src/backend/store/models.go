package store

import (
	"context"
	"time"

	// "github.com/iancoleman/orderedmap"
	"github.com/invopop/jsonschema"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

// TruncateHotels - отчистка таблицы гостиниц
func (db *DB) TruncateHotels(ctx context.Context) error {
	const q = `TRUNCATE TABLE public.hotels CASCADE`
	_, err := db.pool.Exec(ctx, q)
	return err
}

// ServiceLogSource - service_management.source_files
type ServiceLogSource struct {
	SourceType   string    `json:"source_type"`
	SourceLink   string    `json:"source_link"`
	Filename     string    `json:"filename"`
	SHA256Sum    string    `json:"sha256sum"`
	ID           string    `json:"id"`
	TaskDatetime time.Time `json:"task_datetime"`
	Downloaded   bool      `json:"downloaded"`
	Uploaded     bool      `json:"uploaded"`
}

// Hotel - модель описания гостиницы
type Hotel struct {
	FederalNumber       string                `json:"federal_number" jsonschema_description:"номер гостиницы"`
	INN                 string                `json:"inn" jsonschema_description:"ИНН"`
	OGRN                string                `json:"ogrn" jsonschema_description:"ОГРН/ОГРНИП"`
	Fullname            string                `json:"full_name" jsonschema_description:"полное наименование"`
	ShortName           string                `json:"short_name" jsonschema_description:"сокращенное наименование"`
	Type                string                `json:"type" jsonschema_description:"вид"`
	Address             string                `json:"address" jsonschema_description:"адрес"`
	Email               string                `json:"email" jsonschema_description:"email"`
	Region              string                `json:"region" jsonschema_description:"регион"`
	Site                string                `json:"site" jsonschema_description:"сайт"`
	Phone               string                `json:"phone" jsonschema_description:"телефон"`
	Fax                 string                `json:"fax" jsonschema_description:"факс"`
	Owner               string                `json:"owner" jsonschema_description:"владелец"`
	HotelClassification []HotelClassification `json:"classification" jsonschema_description:"классификация гостиницы"`
	HotelRooms          []HotelRoom           `json:"rooms" jsonschema_description:"номера гостиницы"`
}

// HotelClassification - классификация + лицензия
type HotelClassification struct {
	FederalNumber      string `json:"federal_number" jsonschema_description:"номер гостиницы"`
	DateIssued         string `json:"date_issued" jsonschema_description:"дата выдачи"`
	DateEnd            string `json:"date_end" jsonschema_description:"срок действия до"`
	Category           string `json:"category" jsonschema_description:"присвоенная категория"`
	LicenseNumber      string `json:"license_number" jsonschema_description:"регистрационный номер"`
	RegistrationNumber string `json:"registration_number" jsonschema_description:"регистрационный номер свидетельс"`
}

// HotelRoom - номер гостиницы
type HotelRoom struct {
	FederalNumber string `json:"federal_number" jsonschema_description:"номер гостиницы"`
	Category      string `json:"category" jsonschema_description:"категория"`
	Rooms         int64  `json:"rooms" jsonschema_description:"мест"`
	Seats         int64  `json:"seats" jsonschema_description:"номеров"`
}

type StatEGR struct {
	LastDateDischarge string `json:"last_date_discharge"` // Последняя дата выписки
	TotalCount        uint64 `json:"total_count"`         // всего
	EGRUL             uint64 `json:"egrul"`               // кол-во ЕГРЮЛ
	EGRIP             uint64 `json:"egrip"`               // кол-во ЕГРИП
}

type StatBalance struct {
	LastDocDate time.Time `json:"last_doc_date"` // дата последнего документа отчета
	TotalCount  uint64    `json:"total_count"`   // кол-во записей всего
	CountYears  []struct {
		Year  uint32 `json:"year"`
		Count uint64 `json:"count"`
	} `json:"count_years"`
}

type LastUpdateSource struct {
	SourceType string    `json:"source_type"`
	DateTime   time.Time `json:"datetime"`
}

type CountHotelType struct {
	Type  string `json:"type"`
	Count uint64 `json:"count"`
}

type StatsHotels struct {
	TotalCount       uint64           `json:"total_count"`
	CountHotelsTypes []CountHotelType `json:"count_hotels_types"`
}

type RossAccreditationStatus struct {
	CertStatus string `json:"cert_status"`
	Count      uint64 `json:"count"`
}

type StatsRossAccreditation struct {
	TotalCount                uint64                    `json:"total_count"`
	RossAccreditationStatuses []RossAccreditationStatus `json:"ross_accreditation_statuses"`
}

type SumByYear struct {
	Year  uint32  `json:"year"`
	Sum   float64 `json:"sum"`
	Count uint64  `json:"count"`
}

type TaxOffenses struct {
	TotalCount  uint64      `json:"total_count"`
	TotalSum    float64     `json:"total_sum"`
	SumsByYears []SumByYear `json:"sums_by_years"`
}

type StatsRegisterOfTrademarks struct {
	OpenRegistryCount      uint64 `json:"open_registry_count" jsonschema_description:"Открытый реестр товарных знаков"`
	WellKnownRegistryCount uint64 `json:"well_known_registry_count" jsonschema_description:"Общеизвестный реестр товарных знаков"`
}

type StatsFSSP struct {
	IpLegalList         uint64 `json:"ip_legal_list"`
	IpLegalListComplite uint64 `json:"ip_legal_list_complite"`
}

type StatsSMP struct {
	TotalCount uint64 `json:"total_count"`
}

type StatsDEBTAM struct {
	TotalCount uint64  `json:"total_count"`
	TotalSum   float64 `json:"total_sum"`
}

type StatsFGIS struct {
	UnscheduledCount uint64 `json:"unscheduled_count"`
	ScheduledCount   uint64 `json:"scheduled_count"`
}

type StatsTaxRegime struct {
	Count uint64 `json:"count"`
	ЕСХН  uint64 `json:"есхн"`
	УСН   uint64 `json:"усн"`
	ЕНВД  uint64 `json:"енвд"`
	СРП   uint64 `json:"срп"`
}

type StatAvgEmployesNumber struct {
	Count uint64 `json:"count"`
}

type DatabaseStat struct {
	Sources                   map[string]string         `json:"sources" jsonschema_description:"Список сокращений ссылок на источники"`
	LastUpdatesSources        []LastUpdateSource        `json:"last_updates_sources" jsonschema_description:"Последнее обновление источников"`
	StatEGR                   StatEGR                   `json:"stat_egr" jsonschema_description:"ЕГРИП и ЕГРЮЛ"`
	StatBalance               StatBalance               `json:"stat_balance" jsonschema_description:"Бухгалтерская отчетность"`
	StatsHotels               StatsHotels               `json:"stats_hotels" jsonschema_description:"Гостиницы"`
	StatsRossAccreditation    StatsRossAccreditation    `json:"stats_ross_accreditation" jsonschema_description:"Россакредитация"`
	StatsTaxOffenses          TaxOffenses               `json:"stats_tax_offenses" jsonschema_description:"Налоговые правонарушения и штрафы"`
	StatsRegisterOfTrademarks StatsRegisterOfTrademarks `json:"stats_register_of_trademarks" jsonschema_description:"Реестр товарных знаков"`
	StatsFSSP                 StatsFSSP                 `json:"stats_fssp" jsonschema_description:"ФССП"`
	StatsSMP                  StatsSMP                  `json:"stats_smp" jsonschema_description:"СМП"`
	StatsDEBTAM               StatsDEBTAM               `json:"stats_debtam" jsonschema_description:"Сведения о суммах недоимки"`
	StatsFGIS                 StatsFGIS                 `json:"stats_fgis" jsonschema_description:"ФГИС плановые и внеплановые проверки"`
	StatsTaxRegime            StatsTaxRegime            `json:"stats_tax_regime" jsonschema_description:"Режим налогоплательщика"`
	StatAvgEmployesNumber     StatAvgEmployesNumber     `json:"stat_avg_employes_number" jsonschema_description:"Кол-во записей в среднесписочной численности сотрудников"`
}

// CategoriesIP - категории ИП
type CategoriesIP struct {
	Category      string   `json:"category"`
	SubCategories []string `json:"subcategories"`
}

// LegalStatus - статусы ЮЛ
type LegalStatus map[string]string

func (ls LegalStatus) JSONSchema() *jsonschema.Schema {
	var props = orderedmap.New[string, *jsonschema.Schema]()
	props.Set("full_name", &jsonschema.Schema{Type: "string"})
	props.Set("short_name", &jsonschema.Schema{Type: "string"})
	return &jsonschema.Schema{
		Type:        "object",
		Title:       "TypeLegalStatus",
		Description: "Справочник расшифровки статусов ЮЛ ( /api/legal_statuses ) [ Map<string, string> ]",
		Properties:  props,
		Required:    []string{"full_name", "short_name"},
	}
}

// БухОтчет - принимает один из типов: typeБухОтчетностьV503 или typeБухОтчетностьV508
type БухОтчет map[string]any

func (br БухОтчет) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "array",
		Title:       "Бухгалтерская отчетность",
		Description: "Бухгалтерская отчетность [typeБухОтчетностьV503 или typeБухОтчетностьV508]",
		OneOf: []*jsonschema.Schema{
			{
				Ref:         "#/$defs/БухОтчетностьV503",
				Description: "Тип бухотчета 5.03",
			},
			{
				Ref:         "#/$defs/БухОтчетностьV508",
				Description: "Тип бухотчета 5.08",
			},
		},
	}
}

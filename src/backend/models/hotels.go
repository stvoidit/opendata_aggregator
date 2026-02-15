package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

type Hotels []HotelData

func (hs Hotels) GetCategories() []string {
	var unc = make(map[string]struct{})
	for i := range hs {
		for _, room := range hs[i].GetRooms() {
			category := room.RoomCategory
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

type HotelData struct {
	ID         string `json:"_id"`
	Number     string `json:"number"`
	RegistryID string `json:"registryId"`
	Status     string `json:"status"`
	XsdData    struct {
		ProjectOrder struct {
			Number string `json:"number"`
		} `json:"projectOrder"`
	} `json:"xsdData"`
	Objects []struct {
		Address struct {
			Country            string `json:"country"`
			FullAddress        string `json:"fullAddress"`
			IsSpecialAddress   bool   `json:"isSpecialAddress"`
			PostalCode         string `json:"postalCode"`
			UnrecognizablePart string `json:"unrecognizablePart"`
		} `json:"address"`
		GUID   string `json:"guid"`
		Name   string `json:"name"`
		Region struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"region"`
		XsdData struct {
			City                      string `json:"City"`
			ClassificationInformation struct {
				InfoAccredOrganization struct {
					AccredOrganization          string `json:"accredOrganization"`
					AccredOrganizationNumber    string `json:"accredOrganizationNumber"`
					AccredOrganizationShortName string `json:"accredOrganizationShortName"`
					Specialist1                 string `json:"specialist1"`
					Specialist2                 string `json:"specialist2"`
					Specialist3                 string `json:"specialist3"`
				} `json:"InfoAccredOrganization"`
				CategoryStars string `json:"categoryStars"`
				Order         struct {
					DateEnd           string `json:"dateEnd"`
					LicenseDateIssued string `json:"licenseDateIssued"`
					LicenseNumber     string `json:"licenseNumber"`
				} `json:"order"`
			} `json:"ClassificationInformation"`
			Email              string           `json:"Email"`
			Fax                string           `json:"Fax"`
			InfoClassification any              `json:"InfoClassification"`
			InformationRooms   InformationRooms `json:"InformationRooms"`
			Phone              string           `json:"Phone"`
			SiteURL            string           `json:"SiteUrl"`
			View               string           `json:"View"`
			NumberFederalList  string           `json:"numberFederalList"`
			OrderTrack         string           `json:"orderTrack"`
			ShortName          string           `json:"shortName"`
		} `json:"xsdData"`
	} `json:"objects"`
	Subjects []struct {
		Data struct {
			Organization *struct {
				Inn                 string `json:"inn"`
				Name                string `json:"name"`
				Ogrn                string `json:"ogrn"`
				RegistrationAddress struct {
					Country            string `json:"country"`
					FullAddress        string `json:"fullAddress"`
					IsSpecialAddress   bool   `json:"isSpecialAddress"`
					PostalCode         string `json:"postalCode"`
					UnrecognizablePart string `json:"unrecognizablePart"`
				} `json:"registrationAddress"`
				ShortName string `json:"shortName"`
			} `json:"organization"`
			Person *struct {
				FirstName  string `json:"firstName"`
				Inn        string `json:"inn"`
				LastName   string `json:"lastName"`
				MiddleName string `json:"middleName"`
				Ogrn       string `json:"ogrn"`
				Person     struct {
					FirstName  string `json:"firstName"`
					LastName   string `json:"lastName"`
					MiddleName string `json:"middleName"`
					Sex        string `json:"sex"`
				} `json:"person"`
				RegistrationAddress struct {
					UnrecognizablePart string `json:"unrecognizablePart"`
				} `json:"registrationAddress"`
				Sex string `json:"sex"`
			} `json:"person"`
		} `json:"data"`
		GUID          string `json:"guid"`
		Header        string `json:"header"`
		ShortHeader   string `json:"shortHeader"`
		SpecialTypeID string `json:"specialTypeId"`
		XsdData       struct {
			City        string `json:"City"`
			Email       string `json:"Email"`
			Okato       string `json:"Okato"`
			OwnerNameUL string `json:"OwnerNameUL"`
			Phone       string `json:"Phone"`
			Region      string `json:"Region"`
			SiteURL     string `json:"SiteUrl"`
		} `json:"xsdData"`
	} `json:"subjects"`
}

// GetNumber - Порядковый номер в Федеральном перечне
func (h HotelData) GetNumber() string {
	if h.Number != "" {
		return h.Number
	}
	if len(h.Objects) == 0 {
		return ""
	}
	return h.Objects[0].XsdData.NumberFederalList
}

// GetType - Вид
func (h HotelData) GetType() string {
	if len(h.Objects) == 0 {
		return ""
	}
	return strings.TrimSpace(h.Objects[0].XsdData.View)
}

// GetFullName - Полное наименование классифицированного объекта
func (h HotelData) GetFullName() string {
	if len(h.Objects) == 0 {
		return ""
	}
	return strings.TrimSpace(h.Objects[0].Name)
}

// GetShortName - Cокращенное наименование классифицированного объекта
func (h HotelData) GetShortName() string {
	if len(h.Objects) == 0 {
		return ""
	}
	return strings.TrimSpace(h.Objects[0].XsdData.ShortName)
}

// GetOwner - Наименование юридического лица/индивидуального предпринимателя
func (h HotelData) GetOwner() string {
	if len(h.Subjects) == 0 {
		return ""
	}
	if h.Subjects[0].SpecialTypeID == "ulApplicant" {
		return h.Subjects[0].ShortHeader
	}
	return h.Subjects[0].Header
}

// GetRegion - Регион
func (h HotelData) GetRegion() string {
	if len(h.Objects) == 0 {
		return ""
	}
	return h.Objects[0].Region.Name
}

// GetINN - инн
func (h HotelData) GetINN() string {
	if len(h.Subjects) == 0 {
		return ""
	}
	if h.Subjects[0].Data.Organization != nil {
		return h.Subjects[0].Data.Organization.Inn
	}
	return strings.TrimSpace(h.Subjects[0].Data.Person.Inn)
}

// GetINN - огрн
func (h HotelData) GetOGRN() string {
	if len(h.Subjects) == 0 {
		return ""
	}
	if h.Subjects[0].Data.Organization != nil {
		return h.Subjects[0].Data.Organization.Ogrn
	}
	return strings.TrimSpace(h.Subjects[0].Data.Person.Ogrn)
}

// GetAddress - Адрес места нахождения
func (h HotelData) GetAddress() string {
	if len(h.Objects) == 0 {
		return ""
	}
	return h.Objects[0].Address.FullAddress
}

func (h HotelData) GetPhone() string {
	if len(h.Objects) == 0 {
		return ""
	}
	return h.Objects[0].XsdData.Phone
}

func (h HotelData) GetFax() string {
	if len(h.Objects) == 0 {
		return ""
	}
	return h.Objects[0].XsdData.Fax
}

func (h HotelData) GetEmail() string {
	if len(h.Objects) == 0 {
		return ""
	}
	return h.Objects[0].XsdData.Email
}

func (h HotelData) GetSite() string {
	if len(h.Objects) == 0 {
		return ""
	}
	return h.Objects[0].XsdData.SiteURL
}

// GetCategoryStars - Присвоенная категория
func (h HotelData) GetCategoryStars() string {
	if len(h.Objects) == 0 {
		return ""
	}
	return strings.ToUpper(h.Objects[0].XsdData.ClassificationInformation.CategoryStars)
}

// Регистрационный номер
func (h HotelData) GetRegistrationNumber() string {
	return h.XsdData.ProjectOrder.Number
}

// GetLicenseNumber - Регистрационный номер
func (h HotelData) GetLicenseNumber() string {
	if len(h.Objects) == 0 {
		return ""
	}
	return h.Objects[0].XsdData.ClassificationInformation.Order.LicenseNumber
}

// - Дата выдачи свидетельства
func (h HotelData) GetLicenseDateIssued() string {
	if len(h.Objects) == 0 {
		return ""
	}
	t, err := time.Parse("2006-01-02T15:04:05-0700", h.Objects[0].XsdData.ClassificationInformation.Order.LicenseDateIssued)
	if err != nil {
		return ""
	}
	return t.Format("02.01.2006")
}

// - Срок действия свидетельства
func (h HotelData) GetLicenseDateEnd() string {
	if len(h.Objects) == 0 {
		return ""
	}
	t, err := time.Parse("2006-01-02T15:04:05-0700", h.Objects[0].XsdData.ClassificationInformation.Order.DateEnd)
	if err != nil {
		return ""
	}
	return t.Format("02.01.2006")
}

func (h HotelData) GetRooms() []InformationRoomsBlock {
	if len(h.Objects) == 0 {
		return nil
	}
	return h.Objects[0].XsdData.InformationRooms.InformationRoomsBlock
}

type InformationRoomsBlock struct {
	NumberRooms  int    `json:"numberRooms,string"`
	NumberSeats  int    `json:"numberSeats,string"`
	RoomCategory string `json:"roomCategory"`
}

type InformationRooms struct {
	InformationRoomsBlock []InformationRoomsBlock `json:"InformationRoomsBlock"`
}

func (ir *InformationRooms) UnmarshalJSON(b []byte) (err error) {
	ir.InformationRoomsBlock = make([]InformationRoomsBlock, 0)
	if bytes.HasPrefix(b, []byte("[")) {
		var _irb []map[string]InformationRoomsBlock
		err = json.Unmarshal(b, &_irb)
		for i := range _irb {
			for k := range _irb[i] {
				ir.InformationRoomsBlock = append(ir.InformationRoomsBlock, _irb[i][k])
			}
		}
		if err != nil {
			err = json.Unmarshal(b, &ir.InformationRoomsBlock)
		}
	} else {
		if bytes.Count(b, []byte("{")) == 1 {
			var irb InformationRoomsBlock
			err = json.Unmarshal(b, &irb)
			ir.InformationRoomsBlock = append(ir.InformationRoomsBlock, irb)
		} else {
			var _irb map[string]InformationRoomsBlock
			err = json.Unmarshal(b, &_irb)
			for k := range _irb {
				ir.InformationRoomsBlock = append(ir.InformationRoomsBlock, _irb[k])
			}
		}
	}
	return err
}

func (ir InformationRooms) MarshalJSON() ([]byte, error) {
	return json.Marshal(ir.InformationRoomsBlock)
}

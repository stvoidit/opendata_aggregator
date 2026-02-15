// Package csvparser - парсеры csv
package csvparser

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"opendataaggregator/models"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

const chanSize = 1_000

// ParseРосаккредитация - парсинг выгрузки в csv по ParseРосаккредитации
// ! TODO: постоянно меняют формат, в последней выгрузке не было шапки, разделитель был "|"
func ParseРосаккредитация(r *csv.Reader) ([]models.Росаккредитация, error) {
	r.Comma = ','
	r.LazyQuotes = true
	r.ReuseRecord = true
	row, err := r.Read()
	if err != nil {
		return nil, err
	}
	headers := getHeaders(row)
	var list = make([]models.Росаккредитация, chanSize)
	for {
		row, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else if errors.Is(err, csv.ErrFieldCount) {
				log.Warn().Err(err).Msg("ParseIpLegalListComplete")
				continue
			} else if errors.Is(err, csv.ErrQuote) {
				log.Warn().Err(err).Msg("ParseIpLegalListComplete")
				continue
			} else if errors.Is(err, csv.ErrBareQuote) {
				log.Warn().Err(err).Msg("ParseIpLegalListComplete")
				continue
			} else {
				log.Fatal().Err(err).Send()
			}
		}

		doc, err := decode[models.Росаккредитация](headers, row)
		if err != nil {
			return nil, err
		}
		list = append(list, doc)
	}
	return list, nil
}

// ParseRSD - Сведения из Реестра деклараций о соответствии
func ParseRSD(ctx context.Context, r *csv.Reader) <-chan models.RDS {
	var docs = make(chan models.RDS, chanSize)
	go func() {
		defer close(docs)
		r.Comma = ','
		r.LazyQuotes = true
		r.ReuseRecord = true
		r.FieldsPerRecord = 26
		for {
			select {
			case <-ctx.Done():
				return
			default:
				row, err := r.Read()
				if err != nil {
					if errors.Is(err, io.EOF) {
						return
					} else if errors.Is(err, csv.ErrFieldCount) {
						log.Warn().Err(err).Msg("ParseIpLegalListComplete")
						continue
					} else if errors.Is(err, csv.ErrQuote) {
						log.Warn().Err(err).Msg("ParseIpLegalListComplete")
						continue
					} else if errors.Is(err, csv.ErrBareQuote) {
						log.Warn().Err(err).Msg("ParseIpLegalListComplete")
						continue
					} else {
						log.Fatal().Err(err).Send()
					}
				}
				if strings.EqualFold(row[0], "id_decl") {
					continue
				}
				DateBeginning := cleanString(row[4])
				docs <- models.RDS{
					IDDecl:                        row[0],
					RegNumber:                     row[1],
					DeclStatus:                    row[2],
					DeclType:                      row[3],
					DateBeginning:                 &DateBeginning,
					DateFinish:                    cleanString(row[5]),
					DeclarationScheme:             row[6],
					ProductObjectTypeDecl:         row[7],
					ProductType:                   row[8],
					ProductGroup:                  row[9],
					ProductName:                   row[10],
					AsproductInfo:                 row[11],
					ProductTechReg:                row[12],
					OrganToCertificationName:      row[13],
					OrganToCertificationRegNumber: row[14],
					BasisForDecl:                  row[15],
					OldBasisForDecl:               row[16],
					ApplicantType:                 row[17],
					PersonApplicantType:           row[18],
					ApplicantOGRN:                 row[19],
					ApplicantINN:                  row[20],
					ApplicantName:                 row[21],
					ManufacturerType:              row[22],
					ManufacturerOGRN:              row[23],
					ManufacturerINN:               row[24],
					ManufacturerName:              row[25],
				}
			}
		}
	}()
	return docs
}

// ParseИсполнительныеПроизводстваВОтношенииЮридическихЛиц - парсинг выгрузки в csv по ИсполнительныеПроизводстваВОтношенииЮридическихЛиц
func ParseИсполнительныеПроизводстваВОтношенииЮридическихЛиц(r *csv.Reader) (<-chan models.ИсполнительныеПроизводстваВОтношенииЮридическихЛиц, error) {
	r.Comma = ','
	r.LazyQuotes = false
	_, err := r.Read()
	if err != nil {
		return nil, err
	}
	var ch = make(chan models.ИсполнительныеПроизводстваВОтношенииЮридическихЛиц, chanSize)
	go func() {
		defer close(ch)
		for {
			row, err := r.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				} else if errors.Is(err, csv.ErrFieldCount) {
					log.Warn().Err(err).Msg("ParseIpLegalListComplete")
					continue
				} else if errors.Is(err, csv.ErrQuote) {
					log.Warn().Err(err).Msg("ParseIpLegalListComplete")
					continue
				} else if errors.Is(err, csv.ErrBareQuote) {
					log.Warn().Err(err).Msg("ParseIpLegalListComplete")
					continue
				} else {
					log.Fatal().Err(err).Send()
				}
			}
			ch <- models.ИсполнительныеПроизводстваВОтношенииЮридическихЛиц{
				NameOfDebtor:                                        row[0],
				AddressOfDebtorOrganization:                         row[1],
				ActualAddressOfDebtorOrganization:                   row[2],
				NumberOfEnforcementProceeding:                       row[3],
				DateOfInstitutionProceeding:                         row[4],
				TotalNumberOfEnforcementProceedings:                 row[5],
				ExecutiveDocumentType:                               row[6],
				DateOfExecutiveDocument:                             row[7],
				NumberOfExecutiveDocument:                           row[8],
				ObjectOfExecutiveDocuments:                          row[9],
				ObjectOfExecution:                                   row[10],
				AmountDue:                                           parseFloat(row[11]),
				DebtRemainingBalance:                                parseFloat(row[12]),
				DepartmentsOfBailiffs:                               row[13],
				AddressOfDepartmentsOfBailiff:                       row[14],
				DebtorTaxpayerIdentificationNumber:                  row[15],
				TaxpayerIdentificationNumberOfOrganizationCollector: row[16],
			}
		}
	}()
	return ch, nil
}

// ParseIpLegalListComplete -
func ParseIpLegalListComplete(r *csv.Reader) (<-chan models.IpLegalListComplete, error) {
	r.Comma = ','
	r.LazyQuotes = false
	_, err := r.Read()
	if err != nil {
		return nil, err
	}
	var ch = make(chan models.IpLegalListComplete, chanSize)
	go func() {
		defer close(ch)
		for {
			row, err := r.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				} else if errors.Is(err, csv.ErrFieldCount) {
					log.Warn().Err(err).Msg("ParseIpLegalListComplete")
					continue
				} else if errors.Is(err, csv.ErrQuote) {
					log.Warn().Err(err).Msg("ParseIpLegalListComplete")
					continue
				} else if errors.Is(err, csv.ErrBareQuote) {
					log.Warn().Err(err).Msg("ParseIpLegalListComplete")
					continue
				} else {
					log.Fatal().Err(err).Send()
				}
			}
			ch <- models.IpLegalListComplete{
				NameOfDebtor:                                        row[0],
				AddressOfDebtorOrganization:                         row[1],
				ActualAddressOfDebtorOrganization:                   row[2],
				NumberOfEnforcementProceeding:                       row[3],
				DateOfInstitutionProceeding:                         row[4],
				TotalNumberOfEnforcementProceedings:                 row[5],
				ExecutiveDocumentType:                               row[6],
				DateOfExecutiveDocument:                             row[7],
				NumberOfExecutiveDocument:                           row[8],
				ObjectOfExecutiveDocuments:                          row[9],
				ObjectOfExecution:                                   row[10],
				DateCompleteIPreason:                                row[11],
				DepartmentsOfBailiffs:                               row[12],
				AddressOfDepartmentsOfBailiff:                       row[13],
				DebtorTaxpayerIdentificationNumber:                  row[14],
				TaxpayerIdentificationNumberOfOrganizationCollector: row[15],
			}
		}
	}()
	return ch, nil
}

// ParseОткрытыйРеестрТоварныхЗнаков - парсинг выгрузки в csv по ОткрытыйРеестрТоварныхЗнаков
func ParseОткрытыйРеестрТоварныхЗнаков(r *csv.Reader) (<-chan models.ОткрытыйРеестрТоварныхЗнаков, error) {
	r.Comma = ','
	r.LazyQuotes = true
	r.FieldsPerRecord = 58
	r.ReuseRecord = true
	r.TrimLeadingSpace = true
	_, err := r.Read()
	if err != nil {
		return nil, err
	}
	var ch = make(chan models.ОткрытыйРеестрТоварныхЗнаков, chanSize)
	go func() {
		defer close(ch)
		for {
			row, err := r.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				} else if errors.Is(err, csv.ErrFieldCount) {
					log.Warn().Err(err).Msg("ParseIpLegalListComplete")
					continue
				} else if errors.Is(err, csv.ErrQuote) {
					log.Warn().Err(err).Msg("ParseIpLegalListComplete")
					continue
				} else if errors.Is(err, csv.ErrBareQuote) {
					log.Warn().Err(err).Msg("ParseIpLegalListComplete")
					continue
				} else {
					log.Fatal().Err(err).Send()
				}
			}
			row[0] = strings.TrimPrefix(row[0], "\ufeff")
			ch <- models.ОткрытыйРеестрТоварныхЗнаков{
				RegistrationNumber:                                             row[0],
				RegistrationDate:                                               row[1],
				ApplicationNumber:                                              row[2],
				ApplicationDate:                                                row[3],
				PriorityDate:                                                   row[4],
				ExhibitionPriorityDate:                                         row[5],
				ParisConventionPriorityNumber:                                  row[6],
				ParisConventionPriorityDate:                                    row[7],
				ParisConventionPriorityCountryCode:                             row[8],
				InitialApplicationNumber:                                       row[9],
				InitialApplicationOriorityDate:                                 row[10],
				InitialRegistrationNumber:                                      row[11],
				InitialRegistrationDate:                                        row[12],
				InternationalRegistrationNumber:                                row[13],
				InternationalRegistrationDate:                                  row[14],
				InternationalRegistrationPriorityDate:                          row[15],
				InternationalRegistrationEntryDate:                             row[16],
				ApplicationNumberForRecognitionOfTrademarkFromCrimea:           row[17],
				ApplicationDateForRecognitionOfTrademarkFromCrimea:             row[18],
				CrimeanTrademarkApplicationNumberForStateRegistrationInUkraine: row[19],
				CrimeanTrademarkApplicationDateForStateRegistrationInUkraine:   row[20],
				CrimeanTrademarkCertificateNumberInUkraine:                     row[21],
				ExclusiveRightsTransferAgreementRegistrationNumber:             row[22],
				ExclusiveRightsTransferAgreementRegistrationDate:               row[23],
				LegallyRelatedApplications:                                     row[24],
				LegallyRelatedRegistrations:                                    row[25],
				ExpirationDate:                                                 row[26],
				RightHolderName:                                                row[27],
				ForeignRightHolderName:                                         row[28],
				RightHolderAddress:                                             row[29],
				RightHolderCountryCode:                                         row[30],
				RightHolderOgrn:                                                row[31],
				RightHolderInn:                                                 row[32],
				CorrespondenceAddress:                                          row[33],
				Collective:                                                     parseBoolean(row[34]),
				CollectiveUsers:                                                row[35],
				ExtractionFromCharterOfTheCollectiveTrademark:                  row[36],
				ColorSpecification:                                             row[37],
				UnprotectedElements:                                            row[38],
				KindSpecification:                                              row[39],
				Threedimensional:                                               parseBoolean(row[40]),
				ThreedimensionalSpecification:                                  row[41],
				Holographic:                                                    parseBoolean(row[42]),
				HolographicSpecification:                                       row[43],
				Sound:                                                          parseBoolean(row[44]),
				SoundSpecification:                                             row[45],
				Olfactory:                                                      parseBoolean(row[46]),
				OlfactorySpecification:                                         row[47],
				Color:                                                          parseBoolean(row[48]),
				ColorTrademarkSpecification:                                    row[49],
				Light:                                                          parseBoolean(row[50]),
				LightSpecification:                                             row[51],
				Changing:                                                       parseBoolean(row[52]),
				ChangingSpecification:                                          row[53],
				Positional:                                                     parseBoolean(row[54]),
				PositionalSpecification:                                        row[55],
				Actual:                                                         parseBoolean(row[56]),
				PublicationURL:                                                 row[57],
			}
		}
	}()
	return ch, nil
}

// ParseОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков - парсинг выгрузки в csv по ОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков
func ParseОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков(r *csv.Reader) (<-chan models.ОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков, error) {
	r.Comma = ','
	r.LazyQuotes = true
	r.FieldsPerRecord = 35
	r.ReuseRecord = true
	_, err := r.Read()
	if err != nil {
		return nil, err
	}
	var ch = make(chan models.ОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков, chanSize)
	go func() {
		defer close(ch)
		for {
			row, err := r.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				} else if errors.Is(err, csv.ErrFieldCount) {
					log.Warn().Err(err).Msg("ParseIpLegalListComplete")
					continue
				} else if errors.Is(err, csv.ErrQuote) {
					log.Warn().Err(err).Msg("ParseIpLegalListComplete")
					continue
				} else if errors.Is(err, csv.ErrBareQuote) {
					log.Warn().Err(err).Msg("ParseIpLegalListComplete")
					continue
				} else {
					log.Fatal().Err(err).Send()
				}
			}
			row[0] = strings.TrimPrefix(row[0], "\ufeff")
			ch <- models.ОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков{
				RegistrationNumber:                         row[0],
				RegistrationDate:                           row[1],
				WellKnownTrademarkDate:                     row[2],
				LegallyRelatedRegistrations:                row[3],
				RightHolderName:                            row[4],
				ForeignRightHolderName:                     row[5],
				RightHolderAddress:                         row[6],
				RightHolderCountryCode:                     row[7],
				RightHolderOgrn:                            row[8],
				RightHolderInn:                             row[9],
				CorrespondenceAddress:                      row[10],
				Collective:                                 parseBoolean(row[11]),
				CollectiveUsers:                            row[12],
				ExtractionFromCharterOfCollectiveTrademark: row[13],
				ColorSpecification:                         row[14],
				UnprotectedElements:                        row[15],
				KindSpecification:                          row[16],
				Threedimensional:                           parseBoolean(row[17]),
				ThreedimensionalSpecification:              row[18],
				Holographic:                                parseBoolean(row[19]),
				HolographicSpecification:                   row[20],
				Sound:                                      parseBoolean(row[21]),
				SoundSpecification:                         row[22],
				Olfactory:                                  parseBoolean(row[23]),
				OlfactorySpecification:                     row[24],
				Color:                                      parseBoolean(row[25]),
				ColorTrademarkSpecification:                row[26],
				Light:                                      parseBoolean(row[27]),
				LightSpecification:                         row[28],
				Changing:                                   parseBoolean(row[29]),
				ChangingSpecification:                      row[30],
				Positional:                                 parseBoolean(row[31]),
				PositionalSpecification:                    row[32],
				Actual:                                     parseBoolean(row[33]),
				PublicationURL:                             row[34],
			}
		}
	}()
	return ch, nil
}

// ParseОКВЭД - парсинг выгрузки в csv по ОКВЭД
func ParseОКВЭД(r *csv.Reader) (<-chan models.ОКВЭД, error) {
	r.Comma = ';'
	_, err := r.Read()
	if err != nil {
		return nil, err
	}
	var ch = make(chan models.ОКВЭД, chanSize)
	go func() {
		defer close(ch)
		for {
			row, err := r.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				} else if errors.Is(err, csv.ErrFieldCount) {
					log.Warn().Err(err).Msg("ParseIpLegalListComplete")
					continue
				} else if errors.Is(err, csv.ErrQuote) {
					log.Warn().Err(err).Msg("ParseIpLegalListComplete")
					continue
				} else if errors.Is(err, csv.ErrBareQuote) {
					log.Warn().Err(err).Msg("ParseIpLegalListComplete")
					continue
				} else {
					log.Fatal().Err(err).Send()
				}
			}
			if row[1] == "" || row[1] == " " {
				continue
			}
			ch <- models.ОКВЭД{
				КодОКВЭД:  row[1],
				НаимОКВЭД: row[2],
			}
		}
	}()
	return ch, nil
}

// ParseРеестрДисквалифицированныхЛиц - парсинг выгрузки в csv по ОткрытыйРеестрОбщеизвестныхРФТоварныхЗнаков
func ParseРеестрДисквалифицированныхЛиц(r *csv.Reader) (<-chan models.РеестрДисквалифицированныхЛиц, error) {
	r.Comma = ';'
	r.LazyQuotes = true
	_, err := r.Read()
	if err != nil {
		return nil, err
	}
	// headers := getHeaders(row)
	var ch = make(chan models.РеестрДисквалифицированныхЛиц, chanSize)
	go func() {
		defer close(ch)
		for {
			row, err := r.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				} else if errors.Is(err, csv.ErrFieldCount) {
					log.Warn().Err(err).Msg("ParseIpLegalListComplete")
					continue
				} else if errors.Is(err, csv.ErrQuote) {
					log.Warn().Err(err).Msg("ParseIpLegalListComplete")
					continue
				} else if errors.Is(err, csv.ErrBareQuote) {
					log.Warn().Err(err).Msg("ParseIpLegalListComplete")
					continue
				} else {
					log.Fatal().Err(err).Send()
				}
			}
			docID, err := strconv.ParseUint(row[0], 10, 64)
			if err != nil {
				log.Fatal().Err(err).Str("row", row[0]).Send()
			}
			// decode(headers, row, &doc)
			ch <- models.РеестрДисквалифицированныхЛиц{
				ID:                       docID,
				FIO:                      row[1],
				BDate:                    row[2],
				BPlace:                   row[3],
				OrgName:                  row[4],
				INN:                      row[5],
				PositionFL:               row[6],
				NKOAP:                    row[7],
				GOrgName:                 row[8],
				SudFIO:                   row[9],
				SudPosition:              row[10],
				DisqualificationDuration: row[11],
				DisStartDate:             row[12],
				DisEndDate:               row[13],
			}
		}
	}()
	return ch, nil
}

// ParseОКТМО - парсинг выгрузки в csv по ОКТМО
func ParseОКТМО(r *csv.Reader) (<-chan models.OKT, error) {
	r.Comma = ';'
	r.LazyQuotes = true
	var ch = make(chan models.OKT, chanSize)
	go func() {
		defer close(ch)
		for {
			row, err := r.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				} else if errors.Is(err, csv.ErrFieldCount) {
					log.Warn().Err(err).Msg("ParseIpLegalListComplete")
					continue
				} else if errors.Is(err, csv.ErrQuote) {
					log.Warn().Err(err).Msg("ParseIpLegalListComplete")
					continue
				} else if errors.Is(err, csv.ErrBareQuote) {
					log.Warn().Err(err).Msg("ParseIpLegalListComplete")
					continue
				} else {
					log.Fatal().Err(err).Send()
				}
			}
			ch <- models.OKT{
				Ter:      row[0],
				Kod1:     row[1],
				Kod2:     row[2],
				Kod3:     row[3],
				Razdel:   row[4],
				Name:     row[5],
				Centrum:  row[6],
				NomDescr: row[7],
				NomAkt:   row[9],
				Status:   row[10],
				Dateutv:  row[11],
				Datevved: row[12]}
		}
	}()
	return ch, nil
}

// ParseОКТАО - парсинг выгрузки в csv по ОКТАО
func ParseОКТАО(r *csv.Reader) (<-chan models.OKT, error) {
	r.Comma = ';'
	r.LazyQuotes = true
	var ch = make(chan models.OKT, chanSize)
	go func() {
		defer close(ch)
		for {
			row, err := r.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				} else if errors.Is(err, csv.ErrFieldCount) {
					log.Warn().Err(err).Msg("ParseIpLegalListComplete")
					continue
				} else if errors.Is(err, csv.ErrQuote) {
					log.Warn().Err(err).Msg("ParseIpLegalListComplete")
					continue
				} else if errors.Is(err, csv.ErrBareQuote) {
					log.Warn().Err(err).Msg("ParseIpLegalListComplete")
					continue
				} else {
					log.Fatal().Err(err).Send()
				}
			}
			offset := 0
			if len(row) > 12 {
				offset = 1
			}
			ch <- models.OKT{
				Ter:      row[0],
				Kod1:     row[1],
				Kod2:     row[2],
				Kod3:     row[3],
				Razdel:   row[4],
				Name:     row[5],
				Centrum:  row[6],
				NomDescr: row[7],
				NomAkt:   row[8+offset],
				Status:   row[9+offset],
				Dateutv:  row[10+offset],
				Datevved: row[11+offset]}
		}
	}()
	return ch, nil
}

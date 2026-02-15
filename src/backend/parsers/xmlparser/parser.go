// Package xmlparser: пакет для парсинга выгрузок в xml
package xmlparser

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"opendataaggregator/models"
	"opendataaggregator/models/egr"
	"runtime"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/charmap"
)

// KnownTypes - известные типы данных в xml
type KnownTypes interface {
	models.НалоговыйРежимНалогоплательщика |
		models.СведенияОСуммахНедоимки | models.СведенияОбУчастииВКонсГруппе |
		models.СведенияОСреднесписочнойЧисленностиРаботников | models.СведенияОбУплаченныхОрганизациейНалогов |
		models.РеестрСубъектовМалогоИСреднегоПредпринимательства | models.НалоговыхПравонарушенияхИМерахОтветственности |
		models.БухОтчетностьV508 | models.БухОтчетностьV503 | egr.EGRUL | egr.EGRIP | models.InspectionFGIS
}

var badEncodingCP1251 = [...][]byte{
	[]byte(`windows-1251`),
	[]byte(`WINDOWS-1251`),
	[]byte(`Windows-1251`),
}

const (
	baseRootTag = "Документ"
	channelSize = 10_000
)

func decode[T KnownTypes](zf *zip.File, rootTag string) (T, error) {
	var doc T
	rc, err := zf.Open()
	if err != nil {
		return doc, err
	}
	defer rc.Close()
	b, err := io.ReadAll(rc)
	if err != nil {
		return doc, err
	}
	//! Дикий костыль под бухгалтерскую отчетность.
	//! У некоторых есть серьезные ошибки кодировки в xml.
	if !utf8.Valid(b) {
		b, err = io.ReadAll(charmap.Windows1251.NewDecoder().Reader(bytes.NewReader(b)))
		if err != nil {
			return doc, err
		}
		for i := range badEncodingCP1251 {
			b = bytes.ReplaceAll(b, badEncodingCP1251[i], []byte(`utf-8`))
		}
	}
	d := xml.NewDecoder(bytes.NewReader(b))
	d.CharsetReader = charset.NewReaderLabel
	for {
		t, err := d.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return doc, fmt.Errorf("ERROR: %s == %s", err, zf.Name)
		}
		if e, ok := t.(xml.StartElement); ok {
			if e.Name.Local == rootTag {
				if err := d.DecodeElement(&doc, &e); err != nil {
					return doc, nil
				}
				return doc, nil
			}
		}
	}
	return doc, nil
}

// decoder - универсальный декодер
func decoder[T KnownTypes](ctx context.Context, zf *zip.File, rootTag string, ch chan<- T, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	rc, err := zf.Open()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer rc.Close()
	d := xml.NewDecoder(rc)
	d.CharsetReader = charset.NewReaderLabel
	for {
		select {
		case <-ctx.Done():
			return
		default:
			t, err := d.Token()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				log.Error().Err(err).Msg(zf.Name)
				return
			}
			if e, ok := t.(xml.StartElement); ok {
				if e.Name.Local == rootTag {
					var doc T
					if err := d.DecodeElement(&doc, &e); err != nil {
						log.Fatal().Err(err).Send()
					}
					ch <- doc
				}
			}
		}
	}
}

// workersPoolDecode - лимитирование горутин для decoder
func workersPoolDecode[T KnownTypes](ctx context.Context, rch <-chan *zip.File) <-chan T {
	var workersLimit = runtime.NumCPU()
	var ch = make(chan T)
	go func() {
		defer close(ch)
		var workers sync.WaitGroup
		workers.Add(workersLimit)
		for i := 0; i < workersLimit; i++ {
			go func() {
				defer workers.Done()
				for zf := range rch {
					decoder(ctx, zf, baseRootTag, ch, nil)
				}
			}()
		}
		workers.Wait()
	}()
	return ch
}

// ParseБухОтчетность - ...
func ParseБухОтчетность(ctx context.Context, zr *zip.ReadCloser) <-chan models.БухОтчетность {
	const rootTag = "Документ"
	var docsChan = make(chan models.БухОтчетность, 1)
	go func() {
		defer close(docsChan)
		for i := range zr.File {
			if zr.File[i].FileInfo().IsDir() {
				continue
			}
			select {
			case <-ctx.Done():
				return
			default:
				switch {
				case strings.Contains(zr.File[i].Name, "NO_BUHOTCH"):
					doc, err := decode[models.БухОтчетностьV508](zr.File[i], rootTag)
					if err != nil {
						log.Warn().Err(err).Send()
						continue
					}
					doc.ВерсФорм = "5.08"
					docsChan <- doc
				case strings.Contains(zr.File[i].Name, "NO_BOUPR"):
					doc, err := decode[models.БухОтчетностьV503](zr.File[i], rootTag)
					if err != nil {
						log.Warn().Err(err).Send()
						continue
					}
					doc.ВерсФорм = "5.03"
					docsChan <- doc
				default:
					log.Warn().Str("filename", zr.File[i].Name).Msg("unknown filename type")
				}
			}
		}
	}()
	return docsChan
}

// ParseБухОтчетностьV508 - парсинг структур для БухОтчетность
func ParseБухОтчетностьV508(ctx context.Context, zr *zip.ReadCloser) <-chan models.БухОтчетностьV508 {
	var rw = make(chan *zip.File)
	go func() {
		for i := range zr.File {
			// фильтрация по нужному типу отчета
			if !strings.Contains(zr.File[i].Name, "NO_BUHOTCH") {
				continue
			}
			rw <- zr.File[i]
		}
		close(rw)
	}()
	return workersPoolDecode[models.БухОтчетностьV508](ctx, rw)
}

// ParseБухОтчетностьV503 - парсинг структур для БухОтчетность
func ParseБухОтчетностьV503(ctx context.Context, zr *zip.ReadCloser) <-chan models.БухОтчетностьV503 {
	var rw = make(chan *zip.File)
	go func() {
		for i := range zr.File {
			// фильтрация по нужному типу отчета
			if !strings.Contains(zr.File[i].Name, "NO_BOUPR") {
				continue
			}
			rw <- zr.File[i]
		}
		close(rw)
	}()
	return workersPoolDecode[models.БухОтчетностьV503](ctx, rw)
}

// parallelDecodeFunc - унифицированная функция для параллельной работы парсера, пишет в канал.
// рефакторинг чтобы избежать ляпов
func parallelDecodeFunc[T KnownTypes](ctx context.Context, zr *zip.ReadCloser, ch chan<- T) {
	defer close(ch)
	var wg sync.WaitGroup
	wg.Add(len(zr.File))
	for i := range zr.File {
		go decoder(ctx, zr.File[i], baseRootTag, ch, &wg)
	}
	wg.Wait()
}

/*
ParseParallelIterFunc - унифицированная функция парсинга, возвращающая канал как итератор.
Работает со следующими "известными типами":
models.НалоговыйРежимНалогоплательщика;
models.СведенияОСуммахНедоимки;
models.СведенияОбУчастииВКонсГруппе;
models.СведенияОСреднесписочнойЧисленностиРаботников;
models.СведенияОбУплаченныхОрганизациейНалогов;
models.РеестрСубъектовМалогоИСреднегоПредпринимательства;
models.НалоговыхПравонарушенияхИМерахОтветственности;
*/
func ParseParallelIterFunc[T KnownTypes](ctx context.Context, zr *zip.ReadCloser) <-chan T {
	var ch = make(chan T, channelSize)
	go parallelDecodeFunc(ctx, zr, ch)
	return ch
}

// ParseEGRUL - парсинг структуры для models.EGRUL
func ParseEGRUL(ctx context.Context, zr *zip.ReadCloser) <-chan egr.EGRUL {
	const rootTag = "СвЮЛ"
	var ch = make(chan egr.EGRUL)
	go func() {
		defer close(ch)
		for i := range zr.File {
			decoder(ctx, zr.File[i], rootTag, ch, nil)
		}
	}()
	return ch
}

// ParseEGRIP - парсинг структуры для models.EGRIP
func ParseEGRIP(ctx context.Context, zr *zip.ReadCloser) <-chan egr.EGRIP {
	const rootTag = "СвИП"
	var ch = make(chan egr.EGRIP, channelSize)
	go func() {
		defer close(ch)
		for i := range zr.File {
			decoder(ctx, zr.File[i], rootTag, ch, nil)
		}
	}()
	return ch
}

// ParseInspectionsFGIS - ФГИС проверки
func ParseInspectionsFGIS(ctx context.Context, zr *zip.ReadCloser) <-chan models.InspectionFGIS {
	const rootTag = "INSPECTION"
	var ch = make(chan models.InspectionFGIS, channelSize)
	go func() {
		defer close(ch)
		for i := range zr.File {
			decoder(ctx, zr.File[i], rootTag, ch, nil)
		}
	}()
	return ch
}

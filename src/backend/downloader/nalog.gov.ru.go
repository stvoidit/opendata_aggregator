// Ссылки на исходники выгрузок для обработчика:
// https://www.nalog.gov.ru/opendata/7707329152-snr/
// https://www.nalog.gov.ru/opendata/7707329152-debtam/
// https://www.nalog.gov.ru/opendata/7707329152-kgn/
// https://www.nalog.gov.ru/opendata/7707329152-rsmp/
// https://www.nalog.gov.ru/opendata/7707329152-taxoffence/
// https://www.nalog.gov.ru/opendata/7707329152-registerdisqualified/

package downloader

import (
	"context"
	"io"
	"net/http"
	"opendataaggregator/config"

	"github.com/PuerkitoBio/goquery"
)

// TODO: скачивание файла и сохранение в папку
func (dc *DownloadClient) downloadNalogGovRu(ctx context.Context, URL string, sourceType string) (df *config.DownloadedFile, err error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}
	response, err := dc.Do(request)
	if err != nil {
		return nil, err
	}
	readContentEncoding(response)
	href, err := parseNalogGovRu(response.Body)
	if err != nil {
		return nil, err
	}
	if err := response.Body.Close(); err != nil {
		return nil, err
	}
	return dc.downloadFile(ctx, href, sourceType)
}

func parseNalogGovRu(r io.Reader) (hrefDownload string, err error) {
	root, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}
	root.Find("table.border_table").
		Find("td").
		EachWithBreak(func(_ int, td *goquery.Selection) (ok bool) {
			if td.Text() == "Гиперссылка (URL) на набор" {
				hrefDownload, ok = td.Next().Find("a").Attr("href")
				return ok
			}
			return !ok
		})
	if hrefDownload == "" {
		return "", ErrHrefNotFound
	}
	return
}

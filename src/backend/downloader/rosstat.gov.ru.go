// https://rosstat.gov.ru/opendata/7708234640-okved2
// https://rosstat.gov.ru/opendata/7708234640-oktmo
// https://rosstat.gov.ru/opendata/7708234640-okato

package downloader

import (
	"context"
	"io"
	"net/http"
	"opendataaggregator/config"

	"github.com/PuerkitoBio/goquery"
)

func (dc *DownloadClient) downloadRosstatGovRu(ctx context.Context, URL string, sourceType string) (df *config.DownloadedFile, err error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}
	response, err := dc.Do(request)
	if err != nil {
		return nil, err
	}
	readContentEncoding(response)
	href, err := parseRosstatGovRu(response.Body)
	if err != nil {
		return nil, err
	}
	if err := response.Body.Close(); err != nil {
		return nil, err
	}
	return dc.downloadFile(ctx, href, sourceType)
}

func parseRosstatGovRu(r io.Reader) (hrefDownload string, err error) {
	root, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}
	root.Find(`table[typeof="foaf:Document"]`).
		Find("td").EachWithBreak(func(_ int, td *goquery.Selection) (ok bool) {
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

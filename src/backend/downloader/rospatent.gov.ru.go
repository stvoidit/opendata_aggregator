// https://rospatent.gov.ru/opendata/7730176088-otz
// https://rospatent.gov.ru/opendata/7730176088-tz

package downloader

import (
	"context"
	"io"
	"net/http"
	"opendataaggregator/config"

	"github.com/PuerkitoBio/goquery"
)

func (dc *DownloadClient) downloadRospatentGovRu(ctx context.Context, URL string, sourceType string) (df *config.DownloadedFile, err error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}
	response, err := dc.Do(request)
	if err != nil {
		return nil, err
	}
	readContentEncoding(response)
	href, err := parseRospatentGovRu(response.Body)
	if err != nil {
		return nil, err
	}
	if err := response.Body.Close(); err != nil {
		return nil, err
	}
	return dc.downloadFile(ctx, href, sourceType)
}

func parseRospatentGovRu(r io.Reader) (hrefDownload string, err error) {
	root, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}
	root.Find(`div.table`).Find("table").
		Find("td").EachWithBreak(func(_ int, td *goquery.Selection) (ok bool) {
		if td.Text() == "Гиперсылка (URL) на набор" {
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

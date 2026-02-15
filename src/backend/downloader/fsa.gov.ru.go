package downloader

import (
	"context"
	"io"
	"net/http"
	"opendataaggregator/config"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (dc *DownloadClient) downloadFsaGovRu(ctx context.Context, URL string, sourceType string) (df *config.DownloadedFile, err error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}
	response, err := dc.Do(request)
	if err != nil {
		return nil, err
	}
	readContentEncoding(response)
	href, err := parseFsaGovRu(response.Body)
	if err != nil {
		return nil, err
	}
	if err := response.Body.Close(); err != nil {
		return nil, err
	}
	return dc.downloadFile(ctx, href, sourceType)
}

func parseFsaGovRu(r io.Reader) (hrefDownload string, err error) {
	root, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}
	root.Find("div.content.text-primary").Find("table").Find("td").
		EachWithBreak(func(_ int, td *goquery.Selection) (ok bool) {
			if strings.TrimSpace(td.Text()) == "Гиперссылка (URL) на набор" {
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

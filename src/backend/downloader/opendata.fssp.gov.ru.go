package downloader

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"opendataaggregator/config"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (dc *DownloadClient) downloadOpendataFsspGovRu(ctx context.Context, URL, hostname string, sourceType string) (df *config.DownloadedFile, err error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}
	response, err := dc.Do(request)
	if err != nil {
		return nil, err
	}
	readContentEncoding(response)
	href, err := parseOpendataFsspGovRu(response.Body, hostname)
	if err != nil {
		return nil, err
	}
	if err := response.Body.Close(); err != nil {
		return nil, err
	}
	return dc.downloadFile(ctx, href, sourceType)
}

func parseOpendataFsspGovRu(r io.Reader, hostname string) (hrefDownload string, err error) {
	root, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}
	root.Find("table.b-table").Find("td.b-table__cell").
		EachWithBreak(func(_ int, td *goquery.Selection) (ok bool) {
			if strings.TrimSpace(td.Text()) == "8." {
				hrefDownload, ok = td.Next().Next().Find("a").Attr("href")
				if ok && !strings.HasPrefix(hrefDownload, "https://") {
					hrefDownload = (&url.URL{Scheme: "https", Host: hostname, Path: hrefDownload}).String()
				}
				return ok
			}
			return !ok
		})
	if hrefDownload == "" {
		return "", ErrHrefNotFound
	}
	return
}

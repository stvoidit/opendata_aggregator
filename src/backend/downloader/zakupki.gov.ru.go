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

func (dc *DownloadClient) downloadZakupkiGovRu(ctx context.Context, URL string, sourceType string) (df *config.DownloadedFile, err error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}
	response, err := dc.Do(request)
	if err != nil {
		return nil, err
	}
	readContentEncoding(response)
	href, err := parseZakupkiGovRu(response.Body)
	if err != nil {
		return nil, err
	}
	if err := response.Body.Close(); err != nil {
		return nil, err
	}
	return dc.downloadFile(ctx, href, sourceType)
}

func parseZakupkiGovRu(r io.Reader) (hrefDownload string, err error) {
	root, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}
	hrefDownload, ok := root.Find("section.docs-list").Find("a.docs-title.heading-h4").Attr("href")
	if !ok {
		return "", ErrHrefNotFound
	}
	if !strings.HasPrefix(hrefDownload, "https://zakupki.gov.ru/") {
		URL, err := url.Parse(hrefDownload)
		if err != nil {
			return "", err
		}
		URL.Scheme = "https"
		URL.Host = "zakupki.gov.ru"
		hrefDownload = URL.String()
	}
	return hrefDownload, nil
}

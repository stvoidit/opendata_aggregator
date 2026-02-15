package downloader

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"opendataaggregator/config"
	"strconv"
	"strings"
	"time"
)

func prepareLink(link string) string {
	link, _ = url.QueryUnescape(link)
	t := time.Now()
	year := strconv.Itoa(t.Year())
	month := strconv.Itoa(int(t.Month()))
	return strings.NewReplacer(`{y}`, year, `{m}`, month).Replace(link)
}

func (dc *DownloadClient) DownloadUnscheduled(ctx context.Context, URL string, sourceType string) (df *config.DownloadedFile, err error) {
	URL = prepareLink(URL)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}
	response, err := dc.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var body = struct {
		DataZipName string `json:"dataZipName"`
		DataZipUrl  string `json:"dataZipUrl"`
	}{}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		return nil, err
	}
	if err := response.Body.Close(); err != nil {
		return nil, err
	}
	return dc.downloadFile(ctx, body.DataZipUrl, sourceType)
}

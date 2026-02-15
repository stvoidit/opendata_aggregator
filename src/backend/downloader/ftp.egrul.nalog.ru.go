package downloader

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"opendataaggregator/config"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
)

// GetFilesList - link - ссылка на главную папку раздела типа "https://ftp.egrul.nalog.ru/?dir=EGRIP_405"
func (dc *DownloadClient) GetFilesList(ctx context.Context, link string, sourceType string) <-chan *config.DownloadedFile {
	var ch = make(chan *config.DownloadedFile)
	go func() {
		defer close(ch)
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		response, err := dc.Do(request)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		subFolder, err := parseDirs(response)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		var allFiles = make([]string, 0)
		for i := range subFolder {
			log.Info().Str("check folder", subFolder[i]).Send()
			request, err := http.NewRequestWithContext(ctx, http.MethodGet, subFolder[i], nil)
			if err != nil {
				log.Fatal().Err(err).Send()
			}
			response, err := dc.Do(request)
			if err != nil {
				log.Fatal().Err(err).Send()
			}
			files, err := parseDirs(response)
			if err != nil {
				log.Fatal().Err(err).Send()
			}
			allFiles = append(allFiles, files...)
		}
		var filteredFiles = make([]string, 0)
		for i := range allFiles {
			exists, err := dc.db.ExistsDownloadFile(ctx, sourceType, allFiles[i][strings.LastIndex(allFiles[i], `/`)+1:])
			if err != nil {
				log.Fatal().Err(err).Send()
			}
			if !exists {
				filteredFiles = append(filteredFiles, allFiles[i])
			}
		}
		var count = len(filteredFiles)
		log.Info().Int("total files", count).Send()
		for i := range filteredFiles {
			df, err := dc.downloadFile(ctx, filteredFiles[i], sourceType)
			if err != nil {
				if errors.Is(err, ErrExistsFile) {
					log.Err(err).Send()
					continue
				}
				log.Fatal().Err(err).Send()
			}
			log.Info().Int("complite", i+1).Int("total", count).Msg("downloaded files")
			ch <- df
		}
	}()
	return ch
}

func parseDirs(response *http.Response) ([]string, error) {
	defer response.Body.Close()
	root, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	var links = make([]string, 0)
	root.Find(`ul[id="directory-listing"]`).
		Find("a.clearfix").Each(func(_ int, li *goquery.Selection) {
		if href, ok := li.Attr("href"); ok && !strings.HasPrefix(href, response.Request.URL.Scheme) {
			link := url.URL{
				Scheme: response.Request.URL.Scheme,
				Host:   response.Request.URL.Host,
			}
			if strings.HasPrefix(href, "?") {
				link.RawQuery = strings.TrimPrefix(href, "?")
			} else {
				link.Path = href
			}
			links = append(links, link.String())
		}
	})
	return links, nil
}

package downloader

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"opendataaggregator/config"
	"opendataaggregator/models"
	"os"
	"path"
	"time"

	"github.com/rs/zerolog/log"
)

type PaginationData struct {
	Content []models.HotelData `json:"content"`
	Last    bool               `json:"last"`
}

func FetchHotels(ctx context.Context, client *http.Client) ([]models.HotelData, error) {
	var headers = http.Header{
		"Accept":          []string{"application/json, text/plain, */*"},
		"Accept-Encoding": []string{"gzip"},
		"Accept-Language": []string{"ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7,zh-TW;q=0.6,zh;q=0.5"},
		"Content-Type":    []string{"application/json"},
		"Origin":          []string{"https://knd.gov.ru"},
		"Referer":         []string{"https://knd.gov.ru/"},
		"User-Agent":      []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"},
	}
	// jar, _ := cookiejar.New(nil)
	const URL = `https://tor.knd.gov.ru/ext/search/simpleRegistryItems`
	// var client = http.Client{Jar: jar}
	var listHotels = make([]models.HotelData, 0, 10000)
	var currentPage = 0
fetchData:
	var params = map[string]any{
		"search": map[string]any{
			"search": []map[string]string{
				{
					"field":    "registryType.id",
					"operator": "eq",
					"value":    "63ef2fc7a445e900072d7e10",
				},
			},
		},
		"page": currentPage,
		"size": 1000,
	}
	payload, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, URL, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header = headers.Clone()
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	gr, err := gzip.NewReader(res.Body)
	if err != nil {
		return nil, err
	}
	var pd PaginationData
	if err := json.NewDecoder(gr).Decode(&pd); err != nil {
		return nil, err
	}
	gr.Close()
	res.Body.Close()
	listHotels = append(listHotels, pd.Content...)
	log.Info().Int("page", currentPage).Int("count", len(listHotels)).Send()
	if !pd.Last {
		currentPage++
		goto fetchData
	}
	return listHotels, nil
}

func saveToFile(filename string, data []byte) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

func (dc *DownloadClient) downloadHotels(ctx context.Context, sourceType string) (*config.DownloadedFile, error) {
	const SourceLink = "https://tor.knd.gov.ru"
	var filename = fmt.Sprintf("hotels_%s.json", time.Now().Format(time.DateOnly))
	var filepath = path.Join(dc.cnf.FS.DownloadFolder, sourceType, filename)
	if exists, err := dc.db.ExistsDownloadFile(ctx, sourceType, filename); err == nil && exists {
		return nil, ErrExistsFile
	}
	var hotels, err = FetchHotels(ctx, dc.client)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(hotels)
	if err != nil {
		return nil, err
	}
	if err := saveToFile(filepath, b); err != nil {
		return nil, err
	}
	hash := sha256.New()
	if _, err := hash.Write(b); err != nil {
		return nil, err
	}
	return &config.DownloadedFile{
		Filename:   filename,
		Filepath:   filepath,
		SourceType: sourceType,
		SourceLink: SourceLink,
		SHA265SUM:  hex.EncodeToString(hash.Sum(nil)),
	}, nil
}

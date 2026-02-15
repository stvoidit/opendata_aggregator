// Package downloader - скачивание исходников выгрузок
package downloader

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"opendataaggregator/config"
	"opendataaggregator/store"
	"os"
	"path"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
)

var (
	ErrHrefNotFound  = errors.New("href на файл не найден")  // ErrHrefNotFound - ОШИБКА - href на файл не найден
	ErrUnknownHost   = errors.New("unknown host")            // ErrUnknownHost - неизвестный хост для скачивания
	ErrExistsFile    = errors.New("download file is exists") // ErrExistsFile - файл уже скачан
	defaultTransport = &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:          0,
		DisableKeepAlives:     true,
		ResponseHeaderTimeout: time.Minute * 1}
)

// NewDownloadClient - ...
func NewDownloadClient(db *store.DB, cnf *config.Config) *DownloadClient {
	jar, _ := cookiejar.New(nil)
	return &DownloadClient{
		cnf: cnf,
		db:  db,
		client: &http.Client{
			Jar:       jar,
			Transport: defaultTransport}}
}

// DownloadClient - http client download
type DownloadClient struct {
	client *http.Client
	db     *store.DB
	cnf    *config.Config
}

// Close - CloseIdleConnections
func (dc *DownloadClient) Close() { dc.client.CloseIdleConnections() }

// Для скачивания с ftp.egrul.nalog.ru нельзя добавлять в транспорт оба сертификата, только какой то один.
// Поэтому делается сброс транспорта после окончания загрузки из этих источников для дефолтного состояния
func (dc *DownloadClient) resetTransport() { dc.client.Transport = defaultTransport }

// Do - обертка над http.Do с проверкой Content-Encoding
func (dc *DownloadClient) Do(request *http.Request) (*http.Response, error) {
	request.Header.Add("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 Edg/107.0.1418.52`)
	request.Header.Add("Accept-Encoding", "gzip")
	response, err := dc.client.Do(request)
	if err != nil {
		return nil, err
	}
	err = readContentEncoding(response)
	return response, err
}

// Download - ...
func (dc *DownloadClient) Download(ctx context.Context, sourcetype, link string) error {
	URL, err := url.Parse(link)
	if err != nil {
		return err
	}
	link = URL.String()
	hostname := URL.Hostname()
	log.Info().Str("hostname", hostname).Send()
	switch hostname {
	case "nalog.gov.ru", "www.nalog.gov.ru":
		df, err := dc.downloadNalogGovRu(ctx, link, sourcetype)
		if err != nil {
			return err
		}
		if err := dc.db.InsertDownloadedFileInfo(ctx, df); err != nil {
			return err
		}
	case "rosstat.gov.ru", "www.rosstat.gov.ru":
		df, err := dc.downloadRosstatGovRu(ctx, link, sourcetype)
		if err != nil {
			return err
		}
		if err := dc.db.InsertDownloadedFileInfo(ctx, df); err != nil {
			return err
		}
	case "rospatent.gov.ru", "www.rospatent.gov.ru":
		df, err := dc.downloadRospatentGovRu(ctx, link, sourcetype)
		if err != nil {
			return err
		}
		if err := dc.db.InsertDownloadedFileInfo(ctx, df); err != nil {
			return err
		}
	case "opendata.fssp.gov.ru", "www.opendata.fssp.gov.ru", "opendata.old.fssp.gov.ru":
		df, err := dc.downloadOpendataFsspGovRu(ctx, link, hostname, sourcetype)
		if err != nil {
			return err
		}
		if err := dc.db.InsertDownloadedFileInfo(ctx, df); err != nil {
			return err
		}
	case "fsa.gov.ru", "www.fsa.gov.ru":
		df, err := dc.downloadFsaGovRu(ctx, link, sourcetype)
		if err != nil {
			return err
		}
		if err := dc.db.InsertDownloadedFileInfo(ctx, df); err != nil {
			return err
		}
	case "zakupki.gov.ru", "www.zakupki.gov.ru":
		df, err := dc.downloadZakupkiGovRu(ctx, link, sourcetype)
		if err != nil {
			return err
		}
		if err := dc.db.InsertDownloadedFileInfo(ctx, df); err != nil {
			return err
		}
	case "ftp.egrul.nalog.ru", "www.ftp.egrul.nalog.ru":
		var crtname string
		switch URL.Query().Get("dir") {
		case "EGRUL_406", "EGRUL_407":
			crtname = "EGRUL"
		case "EGRIP_405", "EGRIP_406":
			crtname = "EGRIP"
		default:
			return errors.New("unknown dir ftp.egrul.nalog.ru")
		}
		cert, err := tls.LoadX509KeyPair(
			path.Join(dc.cnf.FS.EGRFolder, fmt.Sprintf("%s.crt", crtname)),
			path.Join(dc.cnf.FS.EGRFolder, fmt.Sprintf("%s.key", crtname)),
		)
		if err != nil {
			return fmt.Errorf("certficate load error: %w", err)
		}
		defer dc.resetTransport()
		dc.client.Transport = &http.Transport{
			MaxIdleConns:      0,
			DisableKeepAlives: true,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				Certificates:       []tls.Certificate{cert},
			}}
		docChan := dc.GetFilesList(ctx, link, sourcetype)
		if err != nil {
			return err
		}
		for df := range docChan {
			if err := dc.db.InsertDownloadedFileInfo(ctx, df); err != nil {
				return err
			}
		}
	case "hotels", "tor.knd.gov.ru":
		df, err := dc.downloadHotels(ctx, "hotels")
		if err != nil {
			return err
		}
		if err := dc.db.InsertDownloadedFileInfo(ctx, df); err != nil {
			return err
		}
	case "proverki.gov.ru":
		df, err := dc.DownloadUnscheduled(ctx, link, sourcetype)
		if err != nil {
			return err
		}
		if err := dc.db.InsertDownloadedFileInfo(ctx, df); err != nil {
			return err
		}
	default:
		return ErrUnknownHost
	}
	return nil
}

func (dc *DownloadClient) downloadFile(ctx context.Context, href string, sourceType string) (df *config.DownloadedFile, err error) {
	href = strings.TrimSpace(href)
	log.Info().Str("download href", href).Send()
	filename, err := parseFilenameFromURL(href)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, href, nil)
	if err != nil {
		return nil, err
	}
	response, err := dc.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status not 200: %s", href)
	}
	if contentDisposition := response.Header.Get("Content-Disposition"); len(contentDisposition) > 0 {
		if mt, p, err := mime.ParseMediaType(contentDisposition); err == nil && mt == "attachment" {
			if fn, ok := p["filename"]; ok {
				filename = fn
			}
		}
	}
	return dc.saveToFile(ctx, filename, response, sourceType)
}

func parseFilenameFromURL(link string) (string, error) {
	URL, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	return URL.Path[strings.LastIndex(URL.Path, "/")+1:], nil
}

func readContentEncoding(r *http.Response) error {
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		zr, err := gzip.NewReader(r.Body)
		if err != nil {
			return err
		}
		b, err := io.ReadAll(zr)
		if err != nil {
			return err
		}
		if err := r.Body.Close(); err != nil {
			return err
		}
		if err := zr.Close(); err != nil {
			return err
		}
		r.Body = io.NopCloser(bytes.NewReader(b))
	}
	return nil
}

func (dc *DownloadClient) saveToFile(ctx context.Context, filename string, response *http.Response, sourceType string) (*config.DownloadedFile, error) {
	if exists, err := dc.db.ExistsDownloadFile(ctx, sourceType, filename); err == nil && exists {
		return nil, ErrExistsFile
	}
	log.Info().Str("download file", filename).Send()
	filepath := path.Join(dc.cnf.FS.DownloadFolder, sourceType, filename)
	f, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// f := io.Discard
	hash := sha256.New()
	br := progressbar.NewOptions64(response.ContentLength,
		progressbar.OptionSetDescription("download:"),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowCount())
	if _, err := io.Copy(io.MultiWriter(f, br, hash), response.Body); err != nil {
		os.Remove(filepath) // в случае обрыва соединения нужно удалить недокачанный файл
		return nil, err
	}
	br.Exit()
	fmt.Println()
	return &config.DownloadedFile{
		Filename:   filename,
		Filepath:   filepath,
		SourceType: sourceType,
		SourceLink: response.Request.URL.String(),
		SHA265SUM:  hex.EncodeToString(hash.Sum(nil)),
	}, nil
}

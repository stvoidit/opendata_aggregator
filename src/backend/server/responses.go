package server

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
)

var json = jsoniter.ConfigFastest

const (
	headerContentType        = "Content-Type"
	headerCacheControl       = "Cache-Control"
	headerContentLength      = "Content-Length"
	headerContentDisposition = "Content-Disposition"
	noCacheStore             = `no-cache, no-store`
	applicationJSON          = "application/json"
	multipartFormData        = "multipart/form-data"
	octetStream              = "application/octet-stream"
)

// Jsonify - отправка json response
func Jsonify(w http.ResponseWriter, i any, code int) {
	w.Header().Add(headerContentType, applicationJSON)
	w.WriteHeader(code)
	switch I := i.(type) {
	default:
		json.NewEncoder(w).Encode(i)
	case string:
		w.Write([]byte(I))
	case []byte:
		w.Write(I)
	case io.WriterTo:
		I.WriteTo(w)
	case error:
		log.Err(I).Int("code", code).Msg("response")
		json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{Error: I.Error()})
	}
}

// JSONLoad - преобразование json в структуру
func JSONLoad(rc io.ReadCloser, i interface{}) error {
	defer rc.Close()
	return json.NewDecoder(rc).Decode(i)
}

// FileSystemSPA - обертка над fs.FileSystem для SPA приложения
func FileSystemSPA(dirname string) fs.FS { return &spaFS{os.DirFS(dirname)} }

type spaFS struct {
	hfs fs.FS
}

func (sfs spaFS) Open(name string) (fs.File, error) {
	if _, err := fs.Stat(sfs.hfs, name); err != nil && os.IsNotExist(err) {
		return sfs.hfs.Open("index.html")
	}
	return sfs.hfs.Open(name)
}

// WriteFile - запись файла в ответ
func WriteFile(w http.ResponseWriter, filename string, blob []byte) {
	w.Header().Add(headerContentType, octetStream)
	w.Header().Add(headerCacheControl, noCacheStore)
	w.Header().Add(headerContentLength, strconv.Itoa(len(blob)))
	w.Header().Add(headerContentDisposition, fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.WriteHeader(http.StatusOK)
	w.Write(blob)
}

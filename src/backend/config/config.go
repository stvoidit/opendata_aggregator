// Package config - конфигурация для приложения
package config

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
)

var (
	CommitHash = "devVersion"
	VersionTag = "devTag"
	Branch     = "devBranch"
)

type Megaplan struct {
	Domain       string `toml:"domain"`
	UUID         string `toml:"uuid"`
	Secret       string `toml:"secret"`
	ProviderHost string `toml:"provider_api"`
}

// Config - конфиг приложения
type Config struct {
	Debug    bool     `toml:"debug"`
	LogLevel int8     `toml:"loglevel"`
	Megaplan Megaplan `toml:"mp"`
	DB       struct {
		Host     string `toml:"host"`
		Port     string `toml:"port"`
		Login    string `toml:"user"`
		Password string `toml:"password"`
		DBName   string `toml:"dbname"`
	} `toml:"db"`
	Srv struct {
		Port   string `toml:"port"`
		Domain string `toml:"domain"`
	} `toml:"srv"`
	FS struct {
		DownloadFolder string `toml:"download"`
		EGRFolder      string `toml:"egrfolder"`
	} `toml:"fs"`
	Sources map[string]string `toml:"sources"` // пример перечисления: kgn = "https://www.nalog.gov.ru/opendata/7707329152-kgn/"
}

// func genDefaultConfig() {
// 	log.Info().Msg("GENERATE DEFAULT CONFIG: config.toml")
// 	f, err := os.Create("config.toml")
// 	if err != nil {
// 		log.Fatal().Err(err).Send()
// 	}
// 	defer f.Close()
// 	var cnf = Config{
// 		LogLevel: 0,
// 		DB: {
// 			Host:     "localhost",
// 			Port:     "5432",
// 			Login:    "postgres",
// 			Password: "postgres",
// 			DBName:   "opendata_aggregator",
// 		},
// 		Srv: serverSection{
// 			Port:   ":8080",
// 			Domain: "",
// 		},
// 		FS: filesystemSection{
// 			DownloadFolder: "OPENDATA_SOURCES",
// 			EGRFolder:      "egr_certs",
// 		},
// 	}
// 	if err := toml.NewEncoder(f).Encode(&cnf); err != nil {
// 		log.Fatal().Err(err).Send()
// 	}
// }

func (c Config) String() string {
	var sb strings.Builder
	toml.NewEncoder(&sb).Encode(&c)
	return sb.String()
}

// LoadConfig - ...
func LoadConfig(r io.Reader) (cnf *Config, err error) {
	cnf = new(Config)
	_, err = toml.NewDecoder(r).Decode(cnf)
	if value, ok := os.LookupEnv("DEBUG"); ok && value == "1" {
		cnf.Debug = true
	}
	if value, ok := os.LookupEnv("SERVER_DOMAIN"); ok {
		cnf.Srv.Domain = value
	}
	if cnf.FS.DownloadFolder != "" {
		os.MkdirAll(cnf.FS.DownloadFolder, os.ModePerm)
	} else {
		cnf.FS.DownloadFolder = "DOWNLOADS"
	}
	for k := range cnf.Sources {
		os.MkdirAll(path.Join(cnf.FS.DownloadFolder, k), os.ModePerm)
	}
	return
}

// LoadConfigFromFile - ...
func LoadConfigFromFile(filename string) (cnf *Config, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return LoadConfig(file)
}

// DownloadedFile - загруженный файл
type DownloadedFile struct {
	Filename   string
	Filepath   string
	SHA265SUM  string
	SourceType string
	SourceLink string
}

// Remove - удалить файл
func (df DownloadedFile) Remove() error {
	return os.Remove(df.Filepath)
}

// Sha256Sum - SHA256 файла
func (df *DownloadedFile) Sha256Sum() error {
	f, err := os.Open(df.Filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return err
	}
	df.SHA265SUM = hex.EncodeToString(hash.Sum(nil))
	return nil

}

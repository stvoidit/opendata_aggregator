package csvparser

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"opendataaggregator/models"
	"os"
	"testing"
)

func BenchmarkDecoder_old_version(b *testing.B) {
	filename := `/home/stvoid/rss_20230901-20230930.csv`
	f, err := os.Open(filename)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.Comma = ','
	r.LazyQuotes = true
	r.ReuseRecord = false
	rows, err := r.ReadAll()
	if err != nil {
		b.Fatal(err)
	}
	headers := getHeaders(rows[0])
	rows = rows[1:]
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		index := rand.Intn(len(rows) - 1)
		b.StartTimer()
		var doc models.Росаккредитация
		err = decode_old_version(headers, rows[index], &doc)
		if err != nil {
			b.Fatal(err)
		}
		fmt.Fprint(io.Discard, doc)
	}

}

func BenchmarkDecoder_new_version(b *testing.B) {
	filename := `/home/stvoid/rss_20230901-20230930.csv`
	f, err := os.Open(filename)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.Comma = ','
	r.LazyQuotes = true
	r.ReuseRecord = false
	rows, err := r.ReadAll()
	if err != nil {
		b.Fatal(err)
	}
	headers := getHeaders(rows[0])
	rows = rows[1:]
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		index := rand.Intn(len(rows) - 1)
		b.StartTimer()
		doc, err := decode_new_version[models.Росаккредитация](headers, rows[index])
		if err != nil {
			b.Fatal(err)
		}
		fmt.Fprint(io.Discard, doc)
	}
}

func decode_old_version(headers CSVHeaders, row []string, i any) error {
	var m = make(map[string]string, len(headers))
	for k, v := range headers {
		m[k] = row[v]
	}
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, i); err != nil {
		return err
	}
	return nil
}

func decode_new_version[T any](headers CSVHeaders, row []string) (T, error) {
	var m = make(map[string]string, len(headers))
	for k, v := range headers {
		m[k] = row[v]
	}
	var v T
	b, err := json.Marshal(m)
	if err != nil {
		return v, err
	}
	err = json.Unmarshal(b, &v)
	return v, err
}

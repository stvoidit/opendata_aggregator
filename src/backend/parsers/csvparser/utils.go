package csvparser

import (
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigDefault

type CSVHeaders map[string]int

func getHeaders(row []string) CSVHeaders {
	var m = make(CSVHeaders, len(row))
	for i, v := range row {
		m[v] = i
	}
	return m
}

func decode[T any](headers CSVHeaders, row []string) (T, error) {
	var m = make(map[string]string, len(headers))
	for k, v := range headers {
		m[k] = row[v]
	}
	var v T
	b, err := json.Marshal(&m)
	if err != nil {
		return v, err
	}
	err = json.Unmarshal(b, &v)
	return v, err
}

func parseBoolean(s string) bool {
	n, _ := strconv.ParseBool(s)
	return n
}

func parseFloat(s string) float64 {
	n, _ := strconv.ParseFloat(s, 64)
	return n
}

func cleanString(s string) string {
	var replacerCleanString = strings.NewReplacer(
		"NULL", "",
		"null", "",
	)
	return strings.TrimSpace(replacerCleanString.Replace(s))
}

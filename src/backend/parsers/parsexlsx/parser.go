// Package parsexlsx - парсинг файлов xlsx
package parsexlsx

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/shakinm/xlsReader/xls"
	"github.com/xuri/excelize/v2"
)

func ParseXLS[T any](rc io.ReadSeekCloser) ([]T, error) {
	defer rc.Close()
	workbook, err := xls.OpenReader(rc)
	if err != nil {
		return nil, err
	}
	sheet, err := workbook.GetSheet(0)
	if err != nil {
		return nil, err
	}
	rows := sheet.GetRows()
	if len(rows) == 0 {
		return nil, errors.New("empty table XLS")
	}
	var headers = make(map[string]int)
	var data = make([]T, 0)
	for i, col := range rows[0].GetCols() {
		headers[col.GetString()] = i
	}
	for _, row := range rows[0:] {
		cols := row.GetCols()
		var strRow = make([]string, 0, len(cols))
		for _, col := range cols {
			strRow = append(strRow, col.GetString())
		}
		var v T
		if err := decode(headers, strRow, &v); err != nil {
			return nil, err
		}
		data = append(data, v)
	}
	return data, nil
}

// ParseXLSX - парсинг выгрузки в xlsx
func ParseXLSX[T any](rc io.ReadCloser) ([]T, error) {
	defer rc.Close()
	f, err := excelize.OpenReader(rc)
	if err != nil {
		return nil, err
	}
	var sheetname string
	for _, sheet := range f.WorkBook.Sheets.Sheet {
		sheetname = sheet.Name
		break
	}
	rows, err := f.GetRows(sheetname)
	if err != nil {
		return nil, err
	}
	var headers map[string]int
	var data = make([]T, 0)
	for i, row := range rows {
		if i == 0 {
			headers = getHeaders(row)
			continue
		}
		var v T
		if err := decode(headers, row, &v); err != nil {
			return nil, err
		}
		data = append(data, v)
	}
	return data, nil
}

func getHeaders(row []string) map[string]int {
	var m = make(map[string]int, len(row))
	for i, v := range row {
		m[v] = i
	}
	return m
}

func decode(headers map[string]int, row []string, i interface{}) error {
	var m = make(map[string]string, len(headers))
	for k, v := range headers {
		m[k] = row[v]
	}
	b, _ := json.Marshal(&m)
	if err := json.Unmarshal(b, i); err != nil {
		return err
	}
	return nil
}

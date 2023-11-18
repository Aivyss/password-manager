package csv

import (
	"encoding/csv"
	"fmt"
	"github.com/aivyss/password-manager/pwmErr"
	"reflect"
	"strings"
)

func StructsToCsv[T any](data []T) (*string, error) {
	csvLines := CreateCsvLines(data)

	// convert csv
	var csvBuilder strings.Builder
	writer := csv.NewWriter(&csvBuilder)
	err := writer.WriteAll(csvLines)
	if err != nil {
		return nil, pwmErr.ConvertCsv
	}

	result := csvBuilder.String()
	return &result, nil
}

func CreateCsvLines[T any](data []T) [][]string {
	var csvLines [][]string

	// write header
	var header []string
	var tmp T
	typeOf := reflect.TypeOf(tmp)
	for i := 0; i < typeOf.NumField(); i++ {
		tagHeader := typeOf.Field(i).Tag.Get("csv")
		header = append(header, tagHeader)
	}
	csvLines = append(csvLines, header)

	// write data
	for _, datum := range data {
		valueOf := reflect.ValueOf(datum)
		var line []string
		for i := 0; i < typeOf.NumField(); i++ {
			line = append(line, fmt.Sprintf("%v", valueOf.Field(i).Interface()))
		}

		csvLines = append(csvLines, line)
	}

	return csvLines
}

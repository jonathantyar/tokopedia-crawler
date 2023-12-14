package export

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/shopspring/decimal"
)

func GenerateCsv[T any](filename string, data []T) (string, error) {
	genFilename := fmt.Sprintf("%v.csv", filename)
	csvFile, err := os.Create(genFilename)
	if err != nil {
		return "", err
	}
	csvWriter := csv.NewWriter(csvFile)

	csvData := [][]string{}

	header := []string{}
	writeHeaderCsv[T](&header)
	csvData = append(csvData, header)
	for _, d := range data {
		err := writeRowCsv(&csvData, d)
		if err != nil {
			return "", err
		}
	}

	csvWriter.WriteAll(csvData)

	csvWriter.Flush()
	csvFile.Close()

	return genFilename, nil
}

func writeHeaderCsv[T any](header *[]string) error {
	var headers T
	v := reflect.ValueOf(headers)
	typeOf := v.Type()

	for i := 0; i < v.NumField(); i++ {
		col := ExportTags{}
		col.Convert(typeOf.Field(i).Tag.Get("export"))

		if col.Name != "" {
			*header = append(*header, col.Name)
		}
	}

	return nil
}

func writeRowCsv[T any](rows *[][]string, data T) error {
	v := reflect.ValueOf(data)
	typeOf := v.Type()

	colI := 1
	row := []string{}
	for i := 0; i < v.NumField(); i++ {
		col := ExportTags{}
		col.Convert(typeOf.Field(i).Tag.Get("export"))

		if col.Name != "" {
			fmt.Println(v.Field(i).Type().Name(), i)
			switch col.Type {
			case typeRFC3339:
				loc, _ := time.LoadLocation("Asia/Jakarta")
				value := fmt.Sprintf("%s", v.Field(i).Interface())
				formatedDate, err := time.Parse(time.RFC3339, value)
				if err != nil {
					return err
				}

				row = append(row, formatedDate.In(loc).String())
			default:
				switch v.Field(i).Type().Name() {
				case "uint64":
					row = append(row, fmt.Sprintf("%v", v.Field(i).Interface().(uint64)))
				case "int64":
					row = append(row, fmt.Sprintf("%v", v.Field(i).Interface().(int64)))
				case "uint32":
					row = append(row, fmt.Sprintf("%v", v.Field(i).Interface().(uint32)))
				case "int32":
					row = append(row, fmt.Sprintf("%v", v.Field(i).Interface().(int64)))
				case "float64":
					row = append(row, fmt.Sprintf("%v", v.Field(i).Interface().(float64)))
				case "decimal.Decimal":
					row = append(row, fmt.Sprintf("%v", v.Field(i).Interface().(decimal.Decimal).InexactFloat64()))
				case "float32":
					row = append(row, fmt.Sprintf("%v", v.Field(i).Interface().(float32)))
				case "string":
					row = append(row, fmt.Sprintf("%v", v.Field(i).Interface().(string)))
				}
			}

			colI++
		}
	}
	*rows = append(*rows, row)

	return nil
}

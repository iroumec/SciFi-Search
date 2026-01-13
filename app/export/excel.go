package export

import (
	"io"
	"reflect"

	"github.com/xuri/excelize/v2"
)

type ExcelExporter[T any] struct{}

func (ExcelExporter[T]) ContentType() string {
	return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
}

func (ExcelExporter[T]) FileName() string {
	return "documents.xlsx"
}

func (ExcelExporter[T]) Export(w io.Writer, data []T) error {
	fields := getFields[T]()

	f := excelize.NewFile()
	sheet := "Sheet1"

	// Header
	for i, field := range fields {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, field.Name)
	}

	// Rows
	for r, item := range data {
		val := reflect.ValueOf(item)
		for c, field := range fields {
			cell, _ := excelize.CoordinatesToCellName(c+1, r+2)
			f.SetCellValue(sheet, cell, valueToString(val.Field(field.Index)))
		}
	}

	return f.Write(w)
}

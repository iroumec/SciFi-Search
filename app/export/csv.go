package export

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"encoding/csv"
	"io"
	"reflect"
)

// ------------------------------------------------------------------------------------------------
// Structures
// ------------------------------------------------------------------------------------------------

type CSVExporter[T any] struct{}

// ------------------------------------------------------------------------------------------------

func (CSVExporter[T]) ContentType() string {
	return "text/csv"
}

// ------------------------------------------------------------------------------------------------

func (CSVExporter[T]) FileName() string {
	return "documents.csv"
}

// ------------------------------------------------------------------------------------------------

func (CSVExporter[T]) Export(w io.Writer, data []T) error {
	fields := getFields[T]()
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Header
	header := make([]string, len(fields))
	for i, f := range fields {
		header[i] = f.Name
	}
	writer.Write(header)

	for _, item := range data {
		val := reflect.ValueOf(item)
		row := make([]string, len(fields))

		for i, f := range fields {
			row[i] = valueToString(val.Field(f.Index))
		}
		writer.Write(row)
	}

	return writer.Error()
}

// ------------------------------------------------------------------------------------------------

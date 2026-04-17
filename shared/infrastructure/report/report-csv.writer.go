package report

import (
	"encoding/csv"
	"fmt"
	"os"
)

type CSVWriter struct {
	file   *os.File
	writer *csv.Writer
}

func NewCSVWriter(path string) (*CSVWriter, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return &CSVWriter{
		file:   file,
		writer: csv.NewWriter(file),
	}, nil
}

func (w *CSVWriter) WriteHeader(headers []string) error {
	return w.writer.Write(headers)
}

func (w *CSVWriter) WriteRow(values []any) error {
	record := make([]string, len(values))

	for i, v := range values {
		record[i] = fmt.Sprint(v)
	}

	return w.writer.Write(record)
}

func (w *CSVWriter) Close() error {
	w.writer.Flush()

	if err := w.writer.Error(); err != nil {
		return err
	}

	return w.file.Close()
}


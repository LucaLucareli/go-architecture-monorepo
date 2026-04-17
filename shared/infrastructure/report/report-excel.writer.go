package report

import (
	"github.com/xuri/excelize/v2"
)

type ExcelWriter struct {
	excelFile    *excelize.File
	streamWriter *excelize.StreamWriter
	sheetName    string
	currentRow   int
	filePath     string
}

func NewExcelWriter(filePath string) (*ExcelWriter, error) {
	excelFile := excelize.NewFile()
	sheetName := "Relatorio"

	index := excelFile.GetActiveSheetIndex()
	excelFile.SetSheetName(excelFile.GetSheetName(index), sheetName)

	streamWriter, err := excelFile.NewStreamWriter(sheetName)
	if err != nil {
		return nil, err
	}

	return &ExcelWriter{
		excelFile:    excelFile,
		streamWriter: streamWriter,
		sheetName:    sheetName,
		currentRow:   1,
		filePath:     filePath,
	}, nil
}

func (excelWriter *ExcelWriter) WriteHeader(headers []string) error {
	headerRow := make([]any, len(headers))
	for i, header := range headers {
		headerRow[i] = header
	}

	cell, _ := excelize.CoordinatesToCellName(1, excelWriter.currentRow)
	if err := excelWriter.streamWriter.SetRow(cell, headerRow); err != nil {
		return err
	}

	excelWriter.currentRow++
	return nil
}

func (excelWriter *ExcelWriter) WriteRow(rowValues []any) error {
	cell, _ := excelize.CoordinatesToCellName(1, excelWriter.currentRow)
	if err := excelWriter.streamWriter.SetRow(cell, rowValues); err != nil {
		return err
	}

	excelWriter.currentRow++
	return nil
}

func (excelWriter *ExcelWriter) Close() error {
	if err := excelWriter.streamWriter.Flush(); err != nil {
		return err
	}
	return excelWriter.excelFile.SaveAs(excelWriter.filePath)
}

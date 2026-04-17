package report

import (
	"errors"
	"shared/application/interfaces"
	"shared/domain/types"
)

func NewWriter(t types.ReportType, path string) (interfaces.RowWriterInterface, error) {
	switch t {
	case types.Excel:
		return NewExcelWriter(path)
	case types.CSV:
		return NewCSVWriter(path)
	default:
		return nil, errors.New("tipo de relatório inválido")
	}
}

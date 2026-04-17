package interfaces

type RowWriterInterface interface {
	WriteHeader(headers []string) error

	WriteRow(values []any) error

	Close() error
}

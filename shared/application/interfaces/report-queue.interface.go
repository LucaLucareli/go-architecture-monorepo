package interfaces

import (
	"context"
	"shared/infrastructure/queue/payloads"
)

type ReportQueue interface {
	Enqueue(ctx context.Context, payload payloads.GenerateReportPayload) error
}

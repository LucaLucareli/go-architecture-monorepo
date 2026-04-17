package services

import (
	"context"
	"employee-api/modules/report/dto/io"
	"shared/domain/types"
	"shared/infrastructure/queue/enqueue"
	"shared/infrastructure/queue/payloads"

	"github.com/hibiken/asynq"
)

type RequestReportService struct {
	asynqClient *asynq.Client
}

func NewRequestReportService(asynqClient *asynq.Client) *RequestReportService {
	return &RequestReportService{asynqClient: asynqClient}
}

func (s *RequestReportService) Execute(
	ctx context.Context,
	input io.RequestReportInputDto,
) error {

	return enqueue.EnqueueGenerateReport(
		s.asynqClient,
		payloads.GenerateReportPayload{
			ReportID:   input.ReportID,
			UserID:     input.UserID,
			ReportType: types.ReportType(input.ReportType),
		},
	)
}

package consumer

import (
	"employee-api/modules/report/services"
	"shared/domain/types"
	"shared/infrastructure/queue"
	"shared/infrastructure/queue/payloads"
)

type ReportConsumer struct {
	processExcelReportService *services.ProcessExcelReportService
}

func NewReportConsumer(processExcelReportService *services.ProcessExcelReportService) *ReportConsumer {
	return &ReportConsumer{processExcelReportService: processExcelReportService}
}

func (c *ReportConsumer) GetHandlers() []queue.TaskHandler {
	return []queue.TaskHandler{
		{
			Task:    queue.TaskGenerateReport,
			Handler: c.processExcelReportService.Execute,
		},
	}
}

func (c *ReportConsumer) GetPayload(taskType string) interface{} {
	switch taskType {
	case queue.TaskGenerateReport:
		return payloads.GenerateReportPayload{}
	default:
		return nil
	}
}

func (c *ReportConsumer) GetReportType(taskType string) types.ReportType {
	switch taskType {
	case queue.TaskGenerateReport:
		return types.Excel
	default:
		return 0
	}
}

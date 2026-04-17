package payloads

import "shared/domain/types"

type GenerateReportPayload struct {
	ReportID   int32            `json:"report_id"`
	UserID     string           `json:"user_id"`
	ReportType types.ReportType `json:"reportType"`
}

package request

type RequestReportRequest struct {
	UserID     string `json:"userId" validate:"required,min=1"`
	ReportType int32  `json:"reportType" validate:"required,min=1"`
	ReportID   int32  `json:"reportId" validate:"required,min=1"`
}

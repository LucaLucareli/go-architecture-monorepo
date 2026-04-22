package io

type RequestReportInputDto struct {
	UserID     string `json:"userId"`
	ReportType int32  `json:"reportType"`
	ReportID   int32  `json:"reportId"`
}

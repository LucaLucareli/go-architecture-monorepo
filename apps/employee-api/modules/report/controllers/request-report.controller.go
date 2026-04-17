package controllers

import (
	"employee-api/modules/report/dto/io"
	"employee-api/modules/report/dto/request"
	"employee-api/modules/report/services"
	"shared/application/interfaces"
	exceptionfactory "shared/pkg/validation/exception-factory"

	"github.com/labstack/echo/v4"
)

type RequestReportController struct {
	requestReportService *services.RequestReportService
}

func NewRequestReportController(requestReportService *services.RequestReportService) *RequestReportController {
	return &RequestReportController{requestReportService: requestReportService}
}

func (ctrl *RequestReportController) Handle(c echo.Context) error {
	var req request.RequestReportRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(&req); err != nil {
		return exceptionfactory.CustomExceptionFactory(err)
	}

	err := ctrl.requestReportService.Execute(
		c.Request().Context(),
		io.RequestReportInputDto{
			ReportType: req.ReportType,
			ReportID:   req.ReportID,
			UserID:     req.UserID,
		},
	)

	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	interfaces.Set(c, interfaces.ResponseInterface[any]{
		Message: "Relatório solicitado com sucesso",
	})

	return nil
}

package report

import (
	"employee-api/modules/report/consumer"
	"employee-api/modules/report/controllers"
	"employee-api/modules/report/services"
	"shared"
	"shared/domain/enums"
	"shared/infrastructure/queue"
	"shared/pkg/middlewares"

	"github.com/labstack/echo/v4"
)

type ReportModule struct {
	requestReportController *controllers.RequestReportController
	reportConsumer          *consumer.ReportConsumer
}

func NewReportModule(appState *shared.AppState) *ReportModule {
	requestReportService := services.NewRequestReportService(appState.AsynqClient)
	requestReportController := controllers.NewRequestReportController(requestReportService)

	processExcelReportService := services.NewProcessExcelReportService(appState.UserRepo)
	reportConsumer := consumer.NewReportConsumer(processExcelReportService)

	return &ReportModule{
		requestReportController: requestReportController,
		reportConsumer:          reportConsumer,
	}
}

func (m *ReportModule) RegisterRoutes(e *echo.Echo, appState *shared.AppState) {
	reportGroup := e.Group("/reports")

	reportGroup.Use(middlewares.RequireAccess(
		appState.AuthService,
		enums.AccessGroupAdmin,
		enums.AccessGroupSuperAdmin,
	))

	reportGroup.POST("/request", m.requestReportController.Handle)
}

func (m *ReportModule) GetHandlers() []queue.TaskHandler {
	return m.reportConsumer.GetHandlers()
}

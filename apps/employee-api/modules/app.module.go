package modules

import (
	"employee-api/modules/report"
	"employee-api/modules/users"
	"shared"

	"github.com/labstack/echo/v4"
	"github.com/quantumsheep/plouf"
)

type AppModule struct {
	plouf.Module
	UserModule   *users.UserModule
	ReportModule *report.ReportModule
}

func NewAppModule() *AppModule {
	return &AppModule{
		UserModule:   &users.UserModule{},
		ReportModule: &report.ReportModule{},
	}
}

func (m *AppModule) RegisterAllRoutes(e *echo.Echo, state *shared.AppState) {
	api := e.Group("/api")

	m.UserModule.RegisterRoutes(api, state)
	m.ReportModule.RegisterRoutes(api, state)
}

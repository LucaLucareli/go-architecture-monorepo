package users

import (
	"employee-api/modules/users/controllers"
	"employee-api/modules/users/services"
	"shared"
	"shared/domain/enums"
	"shared/pkg/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/quantumsheep/plouf"
)

type UserModule struct {
	plouf.Module
}

func (m *UserModule) RegisterRoutes(e *echo.Group, appState *shared.AppState) {

	findUserByIdService := services.NewFindUserByIdService(appState.UserRepo())
	findUserByIdController := controllers.NewFindUserByIdController(findUserByIdService)

	userGroup := e.Group("/users")

	userGroup.Use(middlewares.RequireAccess(
		appState.AuthService(),
		enums.AccessGroupEmployee,
		enums.AccessGroupAdmin,
		enums.AccessGroupSuperAdmin,
	))

	userGroup.GET("/:id", findUserByIdController.Handle)
}

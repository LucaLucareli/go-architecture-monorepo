package auth

import (
	"auth-api/modules/auth/controllers"
	"auth-api/modules/auth/services"
	"shared"

	"github.com/labstack/echo/v4"
	"github.com/quantumsheep/plouf"
)

type AuthModule struct {
	plouf.Module
}

func (m *AuthModule) RegisterRoutes(e *echo.Group, state *shared.AppState) {
	loginService := services.NewLoginService(state.AuthService())
	loginController := controllers.NewLoginController(loginService)

	refreshTokenService := services.NewRefreshTokenService(state.AuthService())
	refreshTokenController := controllers.NewRefreshTokenController(refreshTokenService)

	e.POST("/login", loginController.Handle)
	e.POST("/refresh-token", refreshTokenController.Handle)
}

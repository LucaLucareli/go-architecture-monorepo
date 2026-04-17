package controllers

import (
	"auth-api/modules/auth/dto/io"
	"auth-api/modules/auth/dto/request"
	"auth-api/modules/auth/services"
	"shared/application/interfaces"
	exceptionfactory "shared/pkg/validation/exception-factory"

	"github.com/labstack/echo/v4"
)

type RefreshTokenController struct {
	refreshTokenService *services.RefreshTokenService
}

func NewRefreshTokenController(refreshTokenService *services.RefreshTokenService) *RefreshTokenController {
	return &RefreshTokenController{refreshTokenService: refreshTokenService}
}

func (ctrl *RefreshTokenController) Handle(c echo.Context) error {
	var req request.RefreshTokenRequestDTO
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(&req); err != nil {
		return exceptionfactory.CustomExceptionFactory(err)
	}

	authResponse, err := ctrl.refreshTokenService.Execute(
		c.Request().Context(),
		io.RefreshTokenInputDTO{
			RefreshToken: req.RefreshToken,
		},
	)

	if err != nil {
		return echo.NewHTTPError(401, err.Error())
	}

	interfaces.Set(c, interfaces.ResponseInterface[*io.RefreshTokenOutputDTO]{
		Message: "Token atualizado com sucesso",
		Result:  authResponse,
	})

	return nil
}

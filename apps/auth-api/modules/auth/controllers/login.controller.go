package controllers

import (
	"auth-api/modules/auth/dto/io"
	"auth-api/modules/auth/dto/request"
	"auth-api/modules/auth/services"
	"shared/application/interfaces"
	exceptionfactory "shared/pkg/validation/exception-factory"

	"github.com/labstack/echo/v4"
)

type LoginController struct {
	loginService *services.LoginService
}

func NewLoginController(loginService *services.LoginService) *LoginController {
	return &LoginController{loginService: loginService}
}

func (ctrl *LoginController) Handle(c echo.Context) error {
	var req request.LoginRequestDTO
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(&req); err != nil {
		return exceptionfactory.CustomExceptionFactory(err)
	}

	authResponse, err := ctrl.loginService.Execute(
		c.Request().Context(),
		io.LoginInputDTO{
			Document: req.Document,
			Password: req.Password,
		},
	)

	if err != nil {
		return echo.NewHTTPError(401, err.Error())
	}

	interfaces.Set(c, interfaces.ResponseInterface[*io.LoginOutputDTO]{
		Message: "Login realizado com sucesso",
		Result:  authResponse,
	})

	return nil
}

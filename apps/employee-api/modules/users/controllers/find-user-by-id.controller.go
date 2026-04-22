package controllers

import (
	"employee-api/modules/users/services"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type FindUserByIdController struct {
	userService *services.FindUserByIdService
}

func NewFindUserByIdController(s *services.FindUserByIdService) *FindUserByIdController {
	return &FindUserByIdController{
		userService: s,
	}
}

func (ctrl *FindUserByIdController) Handle(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ID inválido")
	}

	user, err := ctrl.userService.Execute(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Usuário não encontrado")
	}

	return c.JSON(http.StatusOK, user)
}

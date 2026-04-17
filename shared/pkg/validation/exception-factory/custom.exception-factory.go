package exceptionfactory

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func CustomExceptionFactory(err error) *echo.HTTPError {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok || len(validationErrors) == 0 {
		return nil
	}

	errors := map[string][]string{}

	for _, e := range validationErrors {
		field := e.Field()
		errors[field] = append(errors[field], parseMessage(e))
	}

	return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
		"message": "Erro na validação dos dados",
		"errors":  errors,
	})
}

func parseMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "campo obrigatório"
	case "min":
		return "tamanho mínimo não atingido"
	case "max":
		return "tamanho máximo excedido"
	case "document":
		return "documento inválido"
	default:
		return "valor inválido"
	}
}

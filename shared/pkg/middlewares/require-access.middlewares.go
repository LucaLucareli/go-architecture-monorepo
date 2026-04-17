package middlewares

import (
	"net/http"
	"slices"
	"strings"

	"shared/application/auth"
	"shared/domain/enums"

	"github.com/labstack/echo/v4"
)

func RequireAccess(
	authService *auth.AuthService,
	allowedGroups ...enums.AccessGroupEnum,
) func(echo.HandlerFunc) echo.HandlerFunc {

	return func(handler echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "É necessário autenticar-se",
				})
			}

			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "Token inválido",
				})
			}

			token := strings.TrimPrefix(authHeader, bearerPrefix)

			claims, err := authService.ValidateAccessToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "Token inválido ou expirado",
				})
			}

			accessGroups := make([]enums.AccessGroupEnum, len(claims.AccessGroups))
			for i, g := range claims.AccessGroups {
				accessGroups[i] = enums.AccessGroupEnum(g)
			}

			for _, g := range accessGroups {
				if slices.Contains(allowedGroups, g) {
					return handler(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]string{
				"message": "Permissão insuficiente",
			})
		}
	}
}

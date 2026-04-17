package validation

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ValidationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		if req.Header.Get(echo.HeaderContentType) == echo.MIMEApplicationJSON {
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "JSON inválido")
			}

			decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
			decoder.DisallowUnknownFields()

			var tmp map[string]interface{}
			if err := decoder.Decode(&tmp); err != nil {
				return echo.NewHTTPError(
					http.StatusBadRequest,
					"Campos inválidos ou não permitidos",
				)
			}

			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		return next(c)
	}
}

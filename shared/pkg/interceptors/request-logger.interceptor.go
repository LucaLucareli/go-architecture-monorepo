package interceptors

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func RequestLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		start := time.Now()

		res.After(func() {
			latency := time.Since(start)
			status := res.Status

			event := log.Info()
			if status >= 400 {
				event = log.Warn()
			}
			if status >= 500 {
				event = log.Error()
			}

			event.
				Str("method", req.Method).
				Str("path", req.URL.Path).
				Int("status", status).
				Dur("latency", latency).
				Str("ip", c.RealIP()).
				Msg("HTTP request")
		})

		return next(c)
	}
}

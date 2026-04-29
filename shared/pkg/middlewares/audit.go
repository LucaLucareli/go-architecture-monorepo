package middlewares

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type AuditEntry struct {
	Timestamp  time.Time
	Method     string
	Path       string
	IP         string
	Status     int
	Latency    time.Duration
	UserAgent  string
	RequestID  string
}

var auditBuffer = make(chan AuditEntry, 1000)

func init() {
	go auditWorker()
}

func auditWorker() {
	for entry := range auditBuffer {
		// Aqui você poderia salvar no banco de dados.
		// Por agora, usaremos o logger estruturado para alta performance.
		log.Info().
			Str("method", entry.Method).
			Str("path", entry.Path).
			Int("status", entry.Status).
			Str("latency", entry.Latency.String()).
			Str("ip", entry.IP).
			Str("request_id", entry.RequestID).
			Msg("Audit Log")
	}
}

func AsyncAuditMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			entry := AuditEntry{
				Timestamp: start,
				Method:    c.Request().Method,
				Path:      c.Path(),
				IP:        c.RealIP(),
				Status:    c.Response().Status,
				Latency:   time.Since(start),
				UserAgent: c.Request().UserAgent(),
				RequestID: c.Response().Header().Get(echo.HeaderXRequestID),
			}

			select {
			case auditBuffer <- entry:
			default:
				log.Warn().Msg("Audit buffer full, dropping entry")
			}

			return err
		}
	}
}

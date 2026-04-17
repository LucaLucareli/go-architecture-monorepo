package interceptors

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

func TaskLogger(next asynq.HandlerFunc) asynq.HandlerFunc {
	return func(ctx context.Context, task *asynq.Task) error {
		start := time.Now()

		err := next(ctx, task)

		latency := time.Since(start)

		event := log.Info()
		if err != nil {
			event = log.Warn()
		}

		event.
			Str("task", task.Type()).
			Dur("latency", latency).
			Bytes("payload", task.Payload()).
			Msg("Asynq task")

		return err
	}
}

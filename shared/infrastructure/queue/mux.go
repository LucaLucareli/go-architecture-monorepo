package queue

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/hibiken/asynq"
)

type TaskHandler struct {
	Task    string
	Handler asynq.HandlerFunc
}

type QueueLifecycleOptions struct {
	QueueName string
	Inspector *asynq.Inspector
}

func NewMux(
	opts QueueLifecycleOptions,
	handlers ...TaskHandler,
) *asynq.ServeMux {
	mux := asynq.NewServeMux()

	for _, h := range handlers {
		mux.Handle(
			h.Task,
			WithLifecycle(h.Handler, opts),
		)
	}

	return mux
}

func WithLifecycle(
	handler asynq.HandlerFunc,
	opts QueueLifecycleOptions,
) asynq.HandlerFunc {
	return func(ctx context.Context, task *asynq.Task) error {
		start := time.Now()
		logger := taskLogger(task, opts)

		logTaskStart(task, opts)

		if err := handler(ctx, task); err != nil {
			logger.Error().
				Err(err).
				Dur("duration", time.Since(start)).
				Msg("task execution failed")

			return err
		}

		logger.Info().
			Dur("duration", time.Since(start)).
			Msg("task execution completed")

		return nil
	}
}

func logTaskStart(task *asynq.Task, opts QueueLifecycleOptions) {
	logger := log.With().
		Str("task_type", task.Type()).
		Str("queue", opts.QueueName).
		Logger()

	event := logger.Info()

	if opts.Inspector != nil {
		if info, err := opts.Inspector.GetQueueInfo(opts.QueueName); err == nil {
			event = event.
				Int("queue_active", info.Active).
				Int("queue_pending", info.Pending).
				Int("queue_scheduled", info.Scheduled)
		} else {
			event = event.
				Err(err).
				Str("queue_info", "unavailable")
		}

		event.Msg("task execution started")

	}
}

func logTaskError(task *asynq.Task, err error) {
	log.Error().
		Err(err).
		Str("task_type", task.Type()).
		Msg("TASK FAILED")
}

func logTaskSuccess(task *asynq.Task) {
	log.Info().
		Str("task_type", task.Type()).
		Msg("TASK COMPLETED")
}

func taskLogger(task *asynq.Task, opts QueueLifecycleOptions) zerolog.Logger {
	l := log.With().
		Str("task_type", task.Type()).
		Str("task_id", task.ResultWriter().TaskID()).
		Str("queue", opts.QueueName).
		Logger()

	return l
}

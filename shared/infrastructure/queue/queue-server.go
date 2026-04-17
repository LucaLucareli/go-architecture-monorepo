package queue

import (
	internalAsynq "shared/infrastructure/queue/asynq"

	"github.com/hibiken/asynq"
)

func NewAsynqServer(cfg internalAsynq.AsynqConfig) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     cfg.Addr,
			Password: cfg.Password,
			DB:       cfg.DB,
		},
		asynq.Config{
			Concurrency: 2,
			Queues: map[string]int{
				"reports": 5,
			},
			Logger: &AsynqLogger{},
		},
	)
}

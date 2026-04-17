package asynq

import (
	"time"

	"github.com/hibiken/asynq"
)

type AsynqConfig struct {
	Addr     string
	Password string
	DB       int
}

func NewAsynqClient(cfg AsynqConfig) *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})
}

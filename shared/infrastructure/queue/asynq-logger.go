package queue

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type AsynqLogger struct{}

func (l *AsynqLogger) Debug(args ...interface{}) {
	log.Debug().Msg(fmt.Sprint(args...))
}

func (l *AsynqLogger) Info(args ...interface{}) {
	log.Info().Msg(fmt.Sprint(args...))
}

func (l *AsynqLogger) Warn(args ...interface{}) {
	log.Warn().Msg(fmt.Sprint(args...))
}

func (l *AsynqLogger) Error(args ...interface{}) {
	log.Error().Msg(fmt.Sprint(args...))
}

func (l *AsynqLogger) Fatal(args ...interface{}) {
	log.Fatal().Msg(fmt.Sprint(args...))
}

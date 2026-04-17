package main

import (
	"log"
	"os"

	"employee-api/modules"
	"shared"
	"shared/infrastructure/queue"
	"shared/pkg/helpers"
	"shared/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	dbPostgresURL := os.Getenv("DATABASE_URL")
	dbRedisURL := os.Getenv("REDIS_URL")

	appState := shared.NewAppState(dbPostgresURL, dbRedisURL)

	logger.Init("WORKER", logger.ColorCyan, "DEV")

	mux := queue.NewMux(
		queue.QueueLifecycleOptions{
			QueueName: "reports",
			Inspector: appState.AsynqInspector,
		},
		modules.NewReportModule(appState).GetHandlers()...,
	)

	if err := appState.AsynqServer.Run(mux); err != nil {
		log.Fatal(err)
	}
}

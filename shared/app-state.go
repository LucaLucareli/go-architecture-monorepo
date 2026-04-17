package shared

import (
	"log"
	"os"

	"shared/application/auth"
	"shared/domain/repositories"
	"shared/infrastructure/persistence/postgres/database"
	postgresRepo "shared/infrastructure/persistence/postgres/repositories"
	"shared/infrastructure/queue"
	"shared/infrastructure/queue/asynq"
	"shared/pkg/helpers"

	hibikenAsynq "github.com/hibiken/asynq"
)

type AppState struct {
	AuthService    *auth.AuthService
	UserRepo       repositories.UsersRepository
	AsynqClient    *hibikenAsynq.Client
	AsynqConfig    asynq.AsynqConfig
	AsynqServer    *hibikenAsynq.Server
	AsynqInspector *hibikenAsynq.Inspector
}

func NewAppState(dbPostgresURL, dbRedisURL string) *AppState {

	clientPostgres, err := database.NewEntClient(dbPostgresURL)
	if err != nil {
		log.Fatalf("erro ao conectar no banco: %v", err)
	}

	userRepo := postgresRepo.NewUsersRepository(clientPostgres)

	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	if accessSecret == "" {
		log.Fatal("JWT_ACCESS_SECRET não definido")
	}

	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if refreshSecret == "" {
		log.Fatal("JWT_REFRESH_SECRET não definido")
	}

	accessExpiryHours := helpers.GetEnv("JWT_ACCESS_EXPIRY_HOURS", 1)
	refreshExpiryDays := helpers.GetEnv("JWT_REFRESH_EXPIRY_DAYS", 7)

	authSvc := &auth.AuthService{
		AccessSecret:      accessSecret,
		RefreshSecret:     refreshSecret,
		AccessExpiryHours: accessExpiryHours,
		RefreshExpiryDays: refreshExpiryDays,
		UserRepo:          userRepo,
	}

	asynqCfg := asynq.AsynqConfig{
		Addr:     helpers.GetEnv("REDIS_ADDR", dbRedisURL),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       helpers.GetEnv("REDIS_DB", 0),
	}

	asynqClient := asynq.NewAsynqClient(asynqCfg)

	asynqServer := queue.NewAsynqServer(asynqCfg)

	inspector := hibikenAsynq.NewInspector(hibikenAsynq.RedisClientOpt{
		Addr: dbRedisURL,
	})

	return &AppState{
		AuthService:    authSvc,
		UserRepo:       userRepo,
		AsynqClient:    asynqClient,
		AsynqConfig:    asynqCfg,
		AsynqServer:    asynqServer,
		AsynqInspector: inspector,
	}
}

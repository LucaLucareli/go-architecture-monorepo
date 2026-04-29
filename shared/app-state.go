package shared

import (
	"log"
	"os"
	"sync"

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
	dbPostgresURL string
	dbRedisURL    string

	authService *auth.AuthService
	authOnce    sync.Once

	userRepo    repositories.UsersRepository
	repoOnce    sync.Once

	asynqClient *hibikenAsynq.Client
	clientOnce  sync.Once

	asynqConfig asynq.AsynqConfig
	configOnce  sync.Once

	asynqServer *hibikenAsynq.Server
	serverOnce  sync.Once

	asynqInspector *hibikenAsynq.Inspector
	inspectorOnce  sync.Once
}

func NewAppState(dbPostgresURL, dbRedisURL string) *AppState {
	return &AppState{
		dbPostgresURL: dbPostgresURL,
		dbRedisURL:    dbRedisURL,
	}
}

func (s *AppState) AuthService() *auth.AuthService {
	s.authOnce.Do(func() {
		userRepo := s.UserRepo()

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

		s.authService = &auth.AuthService{
			AccessSecret:      accessSecret,
			RefreshSecret:     refreshSecret,
			AccessExpiryHours: accessExpiryHours,
			RefreshExpiryDays: refreshExpiryDays,
			UserRepo:          userRepo,
		}
	})
	return s.authService
}

func (s *AppState) UserRepo() repositories.UsersRepository {
	s.repoOnce.Do(func() {
		clientPostgres, err := database.NewEntClient(s.dbPostgresURL)
		if err != nil {
			log.Fatalf("erro ao conectar no banco: %v", err)
		}
		s.userRepo = postgresRepo.NewUsersRepository(clientPostgres)
	})
	return s.userRepo
}

func (s *AppState) AsynqConfig() asynq.AsynqConfig {
	s.configOnce.Do(func() {
		s.asynqConfig = asynq.AsynqConfig{
			Addr:     helpers.GetEnv("REDIS_ADDR", s.dbRedisURL),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       helpers.GetEnv("REDIS_DB", 0),
		}
	})
	return s.asynqConfig
}

func (s *AppState) AsynqClient() *hibikenAsynq.Client {
	s.clientOnce.Do(func() {
		s.asynqClient = asynq.NewAsynqClient(s.AsynqConfig())
	})
	return s.asynqClient
}

func (s *AppState) AsynqServer() *hibikenAsynq.Server {
	s.serverOnce.Do(func() {
		s.asynqServer = queue.NewAsynqServer(s.AsynqConfig())
	})
	return s.asynqServer
}

func (s *AppState) AsynqInspector() *hibikenAsynq.Inspector {
	s.inspectorOnce.Do(func() {
		s.asynqInspector = hibikenAsynq.NewInspector(hibikenAsynq.RedisClientOpt{
			Addr: s.dbRedisURL,
		})
	})
	return s.asynqInspector
}

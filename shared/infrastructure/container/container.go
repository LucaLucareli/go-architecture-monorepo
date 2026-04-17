package container

import (
	"shared/infrastructure/persistence/postgres/database"
	"shared/infrastructure/persistence/postgres/ent"

	"github.com/quantumsheep/plouf"
)

type DatabaseService struct {
	plouf.Service
	Client *ent.Client
}

func NewDatabaseService(dbURL string) (*DatabaseService, error) {
	client, err := database.NewEntClient(dbURL)
	if err != nil {
		return nil, err
	}

	return &DatabaseService{
		Client: client,
	}, nil
}

type MainModule struct {
	plouf.Module
	DB *DatabaseService
}

func Build(dbURL string) (*plouf.Worker, *MainModule, error) {
	dbService, err := NewDatabaseService(dbURL)
	if err != nil {
		return nil, nil, err
	}

	mainModule := &MainModule{
		DB: dbService,
	}

	worker, err := plouf.NewWorker(mainModule)
	if err != nil {
		return nil, nil, err
	}

	if err := worker.Inject(mainModule.DB); err != nil {
		return nil, nil, err
	}

	return worker, mainModule, nil
}

package container

import (
	"shared/infrastructure/persistence/postgres/database"
	"shared/infrastructure/persistence/postgres/ent"
	"sync"

	"github.com/quantumsheep/plouf"
)

type DatabaseService struct {
	plouf.Service
	client *ent.Client
	dbURL  string
	once   sync.Once
	err    error
}

// NewDatabaseService cria o serviço mas não conecta ao banco imediatamente (Lazy).
func NewDatabaseService(dbURL string) *DatabaseService {
	return &DatabaseService{
		dbURL: dbURL,
	}
}

// Client retorna a conexão com o banco, inicializando-a apenas na primeira chamada.
func (s *DatabaseService) Client() (*ent.Client, error) {
	s.once.Do(func() {
		s.client, s.err = database.NewEntClient(s.dbURL)
	})
	return s.client, s.err
}

type MainModule struct {
	plouf.Module
	DB *DatabaseService
}

func Build(dbURL string) (*plouf.Worker, *MainModule, error) {
	dbService := NewDatabaseService(dbURL)

	mainModule := &MainModule{
		DB: dbService,
	}

	worker, err := plouf.NewWorker(mainModule)
	if err != nil {
		return nil, nil, err
	}

	// Injeta o serviço; a conexão real só ocorrerá quando algo chamar s.DB.Client()
	if err := worker.Inject(mainModule.DB); err != nil {
		return nil, nil, err
	}

	return worker, mainModule, nil
}

package database

import (
	"context"
	"log"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"

	"shared/infrastructure/persistence/postgres/ent"
)

func NewEntClient(dbURL string) (*ent.Client, error) {
	drv, err := sql.Open(dialect.Postgres, dbURL)
	if err != nil {
		return nil, err
	}

	client := ent.NewClient(ent.Driver(drv))

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Println("schema create error:", err)
	}

	return client, nil
}

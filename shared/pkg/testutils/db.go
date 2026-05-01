package testutils

import (
	"context"
	"fmt"
	"os"
	"shared/infrastructure/persistence/postgres/ent"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type TestingT interface {
	Error(args ...any)
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	FailNow()
	Helper()
	Log(args ...any)
}

func GetTestClient(t TestingT) *ent.Client {
	t.Helper()

	if os.Getenv("TEST_USE_POSTGRES") == "true" {
		return getPostgresClient(t)
	}

	return getSQLiteClient(t)
}

func getSQLiteClient(t TestingT) *ent.Client {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", uuid.NewString())
	client, err := ent.Open("sqlite3", dsn)
	if err != nil {
		t.Fatalf("failed to open sqlite3: %v", err)
	}
	if err := client.Schema.Create(context.Background()); err != nil {
		t.Fatalf("failed to create sqlite3 schema: %v", err)
	}
	return client
}

func getPostgresClient(t TestingT) *ent.Client {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://sa:YourStrong@Passw0rd@localhost:5432/api-golang-test?sslmode=disable"
	}

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		t.Fatalf("failed to sync postgres schema: %v", err)
	}

	return client
}

func CleanupDatabase(t TestingT, client *ent.Client) {
	t.Helper()
	ctx := context.Background()

	// Using ent's generated delete builders is dialect-agnostic and safe.
	// We delete in order to respect foreign key constraints.
	client.UsersOnAccessGroups.Delete().ExecX(ctx)
	client.User.Delete().ExecX(ctx)
	client.Business.Delete().ExecX(ctx)
	client.UserStatus.Delete().ExecX(ctx)
	client.AccessGroup.Delete().ExecX(ctx)
}

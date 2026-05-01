package testutils

import (
	"fmt"
	"shared/infrastructure/persistence/postgres/ent"
	"shared/infrastructure/persistence/postgres/ent/enttest"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type TestingT interface {
	Error(args ...any)
	FailNow()
	Helper()
}

func GetTestClient(t TestingT) *ent.Client {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", uuid.NewString())
	return enttest.Open(t, "sqlite3", dsn)
}

package postgres

import (
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	dbpg "github.com/kr0uch/studspace/pkg/database/postgres"
	"github.com/stretchr/testify/require"
)

func newSQLMock(t *testing.T) (*dbpg.PostgresDB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	return &dbpg.PostgresDB{DB: sqlx.NewDb(db, "sqlmock")}, mock
}

func matchQuery(q string) string {
	return regexp.QuoteMeta(strings.TrimSpace(q))
}

package pgtest

import (
	"fmt"
	"io/fs"
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/pg/migrations"
)

func TestMigrationsAreReversible(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres integration tests are disabled in short mode")
	}

	dsn, err := shared()
	require.NoError(t, err)

	ctx := t.Context()

	name := fmt.Sprintf("migrate_cycle_%d", counter.Add(1))
	require.NoError(t, createDatabase(ctx, dsn, name, ""))

	db, err := openDB(dsn, name)
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	sub, err := fs.Sub(migrations.FS, "sql")
	require.NoError(t, err)

	p, err := goose.NewProvider(goose.DialectPostgres, db, sub)
	require.NoError(t, err)

	_, err = p.Up(ctx)
	require.NoError(t, err)

	_, err = p.DownTo(ctx, 0)
	require.NoError(t, err)

	var n int
	row := db.QueryRowContext(ctx, `
		SELECT count(*)
		FROM pg_tables
		WHERE schemaname = 'public'
			AND tablename <> 'goose_db_version'
	`)
	require.NoError(t, row.Scan(&n))
	require.Zero(t, n, "after complete rollback tables remained")

	_, err = p.Up(ctx)
	require.NoError(t, err)
}

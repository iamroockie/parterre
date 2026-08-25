package postgrestest

import (
	"database/sql"
	"io/fs"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	// driver pgx for database/sql, required by goose
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/iamroockie/parterre/internal/platform/postgres/migrations"
)

func NewTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	if testing.Short() {
		t.Skip("postgres integration tests are disabled in short mode")
	}

	ctx := t.Context()

	container, err := postgres.Run(
		ctx,
		"postgres:18.4-trixie",
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		),
	)
	testcontainers.CleanupContainer(t, container)
	require.NoError(t, err)

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)

	t.Cleanup(func() {
		pool.Close()
	})

	require.NoError(t, pool.Ping(ctx))

	applyMigrations(t, dsn)

	return pool
}

func applyMigrations(t *testing.T, dsn string) {
	t.Helper()

	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	sub, err := fs.Sub(migrations.FS, "sql")
	require.NoError(t, err)

	p, err := goose.NewProvider(goose.DialectPostgres, db, sub)
	require.NoError(t, err)

	_, err = p.Up(t.Context())
	require.NoError(t, err)
}

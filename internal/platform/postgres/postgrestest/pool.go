package postgrestest

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/iamroockie/parterre/internal/platform/postgres/migrations"
)

const templateDB = "parterre_template"

//nolint:gochecknoglobals // one container per test binary
var (
	shared  = sync.OnceValues(startShared)
	counter atomic.Int64
)

func NewTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	if testing.Short() {
		t.Skip("postgres integration tests are disabled in short mode")
	}

	dsn, err := shared()
	require.NoError(t, err)

	ctx := t.Context()

	name := fmt.Sprintf("test_%d", counter.Add(1))
	require.NoError(t, createDatabase(ctx, dsn, name, templateDB))

	cfg, err := pgxpool.ParseConfig(dsn)
	require.NoError(t, err)
	cfg.ConnConfig.Database = name

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	require.NoError(t, err)
	t.Cleanup(pool.Close)
	require.NoError(t, pool.Ping(ctx))

	return pool
}

func startShared() (string, error) {
	ctx := context.Background()

	dsn, err := createContainer(ctx)
	if err != nil {
		return "", err
	}

	if err := createDatabase(ctx, dsn, templateDB, ""); err != nil {
		return "", err
	}

	if err := migrateTemplate(ctx, dsn); err != nil {
		return "", err
	}

	return dsn, nil
}

func createContainer(ctx context.Context) (string, error) {
	container, err := postgres.Run(
		ctx,
		"postgres:18.4-trixie",
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.
				ForLog("database system is ready to accept connections").
				WithOccurrence(2),
		),
	)
	if err != nil {
		return "", fmt.Errorf("run container: %w", err)
	}

	return container.ConnectionString(ctx, "sslmode=disable")
}

func openDB(dsn, database string) (*sql.DB, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse connection string: %w", err)
	}

	if database != "" {
		cfg.ConnConfig.Database = database
	}

	return stdlib.OpenDB(*cfg.ConnConfig), nil
}

func createDatabase(ctx context.Context, dsn, name, template string) error {
	db, err := openDB(dsn, "")
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	stmt := fmt.Sprintf("CREATE DATABASE %q", name)
	if template != "" {
		stmt += fmt.Sprintf(" TEMPLATE %q", template)
	}

	if _, err := db.ExecContext(ctx, stmt); err != nil {
		return fmt.Errorf("create database %s: %w", name, err)
	}

	return nil
}

func migrateTemplate(ctx context.Context, dsn string) error {
	db, err := openDB(dsn, templateDB)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	sub, err := fs.Sub(migrations.FS, "sql")
	if err != nil {
		return fmt.Errorf("sub fs: %w", err)
	}

	p, err := goose.NewProvider(goose.DialectPostgres, db, sub)
	if err != nil {
		return fmt.Errorf("new provider: %w", err)
	}

	if _, err := p.Up(ctx); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}

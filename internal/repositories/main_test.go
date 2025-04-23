package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	pgdriver "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"testing"
	"time"
)

var testQueries *Queries
var ctx = context.Background()
var dbName = "postgres"
var dbUser = "postgres"
var dbPassword = "postgres"

func TestMain(m *testing.M) {
	postgresContainer := createPostgresContainer()

	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	db := createPostgresConnection(postgresContainer)
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("failed to close database: %s", err)
		}
	}(db)

	runMigration(db)

	testQueries = New(db)
	m.Run()
}

func createPostgresContainer() *postgres.PostgresContainer {
	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}
	return postgresContainer
}

func createPostgresConnection(postgresContainer *postgres.PostgresContainer) *sql.DB {
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open database: %s", err)
	}
	return db
}

func runMigration(db *sql.DB) {
	driver, err := pgdriver.WithInstance(db, &pgdriver.Config{})
	if err != nil {
		log.Fatalf("failed to create postgres driver: %s", err)
	}

	mg, err := migrate.NewWithDatabaseInstance(
		"file://../../database/migrations",
		dbName,
		driver,
	)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %s", err)
	}

	if err := mg.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("failed to run migrations: %s", err)
	}
}

package v1

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/pgsqlc"
	service "github.com/horiondreher/go-web-api-boilerplate/internal/domain/services"
	"github.com/horiondreher/go-web-api-boilerplate/internal/utils"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testUserService *service.ServiceManager

func TestMain(m *testing.M) {
	ctx := context.Background()

	utils.SetConfigFile("../../../../.env")
	if os.Getenv("MIGRATION_URL") == "" {
		os.Setenv("MIGRATION_URL", "file://db/postgres/migration")
	}
	if os.Getenv("DB_SOURCE") == "" {
		os.Setenv("DB_SOURCE", "postgresql://pguser:pgpassword@localhost:5432/go_boilerplate?sslmode=disable")
	}
	if os.Getenv("TOKEN_SYMMETRIC_KEY") == "" {
		os.Setenv("TOKEN_SYMMETRIC_KEY", "V7mQ9xL2pR8kN4tZc1Hf6Wb3sD0yJ5uA")
	}
	config := utils.GetConfig()

	migrationsPath := filepath.Join("..", "..", "..", "..", "db", "postgres", "migration", "*.up.sql")
	upMigrations, err := filepath.Glob(migrationsPath)
	if err != nil {
		log.Fatalf("cannot find up migrations: %v", err)
	}

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16.2"),
		postgres.WithInitScripts(upMigrations...),
		postgres.WithDatabase(config.DBName),
		postgres.WithUsername(config.DBUser),
		postgres.WithPassword(config.DBPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("cannot start postgres container: %v", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("cannot get connection string: %v", err)
	}

	conn, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	testStore := pgsqlc.New(conn)
	testUserService = service.NewServiceManager(testStore)

	os.Exit(m.Run())
}

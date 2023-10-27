package test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/driif/echo-go-starter/internal/server/config"
	dbutils "github.com/driif/echo-go-starter/pkg/db"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// TestDB is a test database instance
type TestDB struct {
	*sql.DB
	Container testcontainers.Container
	t         *testing.T
}

// NewDbInstance creates a new test database instance
func NewDBInstance(t *testing.T, conf config.Server) *TestDB {
	t.Helper()

	db := &TestDB{
		t: t,
	}

	req := testcontainers.ContainerRequest{
		Image: "postgres",
		Env: map[string]string{
			"POSTGRES_PASSWORD": conf.Database.Password,
			"POSTGRES_USER":     conf.Database.Username,
			"POSTGRES_DB":       conf.Database.Database,
		},
		NetworkMode: "bridge",
		// Use WithExposedPorts instead of WithWaitStrategy
		ExposedPorts: []string{"5432/tcp"},
		// Use ForListeningPort instead of ForLog
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	var err error
	db.Container, err = testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}

	mappedPort, err := db.Container.MappedPort(context.Background(), "5432")
	if err != nil {
		t.Fatalf("Could not map port for container: %v", err)
	}

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.Host,
		mappedPort.Port(),
		conf.Database.Database,
	)

	db.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// Close closes the test database instance
func (db *TestDB) Close() {
	db.t.Helper()

	err := db.Container.Terminate(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	err = db.DB.Close()
	if err != nil {
		log.Fatal(err)
	}

	db = nil
}

// ApplyFixtures applies the migrations and test fixtures to the test database instance
func (db *TestDB) ApplyFixtures(t *testing.T) {
	db.t.Helper()

	migrations := &migrate.FileMigrationSource{
		Dir: config.DatabaseMigrationFolder,
	}

	_, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		db.t.Fatal(err)
	}

	ctx := context.Background()

	inserts := Inserts()

	// insert test fixtures in an auto-managed db transaction
	err = dbutils.WithTransaction(ctx, db.DB, func(tx boil.ContextExecutor) error {
		t.Helper()
		for _, fixture := range inserts {
			if err := fixture.Insert(ctx, tx, boil.Infer()); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// WithTestDB creates a new test database instance and applies the migrations and test fixtures to it
func WithTestDB(t *testing.T, closure func(db *sql.DB)) {
	t.Helper()

	conf := config.DefaultServiceConfigFromEnv()

	testDB := NewDBInstance(t, conf)

	testDB.ApplyFixtures(t)

	closure(testDB.DB)

	require.NotPanics(t, func() { testDB.Close() })

}

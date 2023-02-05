package test

import (
	"context"
	"testing"
	"time"

	"github.com/driif/echo-go-starter/internal/db/mongodb"
	"github.com/driif/echo-go-starter/internal/server/config/dbconf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitMongoDB initializes a new mongodb adapter for testing
func InitMongoDB(t *testing.T) *mongodb.MongoAdp {
	t.Helper()
	var (
		err  error
		opts options.ClientOptions
	)

	adapter := mongodb.NewAdapter()

	ctx := context.Background()
	adapter.Context = ctx

	conf := dbconf.ConfigMongo()

	opts.ApplyURI(conf.URI)
	opts.SetTimeout(300 * time.Millisecond)
	if err = opts.Validate(); err != nil {
		t.Fatalf("failed to validate mongo options: %v", err)
	}
	adapter.Conn, err = mongo.Connect(adapter.Context, &opts)
	adapter.DBName = "test"
	adapter.DB = adapter.Conn.Database(adapter.DBName)
	if err != nil {
		t.Fatalf("failed to connect to mongo: %v", err)
	}
	if err = adapter.Conn.Ping(adapter.Context, nil); err != nil {
		t.Fatalf("failed to ping mongo: %v", err)
	}

	return adapter
}

func CloseMongoDB(t *testing.T, adapter *mongodb.MongoAdp) {
	t.Helper()
	// drop the test database
	if err := adapter.DB.Drop(adapter.Context); err != nil {
		t.Fatalf("failed to drop test database: %v", err)
	}
	// close the connection
	if err := adapter.Close(); err != nil {
		t.Fatalf("failed to close mongo: %v", err)
	}
}

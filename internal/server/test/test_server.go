package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/driif/echo-go-starter/internal/api/router"
	"github.com/driif/echo-go-starter/internal/db/mongodb"
	"github.com/driif/echo-go-starter/internal/server"
	"github.com/driif/echo-go-starter/internal/server/config"
)

// E2e is a helper function to run end to end tests
func E2e(t *testing.T, closure func(s *server.Server)) {
	// Code here
	t.Helper()
	conf := config.DefaultServiceConfigFromEnv()
	fmt.Println(conf)

	execClosureNewTestServer(context.Background(), t, &conf, closure)
}

// Executes closure on a new test server with a pre-provided database
func execClosureNewTestServer(ctx context.Context, t *testing.T, config *config.Server, closure func(s *server.Server)) {
	t.Helper()

	// https://stackoverflow.com/questions/43424787/how-to-use-next-available-port-in-http-listenandserve
	// You may use port 0 to indicate you're not specifying an exact port but you want a free, available port selected by the system
	config.Echo.ListenAddress = ":0"

	s := server.New(config)

	// attach test adapter to server.store
	s.Store = InitMongoDB(t)
	if err := s.Initialize(); err != nil {
		t.Fatalf("failed to initialize server: %v", err)
	}

	router.InitGroups(s)
	router.AttachRoutes(s)

	closure(s)

	// close the database connection
	CloseMongoDB(t, s.Store.(*mongodb.MongoAdp))

	// echo is managed and should close automatically after running the test
	if err := s.Echo.Shutdown(ctx); err != nil {
		t.Fatalf("failed to shutdown server: %v", err)
	}

	// disallow any further refs to managed object after running the test
	s = nil
}

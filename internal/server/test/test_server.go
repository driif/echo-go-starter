package test

import (
	"context"
	"testing"

	"github.com/driif/echo-go-starter/internal/api/router"
	"github.com/driif/echo-go-starter/internal/server"
	"github.com/driif/echo-go-starter/internal/server/config"
)

// E2e is a helper function to run end to end tests
func E2e(t *testing.T, closure func(s *server.Server)) {
	// Code here
	t.Helper()
	conf := config.DefaultServiceConfigFromEnv()
	conf.Echo.ListenAddress = ":0"

	testDB := NewDBInstance(t, conf)
	testDB.ApplyFixtures(t)

	s := server.New(conf)
	s.DB = testDB.DB

	router.InitGroups(s)
	router.AttachRoutes(s)

	closure(s)

	// echo is managed and should close automatically after running the test
	if err := s.Echo.Shutdown(context.TODO()); err != nil {
		t.Fatalf("failed to shutdown server: %v", err)
	}

	testDB.Close()
	// disallow any further refs to managed object after running the test
	s = nil

}

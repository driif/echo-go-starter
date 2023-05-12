package test

// import (
// 	"testing"
// 	"time"

// 	"github.com/driif/echo-go-starter/internal/server"
// 	"github.com/google/uuid"
// )

// // FixtureMap represents the main definition which fixtures are available though Fixtures()
// type FixtureMap struct {
// 	User1 *model.User
// }

// // Fixtures returns a function wrapping our fixtures, which tests are allowed to manipulate.
// // Each test (which may run concurrently) receives a fresh copy, preventing side effects between test runs.
// func Fixtures(t *testing.T, s *server.Server) FixtureMap {
// 	t.Helper()

// 	now := time.Now()
// 	f := FixtureMap{}

// 	f.User1 = &model.User{
// 		ID:        uuid.MustParse("9e821ac3-ce61-497f-bbf7-21d04058b044"),
// 		Email:     "declar@dev.de",
// 		LastName:  "Declar",
// 		FirstName: "Tristan",
// 		ShortName: "Teclar",
// 		Active:    true,
// 		UpdatedAt: now.String(),
// 		CreatedAt: now.String(),
// 	}

// 	return f.Insert(s)
// }

// func (f FixtureMap) Insert(s *server.Server) FixtureMap {
// 	return f
// }

package test

import (
	"context"
	"fmt"

	"github.com/driif/echo-go-starter/pkg/structs"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// Insertable represents a common IntFromerface for all model instances so they may be inserted via the Inserts() func
type Insertable interface {
	Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error
}

// FixtureMap definition which fixtures are available through Fixtures().
// Mind the declaration order! The fields get inserted exactly in the order they are declared.
type FixtureMap struct {
}

// Fixtures returns a function wrapping our fixtures, which tests are allowed to manipulate.
// Each test (which may run concurrently) receives a fresh copy, preventing side effects between test runs.
func Fixtures() FixtureMap {
	f := FixtureMap{}

	return f
}

// Inserts defines the order in which the fixtures will be inserted
// IntFromo the test database
func Inserts() []Insertable {
	fix := Fixtures()
	insertableIfc := (*Insertable)(nil)
	inserts, err := structs.GetFieldsImplementing(&fix, insertableIfc)
	if err != nil {
		panic(fmt.Errorf("failed to get insertable fixture fields: %w", err))
	}

	return inserts
}

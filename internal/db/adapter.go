package db

import (
	"github.com/driif/echo-go-starter/internal/db/iface"
	"github.com/driif/echo-go-starter/internal/db/model"
	"github.com/driif/echo-go-starter/internal/db/mongodb"
	"github.com/driif/echo-go-starter/internal/server/config/dbconf"
	"github.com/google/uuid"
)

// Adapter is the interface that wraps the basic methods of a database adapter.
type Adapter interface {
	// General

	// Open and configure the adapter
	Open(conf *dbconf.MongoDB) error
	// Close the adapter
	Close() error
	// IsOpen checks if the adapter is ready for use
	IsOpen() bool
	// GetDbVersion returns current db version
	GetDbVersion() (int, error)
	// CheckDbVersion checks if the db version is compatible with the adapter
	CheckDbVersion() error
	// GetName() returns the name of the adapter
	GetName() string
	// SetMaxResults sets the maximum number of results to return in a single Db call.
	SetMaxResults(max int) error
	// CreateDb creates the database optionally dropping it first.
	CreateDb(dropFirst bool) error
	// UpgradeDb upgrades the database to the current adapter version.
	UpgradeDb() error
	// Version return Adapter version
	Version(conf *dbconf.MongoDB) int
	// Stats return DB connection stats object.
	Stats() interface{}

	// Ultimative
	Insert(model iface.Entity) error
	InsertMany(models []iface.Entity) error
	Update(id uuid.UUID, model iface.Entity) error
	Exists(id uuid.UUID, model iface.Entity) (bool, error)
	Delete(id uuid.UUID, model iface.Entity) error
	SoftDelete(id uuid.UUID, model iface.Entity) error

	// Users
	User(id uuid.UUID) (*model.User, error) // Returns a new User object
	Users() ([]model.User, error)           // Returns all Users
	UserExists(id uuid.UUID) (bool, error)  // Returns true if the User exists
	UserCreate(u *model.User) error         // Create a new User
	UserUpdate(u *model.User) error         // Update an existing User
	UserDelete(id uuid.UUID) error          // Delete an existing User

	// Rooms
	CreateRoom(r *mongodb.Room) error
}

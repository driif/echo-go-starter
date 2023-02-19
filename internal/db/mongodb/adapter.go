package mongodb

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/driif/echo-go-starter/internal/db/iface"
	"github.com/driif/echo-go-starter/internal/db/model"
	"github.com/driif/echo-go-starter/internal/server/config/dbconf"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// adapter holds MongoDB connection data
type MongoAdp struct {
	Conn   *mongo.Client
	DB     *mongo.Database
	DBName string
	// maxResults - max number of results to return
	maxResults int
	// maxMessageResults - max number of messages to return
	maxMessageResults int
	version           int
	Context           context.Context
	userTransactions  bool
}

func NewAdapter() *MongoAdp {
	return &MongoAdp{
		Conn:              nil,
		DB:                nil,
		DBName:            "",
		maxResults:        1024,
		maxMessageResults: 100,
		version:           -1,
		Context:           nil,
		userTransactions:  false,
	}
}

// Open initializes a MongoDB session
func (a *MongoAdp) Open(conf *dbconf.MongoDB) error {
	if a.Conn != nil {
		return errors.New("adapter mongodb is already connected")
	}

	var (
		err  error
		opts options.ClientOptions
	)

	opts.ApplyURI(conf.URI)
	if err = opts.Validate(); err != nil {
		return err
	}
	opts.SetTimeout(300 * time.Millisecond)
	// opts.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	// opts.SetAuth(options.Credential{
	// 	AuthMechanism: "SCRAM-SHA-256",
	// 	AuthSource: "admin",
	// 	Username: conf.Username,
	// 	Password: conf.Password,
	// })

	ctx := context.Background()
	a.Context = ctx
	a.Conn, err = mongo.Connect(a.Context, &opts)
	a.DBName = conf.Database
	a.DB = a.Conn.Database(a.DBName)
	if err != nil {
		return err
	}
	a.version = -1

	if err = a.Conn.Ping(a.Context, nil); err != nil {
		return err
	}

	return nil
}

// Close the MongoDB session
func (a *MongoAdp) Close() error {
	var err error
	if a.Conn != nil {
		err = a.Conn.Disconnect(a.Context)
		a.Conn = nil
		a.DB = nil
		a.DBName = ""
		a.version = -1
	} else {
		err = errors.New("there is no active connection")
	}
	return err
}

// IsOpen checks if the adapter is ready for use
func (a *MongoAdp) IsOpen() bool {
	return a.Conn != nil
}

// GetDbVersion returns current db version
func (a *MongoAdp) GetDbVersion() (int, error) {
	if a.version > 0 {
		return a.version, nil
	}

	var result struct {
		Key string `bson:"_id"`
		Val int    `bson:"val"`
	}
	if err := a.DB.Collection("kvmeta").FindOne(a.Context, bson.M{"_id": "version"}).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			err = errors.New("database not initialized")
		}
		return -1, err
	}
	a.version = result.Val

	return result.Val, nil
}

// CheckDbVersion checks if the db version is compatible with the adapter
func (a *MongoAdp) CheckDbVersion() error {
	version, err := a.GetDbVersion()
	if err != nil {
		return err
	}

	if version != 112 {
		return errors.New("Invalid database version" + strconv.Itoa(version) + ". Expected" + strconv.Itoa(112))
	}
	return nil
}

// Version return Adapter version
func (a *MongoAdp) Version(conf *dbconf.MongoDB) int {
	a.version = conf.AdpVersion
	return conf.AdpVersion
}

// GetName() returns the name of the adapter
func (a *MongoAdp) GetName() string {
	return "mongodb"
}

// SetMaxResults sets the maximum number of results to return in a single Db call.
func (a *MongoAdp) SetMaxResults(max int) error {
	a.maxResults = max
	return nil
}

// Stats return DB connection stats object.
func (a *MongoAdp) Stats() interface{} {
	if a.DB == nil {
		return nil
	}

	var result bson.M
	if err := a.DB.RunCommand(a.Context, bson.D{{Key: "ServerStatus", Value: 1}}, nil).Decode(&result); err != nil {
		return nil
	}

	return result["connections"]
}

// CreateDb creates the database optionally dropping it first.
func (a *MongoAdp) CreateDb(dropFirst bool) error {
	if dropFirst {
		if err := a.DB.Drop(a.Context); err != nil {
			return err
		}
	} else if a.isDbInitialized() {
		return errors.New("database already initialized")
	}
	// Collections (tables) do not need to be created explicitly in MongoDB, will be created on first use, with first write operation.

	indexes := []struct {
		collection string
		field      string
		indexOpts  mongo.IndexModel
	}{
		{
			collection: "rooms",
			field:      "id",
		},
	}

	var err error
	for _, idx := range indexes {
		if idx.field != "" {
			_, err = a.DB.Collection(idx.collection).Indexes().CreateOne(a.Context, mongo.IndexModel{Keys: bson.M{idx.field: 1}})
		} else {
			_, err = a.DB.Collection(idx.collection).Indexes().CreateOne(a.Context, idx.indexOpts)
		}
		if err != nil {
			return err
		}
	}

	// Collection "kvmeta" is used to store key-value pairs.
	// Key in "_id" field.
	// Record current database version.
	if _, err := a.DB.Collection("kvmeta").InsertOne(a.Context, map[string]interface{}{"_id": "version", "val": a.version}); err != nil {
		return err
	}

	return createSystemTopic(a)
}

// UpgradeDb upgrades the database to the current adapter version.
func (a *MongoAdp) UpgradeDb() error {
	return nil
}

func (a *MongoAdp) isDbInitialized() bool {
	var result map[string]int

	findOpts := options.FindOneOptions{Projection: bson.M{"value": 1, "_id": 0}}
	if err := a.DB.Collection("kvmeta").FindOne(a.Context, bson.M{"_id": "version"}, &findOpts).Decode(&result); err != nil {
		return false
	}

	return true
}

func createSystemTopic(a *MongoAdp) error {
	now := time.Now().UTC().Round(time.Millisecond)
	_, err := a.DB.Collection("topics").InsertOne(a.Context, map[string]interface{}{
		"createdat": now,
		"updatedat": now,
		"deletedat": time.Time{},
		"touchedat": now,
		"id":        "sys",
		"owner":     "sys",
		"public":    map[string]interface{}{},
		"private":   map[string]interface{}{},
	})

	return err
}

// == ROOM === //
type Room struct {
	ID      string    `bson:"_id"`
	Created time.Time `bson:"createdat"`
	Updated time.Time `bson:"updatedat"`
	Title   string    `bson:"title, omitempty"`
}

func (a *MongoAdp) CreateRoom(r *Room) error {

	_, err := a.DB.Collection("rooms").InsertOne(a.Context, r)

	return err
}

// === USER === //
func (a *MongoAdp) User(id uuid.UUID) (*model.User, error) {
	cursor := a.DB.Collection("users").FindOne(a.Context, bson.M{"_id": id.String()})
	var user *model.User
	if err := cursor.Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}

func (a *MongoAdp) Users() ([]model.User, error) {
	cursor, err := a.DB.Collection("users").Find(a.Context, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []model.User
	if err := cursor.All(a.Context, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (a *MongoAdp) UserExists(uuid uuid.UUID) (bool, error) {
	count, err := a.DB.Collection("users").CountDocuments(a.Context, bson.M{"_id": uuid})
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (a *MongoAdp) UserCreate(u *model.User) error {
	_, err := a.DB.Collection("users").InsertOne(a.Context, u)
	return err
}

func (a *MongoAdp) UserUpdate(u *model.User) error {
	_, err := a.DB.Collection("users").UpdateOne(a.Context, bson.M{"_id": u.ID.String()}, bson.M{"$set": u})
	return err
}

func (a *MongoAdp) UserDelete(id uuid.UUID) error {
	_, err := a.DB.Collection("users").DeleteOne(a.Context, bson.M{"_id": id})
	return err
}

// === Common Functions === //

// Insert inserts a new record into the database.
func (a *MongoAdp) Insert(model iface.Entity) error {
	_, err := a.DB.Collection(model.GetCollectionName()).InsertOne(a.Context, model)
	return err
}

// InsertMany inserts multiple records into the database.
func (a *MongoAdp) InsertMany(models []iface.Entity) error {
	if len(models) == 0 {
		return nil
	}
	var docs []interface{}
	for _, model := range models {
		docs = append(docs, model)
	}
	_, err := a.DB.Collection(models[0].GetCollectionName()).InsertMany(a.Context, docs)
	return err
}

// Update updates a record in the database.
func (a *MongoAdp) Update(id uuid.UUID, model iface.Entity) error {
	_, err := a.DB.Collection(model.GetCollectionName()).UpdateOne(a.Context, bson.M{"_id": id}, bson.M{"$set": model})
	return err
}

// Exists checks if a record exists in the database.
func (a *MongoAdp) Exists(id uuid.UUID, model iface.Entity) (bool, error) {
	count, err := a.DB.Collection(model.GetCollectionName()).CountDocuments(a.Context, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}

	return false, nil
}

// Delete deletes a record from the database.
func (a *MongoAdp) Delete(id uuid.UUID, model iface.Entity) error {
	_, err := a.DB.Collection(model.GetCollectionName()).DeleteOne(a.Context, bson.M{"_id": id})
	return err
}

// SoftDelete soft deletes a record from the database.
func (a *MongoAdp) SoftDelete(id uuid.UUID, model iface.Entity) error {
	_, err := a.DB.Collection(model.GetCollectionName()).UpdateOne(a.Context, bson.M{"_id": id}, bson.M{"$set": bson.M{"deletedat": time.Now().UTC().Round(time.Millisecond)}})
	return err
}

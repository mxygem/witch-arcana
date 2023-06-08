package witcharcana

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB is a representation of the database in use and provides methods for interaction with stored
// data.
type DB struct {
	ctx    context.Context
	dbloc  string
	conn   *mongo.Client
	dbname string
}

// NewDB returns a pointer to a new DB object with the provided configuration.
func NewDB(ctx context.Context, dbloc, dbname string) *DB {
	return &DB{
		ctx:    ctx,
		dbloc:  dbloc,
		dbname: dbname,
	}
}

// Connect connects to the configured database and stores the returned client.
func (db *DB) Connect() error {
	c, err := connectedClient(db)
	if err != nil {
		return fmt.Errorf("connecting db client: %w", err)
	}

	db.conn = c

	return nil
}

// Disconnect closes all open connections.
func (db *DB) Disconnect() error {
	if err := db.conn.Disconnect(db.ctx); err != nil {
		return fmt.Errorf("disconnecting from db: %w", err)
	}

	return nil
}

// Ping attempts to connect to the configured database to verify connectivity and that the database
// is responding.
func (db *DB) Ping() error {
	var res bson.M
	err := db.conn.Database(db.dbname).
		RunCommand(db.ctx, bson.D{{Key: "ping", Value: 1}}).
		Decode(&res)
	if err != nil {
		return fmt.Errorf("disconnecting db: %w", err)
	}

	return nil
}

func connectedClient(db *DB) (*mongo.Client, error) {
	sAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(db.dbloc).SetServerAPIOptions(sAPI).SetTimeout(3 * time.Second)

	c, err := mongo.Connect(db.ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("connecting to db at: %q: %w", db.dbloc, err)
	}

	return c, nil
}

func (db *DB) Get() error {
	return fmt.Errorf("db get unimplemented")
}

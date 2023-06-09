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
	ctx  context.Context
	cfg  *DBConfig
	conn *mongo.Client
	coll *mongo.Collection
}

type DBConfig struct {
	Loc  string
	Name string
	User string
	Pass string
}

// NewDB returns a pointer to a new DB object with the provided configuration.
func NewDB(ctx context.Context, cfg *DBConfig) *DB {
	return &DB{
		ctx: ctx,
		cfg: cfg,
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
	err := db.conn.Database(db.cfg.Name).
		RunCommand(db.ctx, bson.D{{Key: "ping", Value: 1}}).
		Decode(&res)
	if err != nil {
		return fmt.Errorf("disconnecting db: %w", err)
	}

	return nil
}

func connectedClient(db *DB) (*mongo.Client, error) {
	cred := options.Credential{
		Username: db.cfg.User,
		Password: db.cfg.Pass,
	}
	sAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		SetAuth(cred).
		ApplyURI(db.cfg.Loc).
		SetServerAPIOptions(sAPI).
		SetTimeout(3 * time.Second)

	c, err := mongo.Connect(db.ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("connecting to db at: %q: %w", db.cfg.Loc, err)
	}

	return c, nil
}

// TODO: Generics?
func (db *DB) Get(name string) (any, error) {
	f := bson.D{{Key: "name", Value: name}}

	var res Club
	err := db.coll.FindOne(db.ctx, f, nil).Decode(&res)
	if err != nil {
		return nil, fmt.Errorf("getting club from db: %w", err)
	}

	return &res, nil
}

func (db *DB) Player(club, name string) (*Player, error) {
	f := bson.D{{Key: "name", Value: club}, {Key: "players.$", Value: name}}

	var res Player
	err := db.coll.FindOne(db.ctx, f, nil).Decode(&res)
	if err != nil {
		return nil, fmt.Errorf("getting player from db: %w", err)
	}

	return &res, nil
}

func (db *DB) Create(data any) (any, error) {
	res, err := db.coll.InsertOne(db.ctx, data)
	if err != nil {
		return nil, fmt.Errorf("db insert: %w", err)
	}

	return res.InsertedID, nil
}

func (db *DB) Update(c *Club) error {
	f := bson.D{{Key: "name", Value: c.Name}}
	u := bson.D{{Key: "$set", Value: bson.D{{Key: "players", Value: c.Players}}}}
	o := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var res Club
	err := db.coll.FindOneAndUpdate(db.ctx, f, u, o).Decode(&res)
	if err != nil {
		return fmt.Errorf("db update: %w", err)
	}

	return nil
}

func (db *DB) Delete() error {
	return fmt.Errorf("db delete unimplemented")
}

func (db *DB) setCollection(coll string) {
	db.coll = db.conn.Database(db.cfg.Name).Collection(coll)
}

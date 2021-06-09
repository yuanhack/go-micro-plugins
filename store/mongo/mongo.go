package mongo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
	"unicode"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	log "github.com/macheal/go-micro/v2/logger"
	"github.com/macheal/go-micro/v2/store"
	"github.com/pkg/errors"
)

var (
	// DefaultDatabase is the database that the sql store will use if no database is provided.
	DefaultDatabase = "micro"
	// DefaultTable is the table that the sql store will use if no table is provided.
	DefaultTable = "micro"
)

type MongoDBStore struct {
	store          store.Store
	db             *sql.DB
	client         *mongo.Client
	collection     *mongo.Collection
	listProjection bson.D

	database string
	table    string
	readOne  bool

	options store.Options

	readPrepare, writePrepare, deletePrepare *sql.Stmt
}

func (s *MongoDBStore) Init(opts ...store.Option) error {
	for _, o := range opts {
		o(&s.options)
	}
	// reconfigure
	return s.configure()
}

func (s *MongoDBStore) Options() store.Options {
	return s.options
}

func (s *MongoDBStore) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.client.Disconnect(ctx)
}

// List all the known records
func (s *MongoDBStore) List(opts ...store.ListOption) ([]string, error) {
	var opt store.ListOptions
	for _, o := range opts {
		o(&opt)
	}

	var records []string

	filter := bson.M{}
	if len(opt.Prefix) > 0 && len(opt.Suffix) > 0 {
		filter["key"] = primitive.Regex{Pattern: "^" + opt.Prefix + ".*" + opt.Suffix + "$"}
	} else if len(opt.Prefix) > 0 {
		filter["key"] = primitive.Regex{Pattern: "^" + opt.Prefix}
	} else if len(opt.Suffix) > 0 {
		filter["key"] = primitive.Regex{Pattern: opt.Suffix + "$"}
	}

	//if len(opt.Prefix) > 0 {
	//	filter["key"] = primitive.Regex{Pattern: "^" + opt.Prefix}
	//}
	//if len(opt.Suffix) > 0 {
	//	filter["key"] = primitive.Regex{Pattern: opt.Suffix + "$"}
	//}

	cursor, err := s.collection.Find(
		context.Background(),
		filter,
		options.Find().SetProjection(s.listProjection),
	)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.Background()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		records = append(records, result["key"].(string))
	}
	return records, nil
}

// Read all records with keys
func (s *MongoDBStore) Read(key string, opts ...store.ReadOption) ([]*store.Record, error) {
	var opt store.ReadOptions
	for _, o := range opts {
		o(&opt)
	}

	var records []*store.Record
	filter := bson.M{}

	midSub := ""
	if v := s.options.Context.Value(mid_sub("")); v != nil {
		midSub = v.(string)
	}
	fmt.Println(midSub)
	if opt.Prefix {
		filter["key"] = primitive.Regex{Pattern: "^" + key}
	} else if opt.Suffix {
		filter["key"] = primitive.Regex{Pattern: key + "$"}
	} else {
		filter["key"] = key
	}
	findOpts := []*options.FindOptions{}
	if opt.Limit > 0 {
		findOpts = append(findOpts, options.Find().SetLimit(int64(opt.Limit)))
	}
	if opt.Offset > 0 {
		findOpts = append(findOpts, options.Find().SetSkip(int64(opt.Offset)))
	}

	//options.Find().SetLimit(1),
	cursor, err := s.collection.Find(
		context.Background(),
		filter,
		findOpts...,
	)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.Background()) {
		var result storeModel //bson.M
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		var x time.Duration
		if !result.Expiry.IsZero() {
			x = result.Expiry.Sub(time.Now())
		}
		records = append(records, &store.Record{
			Key:   result.Key,
			Value: []byte(result.Value),
			//Metadata: nil,
			Expiry: x,
		})
	}
	return records, nil
}

type storeModel struct {
	Key    string    `bson:"key,omitempty"`
	Value  string    `bson:"value,omitempty"`
	Expiry time.Time `bson:"expiry,omitempty"`
}

// Write records
func (s *MongoDBStore) Write(r *store.Record, opts ...store.WriteOption) error {
	var opt store.WriteOptions
	for _, o := range opts {
		o(&opt)
	}
	data := storeModel{
		Key:   r.Key,
		Value: string(r.Value),
	}
	if r.Expiry > 0 || opt.TTL > 0 {
		exp := Min(r.Expiry, opt.TTL)
		data.Expiry = time.Now().Local().Add(exp)
	} else if !opt.Expiry.IsZero() {
		data.Expiry = opt.Expiry
	}
	_, err := s.collection.InsertOne(context.Background(), data)
	if err != nil {
		if IsDup(err) {
			filter := bson.D{{"key", r.Key}}
			res1, err1 := s.collection.UpdateOne(context.Background(), filter, bson.D{{"$set", data}})
			if err1 != nil {
				return err1
			}
			if res1.ModifiedCount == 1 {
				return nil
			}
		}
	}
	//ss := s.collection.FindOneAndUpdate(context.Background(), bson.M{"key": r.Key}, data)
	//fmt.Printf("%+v", ss)
	return nil
}

// Delete records with keys
func (s *MongoDBStore) Delete(key string, opts ...store.DeleteOption) error {
	var opt store.DeleteOptions
	for _, o := range opts {
		o(&opt)
	}
	_, err := s.collection.DeleteOne(context.Background(), bson.D{{"key", key}})
	return err
}
func (s *MongoDBStore) initDB() error {
	//_, err := mt.Coll.Indexes().CreateOne(mtest.Background, mongo.IndexModel{Keys: bson.D{{"x", 1}}})
	//Key
	indexView := s.collection.Indexes()
	opts := options.Index().
		SetName("unique_key").
		SetUnique(true)
	_, err := indexView.CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"key", 1}},
		Options: opts,
	})
	if err != nil {
		return err
	}
	//expiry
	opts_exp := options.Index().
		SetName("exp").
		SetExpireAfterSeconds(0)
	_, err = indexView.CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"expiry", 1}},
		Options: opts_exp,
	})
	if err != nil {
		return err
	}
	return nil
}
func (s *MongoDBStore) defaultConfig() {
	s.database = DefaultDatabase
	s.table = DefaultTable
	s.readOne = true
	s.options.Nodes = []string{"localhost:27017"}
}
func (s *MongoDBStore) configure() error {
	s.defaultConfig()

	if v := s.options.Context.Value(uri("")); v != nil {
		s.options.Nodes[0] = v.(string)
	}

	for _, r := range s.database {
		if !unicode.IsLetter(r) {
			return errors.New("store.namespace must only contain letters")
		}
	}
	if len(s.options.Database) > 0 {
		s.database = s.options.Database
	}
	if len(s.options.Table) > 0 {
		s.table = s.options.Table
	}
	source := s.options.Nodes[0]
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(source)) //"mongodb://localhost:27017"))
	if err != nil {
		return err
	}

	//ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	//defer cancel()
	//err = client.Ping(ctx, readpref.Primary())
	//if err != nil {
	//	return err
	//}
	s.client = client
	s.collection = client.Database(s.database).Collection(s.table)

	// initialise the database
	return s.initDB()
}

func (s *MongoDBStore) String() string {
	return "mongodb"
}

// New returns a new micro Store backed by sql
func NewStore(opts ...store.Option) store.Store {
	var opt store.Options
	for _, o := range opts {
		o(&opt)
	}

	// new store
	s := new(MongoDBStore)
	// set the options
	s.options = opt
	s.listProjection = bson.D{
		{"key", 1},
		{"_id", 0},
	}

	// configure the store
	if err := s.configure(); err != nil {
		log.Fatal(err)
	}

	// return store
	return s
}

// IsDup returns whether err informs of a duplicate Key error because
// a primary Key index or a secondary unique index already has an entry
// with the given Value.
func IsDup(err error) bool {
	// Besides being handy, helps with MongoDB bugs SERVER-7164 and SERVER-11493.
	// What follows makes me sad. Hopefully conventions will be more clear over time.
	switch e := err.(type) {
	case mongo.WriteException:
		if len(e.WriteErrors) == 1 {
			er := e.WriteErrors[0]
			return er.Code == 11000 || er.Code == 11001 || er.Code == 12582 || er.Code == 16460 && strings.Contains(er.Error(), " E11000 ")
		}
		//case *QueryError:
		//	return e.Code == 11000 || e.Code == 11001 || e.Code == 12582
		//case *BulkError:
		//	for _, ecase := range e.ecases {
		//		if !IsDup(ecase.Err) {
		//			return false
		//		}
		//	}
		return false
	}
	return false
}

func Min(x, y time.Duration) time.Duration {
	if x < y {
		return x
	}
	return y
}

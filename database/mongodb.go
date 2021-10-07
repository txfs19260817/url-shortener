package database

import (
	"context"
	"time"

	"github.com/txfs19260817/url-shortener/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDB is the database client instance
type MongoDB struct {
	Client          *mongo.Client
	DB              *mongo.Database
	TimeoutDuration time.Duration
	MongoDBConfig
}

// MongoDBConfig loads configs to set up MongoDB instance
type MongoDBConfig struct {
	Uri        string `yaml:"uri"`
	Database   string `yaml:"database"`
	Collection string `yaml:"collection"`
	Timeout    int    `yaml:"timeout"`
}

func NewMongoDB(config MongoDBConfig) (Database, error) {
	// Configure our client to use the correct URI, but we're not yet connecting to it.
	client, err := mongo.NewClient(options.Client().ApplyURI(config.Uri))
	if err != nil {
		return nil, err
	}

	// Define a timeout duration
	timeout := time.Duration(config.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Try to connect using the defined timeout duration
	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	// Ping the cluster to ensure we're already connected
	ctx, cancelPing := context.WithTimeout(context.Background(), timeout)
	defer cancelPing()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	return &MongoDB{Client: client, DB: client.Database(config.Database), TimeoutDuration: time.Duration(config.Timeout) * time.Second, MongoDBConfig: config}, nil
}

func (m *MongoDB) CreateUrl(url *model.Url) error {
	doc, err := bson.Marshal(url)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), m.TimeoutDuration)
	defer cancel()

	if _, err := m.DB.Collection(m.Collection).InsertOne(ctx, doc); err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) ReadUrl(key string) (url *model.Url, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.TimeoutDuration)
	defer cancel()

	err = m.DB.Collection(m.Collection).FindOne(ctx, bson.M{"key": key}).Decode(&url)
	return
}

// KeyExists checks if the given key is already in use
func (m *MongoDB) KeyExists(key string) bool {
	r, err := m.ReadUrl(key)
	if err == mongo.ErrNoDocuments {
		err = nil
	}
	// whether r is nil or not decides the result, however if there is an unexpected type of error, we should return false
	return err == nil && r != nil
}

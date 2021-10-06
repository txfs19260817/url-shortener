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

type MongoDB struct {
	client     *mongo.Client
	database   *mongo.Database
	collection string
	timeout    time.Duration
}

func NewMongoDB(mongoURI, database, collection string, timeoutSec int) (Database, error) {
	// Configure our client to use the correct URI, but we're not yet connecting to it.
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	// Define a timeout duration
	timeout := time.Duration(timeoutSec) * time.Second
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
	return &MongoDB{client, client.Database(database), collection, timeout}, nil
}

func (m *MongoDB) CreateUrl(url *model.Url) error {
	doc, err := bson.Marshal(url)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	if _, err := m.database.Collection(m.collection).InsertOne(ctx, doc); err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) ReadUrl(key string) (url *model.Url, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	err = m.database.Collection(m.collection).FindOne(ctx, bson.M{"key": key}).Decode(&url)
	return
}
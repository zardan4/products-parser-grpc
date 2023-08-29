package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoWrapper struct {
	Client *mongo.Client

	Ctx context.Context
}

func (w *MongoWrapper) Disconnect() error {
	return w.Client.Disconnect(w.Ctx)
}

// init db connection
func InitMongo(connLine string) (MongoWrapper, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Client()
	opts.ApplyURI(connLine)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return MongoWrapper{}, err
	}

	err = client.Ping(ctx, readpref.Primary())

	return MongoWrapper{
		Client: client,
		Ctx:    ctx,
	}, err
}

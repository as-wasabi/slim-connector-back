package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"slim-connector-back/config"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func FetchMongoDB(cfg *config.Config) (*MongoDB, error) {
	uri := cfg.MongoDB.URI

	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// 接続確認
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	database := client.Database(cfg.MongoDB.Database)

	return &MongoDB{
		Client:   client,
		Database: database,
	}, nil
}

func (m *MongoDB) Close() error {
	return m.Client.Disconnect(context.Background())
}

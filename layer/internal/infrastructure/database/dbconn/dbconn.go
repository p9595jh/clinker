package dbconn

import (
	"context"
	"fmt"
	"layer/common/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseConfig struct {
	Host   string
	Port   int
	Schema string
}

func ConnectDB(dbconfig *DatabaseConfig) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf(
		"mongodb://%s:%d",
		dbconfig.Host,
		dbconfig.Port,
	))

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	logger.Info("Database").W("Database connected")

	return client.Database(dbconfig.Schema), nil
}

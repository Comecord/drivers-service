package mongox

import (
	"context"
	"drivers-service/config"
	"drivers-service/pkg/logging"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"sync"
)

var (
	clientInstance      *mongo.Client
	clientDatabase      *mongo.Database
	clientInstanceError error
	mongoOnce           sync.Once
)

func InitMongoClient(conf *config.Config, ctx context.Context, logger logging.Logger) error {
	mongoUrl := fmt.Sprintf(`mongodb://%s:%s@%s:%s/%s?authSource=%s`,
		conf.MongoX.Username, conf.MongoX.Password, conf.MongoX.Host, conf.MongoX.Port,
		conf.MongoX.Database, conf.MongoX.AuthSource)

	mongoconn := options.Client().ApplyURI(mongoUrl)
	var err error
	clientInstance, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		logger.Fatal(logging.MongoDB, logging.Connection, err.Error(), nil)
	}
	return nil
}

func Execute(ctx context.Context, conf *config.Config, operation func(*mongo.Database) error) error {
	if clientInstance == nil {
		return errors.New("MongoDB client is not initialized")
	}

	db := clientInstance.Database(conf.MongoX.Database)
	return operation(db)
}

func Connection(conf *config.Config, ctx context.Context, logger logging.Logger) (*mongo.Database, error) {
	err := InitMongoClient(conf, ctx, logger)
	if err != nil {
		return nil, err
	}
	logger.Info(logging.MongoDB, logging.Connection, "Database connection established.", nil)
	return clientInstance.Database(conf.MongoX.Database), nil
}

func GetMongoClient(conf *config.Config) (*mongo.Database, *mongo.Client, error) {
	mongoUrl := fmt.Sprintf(`mongodb://%s:%s@%s:%s/%s?authSource=%s`,
		conf.MongoX.Username, conf.MongoX.Password, conf.MongoX.Host, conf.MongoX.Port,
		conf.MongoX.Database, conf.MongoX.AuthSource)

	clientOptions := options.Client().ApplyURI(mongoUrl)

	if conf.MongoX.ReplicaSet != "" {
		clientOptions.SetReplicaSet(conf.MongoX.ReplicaSet)
	}
	if conf.MongoX.ReadPreference != "" {
		mode, err := readpref.ModeFromString(conf.MongoX.ReadPreference)
		if err != nil {
			return nil, nil, err
		}
		readPref, err := readpref.New(mode)
		if err != nil {
			return nil, nil, err
		}
		clientOptions.SetReadPreference(readPref)
	}

	mongoOnce.Do(func() {

		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceError = err
			return
		}

		clientDatabase = client.Database(conf.MongoX.Database)
		if clientDatabase == nil {
			clientDatabase = client.Database(conf.MongoX.Database)
		}

		err = client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = err
			return
		}

		clientInstance = client
	})

	return clientDatabase, clientInstance, clientInstanceError
}

func CloseMongoClient() error {
	if clientInstance != nil {
		return clientInstance.Disconnect(context.TODO())
	}
	return nil
}

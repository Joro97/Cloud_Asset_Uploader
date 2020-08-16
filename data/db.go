package data

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Store is an interface for interacting with the applications storage.
type Store interface {
	AddNewAsset(assetName, url string) (id string, err error)
	SetAssetStatus(assetID, status string) (*AssetInfo, error)
	GetAsset(assetID string) (*AssetInfo, error)
}

// DB is a Store implementation, whose underlying DB is MongoDB.
type DB struct {
	Client *mongo.Client
}

// BuildConnectionStringForMongoDB builds a connection string from env vars or returns an error if one of them is missing.
func BuildConnectionStringForMongoDB() (string, error) {
	mongoUser := os.Getenv("MONGO_USERNAME")
	if mongoUser == "" {
		return "", errors.New("no MongoDB username specified")
	}

	mongoPassword := os.Getenv("MONGO_PASSWORD")
	if mongoPassword == "" {
		return "", errors.New("no MongoDB password provided")
	}

	mongoContainer := os.Getenv("MONGO_CONTAINER_NAME")
	if mongoContainer == "" {
		return "", errors.New("no MongoDB container name provided")
	}

	mongoPort := os.Getenv("MONGO_PORT")
	if mongoPort == "" {
		return "", errors.New("no MongoDB port provided")
	}

	connStr := fmt.Sprintf("mongodb://%s:%s@%s:%s", mongoUser, mongoPassword, mongoContainer, mongoPort)
	return connStr, nil
}

// NewDB creates a new DB struct from the given connection string.
func NewDB(connection string) (*DB, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connection))
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 7*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	log.Info().Msg("Successfully connected to MongoDB.")
	return &DB{Client: client}, nil
}

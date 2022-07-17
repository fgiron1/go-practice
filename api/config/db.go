package config

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	once sync.Once
	// Only gain access to this reference through InitDB method
	mongoClient *mongo.Client
)

// Tries to initialize the connection if it's not already up
// If it's up, just returns the client instance
func InitClient(ctx context.Context) *mongo.Client {
	once.Do(func() {
		if possibleClient := connectDB(ctx); possibleClient != nil {
			mongoClient = possibleClient
		}
	})
	return mongoClient
}

func connectDB(ctx context.Context) *mongo.Client {

	// We check whether the connection is already open or not.

	if isClientConnected(ctx, mongoClient) {
		return nil
	}

	// Preparing client configuration

	credentials := options.Credential{
		AuthMechanism: "SCRAM-SHA-256",
		AuthSource:    "admin",
		Username:      "root",
		Password:      "root",
	}

	clientOpts := options.Client().ApplyURI(MongoURI()).
		SetAuth(credentials)

	client, err := mongo.Connect(ctx, clientOpts)

	if err != nil {
		panic(err)
	}

	// Ping the primary to check if the connection is successful.
	isClientConnected(ctx, client)

	return client

}

func isClientConnected(ctx context.Context, client *mongo.Client) bool {

	if mongoClient != nil {
		if err := client.Ping(ctx, readpref.Primary()); err != nil {
			return true
		}
	}
	return false
}

func DisconnectClient() {

	// Panics when trying to disconnect a non-connected client
	if err := mongoClient.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	} else {
		// Shuts down the connection if it's connected
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}
}

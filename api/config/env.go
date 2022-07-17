package config

import (
	"os"
)

func MongoURI() string {
	return os.Getenv("MONGOURI")
}

func MongoDBName() string {

	return os.Getenv("MONGO_DB_NAME")
}

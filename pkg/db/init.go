package db

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/tranhuyducseven/goem-template/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	DB     *mongo.Database

	User *mongo.Collection
)

func Init() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.Env.MONGODB_URI).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	var err error
	client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.M{"ping": 1}).Decode(&result); err != nil {
		panic(err)
	}

	DB = client.Database(config.Env.MONGODB_DATABASE)

	User = DB.Collection("users")
	log.Info().Msg("Connected to MongoDB!")
}

func Close() {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
	log.Info().Msg("Connection to MongoDB closed.")
}

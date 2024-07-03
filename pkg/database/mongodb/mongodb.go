package mongodb

import (
	"context"
	"fmt"
	"log"
	"main/internal/configs"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client, err = mongo.NewClient(options.Client().ApplyURI(configs.MongoURI))

var Collection = client.Database(configs.DBName).Collection(configs.CollectionName)

var Collection2 = client.Database(configs.DBName).Collection(configs.CollectionName2)

func ConnectToMongo() {
	// Create connect
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
}

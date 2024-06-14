package database

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDb struct {
	client      *mongo.Client
	isConnected bool
}

func (mongodb *MongoDb) connect() error {
	mongodbUri := os.Getenv("MONGODB_URI")

	// Connect to the database.
	clientOptions := options.Client().ApplyURI(mongodbUri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	mongodb.client = client

	// Check the connection.
	err = mongodb.client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	fmt.Println("Connected to MongoDB")

	mongodb.isConnected = true

	return nil
}

func (mongodb *MongoDb) BatchSaveFileChunks(fileChunks []FileChunk) error {
	if !mongodb.isConnected {
		err := mongodb.connect()
		if err != nil {
			panic("could not connect to mongodb")
		}
	}

	var documents []interface{}

	for _, fileChunk := range fileChunks {
		documents = append(documents, fileChunk)
	}

	databaseName := os.Getenv("MONGODB_DATABASE_NAME")
	collectionName := os.Getenv("MONGODB_COLLECTION_NAME")

	collection := mongodb.client.Database(databaseName).Collection(collectionName)
	_, err := collection.InsertMany(context.Background(), documents)

	return err
}

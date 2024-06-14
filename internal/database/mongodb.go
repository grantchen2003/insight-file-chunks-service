package database

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
}

func (mongodb *MongoDB) Connect() error {
	mongodbUri := os.Getenv("MONGODB_URI")

	// Connect to the database.
	clientOptions := options.Client().ApplyURI(mongodbUri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	mongodb.client = client

	// Check the connection.
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	fmt.Println("Connected to MongoDB")

	return err
}

func (mongodb *MongoDB) Close() error {
	err := mongodb.client.Disconnect(context.TODO())
	if err != nil {
		return err
	}

	fmt.Println("Connection closed.")

	return err
}

func (mongodb *MongoDB) BatchSaveFileChunks(fileChunks []FileChunk) error {
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

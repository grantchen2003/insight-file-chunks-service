package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDb struct {
	client *mongo.Client
}

func (mongodb *MongoDb) Connect() error {
	mongodbUri := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(mongodbUri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	mongodb.client = client

	log.Println("connected to MongoDB")

	return nil
}

func (mongodb *MongoDb) Close() error {
	if err := mongodb.client.Disconnect(context.TODO()); err != nil {
		return err
	}

	log.Println("connection to MongoDB closed")

	return nil
}

func (mongodb *MongoDb) SaveFileChunk(fileChunk FileChunk) error {
	_, err := mongodb.getCollection().InsertOne(context.Background(), fileChunk)

	return err
}

func (mongodb *MongoDb) BatchSaveFileChunks(fileChunks []FileChunk) error {
	var documents []interface{}

	for _, fileChunk := range fileChunks {
		documents = append(documents, fileChunk)
	}

	_, err := mongodb.getCollection().InsertMany(context.Background(), documents)

	return err
}

func (mongodb *MongoDb) getCollection() *mongo.Collection {
	databaseName := os.Getenv("MONGODB_DATABASE_NAME")
	collectionName := os.Getenv("MONGODB_COLLECTION_NAME")

	return mongodb.client.Database(databaseName).Collection(collectionName)
}
